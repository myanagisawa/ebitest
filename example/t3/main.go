package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/example/t3/ex3"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}
	ebiten.SetRunnableInBackground(true)
	// ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)

	game, _ := ex3.NewGame(1280, 800)
	if err := ebiten.Run(game.Update, game.WindowSize.Width, game.WindowSize.Height, 1.0, "example.t3"); err != nil {
		log.Fatal(err)
	}
}
