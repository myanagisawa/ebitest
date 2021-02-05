package main

import (
	"log"

	"github.com/myanagisawa/ebitest/app"
	"github.com/pkg/profile"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1200
	screenHeight = 1200
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	ebiten.SetScreenTransparent(true)
	// ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowPosition(3360-(screenWidth+100), 100)
	ebiten.SetMaxTPS(15)

	manager := app.NewGameManager(screenWidth, screenHeight)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
