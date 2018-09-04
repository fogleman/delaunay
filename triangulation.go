package delaunay

type Triangulation struct {
	Points     []Point
	ConvexHull []Point
	Triangles  []int
	Halfedges  []int
}

// Triangulate returns a Delaunay triangulation of the provided points.
func Triangulate(points []Point) (*Triangulation, error) {
	t := newTriangulator(points)
	err := t.triangulate()
	return &Triangulation{points, t.convexHull(), t.triangles, t.halfedges}, err
}

func (t *Triangulation) area() float64 {
	var result float64
	points := t.Points
	ts := t.Triangles
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		result += area(p0, p1, p2)
	}
	return result / 2
}

func (t *Triangulation) edgesAroundPoint(start int) []int {
	var result []int
	incoming := start
	for {
		result = append(result, incoming)
		outgoing := nextHalfedge(incoming)
		incoming = t.Halfedges[outgoing]
		if incoming == -1 || incoming == start {
			break
		}
	}
	return result
}

func (t *Triangulation) circumcenters() []Point {
	result := make([]Point, t.NumTriangles())
	for i := range result {
		result[i] = circumcenter(
			t.Points[t.Triangles[i*3+0]],
			t.Points[t.Triangles[i*3+1]],
			t.Points[t.Triangles[i*3+2]])
	}
	return result
}

func (t *Triangulation) NumTriangles() int {
	return len(t.Triangles) / 3
}

func (t *Triangulation) VoronoiEdges() [][2]Point {
	var result [][2]Point
	centers := t.circumcenters()
	for e := range t.Triangles {
		if e < t.Halfedges[e] {
			p0 := centers[triangleOfEdge(e)]
			p1 := centers[triangleOfEdge(t.Halfedges[e])]
			result = append(result, [2]Point{p0, p1})
		}
	}
	return result
}

func (t *Triangulation) VoronoiCells() [][]Point {
	index := make(map[int]int)
	for e := range t.Triangles {
		endpoint := t.Triangles[nextHalfedge(e)]
		_, ok := index[endpoint]
		if !ok || t.Halfedges[e] == -1 {
			index[endpoint] = e
		}
	}
	var result [][]Point
	centers := t.circumcenters()
	for _, p := range index {
		var cell []Point
		for _, e := range t.edgesAroundPoint(p) {
			cell = append(cell, centers[triangleOfEdge(e)])
		}
		result = append(result, cell)
	}
	return result
}
