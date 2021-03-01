package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t7/app"
	"github.com/pkg/profile"
)

const (
	screenWidth  = 1200
	screenHeight = 1200
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowPosition(3360-(screenWidth+100), 100)
	ebiten.SetMaxTPS(60)

	manager := app.NewGameManager(screenWidth, screenHeight)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
