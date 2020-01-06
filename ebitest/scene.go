package ebitest

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// CommonScene ...
	CommonScene struct {
		backgroundImage *ebiten.Image
		count           int
		limit           int
		debug           string
	}
)

var (
	mem runtime.MemStats
)

// NewCommonScene ...
func NewCommonScene(image *ebiten.Image) Scene {
	rand.Seed(time.Now().UnixNano()) //Seed

	return &CommonScene{
		backgroundImage: image,
		count:           0,
		limit:           300,
	}
}

// Update ...
func (s *CommonScene) Update(state *GameState) error {
	s.count++

	runtime.ReadMemStats(&mem)
	s.debug = fmt.Sprintf("\nalloc=%d, talloc=%d, heap-alloc=%d, heap-sys=%d\n", mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.HeapSys)

	change := false
	c, b := state.Input.Control()
	if b {
		s.debug += fmt.Sprintf("control=%s\n", c.String())
		if c.String() == "Z" {
			change = true
		}
	}

	if s.count > s.limit {
		change = true
	}

	if change {
		idx := rand.Intn(len(state.SceneManager.paths) - 1)
		s := NewCommonScene(state.SceneManager.PathToImage(idx))
		state.SceneManager.GoTo(s)
	}

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

	tx := float64(0)
	ty := float64(0)
	if x > y {
		tx = float64(sw - bw)
		ty = (float64(sw) - (float64(bw) * y)) / 2
		op.GeoM.Scale(y, y)
	} else {
		ty = (float64(sh) - (float64(bh) * x)) / 2
		op.GeoM.Scale(x, x)
	}
	op.GeoM.Translate(tx, ty)

	r.DrawImage(s.backgroundImage, op)
	ebitenutil.DebugPrint(r, fmt.Sprintf("x=%0.2f, y=%0.2f, tx=%0.2f, ty=%0.2f, count=%d", x, y, tx, ty, s.count))
	if s.debug != "" {
		ebitenutil.DebugPrint(r, fmt.Sprintf("\n%s", s.debug))
	}
}
