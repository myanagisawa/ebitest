package control

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/event"
)

const (
	// edgeSize 画面の端から何ピクセルを端とするか
	edgeSize = 30
	// edgeSizeOuter Window外の何ピクセルまでを端判定に含めるか
	edgeSizeOuter = 100
)

type (
	// UIControl ...
	UIControl struct {
		t              enum.ControlTypeEnum
		label          string
		bound          g.Bound
		scale          g.Scale
		angle          g.Angle
		vector         g.Vector
		colorScale     g.ColorScale
		scene          interfaces.Scene
		parent         interfaces.UIControl
		children       []interfaces.UIControl
		relations      map[string]interface{}
		eventHandler   interfaces.EventHandler
		drawer         *DrawProps
		updateProc     func(self interfaces.UIControl)
		moving         *g.Point
		_childrenCache []interfaces.UIControl
	}

	// DrawProps ...
	DrawProps struct {
		withoutDraw  bool
		image        *ebiten.Image
		imageSize    *g.Size
		drawSize     *g.Size
		scale        *g.Scale
		angle        *g.Angle
		position     *g.Point
		colorScale   *g.ColorScale
		subImageRect *image.Rectangle
	}
)

// NewDefaultDrawer ...
func NewDefaultDrawer(img *ebiten.Image) *DrawProps {

	return NewDrawer(false, img, g.DefScale(), g.DefAngle(), g.DefPoint(), g.DefCS(), nil)
}

// NewDrawer ...
func NewDrawer(wd bool, img *ebiten.Image, scale *g.Scale, angle *g.Angle, pos *g.Point, cs *g.ColorScale, sr *g.Bound) *DrawProps {
	var ir *image.Rectangle
	if sr != nil {
		ir = sr.ToImageRect()
	}

	return &DrawProps{
		withoutDraw:  wd,
		image:        img,
		imageSize:    g.NewSize(img.Size()),
		scale:        scale,
		angle:        angle,
		position:     pos,
		colorScale:   cs,
		subImageRect: ir,
	}
}

// NewUIControl ...
func NewUIControl(s interfaces.Scene, parent interfaces.UIControl, t enum.ControlTypeEnum, label string, bound *g.Bound, scale *g.Scale, cs *g.ColorScale, img image.Image) interfaces.UIControl {
	o := &UIControl{
		t:            t,
		label:        label,
		bound:        *bound,
		scale:        *scale,
		colorScale:   *cs,
		scene:        s,
		parent:       parent,
		relations:    map[string]interface{}{},
		eventHandler: event.NewEventHandler(),
		drawer:       NewDefaultDrawer(ebiten.NewImageFromImage(img)),
	}

	return o
}

// Type ...
func (o *UIControl) Type() enum.ControlTypeEnum {
	return o.t
}

// Label ...
func (o *UIControl) Label() string {
	return o.label
}

// Scene ...
func (o *UIControl) Scene() interfaces.Scene {
	return o.scene
}

// Parent ...
func (o *UIControl) Parent() interfaces.UIControl {
	return o.parent
}

// Rel ...
func (o *UIControl) Rel(key string) interface{} {
	return o.relations[key]
}

// AddRel ...
func (o *UIControl) AddRel(key string, value interface{}) {
	o.relations[key] = value
}

// SetParent ...
func (o *UIControl) SetParent(parent interfaces.UIControl) {
	o.parent = parent
}

// GetChildren ...
func (o *UIControl) GetChildren() []interfaces.UIControl {
	return o.children
}

// GetControls ...
func (o *UIControl) GetControls() []interfaces.UIControl {
	if o._childrenCache != nil {
		return o._childrenCache
	}
	ret := _getControls(o)
	o._childrenCache = ret
	return ret
}

// AppendChild ...
func (o *UIControl) AppendChild(child interfaces.UIControl) {
	child.SetParent(o)
	o.children = append(o.children, child)

	o.removeChildrenCache()
}

func (o *UIControl) removeChildrenCache() {
	o._childrenCache = nil
	if o.parent != nil {
		o.parent.(*UIControl).removeChildrenCache()
	}
}

// SetChildren ...
func (o *UIControl) SetChildren(children []interfaces.UIControl) {
	o.children = children
}

// Remove ...
func (o *UIControl) Remove() {
	if o.parent == nil {
		return
	}
	o.parent.RemoveChild(o)
}

// RemoveChild ...
func (o *UIControl) RemoveChild(child interfaces.UIControl) {
	if o.children == nil || len(o.children) == 0 {
		return
	}
	newChildren := make([]interfaces.UIControl, len(o.children))
	i := 0
	for _, c := range o.children {
		if c != child {
			newChildren[i] = c
			i++
		}
	}
	newChildren = newChildren[:i]
	o.children = newChildren
	o.removeChildrenCache()
}

