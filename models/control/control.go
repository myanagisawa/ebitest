package control

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/utils"
)

// Base ...
type Base struct {
	label string

	image    *ebiten.Image
	layer    interfaces.Layer
	position *ebitest.Point
	scale    *ebitest.Scale
	angle    int
	moving   *ebitest.Point

	hasHoverAction bool
	hover          bool
}

// Label ...
func (o *Base) Label() string {
	return fmt.Sprintf("%s.%s", o.layer.Label(), o.label)
}

// In ...
func (o *Base) In(x, y int) bool {
	// パーツ位置（左上座標）
	minX, minY := o.Position(enum.TypeGlobal).GetInt()
	// パーツサイズ(オリジナル)
	size := ebitest.NewSize(o.image.Size())
	// スケール
	scale := o.Scale(enum.TypeGlobal)

	// 見かけ上の右下座標を取得
	maxX := int(float64(size.W())*scale.X()) + minX
	maxY := int(float64(size.H())*scale.Y()) + minY

	// フレーム領域
	fPosX, fPosY := o.layer.Frame().Position(enum.TypeGlobal).GetInt()
	fSize := o.layer.Frame().Size()
	fMaxX, fMaxY := fPosX+fSize.W(), fPosY+fSize.H()
	// 座標がフレーム外の場合はフレームのmax座標で置き換え
	if maxX > fMaxX {
		maxX = fMaxX
	}
	if maxY > fMaxY {
		maxY = fMaxY
	}

	// 座標がフレーム外の場合はフレームのmin座標で置き換え
	if minX < fPosX {
		minX = fPosX
	}
	if minY < fPosY {
		minY = fPosY
	}
	// log.Printf("レイヤ座標: {(%d, %d), (%d, %d)}", minX, minY, maxX, maxY)
	return (x >= minX && x <= maxX) && (y > minY && y <= maxY)
}

// SetLayer ...
func (o *Base) SetLayer(l interfaces.Layer) {
	o.layer = l
}

// Position ...
func (o *Base) Position(t enum.ValueTypeEnum) *ebitest.Point {
	dx, dy := 0.0, 0.0
	if o.moving != nil {
		dx, dy = o.moving.Get()
	}
	if t == enum.TypeLocal {
		return ebitest.NewPoint(o.position.X()+dx, o.position.Y()+dy)
	}
	gx, gy := 0.0, 0.0
	if o.layer != nil {
		gx, gy = o.layer.Position(enum.TypeGlobal).Get()
	}
	sx, sy := o.Scale(enum.TypeGlobal).Get()
	gx += (o.position.X() + dx) * sx
	gy += (o.position.Y() + dy) * sy
	return ebitest.NewPoint(gx, gy)
}

// SetPosition ...
func (o *Base) SetPosition(x, y float64) {
	o.position = ebitest.NewPoint(x, y)
}

// Scale ...
func (o *Base) Scale(t enum.ValueTypeEnum) *ebitest.Scale {
	return o.scale
}

// SetScale ...
func (o *Base) SetScale(x, y float64) {
	o.scale = ebitest.NewScale(x, y)
}

// Angle ...
func (o *Base) Angle(t enum.ValueTypeEnum) int {
	return o.angle
}

// SetAngle ...
func (o *Base) SetAngle(a int) {
	o.angle = a
}

// Theta ...
func (o *Base) Theta() float64 {
	return 2 * math.Pi * float64(o.Angle(enum.TypeGlobal)) / 360.0
}

// SetMoving ...
func (o *Base) SetMoving(dx, dy float64) {
	if o.moving == nil {
		o.moving = ebitest.NewPoint(dx, dy)
	} else {
		o.moving.Set(dx, dy)
	}
}

// Moving ...
func (o *Base) Moving() *ebitest.Point {
	return o.moving
}

// Update ...
func (o *Base) Update() error {
	if o.hasHoverAction {
		o.hover = o.In(ebiten.CursorPosition())
		if o.hover {
			log.Printf("hover: %s", o.label)
		}
	}
	return nil
}

// Draw draws the sprite.
func (o *Base) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())

	bgSize := ebitest.NewSize(o.image.Size())
	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(bgSize.W())/2, -float64(bgSize.H())/2)
	// 中心を軸に回転
	op.GeoM.Rotate(o.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(bgSize.W())/2, float64(bgSize.H())/2)

	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	// // フレームからのはみ出し判定
	// cx, cy := o.Position(enum.TypeLocal).GetInt()
	// lx, ly := o.layer.Position(enum.TypeLocal).GetInt()
	// cx += lx
	// cy += cy
	// lw, lh := o.image.Size()
	// fs := o.layer.Frame().Size()

	// x0, y0, x1, y1 := 0, 0, lw, lh
	// // frame外れ判定
	// if cx < 0 {
	// 	// 左にはみ出し
	// 	op.GeoM.Translate(float64(-cx), 0)
	// 	x0 = -cx
	// 	x1 += x0
	// }
	// if cy < 0 {
	// 	// 上にはみ出し
	// 	op.GeoM.Translate(0, float64(-cy))
	// 	y0 = -cy
	// 	y1 += y0
	// }
	// if cx+lw > fs.W() {
	// 	// 右にはみ出し
	// 	x1 = x0 + fs.W()
	// }

	// if cy+lh > fs.H() {
	// 	// 下にはみ出し
	// 	y1 = y0 + fs.H()
	// }

	// fr := image.Rect(x0, y0, x1, y1)
	// log.Printf("%s: pos: %d, %d, fr: %d, %d, %d, %d", o.label, lx, ly, x0, y0, x1, y1)
	// screen.DrawImage(o.image.SubImage(fr).(*ebiten.Image), op)
	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if o.hover {
		r, g, b, a = 0.5, 0.5, 0.5, 1.0
	}
	op.ColorM.Scale(r, g, b, a)
	screen.DrawImage(o.image, op)
}

// NewControlBase ...
func NewControlBase(label string, img image.Image, pos *ebitest.Point, labelColor color.Color) interfaces.UIControl {
	ti := utils.SetTextToCenter(label, img, ebitest.Fonts["btnFont"], labelColor)
	eimg := ebiten.NewImageFromImage(*ti)

	o := &Base{
		label:          label,
		image:          eimg,
		position:       pos,
		scale:          ebitest.NewScale(1.0, 1.0),
		hasHoverAction: true,
	}
	return o
}
