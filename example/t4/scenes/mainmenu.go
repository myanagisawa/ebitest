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
		manager *GameManager
	}
)

var (
	btnBattle UIButton
)

// NewMainMenu ...
func NewMainMenu(m *GameManager) *MainMenu {

	c := NewButton("Battle(dev)", images["btnBase"], fonts["btnFont"], color.Black, 100, 200)
	c.AddEventListener("click", func(target UIController, source *EventSource) {
		log.Printf("btnBattle clicked")
		s := source.scene.(*MainMenu)
		s.manager.TransitionToBattleScene()
	})
	btnBattle = c

	return &MainMenu{
		manager: m,
	}
}

// Update ...
func (g *MainMenu) Update(screen *ebiten.Image) error {
	btnBattle.(*UIButtonImpl).hover = btnBattle.In(ebiten.CursorPosition())

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if btnBattle.In(ebiten.CursorPosition()) {
			log.Printf("btnBattle clicked")
			g.manager.TransitionToBattleScene()
		}
	}

	return nil
}

// Draw ...
func (g *MainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{200, 200, 200, 255})

	btnBattle.Draw(screen)

}

// Layout ...
func (g *MainMenu) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
