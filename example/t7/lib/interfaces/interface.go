package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
)

// Scene ...
type Scene interface {
	Label() string
	GetControls() []UIControl
	TransitionTo(t enum.SceneEnum)
	DidLoad()
	DidActive()
}

// UIControl ...
type UIControl interface {
	Type() enum.ControlTypeEnum
	Label() string
	Update() error
	Draw(screen *ebiten.Image)
	Scene() Scene
	Parent() UIControl
	GetControls() []UIControl
	RemoveChild(UIControl)
	Position(enum.ValueTypeEnum) *g.Point
	Bound() *g.Bound
	SetMoving(*g.Point)
	Scale(enum.ValueTypeEnum) *g.Scale
	Angle(enum.ValueTypeEnum) *g.Angle
	ColorScale() *g.ColorScale
	In() bool
	GetEdgeType() enum.EdgeTypeEnum
	Scroll(et enum.EdgeTypeEnum)
	EventHandler() EventHandler
}

// UIDialog ...
type UIDialog interface {
	UIControl
}

// EventHandler ...
type EventHandler interface {
	AddEventListener(t enum.EventTypeEnum, callback func(o UIControl, params map[string]interface{}))
	Firing(t enum.EventTypeEnum, c UIControl, params map[string]interface{})
	Has(t enum.EventTypeEnum) bool
}

// Event ...
type Event interface {
}
