package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/fogleman/delaunay"
)

type dist func(n int, rnd *rand.Rand) []delaunay.Point

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func uniform(n int, rnd *rand.Rand) []delaunay.Point {
	points := make([]delaunay.Point, n)
	for i := range points {
		x := rnd.Float64()
		y := rnd.Float64()
		points[i] = delaunay.Point{x, y}
	}
	return points
}

func normal(n int, rnd *rand.Rand) []delaunay.Point {
	points := make([]delaunay.Point, n)
	for i := range points {
		x := rnd.NormFloat64()
		y := rnd.NormFloat64()
		points[i] = delaunay.Point{x, y}
	}
	return points
}

func grid(n int, rnd *rand.Rand) []delaunay.Point {
	side := int(math.Floor(math.Sqrt(float64(n))))
	n = side * side
	points := make([]delaunay.Point, 0, n)
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			p := delaunay.Point{float64(x), float64(y)}
			points = append(points, p)
		}
	}
	return points
}

func circle(n int, rnd *rand.Rand) []delaunay.Point {
	points := make([]delaunay.Point, n)
	for i := range points {
		t := float64(i) / float64(n)
		x := math.Cos(t)
		y := math.Sin(t)
		points[i] = delaunay.Point{x, y}
	}
	return points
}

func test(f dist, n int) {
	rnd := rand.New(rand.NewSource(99))
	points := f(n, rnd)
	start := time.Now()
	_, err := delaunay.Triangulate(points)
	elapsed := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n, elapsed)
}

func main() {
	dists := []dist{
		uniform,
		normal,
		grid,
		circle,
	}
	for _, f := range dists {
		fmt.Println(getFunctionName(f))
		for n := 10; n <= 1000000; n *= 10 {
			test(f, n)
		}
		fmt.Println()
	}
}
