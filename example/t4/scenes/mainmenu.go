package scenes

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	// MainMenu ...
	MainMenu struct {
		SceneBase
	}
)

// NewMainMenu ...
func NewMainMenu(m *GameManager) *MainMenu {

	s := &MainMenu{
		SceneBase: SceneBase{
			manager: m,
		},
	}

	l := NewLayerBase(s)
	s.SetLayer(l)

	c := NewButton("Battle(dev)", images["btnBase"], fonts["btnFont"], color.Black, 100, 200)
	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("btnBattle clicked")
		source.scene.Manager().TransitionToBattleScene()
	})

	return s
}

// Update ...
func (s *MainMenu) Update() error {
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
func (s *MainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{200, 200, 200, 255})

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
