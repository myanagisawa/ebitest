package ebitest

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/font"
)

type (

	// LabelFace ...
	LabelFace struct {
		uiFont        font.Face
		uiFontMHeight int
	}

	// Size ...
	Size struct {
		width  int
		height int
	}
)

var (

	// Width ...
	Width int
	// Height ...
	Height int

	// Fonts ...
	Fonts map[string]font.Face
	// Images ...
	Images map[string]draw.Image

	// fontFile
	fontFile = "resources/fonts/GenShinGothic-Light.ttf"

	// EdgeSize 画面の端から何ピクセルを端とするか
	EdgeSize = 30
	// EdgeSizeOuter Window外の何ピクセルまでを端判定に含めるか
	EdgeSizeOuter = 100
)

func init() {
	Fonts = make(map[string]font.Face)
	Fonts["btnFont"] = FontLoad(16)

	Images = make(map[string]draw.Image)
	Images["bgImage"], _ = utils.GetImageByPath("resources/system_images/bg-2.jpg")
	Images["btnBase"], _ = utils.GetImageByPath("resources/system_images/button.png")
	Images["btnBaseHover"], _ = utils.GetImageByPath("resources/system_images/button-hover.png")
	Images["bgFlower"], _ = utils.GetImageByPath("resources/system_images/bg_flower.jpg")

	img := CreateRectImage(10, 10, color.RGBA{0, 0, 0, 255})
	Images["listBase"] = img.(draw.Image)

	img = CreateRectImage(10, 10, color.RGBA{128, 128, 128, 128})
	Images["listScroller"] = img.(draw.Image)

}

// CreateRectImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
func CreateRectImage(w, h int, color color.RGBA) image.Image {
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

// Point ...
type Point struct {
	x float64
	y float64
}

// NewPoint ...
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// X ...
func (p *Point) X() float64 {
	return p.x
}

// Y ...
func (p *Point) Y() float64 {
	return p.y
}

// Get ...
func (p *Point) Get() (float64, float64) {
	return p.x, p.y
}

// GetInt ...
func (p *Point) GetInt() (int, int) {
	return int(p.x), int(p.y)
}

// Set ...
func (p *Point) Set(x, y float64) {
	p.x = x
	p.y = y
}

// SetDelta ...
func (p *Point) SetDelta(dx, dy float64) {
	p.x += dx
	p.y += dy
}

// Scale op.scale設定値
type Scale struct {
	x float64
	y float64
}

// NewScale ...
func NewScale(x, y float64) *Scale {
	return &Scale{x, y}
}

// X ...
func (s *Scale) X() float64 {
	return s.x
}

// Y ...
func (s *Scale) Y() float64 {
	return s.y
}

// Get ...
func (s *Scale) Get() (float64, float64) {
	return s.x, s.y
}

// Set ...
func (s *Scale) Set(x, y float64) {
	s.x = x
	s.y = y
}
