package delaunay

type Point struct {
	X, Y float64
}

func (a Point) squaredDistance(b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func (a Point) sub(b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}
