package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/example/t4/scenes"
)

const (
	screenWidth  = 1200
	screenHeight = 900
)

func main() {

	ebiten.SetScreenTransparent(true)
	// ebiten.SetWindowDecorated(false)
	// ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	game := scenes.NewGameManager(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
