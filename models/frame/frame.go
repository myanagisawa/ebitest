package frame

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/event"
	"github.com/myanagisawa/ebitest/utils"
)

// Base ...
type Base struct {
	scene    interfaces.Scene
	label    string
	image    *ebiten.Image
	parent   interfaces.Scene
	position *g.Point

	layers      []interfaces.Layer
	activeLayer interfaces.Layer

	eventHandler interfaces.EventHandler
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

// Manager ...
func (o *Base) Manager() interfaces.GameManager {
	return o.parent.Manager()
}

// Parent ...
func (o *Base) Parent() interfaces.Scene {
	return o.parent
}

// SetParent ...
func (o *Base) SetParent(parent interfaces.Scene) {
	o.parent = parent
}

// AddLayer ...
func (o *Base) AddLayer(l interfaces.Layer) {
	// l.SetFrame(o)
	o.layers = append(o.layers, l)
}

// LayerAt ...
func (o *Base) LayerAt(x, y int) interfaces.Layer {
	for i := len(o.layers) - 1; i >= 0; i-- {
		l := o.layers[i]
		if l.IsModal() {
			return l
		}
		if l.In(x, y) {
			return l
		}
	}

	return nil
}

// ActiveLayer ...
func (o *Base) ActiveLayer() interfaces.Layer {
	return o.activeLayer
}

// Layers ...
func (o *Base) Layers() []interfaces.Layer {
	return o.layers
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (o *Base) In(x, y int) bool {
	px, py := o.Position(enum.TypeLocal).GetInt()
	return o.image.At(x-px, y-py).(color.RGBA).A > 0
}

// Position ...
func (o *Base) Position(t enum.ValueTypeEnum) *g.Point {
	// 親（scene）は位置を持たないので、常に自分のPositionを返せばOK
	return o.position
}

// Size ...
func (o *Base) Size() *g.Size {
	return g.NewSize(o.image.Size())
}

// GetEdgeType ...
func (o *Base) GetEdgeType(x, y int) enum.EdgeTypeEnum {
	posX, posY := o.position.GetInt()
	frameW, frameH := o.image.Size()

	minX, maxX := posX+g.EdgeSize, posX+frameW-g.EdgeSize
	minY, maxY := posY+g.EdgeSize, posY+frameH-g.EdgeSize

	// 範囲外判定
	if x < posX-g.EdgeSizeOuter || x > posX+frameW+g.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	} else if y < posY-g.EdgeSizeOuter || y > posY+frameH+g.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	}

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

// DoScroll ...
func (o *Base) DoScroll(x, y int) {
	if len(o.layers) > 0 {
		et := o.GetEdgeType(x, y)
		if et != enum.EdgeTypeNotEdge {
			o.layers[0].Scroll(et)
		}
	}
}

// Objects ...
func (o *Base) Objects(lt enum.ListTypeEnum) []interfaces.IListable {
	objs := []interfaces.IListable{}

	switch lt {
	case enum.ListTypeCursor:
		x, y := ebiten.CursorPosition()
		if o.In(x, y) || o.GetEdgeType(x, y) != enum.EdgeTypeNotEdge {
			objs = append(objs, o)
		}
	default:
		objs = append(objs, o)
	}

	for _, layer := range o.layers {
		if c, ok := layer.(interfaces.IListable); ok {
			objs = append(objs, c.Objects(lt)...)
		}
	}

	// log.Printf("SceneBase::GetObjects: %#v", objs)
	return objs
}

// GetObjects ...
func (o *Base) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	for i := len(o.layers) - 1; i >= 0; i-- {
		c := o.layers[i]
		objs = append(objs, c.GetObjects(x, y)...)
	}

	if o.In(x, y) || o.GetEdgeType(x, y) != enum.EdgeTypeNotEdge {
		objs = append(objs, o)
	}
	// log.Printf("FrameBase::GetObjects: %#v", objs)
	return objs
}

// Update ...
func (o *Base) Update() error {
	o.activeLayer = o.LayerAt(ebiten.CursorPosition())

	for _, layer := range o.layers {
		layer.Update()
	}

	return nil
}

// Draw ...
func (o *Base) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.position.Get())
	screen.DrawImage(o.image, op)

	if o.parent.ActiveFrame() != nil && o.parent.ActiveFrame() == o {
		n := "-"
		ac := "-"
		if o.activeLayer != nil {
			x, y := o.activeLayer.Position(enum.TypeGlobal).GetInt()
			n = fmt.Sprintf("%s(%d, %d)", o.activeLayer.Label(), x, y)
			c := o.activeLayer.UIControlAt(ebiten.CursorPosition())
			if c != nil {
				x, y = c.Position(enum.TypeGlobal).GetInt()
				ac = fmt.Sprintf("%s(%d, %d)", c.Label(), x, y)
			}
		}
		g.DebugText += fmt.Sprintf("\n%s / %s / %s", o.label, n, ac)
	}
	// log.Printf("frame.Draw: l=%s, o=%T", o.Label(), o)
}

// // Draw ...
// func (o *Base) Draw(screen *ebiten.Image) {
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(o.position.Get())
// 	screen.DrawImage(o.image, op)

// 	for _, layer := range o.layers {
// 		layer.Draw(screen)
// 	}

// 	if o.parent.ActiveFrame() != nil && o.parent.ActiveFrame() == o {
// 		n := "-"
// 		ac := "-"
// 		if o.activeLayer != nil {
// 			x, y := o.activeLayer.Position(enum.TypeGlobal).GetInt()
// 			n = fmt.Sprintf("%s(%d, %d)", o.activeLayer.Label(), x, y)
// 			c := o.activeLayer.UIControlAt(ebiten.CursorPosition())
// 			if c != nil {
// 				x, y = c.Position(enum.TypeGlobal).GetInt()
// 				ac = fmt.Sprintf("%s(%d, %d)", c.Label(), x, y)
// 			}
// 		}
// 		g.DebugText += fmt.Sprintf("\n%s / %s / %s", o.label, n, ac)
// 	}
// }

// Layout ...
func (o *Base) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// EventHandler ...
func (o *Base) EventHandler() interfaces.EventHandler {
	return o.eventHandler
}

// NewFrame ...
func NewFrame(s interfaces.Scene, label string, pos *g.Point, size *g.Size, c *color.RGBA, scrollable bool) interfaces.Frame {
	img := ebiten.NewImageFromImage(utils.CreateRectImage(size.W(), size.H(), c))

	o := &Base{
		scene:        s,
		label:        label,
		image:        img,
		position:     pos,
		eventHandler: event.NewEventHandler(),
	}
	if scrollable {
		o.eventHandler.AddEventListener(enum.EventTypeScroll, func(o interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
			if t, ok := o.(interfaces.Scrollable); ok {
				t.DoScroll(pos.GetInt())
			}
			// log.Printf("callback::scroll")
		})
	}

	s.AddFrame(o)
	return o
}
