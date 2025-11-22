package eep

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mknyszek/2d/ebiten/text"
	"github.com/mknyszek/2d/geom"
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
	// title, the second is for subtitles, the third
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
		text.AutoBox(screen, t.TextStyles[0].Apply(title).String(), geom.Pt(d.X/8, 3*d.Y/5), geom.Dim(0, 1))
		text.AutoBox(screen, t.TextStyles[1].Apply(subtitle).String(), geom.Pt(d.X/8, 3*d.Y/5), geom.Dim(0, 0))
	})
}

func SectionSlide(t *Theme, title string) Slide {
	return Static(func(screen *ebiten.Image) {
		t.Background(screen)
		text.Box(screen, t.TextStyles[2].Apply(title).String(), geom.ImageAABB(screen.Bounds()), &text.BoxOptions{
			Align:     text.Center,
			VertAlign: text.Middle,
		})
	})
}

func BlankContentSlide(t *Theme, title string) Slide {
	return Static(func(screen *ebiten.Image) {
		d := geom.ImageDim(screen.Bounds())

		t.Background(screen)
		text.AutoBox(screen, t.TextStyles[2].Apply(title).String(), geom.Pt(d.X/16, 3*d.Y/16), geom.Dim(0, 1))
	})
}

func ContentSlide(t *Theme, title string, makeContent func(text.Style) text.String) Slide {
	content := makeContent(t.TextStyles[3])
	return Overlay(
		BlankContentSlide(t, title),
		Static(func(screen *ebiten.Image) {
			d := geom.ImageDim(screen.Bounds())
			text.AutoBox(screen, content, geom.Pt(d.X/16, 4*d.Y/16), geom.Dim(0, 0))
		}),
	)
}

func BasicContentSlide(t *Theme, title, content string) Slide {
	return ContentSlide(t, title, func(style text.Style) text.String {
		return style.Apply(content).String()
	})
}

func BlankSlide(t *Theme) Slide {
	return Static(func(screen *ebiten.Image) {
		t.Background(screen)
	})
}
