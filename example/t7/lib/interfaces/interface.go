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
	DidLoad()
	DidActive()
}

// UIControl ...
type UIControl interface {
	Label() string
	Update() error
	Draw(screen *ebiten.Image)
	Parent() UIControl
	GetControls() []UIControl
	Position(enum.ValueTypeEnum) *g.Point
	Scale(enum.ValueTypeEnum) *g.Scale
}

// UIDialog ...
type UIDialog interface {
	UIControl
}
