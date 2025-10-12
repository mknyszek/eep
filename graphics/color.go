package graphics

import "image/color"

// RGB is a packed 24-bit opaque RGB color.
//
// The packed representation allows for easily
// writing out hex codes in plain Go.
// For example, RGB(0xff0000) is a plain red.
type RGB uint32

// RGBA implements color.Color.
func (c RGB) RGBA() (r, g, b, a uint32) {
	r = (uint32(c) >> 16) & 0xff
	g = (uint32(c) >> 8) & 0xff
	b = (uint32(c) >> 0) & 0xff
	a = 0xff
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a |= a << 8
	return
}

// WithAlpha adds an alpha channel to the color.
func (c RGB) WithAlpha(a uint8) color.RGBA {
	r := uint32(c) >> 16
	g := uint32(c) >> 8
	b := uint32(c) >> 0
	r = r * uint32(a) / 0xff
	g = g * uint32(a) / 0xff
	b = b * uint32(a) / 0xff
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// Black is the color black.
const Black = RGB(0)

// Wong is Bang Wong's colorblind-friendly discrete color palette.
// See https://www.nature.com/articles/nmeth.1618.
var Wong = struct {
	Black     RGB
	Orange    RGB
	LightBlue RGB
	Green     RGB
	Yellow    RGB
	Blue      RGB
	Amber     RGB
	Pink      RGB
}{
	Black:     Black,
	Orange:    RGB(0xE69F00),
	LightBlue: RGB(0x56B4E9),
	Green:     RGB(0x009E73),
	Yellow:    RGB(0xF0E442),
	Blue:      RGB(0x0072B2),
	Amber:     RGB(0xD55E00),
	Pink:      RGB(0xCC79A7),
}

// IBM Design Library's colorblind-friendly discrete color palette.
// See https://www.ibm.com/design/language/resources/color-library.
var IBM = struct {
	Blue   RGB
	Violet RGB
	Pink   RGB
	Orange RGB
	Yellow RGB
}{
	Blue:   RGB(0x648FFF),
	Violet: RGB(0x785EF0),
	Pink:   RGB(0xDC267F),
	Orange: RGB(0xFE6100),
	Yellow: RGB(0xFFB000),
}
