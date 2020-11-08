package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models/scene"
)

type (
	// Manager ...
	Manager struct {
		currentScene interfaces.Scene
		scenes       map[enum.SceneEnum]interfaces.Scene
	}
)

// NewManager ...
func NewManager(screenWidth, screenHeight int) *Manager {
	ebitest.Width, ebitest.Height = screenWidth, screenHeight

	gm := &Manager{}

	scenes := map[enum.SceneEnum]interfaces.Scene{}
	scenes[enum.MainMenuEnum] = scene.NewMainMenu(gm)
	scenes[enum.MapEnum] = scene.NewMap(gm)
	gm.scenes = scenes

	// MainMenuを表示
	gm.TransitionTo(enum.MapEnum)
	return gm
}

// TransitionTo ...
func (g *Manager) TransitionTo(t enum.SceneEnum) {
	// var s interfaces.Scene
	// switch t {
	// case enum.MainMenuEnum:
	// 	s = scene.NewMainMenu(g)
	// case enum.MapEnum:
	// 	s = scene.NewMap(g)
	// default:
	// 	panic(fmt.Sprintf("invalid SceneEnum: %d", t))
	// }
	g.currentScene = g.scenes[t]
}

// SetCurrentScene ...
func (g *Manager) SetCurrentScene(s interfaces.Scene) {
	g.currentScene = s
}

// Update ...
func (g *Manager) Update() error {
	if g.currentScene != nil {
		return g.currentScene.Update()
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
