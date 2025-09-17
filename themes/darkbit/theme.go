package darkbit

import (
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mknyszek/eep"
	"github.com/mknyszek/eep/font"
	"github.com/mknyszek/eep/geom"
	"github.com/mknyszek/eep/graphics"
	"github.com/mknyszek/eep/text"
)

var Theme *eep.Theme

//go:embed m6x11plus.ttf
var titleTTF []byte

//go:embed m5x7.ttf
var textTTF []byte

func init() {
	// Fonts.
	titleFont := must(font.NewSourceFromBytes(titleTTF))
	textFont := must(font.NewSourceFromBytes(textTTF))

	// Palette.
	palette := []color.Color{color.Gray{Y: 32}, color.Gray{Y: 28}}
	Theme = &eep.Theme{
		Palette: palette,
		TextStyles: []text.Style{
			{Face: font.NewFace(titleFont, 108), Color: color.White},
			{Face: font.NewFace(textFont, 64), Color: color.White},
			{Face: font.NewFace(titleFont, 72), Color: color.White},
			{Face: font.NewFace(textFont, 48), Color: color.White},
		},
		Decor: func(screen *ebiten.Image) {
			border := geom.ImageAABB(screen.Bounds())
			padding := border.Dy() / 20
			thickness := padding / 5
			border.Min.X += padding
			border.Min.Y += padding
			border.Max.X -= padding
			border.Max.Y -= padding

			c := graphics.NewContext(screen)
			c.SetColor(palette[1])
			c.SetLineWidth(thickness)
			c.SetLineCap(graphics.LineCapSquare)
			c.SetLineJoin(graphics.LineJoinMiter(10))
			c.Rect(graphics.Stroke, border)
		},
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
