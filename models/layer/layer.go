package layer

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/event"
	"github.com/myanagisawa/ebitest/models/input"
)

// Base ...
type Base struct {
	label    string
	image    *ebiten.Image
	parent   interfaces.Frame
	position *ebitest.Point
	scale    *ebitest.Scale
	moving   *ebitest.Point

	modal        bool
	draggable    bool
	controls     []interfaces.UIControl
	eventHandler interfaces.EventHandler
	stroke       *input.Stroke
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

// Parent ...
func (o *Base) Parent() interfaces.Frame {
	return o.parent
}

// SetParent ...
func (o *Base) SetParent(parent interfaces.Frame) {
	o.parent = parent
}

// Position ...
func (o *Base) Position(t enum.ValueTypeEnum) *ebitest.Point {
	if t == enum.TypeLocal {
		return o.position
	}
	gx, gy := 0.0, 0.0
	if o.parent != nil {
		gx, gy = o.parent.Position(enum.TypeGlobal).Get()
	}
	dx, dy := 0.0, 0.0
	if o.moving != nil {
		dx, dy = o.moving.Get()
		log.Printf("%s, moving: %0.1f,  %0.1f", o.label, dx, dy)
	}
	sx, sy := o.Scale(enum.TypeGlobal).Get()
	gx += (o.position.X() + dx) * sx
	gy += (o.position.Y() + dy) * sy
	return ebitest.NewPoint(gx, gy)
}

// Scale ...
func (o *Base) Scale(t enum.ValueTypeEnum) *ebitest.Scale {
	return o.scale
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

// AddUIControl レイヤに部品を追加します
func (o *Base) AddUIControl(c interfaces.UIControl) {
	c.SetLayer(o)
	o.controls = append(o.controls, c)
}

// UIControlAt (x, y)座標に存在する部品を返します
func (o *Base) UIControlAt(x, y int) interfaces.UIControl {
	for i := len(o.controls) - 1; i >= 0; i-- {
		c := o.controls[i]
		if c.In(x, y) {
			return c
		}
	}
	return nil
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (o *Base) In(x, y int) bool {
	px, py := o.Position(enum.TypeGlobal).GetInt()
	return o.image.At(x-px, y-py).(color.RGBA).A > 0
}

// IsModal ...
func (o *Base) IsModal() bool {
	return o.modal
}

// Scroll ...
func (o *Base) Scroll(et enum.EdgeTypeEnum) {
	// 1フレームあたりの増分値
	dp := 20.0
	switch et {
	case enum.EdgeTypeTopLeft:
		o.position.SetDelta(dp, dp)
	case enum.EdgeTypeTop:
		o.position.SetDelta(0, dp)
	case enum.EdgeTypeTopRight:
		o.position.SetDelta(-dp, dp)
	case enum.EdgeTypeRight:
		o.position.SetDelta(-dp, 0)
	case enum.EdgeTypeBottomRight:
		o.position.SetDelta(-dp, -dp)
	case enum.EdgeTypeBottom:
		o.position.SetDelta(0, -dp)
	case enum.EdgeTypeBottomLeft:
		o.position.SetDelta(dp, -dp)
	case enum.EdgeTypeLeft:
		o.position.SetDelta(dp, 0)
	}

	posX, posY := o.Position(enum.TypeLocal).GetInt()
	layerSize := ebitest.NewSize(o.image.Size())
	frameSize := o.parent.Size()

	if posX > 0 {
		o.position.SetDelta(-dp, 0)
	} else if posX+layerSize.W() < frameSize.W() {
		o.position.SetDelta(dp, 0)
	}

	if posY > 0 {
		o.position.SetDelta(0, -dp)
	} else if posY+layerSize.H() < frameSize.H() {
		o.position.SetDelta(0, dp)
	}
}

// Update ...
func (o *Base) Update() error {

	if o.stroke != nil {
		o.updateStroke(o.stroke)
		if o.stroke.IsReleased() {
			o.updatePositionByDelta()
			o.stroke = nil
			log.Printf("drag end")
		}
	}

	if o.parent.ActiveLayer() != nil && o.parent.ActiveLayer() == o {
		// log.Printf("active layer: %s", l.parent.ActiveLayer().Label())

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// log.Printf("left button push: x=%d, y=%d", x, y)
			if o.In(x, y) {
				stroke := input.NewStroke(&input.MouseStrokeSource{})
				// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
				// o := l.UIControlAt(x, y)
				// if o != nil || l.bg.IsDraggable() {
				// 	l.stroke = stroke
				// 	log.Printf("drag start")
				// }
				if o.draggable {
					o.stroke = stroke
					log.Printf("%s drag start", o.label)
				}
			}
		}

		// for _, c := range o.controls {
		// 	_ = c.Update()
		// }
	}

	return nil
}

// Draw ...
func (o *Base) Draw(screen *ebiten.Image) {
	// log.Printf("LayerBase.Draw")
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())

	op.GeoM.Translate(o.Position(enum.TypeLocal).Get())

	screen.DrawImage(o.image, op)

	// for _, c := range o.controls {
	// 	c.Draw(screen)
	// }

}

// NewLayerBase ...
func NewLayerBase(label string, img image.Image, pos *ebitest.Point, scale *ebitest.Scale, draggable bool) interfaces.Layer {
	eimg := ebiten.NewImageFromImage(img)

	// draggable, ismodal 未実装

	l := &Base{
		label:        label,
		image:        eimg,
		position:     pos,
		scale:        scale,
		draggable:    draggable,
		eventHandler: event.NewEventHandler(),
	}
	return l
}

// updatePositionByDelta ...
func (o *Base) updatePositionByDelta() {

	o.position.SetDelta(o.moving.Get())
	o.moving = nil
}

func (o *Base) updateStroke(stroke *input.Stroke) {
	stroke.Update()
	o.SetMoving(stroke.PositionDiff())
}
