package eep

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func Present(d *Deck) error {
	ebiten.SetWindowSize(d.width, d.height)
	ebiten.SetWindowTitle("Ebitengine Presentation")
	return ebiten.RunGame(d)
}

type Deck struct {
	width, height int
	currSlide     int
	slides        []Slide
}

func NewDeck(width, height int) *Deck {
	return &Deck{width: width, height: height}
}

func (d *Deck) Append(s ...Slide) {
	d.slides = append(d.slides, s...)
}

func (d *Deck) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return d.width, d.height
}

func (d *Deck) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight):
		if d.currSlide > 0 {
			d.currSlide--
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		if d.currSlide < len(d.slides)-1 {
			d.currSlide++
		}
	}
	return d.slides[d.currSlide].Update()
}

func (d *Deck) Draw(screen *ebiten.Image) {
	d.slides[d.currSlide].Draw(screen)
}
