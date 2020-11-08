package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/models/game"
)

const (
	screenWidth  = 1600
	screenHeight = 1200
)

func main() {

	ebiten.SetScreenTransparent(true)
	// ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowPosition(100, 400)
	// ebiten.SetMaxTPS(30)

	manager := game.NewManager(screenWidth, screenHeight)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
