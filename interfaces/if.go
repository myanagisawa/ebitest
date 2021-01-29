package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
)

// EbiObject ...
type EbiObject interface {
	Label() string
	In(x, y int) bool
}

// Focusable ...
type Focusable interface {
	ToggleHover()
}

// Draggable ...
type Draggable interface {
	DidStroke(dx, dy float64)
	FinishStroke()
}

// Wheelable ...
type Wheelable interface {
	DidWheel(dx, dy float64)
}

// Scrollable ...
type Scrollable interface {
	GetEdgeType(x, y int) enum.EdgeTypeEnum
	DoScroll(x, y int)
}

// Positionable ...
type Positionable interface {
	Position(enum.ValueTypeEnum) *g.Point
}

// Scaleable ...
type Scaleable interface {
	Scale(enum.ValueTypeEnum) *g.Scale
}

// Anglable ...
type Anglable interface {
	Angle(enum.ValueTypeEnum) int
}

// Movable ...
type Movable interface {
	Moving() *g.Point
	SetMoving(dx, dy float64)
}

// EventOwner ...
type EventOwner interface {
	EventHandler() EventHandler
}

// GameManager ...
type GameManager interface {
	TransitionTo(enum.SceneEnum)
	DataSet(key enum.DataTypeEnum) DataSet
}

// Scene ...
type Scene interface {
	ebiten.Game
	AddFrame(f Frame)
	ActiveFrame() Frame
	Label() string
	Manager() GameManager
	GetObjects(x, y int) []EbiObject
	SetCustomFunc(t enum.FuncTypeEnum, f interface{})
	ExecDidLoad()
	ExecDidActive()
}

// Frame ...
type Frame interface {
	EbiObject
	Positionable
	EventOwner
	Manager() GameManager
	Size() *g.Size
	Parent() Scene
	SetParent(parent Scene)
	AddLayer(l Layer)
	LayerAt(x, y int) Layer
	ActiveLayer() Layer
	GetObjects(x, y int) []EbiObject
	Update() error
	Draw(screen *ebiten.Image)
}

// Layer ...
type Layer interface {
	EbiObject
	Positionable
	Scaleable
	Movable
	EventOwner
	Size() *g.Size
	Manager() GameManager
	Frame() Frame
	SetFrame(frame Frame)
	IsModal() bool
	Scroll(et enum.EdgeTypeEnum)
	GetObjects(x, y int) []EbiObject
	Update() error
	Draw(screen *ebiten.Image)
	AddUIControl(c UIControl)
	UIControlAt(x, y int) UIControl
}

// UIControl ...
type UIControl interface {
	EbiObject
	Positionable
	Scaleable
	Anglable
	Focusable
	Movable
	EventOwner
	Manager() GameManager
	Layer() Layer
	SetLayer(l Layer)
	SetPosition(x, y float64)
	SetScale(x, y float64)
	SetAngle(a int)
	GetObjects(x, y int) []EbiObject
	Update() error
	Draw(screen *ebiten.Image)
}

// UIScrollView ...
type UIScrollView interface {
	UIControl
	SetDataSource(colNames []interface{}, data [][]interface{})
	SetRowClickFunc(headerfunc func(row interface{}, pos *g.Point, params map[string]interface{}), rowfunc func(idx int, row interface{}, pos *g.Point, params map[string]interface{}))
}

// ListRow ...
type ListRow interface {
	IsHeader() bool
	Index() int
	Parent() UIScrollView
	Source() []string
}

// ListRowClickable ...
type ListRowClickable interface {
	DidClickRowCallBack(idx int, row interface{})
}

// EventHandler ...
type EventHandler interface {
	AddEventListener(t enum.EventTypeEnum, callback func(o EventOwner, pos *g.Point, params map[string]interface{}))
	Firing(t enum.EventTypeEnum, c EventOwner, pos *g.Point, params map[string]interface{})
	Has(t enum.EventTypeEnum) bool
}

// Event ...
type Event interface {
}

// // Stroke ...
// type Stroke interface {
// 	Update()
// 	PositionDiff() (float64, float64)
// }

// DataSet ...
type DataSet interface {
	GetByCode(code string) AppData
}

// AppData ...
type AppData interface {
}
