package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/models/game"
)

const (
	screenWidth  = 1200
	screenHeight = 900
)

func main() {

	ebiten.SetScreenTransparent(true)
	// ebiten.SetWindowDecorated(false)
	// ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowPosition(100, 400)

	manager := game.NewManager(screenWidth, screenHeight)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
