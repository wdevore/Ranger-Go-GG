package engine

import (
	"image"
	"image/color"
)

// DrawRect renderers a rectangle
func DrawRect(x, y float64, w, h int, centered bool, color color.RGBA, pixels *image.RGBA) {
	if centered {
		x -= float64(w) / 2.0
		y -= float64(h) / 2.0
	}

	for xi := x; xi < x+float64(w); xi++ {
		for yi := y; yi < y+float64(h); yi++ {
			pixels.SetRGBA(int(xi), int(yi), color)
		}
	}
}
