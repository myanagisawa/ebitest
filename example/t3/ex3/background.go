package ex3

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type (
	// BackGround ...
	BackGround struct {
		color color.Color
		image ebiten.Image
		size  Size
	}
)

// NewBackGround ...
func NewBackGround(s Size) (*BackGround, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	n := rand.Intn(3) + 1
	eimg := getImage(fmt.Sprintf("bg-%d.jpg", n), s.Width, s.Height)

	scene := &BackGround{
		image: *eimg,
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
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(0.5, 0.5, 0.5, 1.0)

	r.DrawImage(&s.image, op)
}

// GetSize ...
func (s *BackGround) GetSize() (int, int) {
	return s.size.Width, s.size.Height
}
