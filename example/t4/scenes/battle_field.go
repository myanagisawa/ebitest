package scenes

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (

	// BattleScene ...
	BattleScene struct {
		SceneBase
	}
)

// NewBattleScene ...
func NewBattleScene(m *GameManager) Scene {

	s := &BattleScene{
		SceneBase: SceneBase{
			manager: m,
		},
	}

	s.SetLayer(NewLayerBase(s))
	s.SetLayer(NewBattleMap(s))
	s.SetLayer(NewTestWindow(s))

	return s
}

// Update ...
func (s *BattleScene) Update() error {
	s.activeLayer = s.LayerAt(ebiten.CursorPosition())
	if s.activeLayer != nil {
		// log.Printf("activeLayer: %#v", s.activeLayer.Label())

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// click イベントを発火
			s.activeLayer.FiringEvent("click")
		}
	}

	for _, layer := range s.layers {
		layer.Update()
	}

	return nil
}

// Draw ...
func (s *BattleScene) Draw(screen *ebiten.Image) {
	for _, layer := range s.layers {
		layer.Draw(screen)
	}
	active := " - "
	if s.activeLayer != nil {
		active = s.activeLayer.Label()
	}

	dbg := fmt.Sprintf("FPS: %0.2f\nactive: %s", ebiten.CurrentFPS(), active)
	ebitenutil.DebugPrint(screen, dbg)
}

// Layout ...
func (s *BattleScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
