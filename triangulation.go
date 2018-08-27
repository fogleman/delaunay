package delaunay

import (
	"fmt"
	"math"
	"sort"
)

type Triangulation struct {
	Points    []Point
	Triangles []int
	HalfEdges []int
}

func Triangulate(points []Point) (*Triangulation, error) {
	t := newTriangulator(points)
	if err := t.triangulate(); err != nil {
		return nil, err
	}
	return &Triangulation{points, t.triangles, t.halfedges}, nil
}

type triangulator struct {
	points       []Point
	ids          []int
	center       Point
	hash         []*node
	triangles    []int
	halfedges    []int
	trianglesLen int
}

func newTriangulator(points []Point) *triangulator {
	return &triangulator{points: points}
}

// TODO: don't do it this way
func (a *triangulator) Len() int {
	return len(a.points)
}

func (a *triangulator) Swap(i, j int) {
	a.ids[i], a.ids[j] = a.ids[j], a.ids[i]
}

func (a *triangulator) Less(i, j int) bool {
	p1 := a.points[a.ids[i]]
	p2 := a.points[a.ids[j]]
	d1 := p1.distance(a.center)
	d2 := p2.distance(a.center)
	if d1 != d2 {
		return d1 < d2
	}
	if p1.X != p2.X {
		return p1.X < p2.X
	}
	return p1.Y < p2.Y
}
func (tri *triangulator) triangulate() error {
	points := tri.points

	n := len(points)
	if n == 0 { // TODO: < 3?
		return nil
	}

	tri.ids = make([]int, n)

	// compute bounds
	x0 := points[0].X
	y0 := points[0].Y
	x1 := points[0].X
	y1 := points[0].Y
	for i, p := range points {
		if p.X < x0 {
			x0 = p.X
		}
		if p.X > x1 {
			x1 = p.X
		}
		if p.Y < y0 {
			y0 = p.Y
		}
		if p.Y > y1 {
			y1 = p.Y
		}
		tri.ids[i] = i
	}

	var i0, i1, i2 int

	// pick a seed point close to midpoint
	m := Point{(x0 + x1) / 2, (y0 + y1) / 2}
	minDist := infinity
	for i, p := range points {
		d := p.distance(m)
		if d < minDist {
			i0 = i
			minDist = d
		}
	}

	// find point closest to seed point
	minDist = infinity
	for i, p := range points {
		if i == i0 {
			continue
		}
		d := p.distance(points[i0])
		if d > 0 && d < minDist {
			i1 = i
			minDist = d
		}
	}

	// find the third point which forms the smallest circumcircle
	minRadius := infinity
	for i, p := range points {
		if i == i0 || i == i1 {
			continue
		}
		r := circumradius(points[i0], points[i1], p)
		if r < minRadius {
			i2 = i
			minRadius = r
		}
	}
	if minRadius == infinity {
		return fmt.Errorf("No Delaunay triangulation exists for this input.")
	}

	// swap the order of the seed points for counter-clockwise orientation
	if area(points[i0], points[i1], points[i2]) < 0 {
		i1, i2 = i2, i1
	}

	tri.center = circumcenter(points[i0], points[i1], points[i2])

	// sort the points by distance from the seed triangle circumcenter
	// TODO: try precomputing distances
	sort.Sort(tri)

	// initialize a hash table for storing edges of the advancing convex hull
	hashSize := int(math.Ceil(math.Sqrt(float64(n))))
	tri.hash = make([]*node, hashSize)

	// initialize a circular doubly-linked list that will hold an advancing convex hull
	e := newNode(points[i0], i0, nil)
	e.t = 0
	tri.hashEdge(e)

	e = newNode(points[i1], i1, e)
	e.t = 1
	tri.hashEdge(e)

	e = newNode(points[i2], i2, e)
	e.t = 2
	tri.hashEdge(e)

	maxTriangles := 2*n - 5
	tri.trianglesLen = 0
	tri.triangles = make([]int, maxTriangles*3)
	tri.halfedges = make([]int, maxTriangles*3)

	tri.addTriangle(i0, i1, i2, -1, -1, -1)

	pp := Point{infinity, infinity}
	for k := 0; k < n; k++ {
		i := tri.ids[k]
		p := points[i]

		// skip duplicate points
		if p == pp {
			continue
		}
		pp = p

		// skip seed triangle points
		if i == i0 || i == i1 || i == i2 {
			continue
		}

		// find a visible edge on the convex hull using edge hash
		var start *node
		key := tri.hashKey(p)
		for j := 0; j < len(tri.hash); j++ {
			start = tri.hash[key]
			if start != nil && !start.removed {
				break
			}
			key++
			if key >= len(tri.hash) {
				key = 0
			}
		}

		e := start
		for area(p, e.p, e.next.p) >= 0 {
			e = e.next
			if e == start {
				return fmt.Errorf("Something is wrong with the input points.")
			}
		}

		walkBack := e == start

		// add the first triangle from the point
		t := tri.addTriangle(e.i, i, e.next.i, -1, -1, e.t)
		e.t = t // keep track of boundary triangles on the hull
		e = newNode(points[i], i, e)

		// recursively flip triangles from the point until they satisfy the Delaunay condition
		e.t = tri.legalize(t + 2)
		if e.prev.prev.t == tri.halfedges[t+1] {
			e.prev.prev.t = t + 2
		}

		// walk forward through the hull, adding more triangles and flipping recursively
		q := e.next
		for area(p, q.p, q.next.p) < 0 {
			t = tri.addTriangle(q.i, i, q.next.i, q.prev.t, -1, q.t)
			q.prev.t = tri.legalize(t + 2)
			q.remove()
			q = q.next
		}

		if walkBack {
			// walk backward from the other side, adding more triangles and flipping
			q := e.prev
			for area(p, q.prev.p, q.p) < 0 {
				t = tri.addTriangle(q.prev.i, i, q.i, -1, q.t, q.prev.t)
				tri.legalize(t + 2)
				q.prev.t = t
				q.remove()
				q = q.prev
			}
		}

		// save the two new edges in the hash table
		tri.hashEdge(e)
		tri.hashEdge(e.prev)
	}

	tri.triangles = tri.triangles[:tri.trianglesLen]
	tri.halfedges = tri.halfedges[:tri.trianglesLen]

	return nil
}

