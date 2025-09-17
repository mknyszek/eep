package geom

import (
	"math"
)

type QuadraticBezier struct {
	a, b, c Point
}

func Bezier2(a, b, c Point) QuadraticBezier {
	return QuadraticBezier{a, b, c}
}

func (qb QuadraticBezier) At(t float64) Point {
	x, y := quadratic(qb.a.X, qb.a.Y, qb.b.X, qb.b.Y, qb.c.X, qb.c.Y, t)
	return Point{x, y}
}

func (qb QuadraticBezier) Points() []Point {
	l := (math.Hypot(qb.b.X-qb.a.X, qb.b.Y-qb.a.Y) +
		math.Hypot(qb.c.X-qb.b.X, qb.c.Y-qb.b.Y))
	n := int(l + 0.5)
	if n < 4 {
		n = 4
	}
	denom := float64(n) - 1
	result := make([]Point, n)
	for i := 0; i < n; i++ {
		result[i] = qb.At(float64(i) / denom)
	}
	return result
}

func quadratic(x0, y0, x1, y1, x2, y2, t float64) (x, y float64) {
	u := 1 - t
	a := u * u
	b := 2 * u * t
	c := t * t
	x = a*x0 + b*x1 + c*x2
	y = a*y0 + b*y1 + c*y2
	return
}

type CubicBezier struct {
	a, b, c, d Point
}

func Bezier3(a, b, c, d Point) CubicBezier {
	return CubicBezier{a, b, c, d}
}

func (cb CubicBezier) At(t float64) Point {
	x, y := cubic(cb.a.X, cb.a.Y, cb.b.X, cb.b.Y, cb.c.X, cb.c.Y, cb.d.X, cb.d.Y, t)
	return Point{x, y}
}

func (cb CubicBezier) Points() []Point {
	l := (math.Hypot(cb.b.X-cb.a.X, cb.b.Y-cb.a.Y) +
		math.Hypot(cb.c.X-cb.b.X, cb.c.Y-cb.b.Y) +
		math.Hypot(cb.d.X-cb.c.X, cb.d.Y-cb.c.Y))
	n := int(l + 0.5)
	if n < 4 {
		n = 4
	}
	denom := float64(n) - 1
	result := make([]Point, n)
	for i := 0; i < n; i++ {
		result[i] = cb.At(float64(i) / denom)
	}
	return result
}

func cubic(x0, y0, x1, y1, x2, y2, x3, y3, t float64) (x, y float64) {
	u := 1 - t
	a := u * u * u
	b := 3 * u * u * t
	c := 3 * u * t * t
	d := t * t * t
	x = a*x0 + b*x1 + c*x2 + d*x3
	y = a*y0 + b*y1 + c*y2 + d*y3
	return
}
