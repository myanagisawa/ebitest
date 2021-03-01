package menu

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
)

type (
	// UIControl ...
	UIControl struct {
		label    string
		bound    *g.Bound
		scale    *g.Scale
		scene    *Scene
		parent   interfaces.UIControl
		children []interfaces.UIControl
		drawer   *DrawProps
	}

	// DrawProps ...
	DrawProps struct {
		withoutDraw  bool
		image        *ebiten.Image
		scale        *g.Scale
		position     *g.Point
		colorScale   *g.ColorScale
		subImageRect *image.Rectangle
	}
)

// NewDefaultDrawer ...
func NewDefaultDrawer(img *ebiten.Image) *DrawProps {

	return NewDrawer(true, img, g.DefScale(), g.DefPoint(), g.NewCS(1.0, 1.0, 1.0, 1.0), nil)
}

// NewDrawer ...
func NewDrawer(wd bool, img *ebiten.Image, scale *g.Scale, pos *g.Point, cs *g.ColorScale, sr *g.Bound) *DrawProps {
	var ir *image.Rectangle
	if sr != nil {
		ir = sr.ToImageRect()
	}

	return &DrawProps{
		withoutDraw:  wd,
		image:        img,
		scale:        scale,
		position:     pos,
		colorScale:   cs,
		subImageRect: ir,
	}
}

// var (
// 	img = utils.CreateRectImage(1, 1, &color.RGBA{255, 255, 255, 255})
// )

func init() {
	rand.Seed(time.Now().UnixNano())

}

// NewFrame ...
func NewFrame(s *Scene, label string, bound *g.Bound) *UIControl {
	// img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})

	f := &UIControl{
		label: label,
		bound: bound,
		scale: g.DefScale(),
		scene: s,
		// drawer: NewDefaultDrawer(ebiten.NewImageFromImage(img)),
	}

	pos := g.NewPoint(0, 0)
	size := g.NewSize(100, 1200)

	f.children = []interfaces.UIControl{
		NewLayer(s, f, fmt.Sprintf("%s.layer-1", label), g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewLayer(s, f, fmt.Sprintf("%s.layer-2", label), g.NewBoundByPosSize(pos.SetDelta(100, 0), size)),
		NewLayer(s, f, fmt.Sprintf("%s.layer-3", label), g.NewBoundByPosSize(pos.SetDelta(100, 0), size)),
		NewLayer(s, f, fmt.Sprintf("%s.layer-4", label), g.NewBoundByPosSize(pos.SetDelta(100, 0), size)),
	}

	return f
}

// NewLayer ...
func NewLayer(s *Scene, parent interfaces.UIControl, label string, bound *g.Bound) *UIControl {
	// img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 192})

	l := &UIControl{
		label:  label,
		bound:  bound,
		scale:  g.DefScale(),
		scene:  s,
		parent: parent,
		// drawer: NewDefaultDrawer(ebiten.NewImageFromImage(img)),
	}

	pos := g.NewPoint(0, 0)
	size := g.NewSize(100, 100)

	l.children = []interfaces.UIControl{
		NewUIControl(s, l, fmt.Sprintf("%s.control-1", label), g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-2", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-3", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-4", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-5", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-6", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-7", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-8", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-9", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-10", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-11", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-12", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
	}

	return l
}

// NewUIControl ...
func NewUIControl(s *Scene, parent interfaces.UIControl, label string, bound *g.Bound) *UIControl {
	// img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 127})

	c := &UIControl{
		label:  label,
		bound:  bound,
		scale:  g.DefScale(),
		scene:  s,
		parent: parent,
		// drawer: NewDefaultDrawer(ebiten.NewImageFromImage(img)),
	}

	pos := g.NewPoint(0, 0)
	size := g.NewSize(100, 20)

	c.children = []interfaces.UIControl{
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-1", label), g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-2", label), g.NewBoundByPosSize(pos.SetDelta(0, 20), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-3", label), g.NewBoundByPosSize(pos.SetDelta(0, 20), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-4", label), g.NewBoundByPosSize(pos.SetDelta(0, 20), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-5", label), g.NewBoundByPosSize(pos.SetDelta(0, 20), size)),
	}

	return c
}

// NewChildControl ...
func NewChildControl(s *Scene, parent interfaces.UIControl, label string, bound *g.Bound) *UIControl {
	// img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 64})

	c := &UIControl{
		label:  label,
		bound:  bound,
		scale:  g.DefScale(),
		scene:  s,
		parent: parent,
		// drawer: NewDefaultDrawer(ebiten.NewImageFromImage(img)),
	}

	return c
}

// Label ...
func (o *UIControl) Label() string {
	return o.label
}

// Parent ...
func (o *UIControl) Parent() interfaces.UIControl {
	return o.parent
}

// GetControls ...
func (o *UIControl) GetControls() []interfaces.UIControl {
	ret := []interfaces.UIControl{o}
	for _, child := range o.children {
		ret = append(ret, child.GetControls()...)
	}
	return ret
}

// Position ...
func (o *UIControl) Position(t enum.ValueTypeEnum) *g.Point {
	dx, dy := 0.0, 0.0
	// if o.Moving() != nil {
	// 	dx, dy = o.Moving().Get()
	// }
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

// Scale ...
func (o *UIControl) Scale(t enum.ValueTypeEnum) *g.Scale {
	return o.scale
}

// Update ...
func (o *UIControl) Update() error {
	// o.drawer.withoutDraw = false

	// iw, ih := o.drawer.image.Size()
	// _, size := o.bound.ToPosSize()

	// o.drawer.scale.Set(float64(size.W())/float64(iw), float64(size.H())/float64(ih))
	// o.drawer.position = o.Position(enum.TypeGlobal)

	// log.Printf("%s: drawer.pos=%#v", o.label, o.drawer.position)
	return nil
}

// Draw ...
func (o *UIControl) Draw(screen *ebiten.Image) {
	// if o.drawer.withoutDraw {
	// 	return
	// }

	// op := &ebiten.DrawImageOptions{}
	// // 描画位置指定
	// op.GeoM.Scale(o.drawer.scale.Get())
	// op.GeoM.Translate(o.drawer.position.Get())

	// op.ColorM.Scale(o.drawer.colorScale.Get())

	// screen.DrawImage(o.drawer.image, op)

	return
}
