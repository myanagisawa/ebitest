package ex3

import (
	"github.com/hajimehoshi/ebiten"
)

type (
	// BattleScene ...
	BattleScene struct {
		size Size
	}
)

// NewBattleScene ...
func NewBattleScene(s Size) *BattleScene {
	scene := &BattleScene{
		size: s,
	}
	return scene
}

// Update ...
func (s *BattleScene) Update() error {

	//log.Printf("BattleScene.Update")

	return nil
}

// Draw ...
func (s *BattleScene) Draw(r *ebiten.Image) {
}

// GetSize ...
func (s *BattleScene) GetSize() (int, int) {
	return s.size.Width, s.size.Height
}
