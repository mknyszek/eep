package graphics

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mknyszek/eep/geom"
)

// Context is a vector graphics context for drawing vector graphics.
type Context struct {
	ctx
	dst      *ebiten.Image
	vertices []ebiten.Vertex
	indices  []uint16
	stack    []ctx
}

type ctx struct {
	matrix     ebiten.GeoM
	color      color.Color
	path       vector.Path
	strokeOpts vector.StrokeOptions
	fillOpts   vector.FillOptions
	opts       vector.DrawPathOptions
}

// NewContext creates a new Context with the intent to draw to dst.
func NewContext(dst *ebiten.Image) *Context {
	c := &Context{dst: dst}
	c.color = color.Black
	c.opts.AntiAlias = true
	return c
}

// High-level drawing functions.

// Rect draws an axis-aligned rectangle to dst directly.
// It uses the current context, but does not modify the current path.
func (c *Context) Rect(m Method, bounds geom.AABB) {
	c.WithEmpty(func(c *Context) {
		c.MoveTo(bounds.Min)
		c.LineTo(geom.Pt(bounds.Min.X, bounds.Max.Y))
		c.LineTo(bounds.Max)
		c.LineTo(geom.Pt(bounds.Max.X, bounds.Min.Y))
		c.ClosePath()
		c.Draw(m, true) // No need to clear the path.
	})
}

// Arrow draws an arrow from src to dst.
func (c *Context) Arrow(src, dst geom.Point) {
	c.WithEmpty(func(c *Context) {
		// Draw arrow line.
		c.MoveTo(src)
		c.LineTo(src)
		c.Stroke()

		// Compute arrow head points.
		const ahMul = 7                              // Arrow head length multiplier.
		const th = math.Pi / 8                       // Rotation angle (arrow head width).
		ahLen := ahMul * float64(c.strokeOpts.Width) // Arrow head length.
		vec := geom.Vec(dst, src).Normalize().Scale(ahLen)
		ah0 := dst.Add(vec.Rotate(th))
		ah1 := dst.Add(vec.Rotate(-th))

		// Draw arrow head.
		c.MoveTo(dst)
		c.LineTo(ah0)
		c.LineTo(ah1)
		c.ClosePath()
		c.Fill()
	})
}

// Circle draws a circle to dst directly.
// It uses the current context, but does not modify the current path.
func (c *Context) Circle(m Method, center geom.Point, radius float64) {
	c.EllipticalArc(m, center, radius, radius, 0, 2*math.Pi)
}

// Oval draws an axis-aligned oval to dst directly.
// It uses the current context, but does not modify the current path.
func (c *Context) Oval(m Method, center geom.Point, radiusX, radiusY float64) {
	c.EllipticalArc(m, center, radiusX, radiusY, 0, 2*math.Pi)
}

// EllipticalArc draws an axis-aligned elliptical arc to dst directly.
// It uses the current context, but does not modify the current path.
func (c *Context) EllipticalArc(m Method, center geom.Point, radiusX, radiusY, fromRad, toRad float64) {
	c.WithEmpty(func(c *Context) {
		const n = 16
		for i := range n {
			seg0 := float64(i+0) / n
			seg1 := float64(i+1) / n
			angle0 := fromRad + (toRad-fromRad)*seg0
			angle2 := fromRad + (toRad-fromRad)*seg1
			angle1 := (angle0 + angle2) / 2
			p0 := center.Add(geom.Pt(radiusX*math.Cos(angle0), radiusY*math.Sin(angle0)).Vector())
			p1 := center.Add(geom.Pt(radiusX*math.Cos(angle1), radiusY*math.Sin(angle1)).Vector())
			p2 := center.Add(geom.Pt(radiusX*math.Cos(angle2), radiusY*math.Sin(angle2)).Vector())
			ctrl := geom.Pt(2*p1.X-p0.X/2-p2.X/2, 2*p1.Y-p0.Y/2-p2.Y/2)
			if i == 0 {
				c.MoveTo(p0)
			}
			c.QuadTo(ctrl, p2)
		}
		c.Draw(m, true) // No need to clear the path.
	})
}

