package control

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/functions"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/char"
	"github.com/myanagisawa/ebitest/models/event"
	"github.com/myanagisawa/ebitest/utils"
)

// CalcPosition ...
func CalcPosition(o interfaces.IPositionable, t enum.ValueTypeEnum) *g.Point {
	pos := o.Position()
	if t == enum.TypeRaw {
		return pos
	}
	// log.Printf("UIControlBase: Position: %s", o.label)
	dx, dy := 0.0, 0.0
	if o.Moving() != nil {
		dx, dy = o.Moving().Get()
	}
	if t == enum.TypeLocal {
		return g.NewPoint(pos.X()+dx, pos.Y()+dy)
	}
	gx, gy := 0.0, 0.0
	if p := o.Parent(); p != nil {
		gx, gy = CalcPosition(p, enum.TypeGlobal).Get()

		sx, sy := p.Scale().Get()
		gx += (pos.X() + dx) * sx
		gy += (pos.Y() + dy) * sy
	} else {
		gx = pos.X() + dx
		gy = pos.Y() + dy
	}
	return g.NewPoint(gx, gy)
}

// Base ...
type Base struct {
	layer    interfaces.Layer
	label    string
	image    *ebiten.Image
	position *g.Point
	scale    *g.Scale
	angle    float64
	moving   *g.Point
	hover    bool
	visible  bool

	eventHandler interfaces.EventHandler
}

// Children ...
func (o *Base) Children() []interfaces.UIControl {
	// log.Printf("*Base.Children: l=%s, o=%T", o.Label(), o)
	return nil
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

// Size ...
func (o *Base) Size(t enum.SizeTypeEnum) *g.Size {
	switch t {
	case enum.TypeOriginal:
		return g.NewSize(o.image.Size())
	case enum.TypeScaled:
		bx, by := o.image.Size()
		sx := float64(bx) * o.scale.X()
		sy := float64(by) * o.scale.Y()
		return g.NewSize(int(sx), int(sy))
	default:
		return nil
	}
}

// Manager ...
func (o *Base) Manager() interfaces.GameManager {
	return o.Layer().Manager()
}

// Layer ...
func (o *Base) Layer() interfaces.Layer {
	return o.layer
}

// In ...
func (o *Base) In(x, y int) bool {
	if !o.visible {
		return false
	}

	return PointInBound(
		g.NewPoint(float64(x), float64(y)),
		g.NewBoundByPosSize(
			o.Position(enum.TypeGlobal),
			o.Size(enum.TypeScaled),
		),
		g.NewBoundByPosSize(
			o.Layer().Frame().Position(enum.TypeGlobal),
			o.Layer().Frame().Size(),
		),
	)
}

// PointInBound 点pが領域bound内かどうかを返します（親領域areaが指定された時、boundはarea内の範囲に切り詰められます）
func PointInBound(p *g.Point, bound *g.Bound, area *g.Bound) bool {
	minX, minY, maxX, maxY := bound.Min.IntX(), bound.Min.IntY(), bound.Max.IntX(), bound.Max.IntY()
	if area != nil {
		// 座標がフレーム外の場合はフレームのmax座標で置き換え
		if maxX > area.Max.IntX() {
			maxX = area.Max.IntX()
		}
		if maxY > area.Max.IntY() {
			maxY = area.Max.IntY()
		}

		// 座標がフレーム外の場合はフレームのmin座標で置き換え
		if minX < area.Min.IntX() {
			minX = area.Min.IntX()
		}
		if minY < area.Min.IntY() {
			minY = area.Min.IntY()
		}
	}

	x, y := p.GetInt()
	return (x >= minX && x <= maxX) && (y > minY && y <= maxY)
}

// SetLayer ...
func (o *Base) SetLayer(l interfaces.Layer) {
	o.layer = l
}

// Position ...
func (o *Base) Position(t enum.ValueTypeEnum) *g.Point {
	// log.Printf("UIControlBase: Position: %s", o.label)
	dx, dy := 0.0, 0.0
	if o.moving != nil {
		dx, dy = o.moving.Get()
	}
	if t == enum.TypeLocal {
		return g.NewPoint(o.position.X()+dx, o.position.Y()+dy)
	}
	gx, gy := 0.0, 0.0
	if o.Layer() != nil {
		gx, gy = o.Layer().Position(enum.TypeGlobal).Get()

		sx, sy := o.Layer().Scale(enum.TypeGlobal).Get()
		gx += (o.position.X() + dx) * sx
		gy += (o.position.Y() + dy) * sy
	} else {
		gx = o.position.X() + dx
		gy = o.position.Y() + dy
	}
	return g.NewPoint(gx, gy)
}

// SetPosition ...
func (o *Base) SetPosition(x, y float64) {
	o.position = g.NewPoint(x, y)
}

// Scale ...
func (o *Base) Scale(t enum.ValueTypeEnum) *g.Scale {
	return o.scale
}

// SetScale ...
func (o *Base) SetScale(x, y float64) {
	o.scale = g.NewScale(x, y)
}

// Angle ...
func (o *Base) Angle(t enum.ValueTypeEnum) float64 {
	return o.angle
}

// SetAngle ...
func (o *Base) SetAngle(theta float64) {
	o.angle = theta
}

// Theta ...
func (o *Base) Theta() float64 {
	return 2 * math.Pi * float64(o.Angle(enum.TypeGlobal)) / 360.0
}

// Visible ...
func (o *Base) Visible() bool {
	return o.visible
}

// SetVisible ...
func (o *Base) SetVisible(b bool) {
	o.visible = b
}

// Hover ...
func (o *Base) Hover() bool {
	return o.hover
}

// SetHover ...
func (o *Base) SetHover(b bool) {
	o.hover = b
}

// ToggleHover ...
func (o *Base) ToggleHover() {
	o.hover = !o.hover
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

// Objects ...
func (o *Base) Objects(lt enum.ListTypeEnum) []interfaces.IListable {
	objs := []interfaces.IListable{}

	switch lt {
	case enum.ListTypeCursor:
		x, y := ebiten.CursorPosition()
		if o.In(x, y) {
			objs = append(objs, o)
		}
	default:
		objs = append(objs, o)
	}

	// log.Printf("SceneBase::GetObjects: %#v", objs)
	return objs
}

// GetObjects ...
func (o *Base) GetObjects(x, y int) []interfaces.EbiObject {
	if o.In(x, y) {
		return []interfaces.EbiObject{o}
	}
	return nil
}

// Update ...
func (o *Base) Update() error {
	return nil
}

// Draw draws the sprite.
func (o *Base) Draw(screen *ebiten.Image) {
	o.DrawWithOptions(screen, nil)
}

// DrawWithOptions draws the sprite.
func (o *Base) DrawWithOptions(screen *ebiten.Image, in *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	if !o.visible {
		// log.Printf("%sは不可視です。", o.label)
		return in
	}
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	// op.GeoM.Reset()
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())

	bgSize := g.NewSize(o.image.Size())
	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(bgSize.W())/2, -float64(bgSize.H())/2)
	// 中心を軸に回転
	op.GeoM.Rotate(o.Angle(enum.TypeGlobal))
	// ユニットの座標に移動
	op.GeoM.Translate(float64(bgSize.W())/2, float64(bgSize.H())/2)

	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	// log.Printf("%s: pos: %d, %d, fr: %d, %d, %d, %d", o.label, lx, ly, x0, y0, x1, y1)
	// screen.DrawImage(o.image.SubImage(fr).(*ebiten.Image), op)
	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if o.hover {
		r, g, b, a = 0.5, 0.5, 0.5, 1.0
	}
	op.ColorM.Scale(r, g, b, a)
	// if strings.HasPrefix(o.label, "route") {
	// 	log.Printf("Draw: %s: pos=%#v scale=%#v angle=%#v", o.label, o.Position(enum.TypeGlobal), o.Scale(enum.TypeGlobal), o.Angle(enum.TypeGlobal))
	// }
	if in != nil {
		op.GeoM.Concat(in.GeoM)
		op.ColorM.Concat(in.ColorM)
	}
	screen.DrawImage(o.image, op)

	// log.Printf("control.Draw: l=%s, o=%T", o.Label(), o)
	return op
}