func (t *triangulator) hashKey(point Point) int {
	dx := point.X - t.center.X
	dy := point.Y - t.center.Y
	// use pseudo-angle: a measure that monotonically increases
	// with real angle, but doesn't require expensive trigonometry
	p := 1 - dx/(math.Abs(dx)+math.Abs(dy))
	if dy < 0 {
		p = -p
	}
	return int(math.Floor((2 + p) / 4 * float64(len(t.hash))))
}

func (t *triangulator) hashEdge(e *node) {
	t.hash[t.hashKey(e.p)] = e
}

func (t *triangulator) addTriangle(i0, i1, i2, a, b, c int) int {
	i := t.trianglesLen
	t.triangles[i] = i0
	t.triangles[i+1] = i1
	t.triangles[i+2] = i2
	t.link(i, a)
	t.link(i+1, b)
	t.link(i+2, c)
	t.trianglesLen += 3
	return i
}

func (t *triangulator) link(a, b int) {
	t.halfedges[a] = b
	if b >= 0 {
		t.halfedges[b] = a
	}
}

func (t *triangulator) legalize(a int) int {
	b := t.halfedges[a]

	a0 := a - a%3
	b0 := b - b%3

	al := a0 + (a+1)%3
	ar := a0 + (a+2)%3
	bl := b0 + (b+2)%3

	p0 := t.triangles[ar]
	pr := t.triangles[a]
	pl := t.triangles[al]
	p1 := t.triangles[bl]

	illegal := inCircle(t.points[p0], t.points[pr], t.points[pl], t.points[p1])

	if illegal {
		t.triangles[a] = p1
		t.triangles[b] = p0

		t.link(a, t.halfedges[bl])
		t.link(b, t.halfedges[ar])
		t.link(ar, bl)

		br := b0 + (b+1)%3

		t.legalize(a)
		return t.legalize(br)
	}

	return ar
}

type node struct {
	p       Point
	i       int
	t       int
	removed bool
	prev    *node
	next    *node
}

func newNode(p Point, i int, prev *node) *node {
	n := &node{p, i, 0, false, nil, nil}
	if prev == nil {
		n.prev = n
		n.next = n
	} else {
		n.next = prev.next
		n.prev = prev
		prev.next.prev = n
		prev.next = n
	}
	return n
}

func (n *node) remove() *node {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.removed = true
	return n.prev
}
