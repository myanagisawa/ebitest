package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime/pprof"

	"fmt"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/kitchen"
)

const (
	orgImgDir = "resources/images"
	imgDir    = "resources/resized_images"
	objDir    = "resources/object_images"
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
	list := []float64{0, 0.1, 0.2, 0.5, 1.0, 2.0, 3.0, 5.0}
	for _, v := range list {
		log.Printf("tan(%f)=%f", v, math.Tan(v))
	}
	for _, v := range list {
		n := math.Atan(v)
		d := n * 180 / math.Pi
		log.Printf("atan(%f)=%f, deg=%f", v, n, d)
	}
	y, x := -1.0, 1.0
	n := math.Atan2(y, x)
	d := n * 180 / math.Pi
	log.Printf("atan2(%f, %f)=%f, deg=%f", y, x, n, d)

	// 新着イメージの変換
	// imgconv.CreateNewImages(orgImgDir, imgDir)

	ebiten.SetRunnableInBackground(true)

	ebiten.SetWindowDecorated(false)

	game, _ := kitchen.NewGame(1344, 1008)
	if err := ebiten.Run(game.Update, game.WindowSize.Width, game.WindowSize.Height, 0.25, "kitchen sink"); err != nil {
		log.Fatal(err)
	}

	// game, _ := ebitest.NewGame(getResourceNames(), getObjectNames())
	// update := game.Update
	// if err := ebiten.Run(update, ebitest.ScreenWidth, ebitest.ScreenHeight, 1, "ebitest"); err != nil {
	// 	log.Fatal(err)
	// }
}

func getFileNames(dir string) []string {
	files, _ := ioutil.ReadDir(dir)
	fnames := make([]string, len(files))
	for i, f := range files {
		fnames[i] = f.Name()
	}
	return fnames
}

func getResourceNames() []string {
	fnames := getFileNames(imgDir)
	paths := make([]string, len(fnames))
	for i, f := range fnames {
		paths[i] = fmt.Sprintf("%s/%s", imgDir, f)
	}
	return paths
}

func getObjectNames() []string {
	fnames := getFileNames(objDir)
	paths := make([]string, len(fnames))
	for i, f := range fnames {
		paths[i] = fmt.Sprintf("%s/%s", objDir, f)
	}
	return paths
}
