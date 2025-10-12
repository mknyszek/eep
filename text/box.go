package text

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mknyszek/eep/font"
	"github.com/mknyszek/eep/geom"
)

// Box draws a text box containing the provided text within the provided bounds, with
// any layout options applied.
func Box(dst *ebiten.Image, txt String, bounds geom.AABB, boxOpts *BoxOptions) {
	sub := dst.SubImage(bounds.Image()).(*ebiten.Image)
	drawString(sub, txt, bounds.Min, txt.Measure(boxOpts.LineSpacing), boxOpts)
}

// AutoBox draws a text box with an automatically adjusted size to fit all the
// text at position pos.
//
// anchor is a relative anchor point for the box.
// (0, 0) means top-left and (1, 1) means bottom-right of the box.
func AutoBox(dst *ebiten.Image, txt String, pos geom.Point, anchor geom.Dimensions) {
	dim := txt.Measure(0)
	pos = geom.Pt(pos.X-anchor.X*dim.X, pos.Y-anchor.Y*dim.Y)
	sub := dst.SubImage(dim.AABB(pos).Image()).(*ebiten.Image)
	drawString(sub, txt, pos, dim, &BoxOptions{})
}

func drawString(dst *ebiten.Image, txt String, orig geom.Point, txtDim geom.Dimensions, boxOpts *BoxOptions) {
	// Advance our way through the positions and draw.
	m := newBoxDrawMachine(orig, geom.ImageAABB(dst.Bounds()).Dim(), txtDim, txt.direction, boxOpts.Align, boxOpts.VertAlign)
	for segments := range txt.lines() {
		// Figure out the increments for each line.
		primary, secondary := measureLine(segments, txt.direction, boxOpts.LineSpacing)

		// Render the line.
		m.StartLine(primary)
		for _, seg := range segments {
			if seg.text == "\n" {
				continue
			}
			face := seg.style.Face.TextFace()

			var opts text.DrawOptions
			opts.GeoM.Translate(m.X(), m.Y())
			opts.ColorScale.ScaleWithColor(seg.style.Color)
			text.Draw(dst, seg.text, seg.style.Face.TextFace(), &opts)

			m.MoveInLine(text.Advance(seg.text, face))
		}

		// Advance the line.
		m.EndLine(secondary)
	}
}

type BoxOptions struct {
	// Padding sets the distance between the edge of the box and the text inside.
	Padding geom.Dimensions

	// LineSpacing is multiplier proportional to the size of each line of text.
	LineSpacing float64

	// Align sets the alignment of the text.
	Align Alignment

	// VertAlign sets the vertical alignment of the text.
	VertAlign VertAlignment
}

func (o *BoxOptions) set(opts *text.DrawOptions, dir text.Direction) {
	opts.LineSpacing = o.LineSpacing

	if o.Align != AutoAlign {
		switch dir {
		case text.DirectionLeftToRight:
			opts.PrimaryAlign = text.Align(o.Align - 1)
		case text.DirectionRightToLeft:
			opts.PrimaryAlign = text.Align(Right - o.Align - 1)
		case text.DirectionTopToBottomAndLeftToRight:
			opts.SecondaryAlign = text.Align(o.Align - 1)
		case text.DirectionTopToBottomAndRightToLeft:
			opts.SecondaryAlign = text.Align(Right - o.Align - 1)
		}
	}
	if o.VertAlign != AutoVertAlign {
		switch dir {
		case text.DirectionLeftToRight:
			opts.SecondaryAlign = text.Align(o.VertAlign - 1)
		case text.DirectionRightToLeft:
			opts.SecondaryAlign = text.Align(o.VertAlign - 1)
		case text.DirectionTopToBottomAndLeftToRight:
			opts.PrimaryAlign = text.Align(o.VertAlign - 1)
		case text.DirectionTopToBottomAndRightToLeft:
			opts.PrimaryAlign = text.Align(o.VertAlign - 1)
		}
	}
}

// Alignment describes horizontal text alignment.
type Alignment int

const (
	AutoAlign Alignment = iota
	Left
	Center
	Right
)

// VertAlignment describes vertical text alignment.
type VertAlignment int

const (
	AutoVertAlign VertAlignment = iota
	Top
	Middle
	Bottom
)

type boxDrawMachine struct {
	pos        geom.Point // Current position.
	orig       geom.Point // Top-left of text-box.
	boxDim     geom.Dimensions
	txtDim     geom.Dimensions
	init       func(m *boxDrawMachine)
	startLine  func(m *boxDrawMachine, priLen float64)
	moveInLine func(m *boxDrawMachine, priAmt float64)
	endLine    func(m *boxDrawMachine, secLen float64)
}

