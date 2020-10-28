package kitchen

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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
		color: color.RGBA{37, 37, 37, 255},
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
