package models

import (
	"math"

	"github.com/hajimehoshi/ebiten"
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
}

// NewEbiObject ...
func NewEbiObject(tag string, image *ebiten.Image, parent *EbiObject, scale *ebitest.Scale, position *ebitest.Point, angle int, inheritScale bool, inheritAngle bool) *EbiObject {
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
	}
	return eo
}

// EbitenImage ...
func (o *EbiObject) EbitenImage() *ebiten.Image {
	return o.img
}

// Size ...
func (o *EbiObject) Size() (int, int) {
	return o.img.Size()
}

// Transition ...
func (o *EbiObject) Transition() (float64, float64) {
	return o.position.Get()
}

// Scale ...
func (o *EbiObject) Scale() (float64, float64) {
	return o.scale.Get()
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

// GlobalTransition ...
func (o *EbiObject) GlobalTransition() (float64, float64) {
	gx, gy := 0.0, 0.0
	if o.parent != nil {
		gx, gy = o.parent.GlobalTransition()
	}
	gx += o.position.X()
	gy += o.position.Y()
	return gx, gy
}

// GlobalScale ...
func (o *EbiObject) GlobalScale() (float64, float64) {
	if o.parent != nil && o.inheritScale {
		return o.parent.GlobalScale()
	}
	return o.scale.Get()
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

// // SetT ...
// func (o *EbiObject) SetT() {
// 	gx, gy := 0.0, 0.0
// 	if o.parent != nil {
// 		gx, gy = o.parent.Transition()
// 	}
// 	o.img.t.x = gx + o.position.x
// 	o.img.t.y = gy + o.position.y
// 	o.img.isDefault = false
// }

// // SetS ...
// func (o *EbiObject) SetS() {
// 	if o.parent != nil && o.inheritScale {
// 		o.img.s.x, o.img.s.y = o.parent.Scale()
// 	} else {
// 		o.img.s.x = 1.0
// 		o.img.s.y = 1.0
// 	}
// 	o.img.isDefault = false
// }

// // SetA ...
// func (o *EbiObject) SetA() {
// 	if o.parent != nil && o.inheritAngle {
// 		o.img.a = o.parent.img.a
// 	} else {
// 		o.img.a = 0
// 	}
// 	o.img.isDefault = false
// }
