package interfaces

import (
	"github.com/hajimehoshi/ebiten"
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
}

// Layer ...
type Layer interface {
	Label() string
	EbiObjects() []*models.EbiObject
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
	IsModal() bool
	In(x, y int) bool
}

// UIControl ...
type UIControl interface {
	Label() string
	EbiObjects() []*models.EbiObject
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
	In(x, y int) bool
	SetLayer(l Layer)
}

// UIButton ...
type UIButton interface {
	UIControl
}
