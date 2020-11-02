package models

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
)

// EbiObject 画像リソース配置情報を含む構造体
// tag デバッグ用の文字列
// img 画像オブジェクト
// parent 親のリソース
// inheritScale 親のScale情報を引き継ぐかどうか
// inheritAngle 親のAngle情報を引き継ぐかどうか
type EbiObject struct {
	tag          string
	img          *ebiten.Image
	parent       *EbiObject
	scale        *ebitest.Scale
	position     *ebitest.Point
	angle        int
	inheritScale bool
	inheritAngle bool
	moving       *ebitest.Point
	draggable    bool
}

// NewEbiObject ...
func NewEbiObject(tag string, image *ebiten.Image, parent *EbiObject, scale *ebitest.Scale, position *ebitest.Point, angle int, inheritScale, inheritAngle, draggable bool) *EbiObject {
	if scale == nil {
		scale = ebitest.NewScale(1.0, 1.0)
	}
	if position == nil {
		position = ebitest.NewPoint(0.0, 0.0)
	}
	eo := &EbiObject{
		tag:          tag,
		img:          image,
		parent:       parent,
		scale:        scale,
		position:     position,
		angle:        angle,
		inheritScale: inheritScale,
		inheritAngle: inheritAngle,
		draggable:    draggable,
	}
	return eo
}

// EbitenImage ...
func (o *EbiObject) EbitenImage() *ebiten.Image {
	return o.img
}

// Size ...
func (o *EbiObject) Size() *ebitest.Size {
	x, y := o.img.Size()
	return ebitest.NewSize(int(float64(x)*o.scale.X()), int(float64(y)*o.scale.Y()))
}

// Position ...
func (o *EbiObject) Position() *ebitest.Point {
	return o.position
}

// Scale ...
func (o *EbiObject) Scale() *ebitest.Scale {
	return o.scale
}

// Angle ...
func (o *EbiObject) Angle() int {
	return o.angle
}

// Parent ...
func (o *EbiObject) Parent() *EbiObject {
	return o.parent
}

// SetParent ...
func (o *EbiObject) SetParent(parent *EbiObject) {
	o.parent = parent
}

// SetMoving ...
func (o *EbiObject) SetMoving(dx, dy float64) {
	if o.moving == nil {
		o.moving = ebitest.NewPoint(dx, dy)
	} else {
		o.moving.Set(dx, dy)
	}
}

// Moving ...
func (o *EbiObject) Moving() *ebitest.Point {
	return o.moving
}

// UpdatePositionByDelta ...
func (o *EbiObject) UpdatePositionByDelta() {

	o.position.SetDelta(o.moving.Get())
	o.moving = nil
}

// GlobalPosition ...
func (o *EbiObject) GlobalPosition() *ebitest.Point {
	gx, gy := 0.0, 0.0
	if o.parent != nil {
		gx, gy = o.parent.GlobalPosition().Get()
	}
	dx, dy := 0.0, 0.0
	if o.moving != nil {
		dx, dy = o.moving.Get()
		log.Printf("%s, moving: %0.1f,  %0.1f", o.tag, dx, dy)
	}
	sx, sy := o.GlobalScale().Get()
	gx += (o.position.X() + dx) * sx
	gy += (o.position.Y() + dy) * sy
	return ebitest.NewPoint(gx, gy)
}

// GlobalScale ...
func (o *EbiObject) GlobalScale() *ebitest.Scale {
	if o.parent != nil && o.inheritScale {
		return o.parent.GlobalScale()
	}
	return o.scale
}

// GlobalAngle ...
func (o *EbiObject) GlobalAngle() int {
	if o.parent != nil && o.inheritAngle {
		return o.parent.GlobalAngle() + o.angle
	}
	return o.angle
}

// Theta ...
func (o *EbiObject) Theta() float64 {
	return 2 * math.Pi * float64(o.GlobalAngle()) / 360.0
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (o *EbiObject) In(x, y int) bool {
	px, py := o.GlobalPosition().Get()
	return o.img.At(x-int(px), y-int(py)).(color.RGBA).A > 0
}

// SetDraggable ...
func (o *EbiObject) SetDraggable(b bool) {
	o.draggable = b
}

// IsDraggable ...
func (o *EbiObject) IsDraggable() bool {
	return o.draggable
}
