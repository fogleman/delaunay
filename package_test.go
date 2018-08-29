package delaunay

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func uniform(rnd *rand.Rand) Point {
	x := rnd.Float64()
	y := rnd.Float64()
	return Point{x, y}
}

func normal(rnd *rand.Rand) Point {
	x := rnd.NormFloat64()
	y := rnd.NormFloat64()
	return Point{x, y}
}

func circle(rnd *rand.Rand) Point {
	a := rnd.Float64() * math.Pi * 2
	r := 1 + rnd.NormFloat64()*0.1
	x := math.Cos(a) * r
	y := math.Sin(a) * r
	return Point{x, y}
}

func TestSquare(t *testing.T) {
	points := []Point{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}
	tri, err := Triangulate(points)
	if err != nil {
		t.Fatal(err)
	}
	expected := []int{0, 2, 1, 0, 3, 2}
	if !reflect.DeepEqual(tri.Triangles, expected) {
		t.Fatalf("expected: %v, actual: %v", expected, tri.Triangles)
	}
	expected = []int{5, -1, -1, -1, -1, 0}
	if !reflect.DeepEqual(tri.Halfedges, expected) {
		t.Fatalf("expected: %v, actual: %v", expected, tri.Halfedges)
	}
}

func BenchmarkUniform(b *testing.B) {
	rnd := rand.New(rand.NewSource(99))
	points := make([]Point, b.N)
	for i := range points {
		points[i] = uniform(rnd)
	}
	Triangulate(points)
}

func BenchmarkNormal(b *testing.B) {
	rnd := rand.New(rand.NewSource(99))
	points := make([]Point, b.N)
	for i := range points {
		points[i] = normal(rnd)
	}
	Triangulate(points)
}

func BenchmarkCircle(b *testing.B) {
	rnd := rand.New(rand.NewSource(99))
	points := make([]Point, b.N)
	for i := range points {
		points[i] = circle(rnd)
	}
	Triangulate(points)
}
