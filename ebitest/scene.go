package ebitest

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// CommonScene ...
	CommonScene struct {
		backgroundImage *ebiten.Image
		count           int
		debug           string
	}
)

// NewCommonScene ...
func NewCommonScene(fpath string) Scene {
	img, _, err := ebitenutil.NewImageFromFile(fpath, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	return &TitleScene{
		backgroundImage: img,
		count:           0,
	}
}

// Update ...
func (s *CommonScene) Update(state *GameState) error {
	s.count++
	return nil
}

// Draw ...
func (s *CommonScene) Draw(r *ebiten.Image) {
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
