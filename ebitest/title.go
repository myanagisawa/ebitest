package ebitest

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// TitleScene ...
	TitleScene struct {
		backgroundImage *ebiten.Image
		count           int
	}
)

// NewTitleScene ...
func NewTitleScene() Scene {
	img, _, err := ebitenutil.NewImageFromFile("bg1.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	return &TitleScene{
		backgroundImage: img,
		count:           0,
	}
}

// Update ...
func (s *TitleScene) Update(state *GameState) error {
	s.count++
	return nil
}

// Draw ...
func (s *TitleScene) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	r.DrawImage(s.backgroundImage, op)
}
