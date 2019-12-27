package ebitest

import (
	"github.com/hajimehoshi/ebiten"
)

const (
	// ScreenWidth ...
	ScreenWidth = 576
	// ScreenHeight ...
	ScreenHeight = 360
)

type (
	// Game ...
	Game struct {
		input        *Input
		sceneManager *SceneManager
	}
)

// NewGame ...
func NewGame() (*Game, error) {
	g := &Game{
		input:        NewInput(),
		sceneManager: &SceneManager{},
	}
	g.sceneManager.GoTo(NewTitleScene())
	return g, nil
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
	g.input.Update()
	if err := g.sceneManager.Update(g.input); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.sceneManager.Draw(r)
	return nil
}
