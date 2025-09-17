package text

import (
	"image/color"
	"iter"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mknyszek/eep/font"
	"github.com/mknyszek/eep/geom"
)

// Style describes the style of text.
type Style struct {
	Face  *font.Face
	Color color.Color
}

// Basic creates a simple style from the provided font name, a font size,
// and a color. If the font doesn't exist, this falls back to eep's default
// font.
func Basic(fontName string, size float64, c color.Color) Style {
	s, ok := font.FindSource(fontName)
	if !ok {
		s = font.DefaultSource
	}
	return Style{
		Face:  font.NewFace(s, size),
		Color: c,
	}
}

// Apply applies the style to some text to produce a Segment of styled text.
func (s Style) Apply(text string) String {
	var b StringBuilder
	b.Append(text, s)
	return b.String()
}

// String is a string of styled text.
//
// All of the styles must have font face with an identical Direction.
type String struct {
	segments  []segment
	direction font.TextDirection
}

// Concat concatenates the two styled Strings and returns the result.
func (s String) Concat(t String) String {
	if s.direction != t.direction {
		panic("cannot concatenate two styled strings with different text direction")
	}
	segments := make([]segment, 0, len(s.segments)+len(t.segments))
	segments = append(segments, s.segments...)
	segments = append(segments, t.segments...)
	return String{
		direction: s.direction,
		segments:  segments,
	}
}

// segment is a text segment with no line breaks and a single consistent style.
type segment struct {
	text  string
	style Style
}

// Measure returns the dimensions of the String once rendered.
func (s String) Measure(lineSpacing float64) geom.Dimensions {
	var primary, secondary float64 // Totals.
	for segments := range s.lines() {
		// Figure out the increments for each line.
		incPri, incSec := measureLine(segments, s.direction, lineSpacing)

		// Accumulate.
		primary = max(primary, incPri)
		secondary += incSec
	}
	if s.direction == font.DirectionLeftToRight || s.direction == font.DirectionRightToLeft {
		return geom.Dim(primary, secondary)
	}
	return geom.Dim(secondary, primary)
}

func measureLine(segments []segment, direction font.TextDirection, lineSpacing float64) (priLen, secLen float64) {
	var primary, secondary float64
	for _, seg := range segments {
		face := seg.style.Face.TextFace()
		m := face.Metrics()
		if direction == text.DirectionLeftToRight || direction == text.DirectionRightToLeft {
			secondary = max(secondary, (m.HAscent+m.HDescent)*(lineSpacing+1.0))
		} else {
			secondary = max(secondary, (m.VAscent+m.VDescent)*(lineSpacing+1.0))
		}
		if seg.text != "\n" {
			primary += text.Advance(seg.text, face)
		}
	}
	return primary, secondary
}

// lines returns contiguous sections of s.segments with trailing newlines.
func (s String) lines() iter.Seq[[]segment] {
	return func(yield func([]segment) bool) {
		first := 0
		for i := range s.segments {
			if s.segments[i].text == "\n" {
				if !yield(s.segments[first : i+1]) {
					return
				}
				first = i + 1
			}
		}
		if !yield(s.segments[first:len(s.segments)]) {
			return
		}
	}
}

// StringBuilder is a builder for efficiently constructing
// styled strings. The zero value is ready for use.
type StringBuilder struct {
	s       String
	nonZero bool
}

// Append appends the provided text and applies the given style to it.
func (s *StringBuilder) Append(text string, style Style) {
	if !s.nonZero {
		s.s.direction = style.Face.TextFace().Direction
		s.nonZero = true
	}
	if s.s.direction != style.Face.TextFace().Direction {
		panic("cannot append different text direction to builder")
	}
	s.s.segments = appendSegmentsFromText(s.s.segments, text, style)
}

// String returns the builder's accumulated String.
func (s *StringBuilder) String() String {
	return s.s
}

// Reset resets the StringBuilder to be empty.
func (s *StringBuilder) Reset() {
	s.s = String{}
	s.nonZero = false
}

func appendSegmentsFromText(s []segment, text string, style Style) []segment {
	for line := range strings.Lines(text) {
		line, ok := strings.CutSuffix(line, "\n")
		s = append(s, segment{line, style})
		if ok {
			s = append(s, segment{"\n", style})
		}
	}
	return s
}

type noCopy struct{}

func (noCopy) Lock()   {}
func (noCopy) Unlock() {}
