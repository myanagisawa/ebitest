package layer

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/functions"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/event"
	"github.com/myanagisawa/ebitest/utils"
)

// Base ...
type Base struct {
	frame        interfaces.Frame
	label        string
	image        *ebiten.Image
	position     *g.Point
	scale        *g.Scale
	moving       *g.Point
	controls     []interfaces.UIControl
	eventHandler interfaces.EventHandler
	modal        bool
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

// Size ...
func (o *Base) Size() *g.Size {
	return g.NewSize(o.image.Size())
}

// Manager ...
func (o *Base) Manager() interfaces.GameManager {
	return o.frame.Manager()
}

// Frame ...
func (o *Base) Frame() interfaces.Frame {
	return o.frame
}

// SetFrame ...
func (o *Base) SetFrame(frame interfaces.Frame) {
	o.frame = frame
}

// SetPosition ...
func (o *Base) SetPosition(x, y float64) {
	o.position = g.NewPoint(x, y)
}

// Position ...
func (o *Base) Position(t enum.ValueTypeEnum) *g.Point {
	dx, dy := 0.0, 0.0
	if o.moving != nil {
		dx, dy = o.moving.Get()
	}
	if t == enum.TypeLocal {
		return g.NewPoint(o.position.X()+dx, o.position.Y()+dy)
	}
	gx, gy := 0.0, 0.0
	if o.frame != nil {
		gx, gy = o.frame.Position(enum.TypeGlobal).Get()
	}
	sx, sy := o.Scale(enum.TypeGlobal).Get()
	gx += (o.position.X() + dx) * sx
	gy += (o.position.Y() + dy) * sy
	return g.NewPoint(gx, gy)
}

// Scale ...
func (o *Base) Scale(t enum.ValueTypeEnum) *g.Scale {
	return o.scale
}

// SetMoving ...
func (o *Base) SetMoving(dx, dy float64) {
	if o.moving == nil {
		o.moving = g.NewPoint(dx, dy)
	} else {
		o.moving.Set(dx, dy)
	}
}

// Moving ...
func (o *Base) Moving() *g.Point {
	return o.moving
}

// AddUIControl レイヤに部品を追加します
func (o *Base) AddUIControl(c interfaces.UIControl) {
	c.SetLayer(o)
	o.controls = append(o.controls, c)
}

// RemoveUIControl レイヤに登録された部品を削除します
func (o *Base) RemoveUIControl(c interfaces.UIControl) {
	ret := make([]interfaces.UIControl, len(o.controls))
	i := 0
	for _, x := range o.controls {
		if c != x {
			ret[i] = x
			i++
		}
	}
	o.controls = ret[:i]
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
	layerSize := g.NewSize(o.image.Size())
	frameSize := o.frame.Size()

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

	// x, y := o.position.GetInt()
	// log.Printf("%s: %d, %d", o.Label(), x, y)
}

// GetObjects ...
func (o *Base) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	for i := len(o.controls) - 1; i >= 0; i-- {
		c := o.controls[i]
		objs = append(objs, c.GetObjects(x, y)...)
	}

	if o.In(x, y) {
		objs = append(objs, o)
	}
	// log.Printf("LayerBase::GetObjects: %#v", objs)
	return objs
}

// Update ...
func (o *Base) Update() error {

	if o.frame.ActiveLayer() != nil && o.frame.ActiveLayer() == o {
		// log.Printf("active layer: %s", l.frame.ActiveLayer().Label())

		for _, c := range o.controls {
			_ = c.Update()
		}
	}

	// frame外に出ないようにする制御
	o.updatePos()

	return nil
}

// Draw ...
func (o *Base) Draw(screen *ebiten.Image) {
	// log.Printf("LayerBase.Draw")
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())

	// op.GeoM.Translate(o.Position(enum.TypeLocal).Get())
	// screen.DrawImage(o.image, op)

	// log.Printf("%s: translate(%0.2f, %0.2f)", o.label, o.Position(enum.TypeGlobal).X(), o.Position(enum.TypeGlobal).Y())
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	lx, ly := o.Position(enum.TypeLocal).GetInt()
	lw, lh := o.image.Size()
	fs := o.frame.Size()

	x0, y0, x1, y1 := 0, 0, lw, lh
	// frame外れ判定
	if lx < 0 {
		// 左にはみ出し
		op.GeoM.Translate(float64(-lx), 0)
		x0 = -lx
		x1 += x0
	}
	if ly < 0 {
		// 上にはみ出し
		op.GeoM.Translate(0, float64(-ly))
		y0 = -ly
		y1 += y0
	}
	if lx+lw > fs.W() {
		// 右にはみ出し
		x1 = x0 + fs.W()
	}

	if ly+lh > fs.H() {
		// 下にはみ出し
		y1 = y0 + fs.H()
	}

	fr := image.Rect(x0, y0, x1, y1)
	// log.Printf("%s: pos: %d, %d, fr: %d, %d, %d, %d", o.label, lx, ly, x0, y0, x1, y1)
	screen.DrawImage(o.image.SubImage(fr).(*ebiten.Image), op)

	for _, c := range o.controls {
		c.Draw(screen)
	}
}

// EventHandler ...
func (o *Base) EventHandler() interfaces.EventHandler {
	return o.eventHandler
}

// NewLayerBase ...
func NewLayerBase(f interfaces.Frame, label string, pos *g.Point, size *g.Size, c *color.RGBA, draggable bool) interfaces.Layer {
	img := utils.CreateRectImage(size.W(), size.H(), c)

	return NewLayerBaseByImage(f, label, img, pos, draggable)
}

// NewLayerBaseByImage ...
func NewLayerBaseByImage(f interfaces.Frame, label string, img image.Image, pos *g.Point, draggable bool) interfaces.Layer {
	eimg := ebiten.NewImageFromImage(img)

	l := &Base{
		frame:        f,
		label:        label,
		image:        eimg,
		position:     pos,
		scale:        g.NewScale(1.0, 1.0),
		eventHandler: event.NewEventHandler(),
	}
	if draggable {
		l.eventHandler.AddEventListener(enum.EventTypeDragging, functions.CommonEventCallback)
		l.eventHandler.AddEventListener(enum.EventTypeDragDrop, functions.CommonEventCallback)
	}
	f.AddLayer(l)

	return l
}

// FinishStroke ...
func (o *Base) FinishStroke() {
	o.position.SetDelta(o.moving.Get())
	o.moving = nil
}

// DidStroke ...
func (o *Base) DidStroke(dx, dy float64) {
	o.SetMoving(dx, dy)
}

func (o *Base) updatePos() {
	// 自layerがframe外に出ないようにする制御
	// ここでmovingを加算してはいけないので素のpositionで取得する
	px, py := o.position.GetInt()
	frameSize := o.frame.Size()
	lw, lh := o.image.Size()

	//左側に全て隠れてしまう場合
	if px+lw-20 < 0 {
		px = -lw + 20
	}
	// 上に全て隠れてしまう場合
	if py+lh-20 < 0 {
		py = -lh + 20
	}
	// 右に全て隠れてしまう場合
	if px+20 > frameSize.W() {
		px = frameSize.W() - 20
	}
	// 下に全て隠れてしまう場合
	if py+20 > frameSize.H() {
		py = frameSize.H() - 20
	}
	o.position.Set(float64(px), float64(py))
}
