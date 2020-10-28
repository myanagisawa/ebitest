package scenes

import (
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
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
)

func init() {
	fonts = make(map[string]font.Face)
	fonts["btnFont"] = FontLoad(16)

	images = make(map[string]draw.Image)
	images["bgImage"], _ = utils.GetImageByPath("resources/system_images/bg-2.jpg")
	images["btnBase"], _ = utils.GetImageByPath("resources/system_images/button.png")
	images["btnBaseHover"], _ = utils.GetImageByPath("resources/system_images/button-hover.png")

	img := createRectImage(10, 10, color.RGBA{0, 0, 0, 255})
	images["listBase"] = img.(draw.Image)

	img = createRectImage(10, 10, color.RGBA{128, 128, 128, 128})
	images["listScroller"] = img.(draw.Image)

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
func (g *GameManager) Update() error {
	if g.currentScene != nil {
		return g.currentScene.Update()
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
