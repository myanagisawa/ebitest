package main

import (
	_ "image/jpeg"
	"log"

	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/myanagisawa/ebitest/ebitest"
)

const (
	screenWidth  = 320
	screenHeight = 240

	screenScale = 2.0
)

var (
	count = 0

	game *ebitest.Game
)

func update(screen *ebiten.Image) error {
	count++

	if err := game.Update(); err != nil {
		return err
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	game.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("count: %d", count))
	return nil
}

func main() {
	var err error
	game, err = ebitest.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, screenScale, "ebitest"); err != nil {
		log.Fatal(err)
	}
}
