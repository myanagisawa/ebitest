package ebitest

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"sync"

	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/utils"
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
	mask, _ := getImageByPath("resources/system_images/mask.png")
	wg := &sync.WaitGroup{}
	for i, path := range objPaths {
		wg.Add(1)
		go func(list []ebiten.Image, idx int, path string, maskimg image.Image) {
			// 画像読み込み
			img, _ := getImageByPath(path)

			// 画像をマスク
			out := utils.MaskImage(img.(*image.NRGBA), maskimg)

			// 文字列を表示
			setfont(out, "てすと", 48.0)

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

func fontload(fname string) []byte {
	file, err := os.Open(fname)
	defer file.Close()
	if err != nil {
		fmt.Println("error:file\n", err)
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error:fileread\n", err)
		return nil
	}

	return bytes
}

func setfont(out draw.Image, str string, fontsize float64) {
	ft, err := truetype.Parse(fontload("/Library/Fonts/Arial Unicode.ttf"))
	if err != nil {
		fmt.Println("font", err)
		return
	}
	opt := truetype.Options{Size: fontsize}
	face := truetype.NewFace(ft, &opt)

	d := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	// 文字を表示対象の真ん中に表示する
	size := out.Bounds().Size()
	d.Dot.X = (fixed.I(size.X) - d.MeasureString(str)) / 2
	d.Dot.Y = fixed.I((size.Y / 2) + int(fontsize/2))

	d.DrawString(str)
}
