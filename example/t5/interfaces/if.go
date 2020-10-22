package interfaces

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/models"
)

// GameManager ...
type GameManager interface {
	TransitionTo(enum.SceneEnum)
}

// Scene ...
type Scene interface {
	ebiten.Game
	Label() string
	Draw(screen *ebiten.Image)
	SetLayer(l Layer)
	DeleteLayer(l Layer)
	LayerAt(x, y int) Layer
	ActiveLayer() Layer
	GetLayerByLabel(label string) Layer
}

// Layer ...
type Layer interface {
	Label() string
	LabelFull() string
	EbiObjects() []*models.EbiObject
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
	Scroll(t enum.EdgeTypeEnum)
	In(x, y int) bool
	IsModal() bool
	AddUIControl(c UIControl)
	UIControlAt(x, y int) UIControl
	EventHandler() EventHandler
}

// UIControl ...
type UIControl interface {
	Label() string
	EbiObjects() []*models.EbiObject
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
	In(x, y int) bool
	SetLayer(l Layer)
	HasHoverAction() bool
}

// UIButton ...
type UIButton interface {
	UIControl
}

// UIText ...
type UIText interface {
	UIControl
}

// UIColumn ...
type UIColumn interface {
	UIControl
}

// UIScrollView ...
type UIScrollView interface {
	UIControl
}

// EventHandler ...
type EventHandler interface {
	AddEventListener(c UIControl, name string, callback func(UIControl, Scene, *ebitest.Point))
	Firing(s Scene, name string, x, y int)
	Set(name string, ev Event)
}

// Event ...
type Event interface {
}
