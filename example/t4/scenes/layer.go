package scenes

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Layer ...
type Layer interface {
	Label() string
	In(x, y int) bool
	ControlAt(x, y int) interface{}
	ScaleTo(scale float64)
	LocalPosition(x, y int) (int, int)
	IsModal() bool
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
}

// LayerBase ...
type LayerBase struct {
	label      string
	bg         *ebiten.Image
	x          int
	y          int
	scale      float64
	translateX float64
	translateY float64
	parent     *BattleScene
	isModal    bool
}

// Label ...
func (l *LayerBase) Label() string {
	return l.label
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (l *LayerBase) In(x, y int) bool {
	// レイヤ位置（左上座標）
	tx, ty := int(l.translateX), int(l.translateY)

	// レイヤサイズ(オリジナル)
	w, h := l.bg.Size()

	// 見かけ上の右下座標を取得
	maxX := int(float64(w)*l.scale) + tx
	maxY := int(float64(h)*l.scale) + ty
	if maxX > width {
		maxX = width
	}
	if maxY > height {
		maxY = height
	}

	// 見かけ上の左上座標を取得
	minX, minY := tx, ty
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	log.Printf("レイヤ座標: {(%d, %d), (%d, %d)}", minX, minY, maxX, maxY)
	return (x >= minX && x <= maxX) && (y > minY && y <= maxY)
	// return l.bg.At(x-l.x, y-l.y).(color.RGBA).A > 0
}

// ControlAt (x, y)座標に存在する部品を返します
func (l *LayerBase) ControlAt(x, y int) interface{} {
	return l.bg
}

// ScaleTo ...
func (l *LayerBase) ScaleTo(scale float64) {
	l.scale = scale

	w, h := l.bg.Size()
	w = int(float64(w) * l.scale)
	h = int(float64(h) * l.scale)

	if l.x < width-w {
		l.x = width - w
	}
	if l.y < height-h {
		l.y = height - h
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)
	// log.Printf("MoveBy: s.x=%0.2f, s.y=%0.2f", s.translateX, s.translateY)
}

// LocalPosition スクリーン上の座標をシーンオブジェクト上の座標に変換します
func (l *LayerBase) LocalPosition(x, y int) (int, int) {
	sx := float64(l.x) * l.scale * -1
	sy := float64(l.y) * l.scale * -1
	localX := int((float64(x) + sx) / l.scale)
	localY := int((float64(y) + sy) / l.scale)

	// log.Printf("scale: %0.2f [x: %d, s.x: %d = %d] [y: %d, s.y: %d = %d]", s.scale, x, s.x, localX, y, s.y, localY)
	return localX, localY
}

// IsModal ...
func (l *LayerBase) IsModal() bool {
	return l.isModal
}

// Update ...
func (l *LayerBase) Update(screen *ebiten.Image) error {

	return nil
}

// Draw ...
func (l *LayerBase) Draw(screen *ebiten.Image) {
	// log.Printf("LayerBase.Draw")
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.translateX, l.translateY)

	screen.DrawImage(l.bg, op)
}

// NewLayerBase ...
func NewLayerBase() *LayerBase {
	img := createRectImage(width, height, color.RGBA{32, 32, 32, 255})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &LayerBase{
		label:   "layer base",
		bg:      eimg,
		x:       0,
		y:       0,
		scale:   1.0,
		isModal: false,
	}

	l.translateX = float64(l.x)
	l.translateY = float64(l.y)
	return l
}

// TestWindow ...
type TestWindow struct {
	LayerBase
}

// NewTestWindow ...
func NewTestWindow() *TestWindow {
	img := createRectImage(200, 400, color.RGBA{0, 0, 0, 64})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &TestWindow{
		LayerBase: LayerBase{
			label:   "test window",
			bg:      eimg,
			x:       50,
			y:       100,
			scale:   1.0,
			isModal: false,
		},
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)
	return l
}

// Draw ...
func (l *TestWindow) Draw(screen *ebiten.Image) {
	// log.Printf("TestWindow.Draw")

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.translateX, l.translateY)

	screen.DrawImage(l.bg, op)
}
