package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/delaunay"
	"github.com/fogleman/gg"
)

const (
	W = 1024
	H = 1024
	N = 1000
)

func uniform(rnd *rand.Rand) delaunay.Point {
	x := rnd.Float64()
	y := rnd.Float64()
	return delaunay.Point{x, y}
}

func normal(rnd *rand.Rand) delaunay.Point {
	x := rnd.NormFloat64()
	y := rnd.NormFloat64()
	return delaunay.Point{x, y}
}

func circle(rnd *rand.Rand) delaunay.Point {
	a := rnd.Float64() * math.Pi * 2
	r := 1 + rnd.NormFloat64()*0.1
	x := math.Cos(a) * r
	y := math.Sin(a) * r
	return delaunay.Point{x, y}
}

func main() {
	// generate points
	rnd := rand.New(rand.NewSource(99))
	points := make([]delaunay.Point, N)
	for i := range points {
		points[i] = circle(rnd)
	}

	// triangulate
	start := time.Now()
	triangulation, err := delaunay.Triangulate(points)
	elapsed := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(elapsed)
	fmt.Println(len(triangulation.Triangles) / 3)

	// compute point bounds for rendering
	min := points[0]
	max := points[0]
	for _, p := range points {
		min.X = math.Min(min.X, p.X)
		min.Y = math.Min(min.Y, p.Y)
		max.X = math.Max(max.X, p.X)
		max.Y = math.Max(max.Y, p.Y)
	}

	size := delaunay.Point{max.X - min.X, max.Y - min.Y}
	center := delaunay.Point{min.X + size.X/2, min.Y + size.Y/2}
	scale := math.Min(W/size.X, H/size.Y) * 0.9

	// render points and edges
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	dc.Translate(W/2, H/2)
	dc.Scale(scale, scale)
	dc.Translate(-center.X, -center.Y)

	ts := triangulation.Triangles
	hs := triangulation.Halfedges
	for i, h := range hs {
		if i > h {
			p := points[ts[i]]
			q := points[ts[nextHalfEdge(i)]]
			dc.DrawLine(p.X, p.Y, q.X, q.Y)
		}
	}
	dc.Stroke()

	for _, p := range points {
		dc.DrawPoint(p.X, p.Y, 3)
	}
	dc.Fill()

	dc.SavePNG("out.png")
}

func nextHalfEdge(e int) int {
	if e%3 == 2 {
		return e - 2
	}
	return e + 1
}
