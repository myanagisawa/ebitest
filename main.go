package main

import (
	"log"

	"github.com/pkg/profile"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/game"
)

const (
	screenWidth  = 800
	screenHeight = 1200
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	ebiten.SetScreenTransparent(true)
	// ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowPosition(100, 100)
	// ebiten.SetMaxTPS(30)

	manager := game.NewManager(screenWidth, screenHeight)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
