package scenes

import (
	"image/draw"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/font"
)

type (
	// GameManager ...
	GameManager struct {
		currentScene Scene
	}
)

var (
	fonts  map[string]font.Face
	images map[string]draw.Image

	eventHandler *EventHandler
)

func init() {
	fonts = make(map[string]font.Face)
	fonts["btnFont"] = FontLoad(16)

	images = make(map[string]draw.Image)
	images["bgImage"], _ = utils.GetImageByPath("resources/system_images/bg-2.jpg")
	images["btnBase"], _ = utils.GetImageByPath("resources/system_images/button.png")
	images["btnBaseHover"], _ = utils.GetImageByPath("resources/system_images/button-hover.png")

	eventHandler = &EventHandler{
		events: map[string]map[*Event]struct{}{},
	}

}

// NewGameManager ...
func NewGameManager(screenWidth, screenHeight int) *GameManager {
	width, height = screenWidth, screenHeight

	gm := &GameManager{}

	// MainMenuを表示
	gm.TransitionToMainMenu()
	return gm
}

// TransitionToMainMenu ...
func (g *GameManager) TransitionToMainMenu() {
	s := NewMainMenu(g)
	g.currentScene = s
}

// TransitionToBattleScene ...
func (g *GameManager) TransitionToBattleScene() {
	s := NewBattleScene(g)
	g.currentScene = s
}

// SetCurrentScene ...
func (g *GameManager) SetCurrentScene(s Scene) {
	g.currentScene = s
}

// Update ...
func (g *GameManager) Update(screen *ebiten.Image) error {
	if g.currentScene != nil {
		return g.currentScene.Update(screen)
	}
	return nil
}

// Draw ...
func (g *GameManager) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

// Layout ...
func (g *GameManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}
