package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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

	// game := scenes.NewGameManager(screenWidth, screenHeight)

	// if err := ebiten.RunGame(game); err != nil {
	// 	log.Fatal(err)
	// }
}
