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

// StrokeTarget ...
type StrokeTarget interface {
	UpdateStroke(stroke Stroke)
	UpdatePositionByDelta()
}

// GameManager ...
type GameManager interface {
	TransitionTo(enum.SceneEnum)
	SetStroke(target StrokeTarget)
}

// Scene ...
type Scene interface {
	ebiten.Game
	AddFrame(f Frame)
	ActiveFrame() Frame
	Label() string
	Manager() GameManager
	// SetLayer(l Layer)
	// DeleteLayer(l Layer)
	// LayerAt(x, y int) Layer
	// ActiveLayer() Layer
	// GetLayerByLabel(label string) Layer
}

// Frame ...
type Frame interface {
	Positionable
	Label() string
	Manager() GameManager
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
	Manager() GameManager
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
	Manager() GameManager
	Layer() Layer
	In(x, y int) bool
	SetLayer(l Layer)
	SetPosition(x, y float64)
	SetScale(x, y float64)
	SetAngle(a int)
	Update() error
	Draw(screen *ebiten.Image)
}

// UIScrollView ...
type UIScrollView interface {
	UIControl
	SetDataSource(colNames []interface{}, data [][]interface{})
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

// Stroke ...
type Stroke interface {
	Update()
	PositionDiff() (float64, float64)
}
