package eep

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Slide represents a single slide in the slide deck.
//
// Slide is a subset of the ebiten.Game interface. This is intentional,
// since it has the same split between update logic and draw logic.
//
// See [github.com/hajimehoshi/ebiten/v2.Game] for more context and
// details.
// See https://ebitencookbook.vercel.app/blog for an explanation of
// the split between Update and Draw.
type Slide interface {
	// Draw renders the slide deck to the screen.
	//
	// Simple, static slides need only implement Draw; Update can be a no-op.
	//
	// Draw is called in lock-step with the refresh rate of the screen.
	// It does not guarantee a fixed time delta between calls, and so
	// is not ideal for any logic that causes your slide to evolve over
	// time, especially if your logic includes any kind of physical
	// simulation.
	//
	// Note that it's still feasible to do so if you wish to compute the
	// time delta between frames, but this is likely much more complicated
	// than you need for a slideshow presentation.
	Draw(*ebiten.Image)

	// Update is called once per tick.
	//
	// Update is where to put logic that influences what is drawn
	// to the screen. For example, anything depending on time evolution,
	// like a physical simulation.
	//
	// Ebiten will arrange for update to be called, in the long run,
	// on average, 60 times per second. As a result, Update may
	// assume that 1/60th of a second passes between each call.
	//
	// If Update returns a non-nil error, the presentation will exit.
	// Update may return ebiten.Termination for a clean exit.
	Update() error
}

// Static returns a slide that only has a Draw command.
//
// Without an Update function, the Slide is not expected
// to evolve over time, hence the name.
func Static(f func(screen *ebiten.Image)) Slide {
	return staticSlide{f}
}

type staticSlide struct {
	draw func(screen *ebiten.Image)
}

func (s staticSlide) Draw(screen *ebiten.Image) {
	s.draw(screen)
}

func (s staticSlide) Update() error {
	return nil
}

// Dynamic returns a slide that has both a Draw and Update command.
func Dynamic(draw func(screen *ebiten.Image), update func() error) Slide {
	return dynamicSlide{draw, update}
}

type dynamicSlide struct {
	draw   func(screen *ebiten.Image)
	update func() error
}

func (s dynamicSlide) Draw(screen *ebiten.Image) {
	s.draw(screen)
}

func (s dynamicSlide) Update() error {
	return s.update()
}

// Overlay returns a slide that draws dst, then src0 on top
// of it, then src on top of that, etc.
//
// The returned slide will also execute all slides' Update
// methods in its Update method.
func Overlay(dst, src0 Slide, src ...Slide) Slide {
	return slideStack{slides: append([]Slide{dst, src0}, src...)}
}

type slideStack struct {
	slides []Slide
}

func (s slideStack) Draw(screen *ebiten.Image) {
	for _, slide := range s.slides {
		slide.Draw(screen)
	}
}

func (s slideStack) Update() error {
	for _, slide := range s.slides {
		if err := slide.Update(); err != nil {
			return err
		}
	}
	return nil
}
