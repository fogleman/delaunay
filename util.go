package delaunay

import "math"

var infinity = math.Inf(1)

func area(a, b, c Point) float64 {
	return (b.Y-a.Y)*(c.X-b.X) - (b.X-a.X)*(c.Y-b.Y)
}

func inCircle(a, b, c, p Point) bool {
	dx := a.X - p.X
	dy := a.Y - p.Y
	ex := b.X - p.X
	ey := b.Y - p.Y
	fx := c.X - p.X
	fy := c.Y - p.Y

	ap := dx*dx + dy*dy
	bp := ex*ex + ey*ey
	cp := fx*fx + fy*fy

	return dx*(ey*cp-bp*fy)-dy*(ex*cp-bp*fx)+ap*(ex*fy-ey*fx) < 0
}

func circumradius(a, b, c Point) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	ex := c.X - a.X
	ey := c.Y - a.Y

	bl := dx*dx + dy*dy
	cl := ex*ex + ey*ey
	d := dx*ey - dy*ex

	x := (ey*bl - dy*cl) * 0.5 / d
	y := (dx*cl - ex*bl) * 0.5 / d

	r := x*x + y*y

	if bl == 0 || cl == 0 || d == 0 || r == 0 {
		return infinity
	}

	return r
}

func circumcenter(a, b, c Point) Point {
	dx := b.X - a.X
	dy := b.Y - a.Y
	ex := c.X - a.X
	ey := c.Y - a.Y

	bl := dx*dx + dy*dy
	cl := ex*ex + ey*ey
	d := dx*ey - dy*ex

	x := a.X + (ey*bl-dy*cl)*0.5/d
	y := a.Y + (dx*cl-ex*bl)*0.5/d

	return Point{x, y}
}
