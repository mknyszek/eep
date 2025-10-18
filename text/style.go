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

// Apply applies the style to some text to produce a Piece of styled text.
func (s Style) Apply(text string) Piece {
	return Piece{text, s}
}

// Recolor returns a new style that's the same as this one but with a different color.
func (s Style) Recolor(c color.Color) Style {
	newStyle := s
	newStyle.Color = c
	return newStyle
}

// String is a string of styled text.
//
// All of the styles must have font face with an identical Direction.
type String struct {
	segments []segment
}

// direction returns the canonical direction of the String.
func (s String) direction() font.TextDirection {
	if len(s.segments) == 0 {
		return 0
	}
	return s.segments[0].style.Face.TextFace().Direction
}

// Concat concatenates the two styled Strings and returns the result.
func (s String) Concat(t String) String {
	if s.direction() != t.direction() {
		panic("cannot concatenate two styled strings with different text direction")
	}
	segments := make([]segment, 0, len(s.segments)+len(t.segments))
	segments = append(segments, s.segments...)
	segments = append(segments, t.segments...)
	return String{
		segments: segments,
	}
}

// segment is a text segment with no line breaks and a single consistent style.
//
// Distinct from Piece because Piece doesn't enforce the "no line breaks" rule.
type segment struct {
	text  string
	style Style
}

// Measure returns the dimensions of the String once rendered.
func (s String) Measure(lineSpacing float64) geom.Dimensions {
	var primary, secondary float64 // Totals.
	for segments := range s.lines() {
		// Figure out the increments for each line.
		incPri, incSec := measureLine(segments, lineSpacing)

		// Accumulate.
		primary = max(primary, incPri)
		secondary += incSec
	}
	if dir := s.direction(); dir == font.DirectionLeftToRight || dir == font.DirectionRightToLeft {
		return geom.Dim(primary, secondary)
	}
	return geom.Dim(secondary, primary)
}

// measureLine measures the primary and secondary dimensions of the segments, which
// must represent a single line. All the segments must have the same direction.
//
// The segments must contain only a single newline: at the end.
func measureLine(segments []segment, lineSpacing float64) (priLen, secLen float64) {
	var primary, secondary float64
	for _, seg := range segments {
		secondary = max(secondary, seg.style.Face.LineSize(lineSpacing))
		if seg.text != "\n" {
			primary += text.Advance(seg.text, seg.style.Face.TextFace())
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

// Piece represents a piece of text with a consistent style.
type Piece struct {
	Text  string
	Style Style
}

// String creates a new String containing only the piece.
func (p Piece) String() String {
	var b StringBuilder
	b.Append(p)
	return b.String()
}

// Concat concatenates a series of Pieces to form a String.
func Concat(pieces ...Piece) String {
	var b StringBuilder
	for _, p := range pieces {
		b.Append(p)
	}
	return b.String()
}

// StringBuilder is a builder for efficiently constructing
// styled strings. The zero value is ready for use.
type StringBuilder struct {
	segments  []segment
	direction font.TextDirection
	nonZero   bool
}

// Append appends the provided piece.
func (s *StringBuilder) Append(piece Piece) {
	if !s.nonZero {
		s.direction = piece.Style.Face.TextFace().Direction
		s.nonZero = true
	}
	if s.direction != piece.Style.Face.TextFace().Direction {
		panic("cannot append different text direction to builder")
	}
	s.segments = appendSegmentsFromText(s.segments, piece.Text, piece.Style)
}

// String returns the builder's accumulated String.
func (s *StringBuilder) String() String {
	return String{s.segments}
}

// Reset resets the StringBuilder to be empty.
func (s *StringBuilder) Reset() {
	s.segments = nil
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
