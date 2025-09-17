package eep

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Slide represents a single slide in the slide deck.
type Slide interface {
	Draw(*ebiten.Image)
	Update() error
}

// Static returns a slide that only has a Draw command.
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
