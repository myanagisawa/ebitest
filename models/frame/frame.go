package frame

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Base ...
type Base struct {
	label    string
	image    *ebiten.Image
	parent   interfaces.Scene
	position *ebitest.Point

	layers      []interfaces.Layer
	activeLayer interfaces.Layer
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

// Parent ...
func (o *Base) Parent() interfaces.Scene {
	return o.parent
}

// SetParent ...
func (o *Base) SetParent(parent interfaces.Scene) {
	o.parent = parent
}

/// AddLayer ...
func (o *Base) AddLayer(l interfaces.Layer) {
	l.SetParent(o)
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

// In returns true if (x, y) is in the sprite, and false otherwise.
func (o *Base) In(x, y int) bool {
	px, py := o.Position(enum.TypeLocal).GetInt()
	return o.image.At(x-px, y-py).(color.RGBA).A > 0
}

// Position ...
func (o *Base) Position(t enum.ValueTypeEnum) *ebitest.Point {
	// 親（scene）は位置を持たないので、常に自分のPositionを返せばOK
	return o.position
}

// Size ...
func (o *Base) Size() *ebitest.Size {
	return ebitest.NewSize(o.image.Size())
}

// GetEdgeType ...
func (o *Base) GetEdgeType(x, y int) enum.EdgeTypeEnum {
	posX, posY := o.position.GetInt()
	frameW, frameH := o.image.Size()

	minX, maxX := posX+ebitest.EdgeSize, posX+frameW-ebitest.EdgeSize
	minY, maxY := posY+ebitest.EdgeSize, posY+frameH-ebitest.EdgeSize

	// 範囲外判定
	if x < posX-ebitest.EdgeSizeOuter || x > posX+frameW+ebitest.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	} else if y < posY-ebitest.EdgeSizeOuter || y > posY+frameH+ebitest.EdgeSizeOuter {
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

// Update ...
func (o *Base) Update() error {
	if len(o.layers) > 0 {
		et := o.GetEdgeType(ebiten.CursorPosition())
		if et != enum.EdgeTypeNotEdge {
			o.layers[0].Scroll(et)
		}
	}

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

	i := ebiten.NewImage(o.image.Size())
	for _, layer := range o.layers {
		layer.Draw(i)
	}
	screen.DrawImage(i, op)

	if o.parent.ActiveFrame() != nil && o.parent.ActiveFrame() == o {
		n := "-"
		if o.activeLayer != nil {
			x, y := o.activeLayer.Position(enum.TypeLocal).GetInt()
			n = fmt.Sprintf("%s(%d, %d)", o.activeLayer.Label(), x, y)
		}
		dbg := fmt.Sprintf("\n%s / %s", o.label, n)
		ebitenutil.DebugPrint(screen, dbg)
	}
}

// Layout ...
func (o *Base) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// NewFrame ...
func NewFrame(label string, pos *ebitest.Point, size *ebitest.Size, c color.RGBA) interfaces.Frame {
	img := ebiten.NewImageFromImage(ebitest.CreateRectImage(size.W(), size.H(), c))

	s := &Base{
		label:    label,
		image:    img,
		position: pos,
	}

	return s
}
