package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/ebitest"
)

// GlobalPositionable ...
type GlobalPositionable interface {
	GlobalPosition() *ebitest.Point
}

// GlobalScaleable ...
type GlobalScaleable interface {
	GlobalScale() *ebitest.Scale
}

// GlobalAnglable ...
type GlobalAnglable interface {
	GlobalAngle() int
}

// Scene ...
type Scene interface {
	ebiten.Game
	// Label() string
	// SetLayer(l Layer)
	// DeleteLayer(l Layer)
	// LayerAt(x, y int) Layer
	// ActiveLayer() Layer
	// GetLayerByLabel(label string) Layer
}

// // Layer ...
// type Layer interface {
// 	Label() string
// 	LabelFull() string
// 	EbiObjects() []*models.EbiObject
// 	Update() error
// 	Draw(screen *ebiten.Image)
// 	Scroll(t enum.EdgeTypeEnum)
// 	In(x, y int) bool
// 	IsModal() bool
// 	AddUIControl(c UIControl)
// 	UIControlAt(x, y int) UIControl
// 	EventHandler() EventHandler
// }

// // UIControl ...
// type UIControl interface {
// 	Label() string
// 	EbiObjects() []*models.EbiObject
// 	Update() error
// 	Draw(screen *ebiten.Image)
// 	In(x, y int) bool
// 	SetLayer(l Layer)
// 	HasHoverAction() bool
// }

// // UIButton ...
// type UIButton interface {
// 	UIControl
// }

// // UIText ...
// type UIText interface {
// 	UIControl
// }

// // UIColumn ...
// type UIColumn interface {
// 	UIControl
// }

// // UIScrollView ...
// type UIScrollView interface {
// 	UIControl
// }

// // EventHandler ...
// type EventHandler interface {
// 	AddEventListener(c UIControl, name string, callback func(UIControl, Scene, *ebitest.Point))
// 	Firing(s Scene, name string, x, y int)
// 	Set(name string, ev Event)
// }

// // Event ...
// type Event interface {
// }
