package ebitest

import (
	"log"

	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// Game ...
	Game struct {
		input      *Input
		scene      *Scene
		sceneImage *ebiten.Image
	}
)

// NewGame ...
func NewGame() (*Game, error) {
	g := &Game{
		input: NewInput(),
	}

	var err error
	g.scene, err = NewScene()
	if err != nil {
		return nil, err
	}

	return g, nil
}

// Update ...
func (g *Game) Update() error {
	g.input.Update()
	if err := g.scene.Update(g.input); err != nil {
		return err
	}

	//	log.Printf("Game: Update: Input: %#v", g.input)
	return nil
}

// Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	log.Printf("Game: Draw: ")

	if g.sceneImage == nil {
		w, h := g.scene.Size()
		g.boardImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n%d, %d, %d, %d", g.input.mouseState, g.input.mouseInitPosX, g.input.mouseInitPosY, g.input.mouseDir))
	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n%d, %d", x, y))
}
