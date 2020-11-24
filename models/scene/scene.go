package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Base ...
type Base struct {
	label       string
	frames      []interfaces.Frame
	activeFrame interfaces.Frame
}

// Label ...
func (o *Base) Label() string {
	return o.label
}

/// AddFrame ...
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

/// Update ...
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
func NewScene(m interfaces.GameManager) interfaces.Scene {

	s := &Base{
		label: "BaseScene",
	}

	return s
}