// Context stack functions.

// WithEmpty temporarily swaps out the context's path with a new empty path, for the duration of f.
func (c *Context) WithEmpty(f func(c *Context)) {
	var tmp vector.Path
	c.path, tmp = tmp, c.path
	f(c)
	c.path, tmp = tmp, c.path
}

// Push clones the current context and pushes it onto the internal stack.
func (c *Context) Push() {
	old := c.ctx
	c.stack = append(c.stack, old)
	c.path = vector.Path{}
	c.path.AddPath(&old.path, &vector.AddPathOptions{})
}

// Pop restores the previously pushed context.
func (c *Context) Pop() {
	c.ctx = c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
}

// Draw functions.

// Method describes the draw method, either fill (area filling) or stroke (boundary drawing).
type Method int

const (
	Fill Method = iota
	Stroke
)

// Draw draws the current path with the provided method. If preserve is true, preserves the current path.
func (c *Context) Draw(m Method, preserve bool) {
	switch m {
	case Fill:
		vector.FillPath(c.dst, &c.path, &c.fillOpts, &c.opts)
	case Stroke:
		vector.StrokePath(c.dst, &c.path, &c.strokeOpts, &c.opts)
	}
	if !preserve {
		c.ClearPath()
	}
}

// Stroke draws the stroke of the current path to the destination image.
// It clears the current path.
func (c *Context) Stroke() {
	c.Draw(Stroke, false)
}

// StrokePreserve fills the boundary of the current path with the current color and draws it to the destination image.
// It does not clear (preserves) the current path.
func (c *Context) StrokePreserve() {
	c.Draw(Stroke, true)
}

// Fill fills the boundary of the current path with the current color and draws it to the destination image.
// It clears the current path.
func (c *Context) Fill() {
	c.Draw(Fill, false)
}

// FillPreserve fills the boundary of the current path with the current color and draws it to the destination image.
// It does not clear (preserves) the current path.
func (c *Context) FillPreserve() {
	c.Draw(Fill, true)
}

// Styling functions.

// SetColor sets the current color.
func (c *Context) SetColor(clr color.Color) {
	c.opts.ColorScale.Reset()
	c.opts.ColorScale.ScaleWithColor(clr)
}

// SetLineWidth sets the width of the stroke to draw.
func (c *Context) SetLineWidth(w float64) {
	c.strokeOpts.Width = float32(w)
}

type LineCap vector.LineCap

const (
	LineCapButt   = LineCap(vector.LineCapButt)
	LineCapRound  = LineCap(vector.LineCapRound)
	LineCapSquare = LineCap(vector.LineCapSquare)
)

// SetLineCap sets style of line cap.
func (c *Context) SetLineCap(cap LineCap) {
	c.strokeOpts.LineCap = vector.LineCap(cap)
}

type LineJoin struct {
	join  vector.LineJoin
	param float64
}

var (
	LineJoinMiterDefault = LineJoin{vector.LineJoinMiter, 0}
	LineJoinBevel        = LineJoin{vector.LineJoinBevel, 0}
	LineJoinRound        = LineJoin{vector.LineJoinRound, 0}
)

func LineJoinMiter(limit float64) LineJoin {
	return LineJoin{vector.LineJoinMiter, limit}
}

// SetLineJoin sets style of line join.
func (c *Context) SetLineJoin(j LineJoin) {
	c.strokeOpts.LineJoin = j.join
	c.strokeOpts.MiterLimit = float32(j.param)
}

// Set Ebiten drawing controls.

// SetFillRule sets the fill rule for drawing. The default is FillRuleFillAll.
func (c *Context) SetFillRule(r vector.FillRule) {
	c.fillOpts.FillRule = r
}

