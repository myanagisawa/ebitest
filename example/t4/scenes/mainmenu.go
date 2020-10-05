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
		manager      *GameManager
		eventHandler *EventHandler
	}
)

var (
	btnBattle UIButton
)

// NewMainMenu ...
func NewMainMenu(m *GameManager) *MainMenu {

	s := &MainMenu{
		manager: m,
		eventHandler: &EventHandler{
			events: map[string]map[*Event]struct{}{},
		},
	}

	c := NewButton("Battle(dev)", images["btnBase"], fonts["btnFont"], color.Black, 100, 200)

	c.AddEventListener(s, "click", func(target UIController, source *EventSource) {
		log.Printf("btnBattle clicked")
		s := source.scene.(*MainMenu)
		s.manager.TransitionToBattleScene()
	})

	btnBattle = c

	return s
}

// SetEvent ...
func (s *MainMenu) SetEvent(name string, e *Event) {
	if s.eventHandler.events[name] != nil {
		s.eventHandler.events[name][e] = struct{}{}
	} else {
		m := map[*Event]struct{}{e: {}}
		s.eventHandler.events[name] = m
	}
}

// Update ...
func (s *MainMenu) Update(screen *ebiten.Image) error {
	btnBattle.(*UIButtonImpl).hover = btnBattle.In(ebiten.CursorPosition())

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if btnBattle.In(ebiten.CursorPosition()) {
			log.Printf("btnBattle clicked")
			s.manager.TransitionToBattleScene()
		}
	}

	return nil
}

// Draw ...
func (s *MainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{200, 200, 200, 255})

	btnBattle.Draw(screen)

}

// Layout ...
func (s *MainMenu) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
