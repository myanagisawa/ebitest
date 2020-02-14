package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

var (
	img *ebiten.Image
)

func update(r *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}

	r.DrawImage(img, op)
	return nil
}

func main() {

	colors := []color.RGBA{
		color.RGBA{255, 100, 100, 120},
		color.RGBA{100, 255, 100, 120},
		color.RGBA{100, 100, 255, 120},
		color.RGBA{220, 220, 220, 120},
	}
	lines := []color.RGBA{
		color.RGBA{255, 200, 200, 255},
		color.RGBA{200, 255, 200, 255},
		color.RGBA{200, 200, 255, 255},
		color.RGBA{220, 220, 220, 255},
	}

	w, h := 160, 40
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		tilenum := x / h
		tilewidth := x % h
		log.Printf("x: %d (num=%d, w=%d)", x, tilenum, tilewidth)
		for y := 0; y < h; y++ {
			if (y == 2 || y == 37) && (tilewidth >= 2 && tilewidth <= 37) {
				m.Set(x, y, lines[tilenum])
			} else if (tilewidth == 2 || tilewidth == 37) && (y >= 2 && y <= 37) {
				m.Set(x, y, lines[tilenum])
			} else {
				m.Set(x, y, colors[tilenum])
			}
		}
	}
	img, _ = ebiten.NewImageFromImage(m, ebiten.FilterDefault)

	if err := ebiten.Run(update, 200, 200, 1.0, "t4"); err != nil {
		log.Fatal(err)
	}
}