// SetBlend sets the blend rule for drawing. The default is regular alpha blending.
func (c *Context) SetBlend(b ebiten.Blend) {
	c.opts.Blend = b
}

// Basic path primitives.

// MoveTo moves the current point in the path to (x, y).
func (c *Context) MoveTo(pt geom.Point) {
	pt = c.TransformPoint(pt)
	c.path.MoveTo(float32(pt.X), float32(pt.Y))
}

// LineTo appends the current path with a line from the current point to the provided point, and sets
// (x, y) as the new current point.
func (c *Context) LineTo(pt geom.Point) {
	pt = c.TransformPoint(pt)
	c.path.LineTo(float32(pt.X), float32(pt.Y))
}

// QuadTo appends the current path with a quadratic Bezier curve starting at the current point through to dst,
// using ctrl as the control point.
func (c *Context) QuadTo(ctrl, dst geom.Point) {
	ctrl = c.TransformPoint(ctrl)
	dst = c.TransformPoint(dst)
	c.path.QuadTo(float32(ctrl.X), float32(ctrl.Y), float32(dst.X), float32(dst.Y))
}

// CubicTo appends the current path with a cubic Bezier curve starting at the current point through to dst,
// using ctrl0 and ctrl1 as the control points.
func (c *Context) CubicTo(ctrl0, ctrl1, dst geom.Point) {
	ctrl0 = c.TransformPoint(ctrl0)
	ctrl1 = c.TransformPoint(ctrl1)
	dst = c.TransformPoint(dst)
	c.path.CubicTo(float32(ctrl0.X), float32(ctrl0.Y), float32(ctrl1.X), float32(ctrl1.Y), float32(dst.X), float32(dst.Y))
}

// ClosePath closes the current path.
func (c *Context) ClosePath() {
	c.path.Close()
}

// ClearPath clears the current path.
func (c *Context) ClearPath() {
	c.path.Reset()
}

// Transformation matrix primitives.

// TransformPoint applies the context's transformation.
func (c *Context) TransformPoint(pt geom.Point) geom.Point {
	x, y := c.matrix.Apply(pt.X, pt.Y)
	return geom.Pt(x, y)
}

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func (c *Context) Identity() {
	c.matrix.Reset()
}

// Translate updates the current matrix with a translation.
func (c *Context) Translate(x, y float64) {
	c.matrix.Translate(x, y)
}

// Scale updates the current matrix with a scaling factor.
// Scaling occurs about the origin.
func (c *Context) Scale(x, y float64) {
	c.matrix.Scale(x, y)
}

// ScaleAbout updates the current matrix with a scaling factor.
// Scaling occurs about the specified point.
func (c *Context) ScaleAbout(sx, sy, x, y float64) {
	c.Translate(x, y)
	c.Scale(sx, sy)
	c.Translate(-x, -y)
}

// Rotate updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the origin. Angle is specified in radians.
func (c *Context) Rotate(angle float64) {
	c.matrix.Rotate(angle)
}

// RotateAbout updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the specified point. Angle is specified in radians.
func (c *Context) RotateAbout(angle float64, pt geom.Point) {
	c.Translate(-pt.X, -pt.Y)
	c.Rotate(angle)
	c.Translate(pt.X, pt.Y)
}

// Skew updates the current matrix with a shearing angle.
// Skewing occurs about the origin.
func (c *Context) Skew(sx, sy float64) {
	c.matrix.Skew(sx, sy)
}

// SkewAbout updates the current matrix with a shearing angle.
// Skewing occurs about the specified point.
func (c *Context) SkewAbout(sx, sy float64, pt geom.Point) {
	c.Translate(-pt.X, -pt.Y)
	c.Skew(sx, sy)
	c.Translate(pt.X, pt.Y)
}

// InvertY flips the Y axis so that Y grows from bottom to top and Y=0 is at
// the bottom of the image.
func (c *Context) InvertY() {
	c.Translate(0, float64(c.dst.Bounds().Dx()))
	c.Scale(1, -1)
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}
