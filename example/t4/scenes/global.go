package scenes

import (
	"image"
	"image/color"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type (

	// LabelFace ...
	LabelFace struct {
		uiFont        font.Face
		uiFontMHeight int
	}
)

var (
	width, height int

	fontFile = "resources/fonts/GenShinGothic-Light.ttf"
)

// createCircleImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
func createRectImage(w, h int, color color.RGBA) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	// 横ループ、半径*2＝直径まで
	for x := 0; x < w; x++ {
		// 縦ループ、半径*2＝直径まで
		for y := 0; y < h; y++ {
			m.Set(x, y, color)
		}
	}
	return m
}

// FontLoad ...
func FontLoad(size int) font.Face {
	// ebitenフォントのテスト
	ftBinary, err := ioutil.ReadFile(fontFile)
	if err != nil {
		panic(err)
	}
	tt, err := truetype.Parse(ftBinary)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(tt, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