// Position ...
func (o *UIControl) Position(t enum.ValueTypeEnum) *g.Point {
	dx, dy := 0.0, 0.0
	if o.moving != nil {
		dx, dy = o.moving.Get()
	}
	pos, _ := o.bound.ToPosSize()
	if t == enum.TypeLocal {
		return g.NewPoint(pos.X()+dx, pos.Y()+dy)
	}
	gx, gy := 0.0, 0.0
	if p := o.Parent(); p != nil {
		gx, gy = p.Position(enum.TypeGlobal).Get()

		sx, sy := p.Scale(enum.TypeGlobal).Get()
		gx += (pos.X() + dx) * sx
		gy += (pos.Y() + dy) * sy
	} else {
		gx = pos.X() + dx
		gy = pos.Y() + dy
	}
	return g.NewPoint(gx, gy)
}

// Bound ...
func (o *UIControl) Bound() *g.Bound {
	return &o.bound
}

// Moving ...
func (o *UIControl) Moving() *g.Point {
	return o.moving
}

// SetMoving ...
func (o *UIControl) SetMoving(dp *g.Point) {
	if dp == nil {
		o.moving = nil
		return
	}
	if o.moving == nil {
		o.moving = dp
	} else {
		o.moving.Set(dp.Get())
	}
}

// Scale ...
func (o *UIControl) Scale(t enum.ValueTypeEnum) *g.Scale {
	return &o.scale
}

// Angle ...
func (o *UIControl) Angle(t enum.ValueTypeEnum) *g.Angle {
	return &o.angle
}

// SetAngle ...
func (o *UIControl) SetAngle(a *g.Angle) {
	o.angle = *a
}

// Vector ...
func (o *UIControl) Vector() *g.Vector {
	return &o.vector
}

// SetVector ...
func (o *UIControl) SetVector(v *g.Vector) {
	o.vector = *v
}

// ColorScale ...
func (o *UIControl) ColorScale() *g.ColorScale {
	return &o.colorScale
}

// SetUpdateProc ...
func (o *UIControl) SetUpdateProc(f func(self interfaces.UIControl)) {
	o.updateProc = f
}

// In ...
func (o *UIControl) In() bool {
	x, y := ebiten.CursorPosition()

	// オブジェクトの位置を取得する
	pos := o.Position(enum.TypeGlobal)

	// 無回転当たり判定オブジェクトを取得
	_, bsize := o.bound.ToPosSize()

	// あたり判定objの原点を取得(だいたい中心を原点としているが描画時の回転軸に合わせてlefttopを原点にする)
	// height側を実線より太くしている関係で、Y側の原点は
	rectCenter := g.NewPoint(pos.X(), pos.Y()+float64(bsize.H())/2)

	// ポインタ座標を当たり判定objとの相対座標に変換
	relativeX := float64(x) - rectCenter.X()
	relativeY := float64(y) - rectCenter.Y()

	// ポインタ座標を座標変換する（回転を打ち消す）
	rad := o.Angle(enum.TypeLocal).Get()
	transformPos := g.NewPoint(
		math.Cos(rad)*relativeX+math.Sin(rad)*relativeY,
		-math.Sin(rad)*relativeX+math.Cos(rad)*relativeY)

	// 当たり判定objと変換したポインタ座標で当たり判定を行う
	objWidth := float64(bsize.W())
	objHeight := float64(bsize.H())
	if 0 <= transformPos.X() && objWidth >= transformPos.X() &&
		-objHeight/2 <= transformPos.Y() && objHeight/2 >= transformPos.Y() {
		// log.Printf("%s にフォーカス！", o.Label())
		return true
	}
	return false
}

// GetEdgeType ...
func (o *UIControl) GetEdgeType() enum.EdgeTypeEnum {
	if o.t != enum.ControlTypeFrame {
		// スクロール対応はひとまずframeのみ
		return enum.EdgeTypeNotEdge
	}
	x, y := ebiten.CursorPosition()

	// 範囲外判定
	if x < o.bound.Min.IntX()-edgeSizeOuter || x > o.bound.Max.IntX()+edgeSizeOuter {
		return enum.EdgeTypeNotEdge
	} else if y < o.bound.Min.IntY()-edgeSizeOuter || y > o.bound.Max.IntY()+edgeSizeOuter {
		return enum.EdgeTypeNotEdge
	}

	minX, maxX := o.bound.Min.IntX()+edgeSize, o.bound.Max.IntX()-edgeSize
	minY, maxY := o.bound.Min.IntY()+edgeSize, o.bound.Max.IntY()-edgeSize

	// 判定
	if x <= minX && y <= minY {
		return enum.EdgeTypeTopLeft
	} else if x > minX && x < maxX && y <= minY {
		return enum.EdgeTypeTop
	} else if x >= maxX && y <= minY {
		return enum.EdgeTypeTopRight
	} else if x >= maxX && y > minY && y < maxY {
		return enum.EdgeTypeRight
	} else if x >= maxX && y >= maxY {
		return enum.EdgeTypeBottomRight
	} else if x > minX && x < maxX && y >= maxY {
		return enum.EdgeTypeBottom
	} else if x <= minX && y >= maxY {
		return enum.EdgeTypeBottomLeft
	} else if x <= minX && y > minY && y < maxY {
		return enum.EdgeTypeLeft
	}
	return enum.EdgeTypeNotEdge
}

