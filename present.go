package eep

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Present launches a new window and starts the presentation.
//
// The launched presentation provides basic keyboard and mouse
// controls for navigating the SlideDeck. The right arrow key
// and left mouse button will advance the SlideDeck forward by
// calling Next, while the left arrow key and right mouse button
// will advance the SlideDeck backward by calling Prev.
//
// Present must be called from the program's main thread. Since
// eep is imported, the main goroutine will be locked to the
// main thread, so this must be called from the main goroutine.
//
// It does not return unless the window has exited or deck's
// Update returns a non-nil error.
func Present(width, height int, deck SlideDeck) error {
	p := &presentation{width, height, deck}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebitengine Presentation")
	return ebiten.RunGame(p)
}

type presentation struct {
	width, height int
	deck          SlideDeck
}

func (p *presentation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return p.width, p.height
}

func (p *presentation) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight):
		p.deck.Prev()
	case inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		p.deck.Next()
	}
	return p.deck.Update()
}

func (p *presentation) Draw(screen *ebiten.Image) {
	p.deck.Draw(screen)
}

// SlideDeck is a Slide that can advance forward or backward.
type SlideDeck interface {
	Slide

	// Next advances the deck forward in response to mouse/keyboard input.
	//
	// Returns true on success. A call to Prev might fail to advance if,
	// for example, there's no next slide (already on the last slide).
	Next() DeckStatus

	// Prev advances the deck backward in response to mouse/keyboard input.
	//
	// Returns true on success. A call to Prev might fail to advance if,
	// for example, there's no previous slide (already on the first slide).
	Prev() DeckStatus
}

// DeckStatus is the status of the SlideDeck after calling Next or Prev.
type DeckStatus int

const (
	// DeckOK indicates that the Next or Prev call succeeded.
	DeckOK DeckStatus = iota

	// DeckBusy indicates that the Next or Prev call failed because
	// the Deck is currently busy showing the viewer something.
	// Typically, this is used to indicate that some animated
	// transition is playing back.
	DeckBusy

	// DeckEnd indicates that the Next or Prev call failed because there
	// are no more slides after or before respectively.
	DeckEnd
)

// StaticDeck is a SlideDeck that contains a static set of Slides.
type StaticDeck struct {
	currSlide int
	slides    []Slide
}

// Append adds a Slide to the deck.
func (d *StaticDeck) Append(s ...Slide) {
	d.slides = append(d.slides, s...)
}

// Next advances the slide deck forward.
func (d *StaticDeck) Next() DeckStatus {
	if d.currSlide < len(d.slides)-1 {
		d.currSlide++
		return DeckOK
	}
	return DeckEnd
}

// Prev advances the slide deck backward.
func (d *StaticDeck) Prev() DeckStatus {
	if d.currSlide > 0 {
		d.currSlide--
		return DeckOK
	}
	return DeckEnd
}

// Update implements Slide by updating the current Slide.
func (d *StaticDeck) Update() error {
	return d.slides[d.currSlide].Update()
}

// Draw implements Slide by drawing the current Slide.
func (d *StaticDeck) Draw(screen *ebiten.Image) {
	d.slides[d.currSlide].Draw(screen)
}

// ChainDeck chains together one or more slide decks.
type ChainDeck struct {
	decks    []SlideDeck
	currDeck int
}

// Append adds a Slide to the deck.
func (d *ChainDeck) Append(s ...SlideDeck) {
	d.decks = append(d.decks, s...)
}

// Next advances the slide deck forward.
func (d *ChainDeck) Next() DeckStatus {
	if status := d.decks[d.currDeck].Next(); status != DeckEnd {
		return status
	}
	if d.currDeck < len(d.decks)-1 {
		d.currDeck++
		return DeckOK
	}
	return DeckEnd
}

// Prev advances the slide deck backward.
func (d *ChainDeck) Prev() DeckStatus {
	if status := d.decks[d.currDeck].Prev(); status != DeckEnd {
		return status
	}
	if d.currDeck > 0 {
		d.currDeck--
		return DeckOK
	}
	return DeckEnd
}

// Update implements Slide by updating the current Slide.
func (d *ChainDeck) Update() error {
	return d.decks[d.currDeck].Update()
}

// Draw implements Slide by drawing the current Slide.
func (d *ChainDeck) Draw(screen *ebiten.Image) {
	d.decks[d.currDeck].Draw(screen)
}
