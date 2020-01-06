package ebitest

import (
	"fmt"
	"sync"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	// ScreenWidth ...
	ScreenWidth = 600
	// ScreenHeight ...
	ScreenHeight = 800
)

type (
	// Game ...
	Game struct {
		input        *Input
		sceneManager *SceneManager
	}
)

// NewGame ...
func NewGame(paths, objPaths []string) (*Game, error) {
	fmt.Print("loading images")
	images := make([]ebiten.Image, len(objPaths))

	wg := &sync.WaitGroup{}
	for i, path := range objPaths {
		wg.Add(1)
		go func(list []ebiten.Image, idx int, path string) {
			img, _, err := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
			if err != nil {
				panic(err)
			}
			list[idx] = *img
			fmt.Print(".")
			wg.Done()
		}(images, i, path)
	}
	wg.Wait()
	fmt.Println("complete!")

	g := &Game{
		input: NewInput(),
		sceneManager: &SceneManager{
			paths:        paths,
			objectImages: images,
		},
	}
	g.sceneManager.GoTo(NewCommonScene(g.sceneManager.PathToImage(0)))
	return g, nil
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
	g.input.Update()
	if err := g.sceneManager.Update(g.input); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.sceneManager.Draw(r)

	return nil
}
