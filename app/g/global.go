package g

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/myanagisawa/ebitest/utils"
)

var (
	// DebugText ...
	DebugText string

	// Width ...
	Width int
	// Height ...
	Height int

	// Master マスタデータ
	Master *MasterData

	// Images ...
	Images map[string]image.Image

	// EdgeSize 画面の端から何ピクセルを端とするか
	EdgeSize = 30
	// EdgeSizeOuter Window外の何ピクセルまでを端判定に含めるか
	EdgeSizeOuter = 100
)

func init() {
	// masterロード
	{
		log.Printf("MASTERデータ読み込み...")
		Master = LoadMaster()
		log.Printf("MASTERデータ読み込み...done")
	}

	// 画像リソース読み込み
	{
		log.Printf("Resource読み込み...")
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
		log.Printf("Resource読み込み...done")
	}

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

// DefPoint ...
func DefPoint() *Point {
	return &Point{0.0, 0.0}
}

// X ...
func (p *Point) X() float64 {
	return p.x
}

// Y ...
func (p *Point) Y() float64 {
	return p.y
}

// IntX ...
func (p *Point) IntX() int {
	return int(p.x)
}

// IntY ...
func (p *Point) IntY() int {
	return int(p.y)
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

// DefScale ...
func DefScale() *Scale {
	return &Scale{1.0, 1.0}
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

// Bound ...
type Bound struct {
	Min *Point
	Max *Point
}

// NewBound ...
func NewBound(min, max *Point) *Bound {
	return &Bound{Min: min, Max: max}
}

// NewBoundByPosSize ...
func NewBoundByPosSize(pos *Point, size *Size) *Bound {
	return &Bound{Min: pos, Max: NewPoint(pos.X()+float64(size.W()), pos.Y()+float64(size.H()))}
}
