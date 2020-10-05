package scenes

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (
	// MainMenu ...
	MainMenu struct {
		SceneBase
	}
)

var (
	btnBattle UIButton
)

// NewMainMenu ...
func NewMainMenu(m *GameManager) *MainMenu {

	s := &MainMenu{
		SceneBase: SceneBase{
			manager: m,
			eventHandler: &EventHandler{
				events: map[string]map[*Event]struct{}{},
			},
		},
	}

	l := NewLayerBase(s)
	s.SetLayer(l)

	c := NewButton("Battle(dev)", images["btnBase"], fonts["btnFont"], color.Black, 100, 200)

	c.AddEventListener(s, "click", func(target UIController, source *EventSource) {
		log.Printf("btnBattle clicked")
		source.scene.Manager().TransitionToBattleScene()
	})

	btnBattle = c

	return s
}

// Update ...
func (s *MainMenu) Update(screen *ebiten.Image) error {
	s.activeLayer = s.LayerAt(ebiten.CursorPosition())
	if s.activeLayer != nil {
		// log.Printf("activeLayer: %#v", s.activeLayer.Label())

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// click イベントを発火
			x, y := s.activeLayer.LocalPosition(ebiten.CursorPosition())
			s.eventHandler.Firing(s, "click", x, y)
		}
	}

	btnBattle.(*UIButtonImpl).hover = btnBattle.In(ebiten.CursorPosition())

	return nil
}

// Draw ...
func (s *MainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{200, 200, 200, 255})

	btnBattle.Draw(screen)

}
