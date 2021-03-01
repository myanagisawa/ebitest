package scene

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Base ...
type Base struct {
	label       string
	manager     interfaces.GameManager
	frames      []interfaces.Frame
	activeFrame interfaces.Frame

	customFuncDidLoad   func()
	customFuncDidActive func()
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

// Frames ...
func (o *Base) Frames() []interfaces.Frame {
	return o.frames
}

// Objects ...
func (o *Base) Objects(lt enum.ListTypeEnum) []interfaces.IListable {
	objs := []interfaces.IListable{}

	for _, frame := range o.frames {
		if c, ok := frame.(interfaces.IListable); ok {
			objs = append(objs, c.Objects(lt)...)
		}
	}

	switch lt {
	case enum.ListTypeCursor:
		// カーソルの場合は逆順リストに変換
		// （上に重なってる方を優先したいから）
		sort.Slice(objs, func(i, j int) bool {
			return i > j
		})
	}
	// log.Printf("SceneBase::GetObjects: %#v", objs)
	return objs
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
	return
}

// // Draw ...
// func (o *Base) Draw(screen *ebiten.Image) {

// 	for _, frame := range o.frames {
// 		frame.Draw(screen)
// 	}
// }

// Layout ...
func (o *Base) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// SetCustomFunc ...
func (o *Base) SetCustomFunc(t enum.FuncTypeEnum, f interface{}) {
	switch t {
	case enum.FuncTypeDidLoad:
		if v, ok := f.(func()); ok {
			o.customFuncDidLoad = v
		}
	case enum.FuncTypeDidActive:
		if v, ok := f.(func()); ok {
			o.customFuncDidActive = v
		}
	}
}

// ExecDidLoad ...
func (o *Base) ExecDidLoad() {
	if o.customFuncDidLoad != nil {
		o.customFuncDidLoad()
	}
}

// ExecDidActive ...
func (o *Base) ExecDidActive() {
	if o.customFuncDidActive != nil {
		o.customFuncDidActive()
	}
}

// NewScene ...
func NewScene(label string, m interfaces.GameManager) interfaces.Scene {

	s := &Base{
		label:   label,
		manager: m,
	}

	return s
}
