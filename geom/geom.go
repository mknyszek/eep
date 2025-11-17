package geom

import (
	"image"
	"math"
)

var Origin = Point{0, 0}

// Point represents a position in R^2.
//
// A Point is a special case of a Vector: a distance vector
// from Origin. The distinction is useful at the type level,
// but functionally the two are essentially equivalent.
type Point struct {
	X, Y float64
}

// Point is shorthand for a new point.
func Pt(x, y float64) Point {
	return Point{x, y}
}

// Add moves a point by a vector.
func (p Point) Add(v Vector) Point {
	return Point{p.X + v.X, p.Y + v.Y}
}

// Sub subtracts one point's elements from another.
func (p0 Point) Sub(p1 Point) Point {
	return Point{p0.X - p1.X, p0.Y - p1.Y}
}

// Vector converts the point to a vector. This is equivalent to Vec(Origin, p).
func (p Point) Vector() Vector {
	return Vector{p.X, p.Y}
}

// Image returns a rounded point, suitable for use in image processing.
func (p Point) Image() image.Point {
	return image.Point{int(math.Round(p.X)), int(math.Round(p.Y))}
}

// ImagePoint produces a Point from an image.Point.
func ImagePoint(pt image.Point) Point {
	return Point{float64(pt.X), float64(pt.Y)}
}

// Segment represents a line segment between two points.
//
// Segment implements Curve.
type Segment struct {
	Start, End Point
}

// Seg returns a new Segment between the two provided points.
func Seg(a, b Point) Segment {
	return Segment{a, b}
}

// At interpolates between the two points in the segment according to t, a parameter
// from 0 to 1.
func (s Segment) At(t float64) Point {
	return Point{s.Start.X + t*(s.End.X-s.Start.X), s.Start.Y + t*(s.End.Y-s.Start.Y)}
}

// Line represents an infinite line of the form y=mx+b where m is the slope
// and b is the y-intercept.
type Line struct {
	M, B float64
}

// LineFromPoints computes the Line that is intersects the two provided points.
func LineFromPoints(p0, p1 Point) Line {
	r := 1 / (p1.X - p0.X)
	return Line{(p1.Y - p0.Y) * r, (p1.X*p0.Y - p0.X*p1.Y) * r}
}

// Intercept returns the intersection point of the two lines.
//
// Returns false if the lines do not intersect, or if they're identical (intersect at every point).
func (l0 Line) Intercept(l1 Line) (Point, bool) {
	if l1.M == l0.M {
		return Point{}, false
	}
	r := 1 / (l1.M - l0.M)
	return Point{(l0.B - l1.B) * r, (l1.M*l0.B - l0.M*l1.B) * r}, true
}

// Vector is a two-dimensional vector.
type Vector struct {
	X, Y float64
}

// Vec create a new Vector in R^2 from a pair of points.
func Vec(origin, p Point) Vector {
	v := p.Sub(origin)
	return Vector{v.X, v.Y}
}

func (v0 Vector) Add(v1 Vector) Vector {
	return Vector{v0.X + v1.X, v0.Y + v1.Y}
}

// Dot computes the dot product of two vectors.
func (v0 Vector) Dot(v1 Vector) float64 {
	return v0.X*v1.X + v0.Y*v1.Y
}

// Scale scales the vector by a.
func (v Vector) Scale(a float64) Vector {
	return Vector{a * v.X, a * v.Y}
}

// Point produces a Point from the Vector, given an origin.
func (v Vector) Point(origin Point) Point {
	return Point{v.X + origin.X, v.Y + origin.Y}
}

// Length2 returns the square of the length of the vector.
func (v Vector) Length2() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Length returns the length of the vector.
func (v Vector) Length() float64 {
	return math.Sqrt(v.Length2())
}

// Normalize returns a unit vector copy of v pointing in the same direction.
func (v Vector) Normalize() Vector {
	l := v.Length()
	return Vector{v.X / l, v.Y / l}
}

// Rotate rotates the vector about the origin by rad radians.
func (v Vector) Rotate(rad float64) Vector {
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	return Vector{
		v.X*cos - v.Y*sin,
		v.X*sin + v.Y*cos,
	}
}

// RightNormal computes the right-normal vector of this vector.
func (v Vector) RightNormal() Vector {
	return Vector{v.Y, -v.X}
}

// Dim is a set of 2D dimensions, a width and a height, without a location.
type Dimensions struct {
	X, Y float64
}

// Dim creates a new set of 2D dimensions.
func Dim(x, y float64) Dimensions {
	return Dimensions{x, y}
}

// ImageDim returns the dimensions of an image.Rectangle.
func ImageDim(r image.Rectangle) Dimensions {
	return Dimensions{float64(r.Dx()), float64(r.Dy())}
}

// AABB gives the dimensions a starting location, producing an AABB.
func (d Dimensions) AABB(start Point) AABB {
	return AABB{start, start.Add(d.Vector())}
}

// Vector returns a vector that represents the dimensions.
func (d Dimensions) Vector() Vector {
	return Vector{d.X, d.Y}
}

// AABB describes an axis-aligned bounding box in R^2.
type AABB struct {
	Min, Max Point
}

// ImageAABB returns an AABB for the image.Rectangle.
func ImageAABB(r image.Rectangle) AABB {
	return AABB{ImagePoint(r.Min), ImagePoint(r.Max)}
}

// Image returns an image.Rectangle, rounding the AABB for use in image processing.
func (a AABB) Image() image.Rectangle {
	return image.Rectangle{a.Min.Image(), a.Max.Image()}
}

// Bound creates a new AABB from two points.
func Bound(x0, y0, x1, y1 float64) AABB {
	return AABB{Point{x0, y0}, Point{x1, y1}}
}

func (a AABB) Dx() float64 {
	return a.Max.X - a.Min.X
}

func (a AABB) Dy() float64 {
	return a.Max.Y - a.Min.Y
}

func (a AABB) Dim() Dimensions {
	return Dim(a.Dx(), a.Dy())
}

// Translate moves the AABB in the direction of the provided vector.
func (a AABB) Translate(v Vector) AABB {
	return AABB{a.Min.Add(v), a.Max.Add(v)}
}

// Intersects returns true if the two AABBs intersect.
func (a AABB) Intersects(b AABB) bool {
	return !(a.Max.X < b.Min.X || a.Min.X > b.Max.X || a.Max.Y < b.Min.Y || a.Min.Y > b.Max.Y)
}

// Rad converts degrees to radians.
func Rad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Deg converts radiasn to degrees.
func Deg(radians float64) float64 {
	return radians * 180 / math.Pi
}
