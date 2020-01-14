package ebitest

import (
	"fmt"
	"image"
	"os"
	"sync"

	"log"

	"golang.org/x/image/draw"

	"github.com/hajimehoshi/ebiten"
)

const (
	// ScreenWidth ...
	ScreenWidth = 300
	// ScreenHeight ...
	ScreenHeight = 400
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
	// mask画像読み込み
	mask, _ := getImageByPath("system_images/mask.png")

	wg := &sync.WaitGroup{}
	for i, path := range objPaths {
		wg.Add(1)
		go func(list []ebiten.Image, idx int, path string, maskimg image.Image) {
			// 画像読み込み
			img, _ := getImageByPath(path)
			// 画像サイズに合わせたマスクの作成
			rmask := image.NewRGBA(img.Bounds())
			draw.BiLinear.Scale(rmask, rmask.Bounds(), maskimg, maskimg.Bounds(), draw.Over, nil)
			// 円形maskの適用
			// log.Printf("img.bounds: %#v", img.Bounds())
			// log.Printf("mask.bounds: %#v", rmask.Bounds())
			out := image.NewRGBA(img.Bounds())
			draw.DrawMask(out, out.Bounds(), img, image.Point{0, 0}, rmask, image.Point{0, 0}, draw.Over)

			// 画像からebiten.imageを作成
			// img, _, err := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
			// if err != nil {
			// 	panic(err)
			// }
			ebimg, err := ebiten.NewImageFromImage(out, ebiten.FilterDefault)
			if err != nil {
				panic(err)
			}
			list[idx] = *ebimg
			fmt.Print(".")
			wg.Done()
		}(images, i, path, mask)
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

func getImageByPath(path string) (image.Image, string) {
	//画像読み込み
	fileIn, err := os.Open(path)
	defer fileIn.Close()
	if err != nil {
		fmt.Println("error:file\n", err)
		log.Panic(err.Error())
	}

	//画像をimage型として読み込む
	img, format, err := image.Decode(fileIn)
	if err != nil {
		fmt.Println("error:decode\n", format, err)
		log.Panic(err.Error())
	}
	return img, format
}
