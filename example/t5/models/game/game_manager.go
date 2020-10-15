package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models/scene"
)

type (
	// Manager ...
	Manager struct {
		currentScene interfaces.Scene
	}
)

// NewManager ...
func NewManager(screenWidth, screenHeight int) *Manager {
	ebitest.Width, ebitest.Height = screenWidth, screenHeight

	gm := &Manager{}

	// MainMenuを表示
	gm.TransitionTo(enum.MainMenuEnum)
	return gm
}

// TransitionTo ...
func (g *Manager) TransitionTo(t enum.SceneEnum) {
	var s interfaces.Scene
	switch t {
	case enum.MainMenuEnum:
		s = scene.NewMainMenu(g)
	default:
		panic(fmt.Sprintf("invalid SceneEnum: %d", t))
	}
	g.currentScene = s
}

// SetCurrentScene ...
func (g *Manager) SetCurrentScene(s interfaces.Scene) {
	g.currentScene = s
}

// Update ...
func (g *Manager) Update(screen *ebiten.Image) error {
	if g.currentScene != nil {
		return g.currentScene.Update(screen)
	}
	return nil
}

// Draw ...
func (g *Manager) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

// Layout ...
func (g *Manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebitest.Width, ebitest.Height
}