// EventHandler ...
func (o *Base) EventHandler() interfaces.EventHandler {
	return o.eventHandler
}

// NewControlBase ...
func NewControlBase(l interfaces.Layer, label string, eimg *ebiten.Image, pos *g.Point, scale *g.Scale, visible bool) interfaces.UIControl {
	o := &Base{
		layer:        l,
		label:        label,
		image:        eimg,
		position:     pos,
		scale:        scale,
		visible:      visible,
		eventHandler: event.NewEventHandler(),
	}

	return o
}

// NewUIControl ...
func NewUIControl(l interfaces.Layer, label string, eimg *ebiten.Image, pos *g.Point, scale *g.Scale, visible bool) interfaces.UIControl {
	o := NewControlBase(l, label, eimg, pos, scale, true).(*Base)
	l.AddUIControl(o)
	return o
}

// NewSimpleLabel ...
func NewSimpleLabel(l interfaces.Layer, label string, pos *g.Point, pt int, c *color.RGBA, family enum.FontStyleEnum) interfaces.UIControl {
	fset := char.Res.Get(pt, family)
	ti := fset.GetStringImage(label)
	ti2 := utils.TextColorTo(ti.(draw.Image), c)
	eimg := ebiten.NewImageFromImage(ti2)

	o := NewControlBase(l, label, eimg, pos, g.DefScale(), true).(*Base)
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)

	l.AddUIControl(o)
	return o
}

// NewSimpleButton ...
func NewSimpleButton(l interfaces.Layer, label string, img image.Image, pos *g.Point, pt int, c *color.RGBA) interfaces.UIControl {
	fset := char.Res.Get(pt, enum.FontStyleGenShinGothicBold)
	ti := fset.GetStringImage(label)
	ti = utils.TextColorTo(ti.(draw.Image), c)
	ti = utils.ImageOnTextToCenter(img.(draw.Image), ti)
	eimg := ebiten.NewImageFromImage(ti)

	o := NewControlBase(l, label, eimg, pos, g.DefScale(), true).(*Base)
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)

	l.AddUIControl(o)
	return o
}