// Scroll ...
func (o *UIControl) Scroll(et enum.EdgeTypeEnum) {
	// 1フレームあたりの増分値
	dp := 20.0
	switch et {
	case enum.EdgeTypeTopLeft:
		o.bound.SetDelta(g.NewPoint(dp, dp), nil)
	case enum.EdgeTypeTop:
		o.bound.SetDelta(g.NewPoint(0, dp), nil)
	case enum.EdgeTypeTopRight:
		o.bound.SetDelta(g.NewPoint(-dp, dp), nil)
	case enum.EdgeTypeRight:
		o.bound.SetDelta(g.NewPoint(-dp, 0), nil)
	case enum.EdgeTypeBottomRight:
		o.bound.SetDelta(g.NewPoint(-dp, -dp), nil)
	case enum.EdgeTypeBottom:
		o.bound.SetDelta(g.NewPoint(0, -dp), nil)
	case enum.EdgeTypeBottomLeft:
		o.bound.SetDelta(g.NewPoint(dp, -dp), nil)
	case enum.EdgeTypeLeft:
		o.bound.SetDelta(g.NewPoint(dp, 0), nil)
	}

	pos, size := o.bound.ToPosSize()
	_, frameSize := o.parent.Bound().ToPosSize()

	if pos.X() > 0 {
		o.bound.SetDelta(g.NewPoint(-dp, 0), nil)
	} else if pos.IntX()+size.W() < frameSize.W() {
		o.bound.SetDelta(g.NewPoint(dp, 0), nil)
	}

	if pos.Y() > 0 {
		o.bound.SetDelta(g.NewPoint(0, -dp), nil)
	} else if pos.IntY()+size.H() < frameSize.H() {
		o.bound.SetDelta(g.NewPoint(0, dp), nil)
	}
}

// EventHandler ...
func (o *UIControl) EventHandler() interfaces.EventHandler {
	return o.eventHandler
}

// Update ...
func (o *UIControl) Update() error {
	if o.vector.GetAmount() != 0.0 {
		o.bound.SetDelta(o.vector.GetDelta(), nil)
		// 跳ね返り判定
		o.bounce()
	}

	if o.updateProc != nil {
		o.updateProc(o)
	}

	// drawer設定
	o.drawer.withoutDraw = false
	iw, ih := o.drawer.imageSize.Get()
	_, size := o.bound.ToPosSize()
	o.drawer.drawSize = g.NewSize(int(float64(size.W())*o.scale.X()), int(float64(size.H())*o.scale.Y()))

	o.drawer.scale.Set(float64(size.W())/float64(iw), float64(size.H())/float64(ih))
	o.drawer.angle = &o.angle
	o.drawer.position = o.Position(enum.TypeGlobal)
	o.drawer.colorScale = &o.colorScale

	o.drawer.subImageRect = nil
	return nil
}

// Draw ...
func (o *UIControl) Draw(screen *ebiten.Image) {
	if o.drawer.withoutDraw {
		return
	}

	op := &ebiten.DrawImageOptions{}
	// 拡大・縮小
	op.GeoM.Scale(o.drawer.scale.Get())

	// 回転
	imgSize := o.drawer.drawSize
	op.GeoM.Translate(-float64(imgSize.W())/2, -float64(imgSize.H())/2)
	op.GeoM.Rotate(o.drawer.angle.Get())
	op.GeoM.Translate(float64(imgSize.W())/2, float64(imgSize.H())/2)

	// 位置
	op.GeoM.Translate(o.drawer.position.Get())

	// 色
	op.ColorM.Scale(o.drawer.colorScale.Get())

	// log.Printf("draw: %#v", o.drawer)
	if o.drawer.subImageRect != nil {
		screen.DrawImage(o.drawer.image.SubImage(*o.drawer.subImageRect).(*ebiten.Image), op)
		// log.Printf("SubImage描画: %#v", o.drawer.subImageRect)
	} else {
		screen.DrawImage(o.drawer.image, op)
	}
	return
}

// bounce ...
func (o *UIControl) bounce() {
	// はみ出し判定
	min := o.Position(enum.TypeGlobal)
	_, size := o.bound.ToPosSize()
	max := g.NewPoint(min.X()+float64(size.W()), min.Y()+float64(size.H()))

	bounced := false
	if min.X() < 0 || max.X() > float64(g.Width) {
		// 180 - angle
		o.vector.Angle().Set(math.Pi - o.vector.Angle().Get())
		bounced = true
	}
	if min.Y() < 0 || max.Y() > float64(g.Height) {
		// 360 - angle
		o.vector.Angle().Set((2 * math.Pi) - o.vector.Angle().Get())
		bounced = true
	}
	if bounced {
		o.bound.SetDelta(o.vector.GetDelta(), nil)
	}
}

func _getControls(o interfaces.UIControl) []interfaces.UIControl {
	ret := []interfaces.UIControl{o}
	for _, child := range o.GetChildren() {
		ret = append(ret, child.GetControls()...)
	}
	return ret
}
