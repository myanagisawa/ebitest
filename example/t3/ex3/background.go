package ex3

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type (
	// BackGround ...
	BackGround struct {
		color color.Color
		size  Size
	}
)

// NewBackGround ...
func NewBackGround(s Size) (*BackGround, error) {
	scene := &BackGround{
		color: color.RGBA{37, 37, 37, 192},
		size:  s,
	}

	return scene, nil
}

// Update ...
func (s *BackGround) Update() error {
	return nil
}

// Draw ...
func (s *BackGround) Draw(r *ebiten.Image) {
	r.Fill(s.color)
}

// GetSize ...
func (s *BackGround) GetSize() (int, int) {
	return s.size.Width, s.size.Height
}
