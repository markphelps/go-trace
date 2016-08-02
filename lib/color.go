package primitives

import (
	"math"
)

type Color struct {
	R, G, B float64
}

var (
	Black = Color{}
	White = Color{1.0, 1.0, 1.0}
	Blue  = Color{0.5, 0.7, 1.0}
)

// For compatibility with image.Color
func (c Color) RGBA() (r, g, b, a uint32) {
	// Sqrt() for gamma-2 correction
	r = uint32(math.Sqrt(c.R) * 0xffff)
	g = uint32(math.Sqrt(c.G) * 0xffff)
	b = uint32(math.Sqrt(c.B) * 0xffff)
	a = 0xffff
	return
}

func (c Color) Add(o Color) Color {
	return Color{c.R + o.R, c.G + o.G, c.B + o.B}
}

func (c Color) Multiply(o Color) Color {
	return Color{c.R * o.R, c.G * o.G, c.B * o.B}
}

func (c Color) AddScalar(f float64) Color {
	return Color{c.R + f, c.G + f, c.B + f}
}

func (c Color) MultiplyScalar(f float64) Color {
	return Color{c.R * f, c.G * f, c.B * f}
}

func (c Color) DivideScalar(f float64) Color {
	return Color{c.R / f, c.G / f, c.B / f}
}

func Gradient(a, b Color, f float64) Color {
	// scale between 0.0 and 1.0
	f = 0.5 * (f + 1.0)

	// linear blend: blended_value = (1 - f) * a + f * b
	return a.MultiplyScalar(1.0 - f).Add(b.MultiplyScalar(f))
}
