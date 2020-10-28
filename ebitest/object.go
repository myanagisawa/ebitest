package ebitest

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type (
	// Object ...
	Object struct {
		image    *ebiten.Image
		op       *ebiten.DrawImageOptions
		lifetime int
	}
)

// NewObject ...
func NewObject(i *ebiten.Image, bw, bh int) Scene {
	rand.Seed(time.Now().UnixNano()) //Seed

	scale := float64(rand.Intn(4)+1) / 20
	//log.Printf("scale=%f", scale)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scale, scale)

	// 描画オブジェクトのサイズを取得
	w, h := i.Size()
	iw := float64(w) * scale
	ih := float64(h) * scale
	// 描画位置を計算
	x := rand.Float64()*float64(bw) - iw
	if x < 0 {
		x = 0
	}
	y := rand.Float64()*float64(bh) - ih
	if y < 0 {
		y = 0
	}

	op.GeoM.Translate(x, y)

	return &Object{
		image:    i,
		op:       op,
		lifetime: rand.Intn(500) + transitionMaxCount,
	}
}

// Update ...
func (s *Object) Update(state *GameState) error {
	s.lifetime--
	return nil
}

// Draw ...
func (s *Object) Draw(r *ebiten.Image) {
	if s.lifetime > transitionMaxCount {
		r.DrawImage(s.image, s.op)
		return
	}

	// transition
	alpha := 1 - float64(transitionMaxCount-s.lifetime)/float64(transitionMaxCount)
	op := s.op
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(s.image, op)

}
