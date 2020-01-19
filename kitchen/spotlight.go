package kitchen

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type (
	// Spotlight ...
	Spotlight struct {
		image  *ebiten.Image
		circle Circle
		color  int
		x      float64
		y      float64
	}
)

// NewSpotlight ...
func NewSpotlight(x, y, r float64, c int) (*Spotlight, error) {
	circle := Circle{r, r, r}
	img := image.NewRGBA(image.Rect(0, 0, int(r*2), int(r*2)))
	for x := 0; x < int(r*2); x++ {
		for y := 0; y < int(r*2); y++ {
			a := circle.Brightness(float64(x), float64(y))
			c := color.RGBA{0, 0, 0, 255 - a}
			img.Set(x, y, c)
		}
	}

	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	s := &Spotlight{
		image:  eimg,
		circle: circle,
		color:  c,
		x:      x,
		y:      y,
	}
	return s, nil
}

// Update ...
func (s *Spotlight) Update() error {

	return nil
}

// Draw ...
func (s *Spotlight) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.x, s.y)
	r.DrawImage(s.image, op)
}
