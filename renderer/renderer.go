package renderer

import (
	"image"
	"image/color"
)

func Render(scale int, pixels [][]color.Color) image.Image {
	size := len(pixels) * scale
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := range pixels {
		for x := range pixels[y] {
			for i := 0; i < scale; i++ {
				for j := 0; j < scale; j++ {
					img.Set(x*scale+i, y*scale+j, pixels[y][x])
				}
			}
		}
	}

	return img
}
