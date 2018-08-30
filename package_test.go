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

func testHalfedges(t *testing.T, points []Point) {
	tri, err := Triangulate(points)
	if err != nil {
		t.Fatal(err)
	}
	for i1, i2 := range tri.Halfedges {
		if i1 != -1 && tri.Halfedges[i1] != i2 {
			t.Fatal("invalid halfedge connection")
		}
		if i2 != -1 && tri.Halfedges[i2] != i1 {
			t.Fatal("invalid halfedge connection")
		}
	}
}

func TestHalfedges(t *testing.T) {
	testHalfedges(t, []Point{{516, 661}, {369, 793}, {426, 539}, {273, 525}, {204, 694}, {747, 750}, {454, 390}})
	testHalfedges(t, []Point{{382, 302}, {382, 328}, {382, 205}, {623, 175}, {382, 188}, {382, 284}, {623, 87}, {623, 341}, {141, 227}})
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
