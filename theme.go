package eep

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mknyszek/eep/geom"
	"github.com/mknyszek/eep/text"
)

type Theme struct {
	// Palette is the color palette for the theme.
	//
	// By convention, the first color is the background
	// color. The rest are all accent colors.
	Palette []color.Color

	// TextStyles is the set of text styles for the theme.
	//
	// By convention, the first text style is for the main
	// title, the second is for subtitles, and the third
	// is for secondary titles, and the fourth is for text
	// content.
	TextStyles []text.Style

	// Decor is a function that draws any theme-specific decor
	// to the slide. It is always drawn first.
	Decor func(*ebiten.Image)
}

func (t *Theme) Background(dst *ebiten.Image) {
	dst.Fill(t.Palette[0])
	t.Decor(dst)
}

func TitleSlide(t *Theme, title, subtitle string) Slide {
	return Static(func(screen *ebiten.Image) {
		d := geom.ImageDim(screen.Bounds())

		t.Background(screen)
		text.AutoBox(screen, t.TextStyles[0].Apply(title), geom.Pt(d.X/8, 3*d.Y/5), geom.Dim(0, 1))
		text.AutoBox(screen, t.TextStyles[1].Apply(subtitle), geom.Pt(d.X/8, 3*d.Y/5), geom.Dim(0, 0))
	})
}

func ContentSlide(t *Theme, title, content string) Slide {
	return Static(func(screen *ebiten.Image) {
		d := geom.ImageDim(screen.Bounds())

		t.Background(screen)
		text.AutoBox(screen, t.TextStyles[2].Apply(title), geom.Pt(d.X/16, 3*d.Y/16), geom.Dim(0, 1))
		text.AutoBox(screen, t.TextStyles[3].Apply(content), geom.Pt(d.X/16, 5*d.Y/16), geom.Dim(0, 0))
	})
}