func (m *boxDrawMachine) StartLine(priLen float64) {
	m.startLine(m, priLen)
}

func (m *boxDrawMachine) MoveInLine(priAmt float64) {
	m.moveInLine(m, priAmt)
}

func (m *boxDrawMachine) EndLine(secLen float64) {
	m.endLine(m, secLen)
}

func (m *boxDrawMachine) X() float64 {
	return m.pos.X
}

func (m *boxDrawMachine) Y() float64 {
	return m.pos.Y
}

func newBoxDrawMachine(orig geom.Point, boxDim, txtDim geom.Dimensions, dir font.TextDirection, align Alignment, valign VertAlignment) *boxDrawMachine {
	m := new(boxDrawMachine)
	m.orig = orig
	m.boxDim = boxDim
	m.txtDim = txtDim

	switch dir {
	case font.DirectionLeftToRight, font.DirectionRightToLeft:
		switch valign {
		case Top, AutoVertAlign:
			m.init = func(m *boxDrawMachine) { m.pos.Y = m.orig.Y }
		case Middle:
			m.init = func(m *boxDrawMachine) { m.pos.Y = m.orig.Y + m.boxDim.Y/2 - m.txtDim.Y/2 }
		case Bottom:
			m.init = func(m *boxDrawMachine) { m.pos.Y = m.orig.Y + m.boxDim.Y - m.txtDim.Y }
		}
	case font.DirectionTopToBottomAndLeftToRight:
		switch align {
		case Left, AutoAlign:
			m.init = func(m *boxDrawMachine) { m.pos.X = m.orig.X }
		case Center:
			m.init = func(m *boxDrawMachine) { m.pos.X = m.orig.X + m.boxDim.X/2 - m.txtDim.X/2 }
		case Right:
			m.init = func(m *boxDrawMachine) { m.pos.X = m.orig.X + m.boxDim.X - m.txtDim.X }
		}
	case font.DirectionTopToBottomAndRightToLeft:
		switch align {
		case Left:
			m.init = func(m *boxDrawMachine) { m.pos.X = m.orig.X + m.txtDim.X }
		case Center:
			m.init = func(m *boxDrawMachine) { m.pos.X = m.orig.X + m.boxDim.X - (m.boxDim.X-m.txtDim.X)/2 }
		case Right, AutoAlign:
			m.init = func(m *boxDrawMachine) { m.pos.X = m.orig.X + m.boxDim.X }
		}
	}
	switch dir {
	case font.DirectionLeftToRight:
		m.moveInLine = func(m *boxDrawMachine, a float64) { m.pos.X += a }
	case font.DirectionRightToLeft:
		m.moveInLine = func(m *boxDrawMachine, a float64) { m.pos.X -= a }
	case font.DirectionTopToBottomAndLeftToRight, font.DirectionTopToBottomAndRightToLeft:
		m.moveInLine = func(m *boxDrawMachine, a float64) { m.pos.Y += a }
	}
	switch dir {
	case font.DirectionLeftToRight, font.DirectionRightToLeft:
		m.endLine = func(m *boxDrawMachine, a float64) { m.pos.Y += a }
	case font.DirectionTopToBottomAndLeftToRight:
		m.endLine = func(m *boxDrawMachine, a float64) { m.pos.X += a }
	case font.DirectionTopToBottomAndRightToLeft:
		m.endLine = func(m *boxDrawMachine, a float64) { m.pos.X -= a }
	}
	switch dir {
	case font.DirectionLeftToRight:
		switch align {
		case Left, AutoAlign:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.X = m.orig.X }
		case Center:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.X = m.orig.X + (m.boxDim.X-primary)/2 }
		case Right:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.X = m.orig.X + m.boxDim.X - primary }
		}
	case font.DirectionRightToLeft:
		switch align {
		case Left:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.X = m.orig.X + primary }
		case Center:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.X = m.orig.X + m.boxDim.X - (m.boxDim.X-primary)/2 }
		case Right, AutoAlign:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.X = m.orig.X + m.boxDim.X }
		}
	case font.DirectionTopToBottomAndLeftToRight, font.DirectionTopToBottomAndRightToLeft:
		switch valign {
		case Top, AutoVertAlign:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.Y = m.orig.Y }
		case Middle:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.Y = m.orig.Y + (m.boxDim.Y-primary)/2 }
		case Bottom:
			m.startLine = func(m *boxDrawMachine, primary float64) { m.pos.Y = m.orig.Y + m.boxDim.Y - primary }
		}
	}
	m.init(m)
	return m
}
