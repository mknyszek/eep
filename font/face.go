package font

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"sync"
	"weak"

	"github.com/adrg/sysfont"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

// Face represents a configuration of the display of text.
type Face struct {
	face text.GoTextFace
}

// NewFace creates a new font face that can be used to draw text.
func NewFace(src Source, size float64, opts ...Option) *Face {
	f := new(Face)
	f.face.Source = src.src
	f.face.Size = size
	return f
}

// Option represents additional optional configuration for a Face.
type Option struct {
	f func(*Face)
}

type TextDirection = text.Direction

const (
	DirectionLeftToRight               = text.DirectionLeftToRight
	DirectionRightToLeft               = text.DirectionRightToLeft
	DirectionTopToBottomAndLeftToRight = text.DirectionTopToBottomAndLeftToRight
	DirectionTopToBottomAndRightToLeft = text.DirectionTopToBottomAndRightToLeft
)

// Direction sets the Face's rendering direction.
func Direction(d TextDirection) Option {
	return Option{func(f *Face) {
		f.face.Direction = d
	}}
}

// Language sets the Face's language hint.
func Language(t language.Tag) Option {
	return Option{func(f *Face) {
		f.face.Language = t
	}}
}

// Source returns the Source for the Face.
func (f *Face) Source() Source {
	return Source{f.face.Source}
}

// Size returns the size of the Face.
func (f *Face) Size() float64 {
	return f.face.Size
}

// Resize returns a new Face with all the same features except with the
// font size changed to the provided size.
func (f *Face) Resize(size float64) *Face {
	g := new(Face)
	*g = *f
	g.face.Size = size
	return g
}

// LineSize is the amount of vertical or horizontal space (depending on
// the face's Direction) takes up on a line. lineSpacing is the relative
// amount of additional spacing to provide the line. For example, 0.0 is
// single-spaced and 1.0 is double-spaced. Negative lineSpacing reduces
// the line size.
func (f *Face) LineSize(lineSpacing float64) float64 {
	m := f.face.Metrics()
	if f.face.Direction == text.DirectionLeftToRight || f.face.Direction == text.DirectionRightToLeft {
		return (m.HAscent + m.HDescent) * (lineSpacing + 1.0)
	}
	return (m.VAscent + m.VDescent) * (lineSpacing + 1.0)
}

// TextFace returns the underlying Ebiten GoTextFace.
//
// Mutating the result will also mutate Face.
func (f *Face) TextFace() *text.GoTextFace {
	return &f.face
}

var sourceRegistry sync.Map // string -> weak.Pointer[text.GoTextFaceSource]

// FindSource looks first for a pre-registered font source, registered by RegisterSource,
// and if that fails, searches the system for related fonts (via fuzzy match), then
// registers and returns that font source.
//
// Safe to call from multiple goroutines simultaneously.
func FindSource(name string) (Source, bool) {
	// Try the sourceRegistry.
	if a, ok := sourceRegistry.Load(name); ok {
		if s := a.(weak.Pointer[text.GoTextFaceSource]).Value(); s != nil {
			return Source{s}, true
		}
	}

	// Try to look up a system font.
	sysFonts.mu.Lock()
	sf := sysFonts.finder.Match(name)
	if sf == nil {
		sysFonts.mu.Unlock()
		return Source{}, false
	}
	sysFonts.mu.Unlock()

	// Try the sourceRegistry for the full name.
	if a, ok := sourceRegistry.Load(sf.Name); ok {
		if s := a.(weak.Pointer[text.GoTextFaceSource]).Value(); s != nil {
			return Source{s}, true
		}
	}

	// Load the system font.
	s, err := NewSourceFromFile(sf.Filename)
	if err != nil {
		return Source{}, false
	}
	RegisterSource(sf.Name, s)
	return s, true
}

// RegisterSource adds a source to the registry under the provided name.
//
// Overrides any system fonts.
// Safe to call from multiple goroutines simultaneously.
func RegisterSource(name string, s Source) {
	wp := weak.Make(s.src)
	type entry struct {
		name string
		wp   weak.Pointer[text.GoTextFaceSource]
	}
	runtime.AddCleanup(s.src, func(e entry) {
		sourceRegistry.CompareAndDelete(e.name, e.wp)
	}, entry{name, wp})
	sourceRegistry.Store(name, wp)
}

// Source is a font source used to create Faces, which are used to draw text.
type Source struct {
	src *text.GoTextFaceSource
}

// NewSourceFromBytes creates a new font source from the bytes of an OTF or TTF file.
func NewSourceFromBytes(ttf []byte) (Source, error) {
	return NewSource(bytes.NewReader(ttf))
}

// NewSourceFromFile create a new font source from an OTF or TTF file.
func NewSourceFromFile(filename string) (Source, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Source{}, err
	}
	defer f.Close()
	return NewSource(f)
}

// NewSource creates a new font source from an io.Reader whose stream must be an OTF or TTF-formatted file.
func NewSource(r io.Reader) (Source, error) {
	src, err := text.NewGoTextFaceSource(r)
	if err != nil {
		return Source{}, err
	}
	return Source{src: src}, nil
}

// DefaultSource is a Source guaranteed to exist that may be used as a fallback.
var DefaultSource Source

var sysFonts struct {
	mu     sync.Mutex
	finder *sysfont.Finder
}

func init() {
	// Set up a font finder.
	sysFonts.finder = sysfont.NewFinder(&sysfont.FinderOpts{Extensions: []string{".ttf", ".otf"}})

	// Try to set a default font.
	for _, name := range []string{"Arial", "Helvetica", "Times New Roman", "Times", "Courier New", "Courier"} {
		def, ok := FindSource(name)
		if ok {
			DefaultSource = def
			break
		}
	}
	if DefaultSource.src == nil {
		println("failed to set a default font")
	}
}
