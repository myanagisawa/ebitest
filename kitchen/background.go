package kitchen

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type (
	// BackGround ...
	BackGround struct {
		color color.Color
	}
)

// NewBackGround ...
func NewBackGround() (*BackGround, error) {
	s := &BackGround{
		color: color.RGBA{0x90, 0x7e, 0xb4, 0xff},
	}

	return s, nil
}

// Update ...
func (s *BackGround) Update() error {
	return nil
}

// Draw ...
func (s *BackGround) Draw(r *ebiten.Image) {
	r.Fill(s.color)
}
