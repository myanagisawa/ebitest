package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Base ...
type Base struct {
	label       string
	manager     interfaces.GameManager
	frames      []interfaces.Frame
	activeFrame interfaces.Frame
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

// Manager ...
func (o *Base) Manager() interfaces.GameManager {
	return o.manager
}

// AddFrame ...
func (o *Base) AddFrame(f interfaces.Frame) {
	f.SetParent(o)
	o.frames = append(o.frames, f)
}

// FrameAt ...
func (o *Base) FrameAt(x, y int) interfaces.Frame {
	for i := len(o.frames) - 1; i >= 0; i-- {
		f := o.frames[i]
		if f.In(x, y) {
			return f
		}
	}
	return nil
}

// ActiveFrame ...
func (o *Base) ActiveFrame() interfaces.Frame {
	return o.activeFrame
}

// GetObjects ...
func (o *Base) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	for i := len(o.frames) - 1; i >= 0; i-- {
		c := o.frames[i]
		objs = append(objs, c.GetObjects(x, y)...)
	}
	// log.Printf("SceneBase::GetObjects: %#v", objs)
	return objs
}

// Update ...
func (o *Base) Update() error {

	o.activeFrame = o.FrameAt(ebiten.CursorPosition())

	for _, frame := range o.frames {
		frame.Update()
	}

	return nil
}

// Draw ...
func (o *Base) Draw(screen *ebiten.Image) {

	for _, frame := range o.frames {
		frame.Draw(screen)
	}
}

// Layout ...
func (o *Base) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// NewScene ...
func NewScene(label string, m interfaces.GameManager) interfaces.Scene {

	s := &Base{
		label:   label,
		manager: m,
	}

	return s
}
