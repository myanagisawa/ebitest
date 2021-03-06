package _ebitest

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/myanagisawa/ebitest/utils"
)

// type (

// 	// LabelFace ...
// 	LabelFace struct {
// 		uiFont        font.Face
// 		uiFontMHeight int
// 	}
// )

var (
	// DebugText ...
	DebugText string

	// Width ...
	Width int
	// Height ...
	Height int

	// Fonts ...
	// Fonts map[string]font.Face

	// Images ...
	Images map[string]image.Image

	// // SpotLightImage ...
	// SpotLightImage *ebiten.Image

	// fontFile
	// fontFile = "resources/fonts/GenShinGothic-Regular.ttf"

	// EdgeSize 画面の端から何ピクセルを端とするか
	EdgeSize = 30
	// EdgeSizeOuter Window外の何ピクセルまでを端判定に含めるか
	EdgeSizeOuter = 100
)

func init() {
	// Fonts = make(map[string]font.Face)
	// Fonts["btnFont"] = FontLoad(16)

	Images = make(map[string]image.Image)
	Images["btnBase"], _ = utils.GetImageByPath("resources/system_images/button.png")
	Images["world"], _ = utils.GetImageByPath("resources/system_images/world.jpg")

	Images["site_1"], _ = utils.GetImageByPath("resources/system_images/site_1.png")
	Images["site_2"], _ = utils.GetImageByPath("resources/system_images/site_2.png")
	Images["site_3"], _ = utils.GetImageByPath("resources/system_images/site_3.png")
	Images["site_4"], _ = utils.GetImageByPath("resources/system_images/site_4.png")
	Images["site_5"], _ = utils.GetImageByPath("resources/system_images/site_5.png")

	Images["route_1"], _ = utils.GetImageByPath("resources/system_images/route_1.png")
	Images["route_2"], _ = utils.GetImageByPath("resources/system_images/route_2.png")
	Images["route_3"], _ = utils.GetImageByPath("resources/system_images/route_3.png")

	img := utils.CreateRectImage(1, 1, &color.RGBA{200, 200, 200, 255})
	Images["routeLine"] = img.(draw.Image)

}

// CreateRectImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
// func CreateRectImage(w, h int, color *color.RGBA) image.Image {
// 	m := image.NewRGBA(image.Rect(0, 0, w, h))
// 	// 横ループ、半径*2＝直径まで
// 	for x := 0; x < w; x++ {
// 		// 縦ループ、半径*2＝直径まで
// 		for y := 0; y < h; y++ {
// 			m.Set(x, y, color)
// 		}
// 	}
// 	return m
// }

// // CreateBorderedRectImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
// func CreateBorderedRectImage(w, h int, c *color.RGBA) image.Image {
// 	m := image.NewRGBA(image.Rect(0, 0, w, h))
// 	// 横ループ、半径*2＝直径まで
// 	for x := 0; x < w; x++ {
// 		// 縦ループ、半径*2＝直径まで
// 		for y := 0; y < h; y++ {
// 			if (x == 0 || x == (w-1)) || (y == 0 || y == (h-1)) {
// 				m.Set(x, y, color.RGBA{c.R / 2, c.G / 2, c.B / 2, c.A})
// 			} else {
// 				m.Set(x, y, c)
// 			}
// 		}
// 	}
// 	return m
// }

// // FontLoad ...
// func FontLoad(size int) font.Face {
// 	// ebitenフォントのテスト
// 	ftBinary, err := ioutil.ReadFile(fontFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	tt, err := truetype.Parse(ftBinary)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return truetype.NewFace(tt, &truetype.Options{
// 		Size:    float64(size),
// 		DPI:     72,
// 		Hinting: font.HintingFull,
// 	})
// }

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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Size ...
type Size struct {
	width  int
	height int
}

// NewSize ...
func NewSize(width, height int) *Size {
	return &Size{width, height}
}

// W ...
func (s *Size) W() int {
	return s.width
}

// H ...
func (s *Size) H() int {
	return s.height
}

// Get ...
func (s *Size) Get() (int, int) {
	return s.width, s.height
}

// Set ...
func (s *Size) Set(width, height int) {
	s.width = width
	s.height = height
}
