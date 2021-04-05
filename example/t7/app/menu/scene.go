package menu

import (
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/game"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
)

type (
	// Scene ...
	Scene struct {
		children []interfaces.UIControl
		manager  *game.Manager
	}
)

// NewScene ...
func NewScene(manager *game.Manager) *Scene {
	s := &Scene{
		manager: manager,
	}

	pos := g.NewPoint(0, 0)
	size := g.NewSize(300, 1200)

	s.children = []interfaces.UIControl{
		NewFrame(s, "frame-1", g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewFrame(s, "frame-2", g.NewBoundByPosSize(pos.SetDelta(300, 0), size)),
		NewFrame(s, "frame-3", g.NewBoundByPosSize(pos.SetDelta(300, 0), size)),
		NewFrame(s, "frame-4", g.NewBoundByPosSize(pos.SetDelta(300, 0), size)),
		NewButtonControl(s, "シーン切り替え", g.NewBoundByPosSize(g.NewPoint(10, 10), g.NewSize(199, 40))),
	}

	return s
}

// Label ...
func (o *Scene) Label() string {
	return "menu"
}

// TransitionTo ...
func (o *Scene) TransitionTo(t enum.SceneEnum) {
	o.manager.TransitionTo(t)
}

// GetControls ...
func (o *Scene) GetControls() []interfaces.UIControl {
	ret := []interfaces.UIControl{}
	for _, child := range o.children {
		ret = append(ret, child.GetControls()...)
	}
	return ret
}

// Update ...
func (o *Scene) Update() error {
	return nil
}

// DidLoad ...
func (o *Scene) DidLoad() {
	log.Printf("menu.DidLoad")
}

// DidActive ...
func (o *Scene) DidActive() {
	log.Printf("menu.DidActive")
}
