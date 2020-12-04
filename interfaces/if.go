package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
)

// Positionable ...
type Positionable interface {
	Position(enum.ValueTypeEnum) *ebitest.Point
}

// Scaleable ...
type Scaleable interface {
	Scale(enum.ValueTypeEnum) *ebitest.Scale
}

// Anglable ...
type Anglable interface {
	Angle(enum.ValueTypeEnum) int
}

// Movable ...
type Movable interface {
	Moving() *ebitest.Point
	SetMoving(dx, dy float64)
}

// GameManager ...
type GameManager interface {
	TransitionTo(enum.SceneEnum)
}

// Scene ...
type Scene interface {
	ebiten.Game
	AddFrame(f Frame)
	ActiveFrame() Frame
	// Label() string
	// SetLayer(l Layer)
	// DeleteLayer(l Layer)
	// LayerAt(x, y int) Layer
	// ActiveLayer() Layer
	// GetLayerByLabel(label string) Layer
}

// Frame ...
type Frame interface {
	Positionable
	Size() *ebitest.Size
	Parent() Scene
	SetParent(parent Scene)
	AddLayer(l Layer)
	In(x, y int) bool
	LayerAt(x, y int) Layer
	ActiveLayer() Layer
	Update() error
	Draw(screen *ebiten.Image)
}

// Layer ...
type Layer interface {
	Positionable
	Scaleable
	Movable
	Label() string
	Frame() Frame
	SetFrame(frame Frame)
	In(x, y int) bool
	IsModal() bool
	Scroll(et enum.EdgeTypeEnum)
	Update() error
	Draw(screen *ebiten.Image)
	AddUIControl(c UIControl)
	UIControlAt(x, y int) UIControl
	EventHandler() EventHandler
}

// UIControl ...
type UIControl interface {
	Positionable
	Scaleable
	Anglable
	Movable
	Label() string
	In(x, y int) bool
	SetLayer(l Layer)
	SetPosition(x, y float64)
	SetScale(x, y float64)
	SetAngle(a int)
	Update() error
	Draw(screen *ebiten.Image)
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

// UIScrollView ...
type UIScrollView interface {
	UIControl
}

// EventHandler ...
type EventHandler interface {
	AddEventListener(c UIControl, name string, callback func(UIControl, *ebitest.Point))
	Firing(s Scene, name string, x, y int)
	Set(name string, ev Event)
}

// Event ...
type Event interface {
}
