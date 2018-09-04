package delaunay

import (
	"math"
	"math/big"
)

var floatErr = math.Nextafter(1, 2) - 1
var inCircleErr = 12 * floatErr
var areaErr = 12 * floatErr // TODO: correct value?

func makeRobustPoints(points []Point) []Point {
	if len(points) == 0 {
		return nil
	}
	result := make([]Point, len(points))
	x0, y0, x1, y1 := bounds(points)
	if x0 >= 1 && y0 >= 1 && x1 < 2 && y1 < 2 {
		copy(result, points)
	} else {
		mx := 1 / (x1 - x0 + 1e-9)
		my := 1 / (y1 - y0 + 1e-9)
		m := math.Min(mx, my)
		for i, p := range points {
			rx := 1 + (p.X-x0)*m
			ry := 1 + (p.Y-y0)*m
			result[i] = Point{rx, ry}
		}
	}
	return result
}

func extractBig(x float64) *big.Int {
	if x < 1 || x >= 2 {
		// panic("invalid argument for extractBig")
	}
	return big.NewInt(int64(math.Float64bits(x) & 0x000fffffffffffff))
}

func add(x, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

func sub(x, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

func mul(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

func mul3(x, y, z *big.Int) *big.Int {
	b := new(big.Int)
	b.Mul(x, y)
	return b.Mul(b, z)
}

func dot(x, y *big.Int) *big.Int {
	a := new(big.Int)
	b := new(big.Int)
	a.Mul(x, x)
	b.Mul(y, y)
	return a.Add(a, b)
}

func exactInCircle(ax, ay, bx, by, cx, cy, px, py *big.Int) bool {
	bx.Sub(bx, ax)
	by.Sub(by, ay)
	cx.Sub(cx, ax)
	cy.Sub(cy, ay)
	px.Sub(px, ax)
	py.Sub(py, ay)
	br2 := dot(bx, by)
	cr2 := dot(cx, cy)
	pr2 := dot(px, py)
	result := new(big.Int)
	result = add(result, mul3(br2, cy, px))
	result = sub(result, mul3(br2, cx, py))
	result = add(result, mul3(bx, cr2, py))
	result = sub(result, mul3(bx, cy, pr2))
	result = add(result, mul3(by, cx, pr2))
	result = sub(result, mul3(by, cr2, px))
	return result.Sign() < 0
}

func robustInCircle(robust bool, a, b, c, p Point) bool {
	dx := a.X - p.X
	dy := a.Y - p.Y
	ex := b.X - p.X
	ey := b.Y - p.Y
	fx := c.X - p.X
	fy := c.Y - p.Y
	ap := dx*dx + dy*dy
	bp := ex*ex + ey*ey
	cp := fx*fx + fy*fy
	d := dx*(ey*cp-bp*fy) - dy*(ex*cp-bp*fx) + ap*(ex*fy-ey*fx)
	if !robust {
		return d < 0
	}
	if d < -inCircleErr {
		return true
	}
	if d > inCircleErr {
		return false
	}
	return exactInCircle(
		extractBig(a.X),
		extractBig(a.Y),
		extractBig(b.X),
		extractBig(b.Y),
		extractBig(c.X),
		extractBig(c.Y),
		extractBig(p.X),
		extractBig(p.Y),
	)
}

func robustDistanceComparison(robust bool, p, a, b Point) int {
	da := p.squaredDistance(a)
	db := p.squaredDistance(b)
	d := da - db

	if !robust || d < -areaErr || d > areaErr {
		return int(math.Copysign(1, d))
	}

	px := extractBig(p.X)
	py := extractBig(p.Y)
	ax := extractBig(a.X)
	ay := extractBig(a.Y)
	bx := extractBig(b.X)
	by := extractBig(b.Y)

	ad := dot(sub(ax, px), sub(ay, py))
	bd := dot(sub(bx, px), sub(by, py))
	return sub(ad, bd).Sign()
}

func robustOrient(robust bool, a, b, c Point) bool {
	d := (b.Y-a.Y)*(c.X-b.X) - (b.X-a.X)*(c.Y-b.Y)
	if !robust {
		return d < 0
	}
	if d < -areaErr {
		return true
	}
	if d > areaErr {
		return false
	}
	ax := extractBig(a.X)
	ay := extractBig(a.Y)
	bx := extractBig(b.X)
	by := extractBig(b.Y)
	cx := extractBig(c.X)
	cy := extractBig(c.Y)
	bay := sub(by, ay)
	cbx := sub(cx, bx)
	bax := sub(bx, ax)
	cby := sub(cy, by)
	p1 := mul(bay, cbx)
	p2 := mul(bax, cby)
	sign := sub(p1, p2).Sign()
	return sign < 0
}

func robustCross2D(robust bool, p, a, b Point) int {
	d := (a.X-p.X)*(b.Y-p.Y) - (a.Y-p.Y)*(b.X-p.X)
	if d < -areaErr {
		return -1
	}
	if d > areaErr {
		return 1
	}
	px := extractBig(p.X)
	py := extractBig(p.Y)
	ax := extractBig(a.X)
	ay := extractBig(a.Y)
	bx := extractBig(b.X)
	by := extractBig(b.Y)
	apx := sub(ax, px)
	bpy := sub(by, py)
	apy := sub(ay, py)
	bpx := sub(bx, px)
	p1 := mul(apx, bpy)
	p2 := mul(apy, bpx)
	return sub(p1, p2).Sign()
}
