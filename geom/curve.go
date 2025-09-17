package geom

// Curve is an interface representing a parametric curve.
type Curve interface {
	// Point returns the geom.Point along the curve for 0 <= t <= 1.
	At(t float64) Point
}
