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
	SetParent(parent UIControl)
	Rel(key string) interface{}
	AddRel(key string, value interface{})
	GetChildren() []UIControl
	GetControls() []UIControl
	AppendChild(child UIControl)
	SetChildren(children []UIControl)
	Remove()
	RemoveChild(child UIControl)
	Position(enum.ValueTypeEnum) *g.Point
	Bound() *g.Bound
	Moving() *g.Point
	SetMoving(*g.Point)
	Scale(enum.ValueTypeEnum) *g.Scale
	Angle(enum.ValueTypeEnum) *g.Angle
	SetAngle(a *g.Angle)
	Vector() *g.Vector
	SetVector(v *g.Vector)
	ColorScale() *g.ColorScale
	SetUpdateProc(f func(self UIControl))
	In() bool
	GetEdgeType() enum.EdgeTypeEnum
	Scroll(et enum.EdgeTypeEnum)
	EventHandler() EventHandler
}

// UIScrollView ...
type UIScrollView interface {
	UIControl
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
