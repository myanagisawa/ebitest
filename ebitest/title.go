package ebitest

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type (
	// TitleScene ...
	TitleScene struct {
		backgroundImage *ebiten.Image
		count           int
		debug           string
	}
)

// NewTitleScene ...
func NewTitleScene() Scene {
	img, _, err := ebitenutil.NewImageFromFile("bg2.png")
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

	s.debug = ""

	d, b := state.Input.Dir()
	if b {
		s.debug += fmt.Sprintf("dir=%s\n", d.String())
	}

	// c, b := state.Input.Control()
	// if b {
	// 	s.debug += fmt.Sprintf("control=%s\n", c.String())
	// 	if c.String() == "A" {
	// 		state.SceneManager.GoTo(NewCommonScene("bg3.png"))
	// 	}
	// 	if c.String() == "Z" {
	// 		state.SceneManager.GoTo(NewCommonScene("bg4.png"))
	// 	}
	// }
	return nil
}

// Draw ...
func (s *TitleScene) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	sw, sh := r.Size()
	bw, bh := s.backgroundImage.Size()
	// x := (sw - bw) / 2
	// y := (sh - bh) / 2
	//	op.GeoM.Translate(float64(sw), float64(sh))
	x := float64(sw) / float64(bw)
	y := float64(sh) / float64(bh)
	op.GeoM.Scale(x, y)

	r.DrawImage(s.backgroundImage, op)
	ebitenutil.DebugPrint(r, fmt.Sprintf("x=%0.2f, y=%0.2f", x, y))
	if s.debug != "" {
		ebitenutil.DebugPrint(r, fmt.Sprintf("\n%s", s.debug))
	}
}
