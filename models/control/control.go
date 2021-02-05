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

// Base ...
type Base struct {
	label string

	image    *ebiten.Image
	layer    interfaces.Layer
	position *g.Point
	scale    *g.Scale
	angle    float64
	moving   *g.Point
	hover    bool

	eventHandler interfaces.EventHandler
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
	return controlIn(x, y,
		o.Position(enum.TypeGlobal),
		g.NewSize(o.image.Size()),
		o.Scale(enum.TypeGlobal),
		o.Layer().Frame().Position(enum.TypeGlobal),
		o.Layer().Frame().Size())
}

// controlIn
func controlIn(x, y int, pos *g.Point, size *g.Size, scale *g.Scale, framePos *g.Point, frameSize *g.Size) bool {
	// パーツ位置（左上座標）
	minX, minY := pos.GetInt()

	// 見かけ上の右下座標を取得
	maxX := int(float64(size.W())*scale.X()) + minX
	maxY := int(float64(size.H())*scale.Y()) + minY

	// フレーム領域
	fPosX, fPosY := framePos.GetInt()
	fMaxX, fMaxY := fPosX+frameSize.W(), fPosY+frameSize.H()
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

// // Draw draws the sprite.
// func (o *Base) Draw(screen *ebiten.Image) *ebiten.DrawImageOptions {
// 	op := &ebiten.DrawImageOptions{}

// 	// 描画位置指定
// 	op.GeoM.Reset()
// 	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())

// 	bgSize := g.NewSize(o.image.Size())
// 	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
// 	op.GeoM.Translate(-float64(bgSize.W())/2, -float64(bgSize.H())/2)
// 	// 中心を軸に回転
// 	op.GeoM.Rotate(o.Theta())
// 	// ユニットの座標に移動
// 	op.GeoM.Translate(float64(bgSize.W())/2, float64(bgSize.H())/2)

// 	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

// 	// log.Printf("%s: pos: %d, %d, fr: %d, %d, %d, %d", o.label, lx, ly, x0, y0, x1, y1)
// 	// screen.DrawImage(o.image.SubImage(fr).(*ebiten.Image), op)
// 	r, g, b, a := 1.0, 1.0, 1.0, 1.0
// 	if o.hover {
// 		r, g, b, a = 0.5, 0.5, 0.5, 1.0
// 	}
// 	op.ColorM.Scale(r, g, b, a)
// 	screen.DrawImage(o.image, op)

// 	return op
// }

// Draw draws the sprite.
func (o *Base) Draw(screen *ebiten.Image) {
	_ = o.DrawWithOptions(screen, nil)
}

// DrawWithOptions draws the sprite.
func (o *Base) DrawWithOptions(screen *ebiten.Image, op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	if op == nil {
		op = &ebiten.DrawImageOptions{}
	}

	// 描画位置指定
	op.GeoM.Reset()
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
	screen.DrawImage(o.image, op)

	return op
}

// EventHandler ...
func (o *Base) EventHandler() interfaces.EventHandler {
	return o.eventHandler
}

// NewControlBase ...
func NewControlBase(label string, eimg *ebiten.Image, pos *g.Point) interfaces.UIControl {
	o := &Base{
		label:        label,
		image:        eimg,
		position:     pos,
		scale:        g.NewScale(1.0, 1.0),
		eventHandler: event.NewEventHandler(),
	}

	return o
}

// NewSimpleLabel ...
func NewSimpleLabel(label string, pos *g.Point, pt int, c *color.RGBA, family enum.FontStyleEnum) interfaces.UIControl {
	fset := char.Res.Get(pt, family)
	ti := fset.GetStringImage(label)
	ti2 := utils.TextColorTo(ti.(draw.Image), c)
	eimg := ebiten.NewImageFromImage(ti2)

	o := &Base{
		label:        label,
		image:        eimg,
		position:     pos,
		scale:        g.NewScale(1.0, 1.0),
		eventHandler: event.NewEventHandler(),
	}
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)

	return o
}

// NewSimpleButton ...
func NewSimpleButton(label string, img image.Image, pos *g.Point, pt int, c *color.RGBA) interfaces.UIControl {
	fset := char.Res.Get(pt, enum.FontStyleGenShinGothicBold)
	ti := fset.GetStringImage(label)
	ti = utils.TextColorTo(ti.(draw.Image), c)
	ti = utils.ImageOnTextToCenter(img.(draw.Image), ti)
	eimg := ebiten.NewImageFromImage(ti)

	o := &Base{
		label:        label,
		image:        eimg,
		position:     pos,
		scale:        g.NewScale(1.0, 1.0),
		eventHandler: event.NewEventHandler(),
	}
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)

	return o
}
