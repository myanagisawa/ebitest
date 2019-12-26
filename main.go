package main

import (
	"flag"
	_ "image/jpeg"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/ebitest"
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

	game, _ := ebitest.NewGame()
	update := game.Update
	if err := ebiten.Run(update, ebitest.ScreenWidth, ebitest.ScreenHeight, 2, "ebitest"); err != nil {
		log.Fatal(err)
	}
}
