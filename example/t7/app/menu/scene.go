package menu

import (
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
)

type (
	// Scene ...
	Scene struct {
		children []interfaces.UIControl
	}
)

// NewScene ...
func NewScene() *Scene {
	s := &Scene{}

	pos := g.NewPoint(0, 0)
	size := g.NewSize(300, 1200)

	s.children = []interfaces.UIControl{
		NewFrame(s, "frame-1", g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewFrame(s, "frame-2", g.NewBoundByPosSize(pos.SetDelta(300, 0), size)),
		NewFrame(s, "frame-3", g.NewBoundByPosSize(pos.SetDelta(300, 0), size)),
		NewFrame(s, "frame-4", g.NewBoundByPosSize(pos.SetDelta(300, 0), size)),
	}

	return s
}

// Label ...
func (o *Scene) Label() string {
	return "menu"
}

// GetControls ...
func (o *Scene) GetControls() []interfaces.UIControl {
	ret := []interfaces.UIControl{}
	for _, child := range o.children {
		ret = append(ret, child.GetControls()...)
	}
	return ret
}

// DidLoad ...
func (o *Scene) DidLoad() {
	log.Printf("menu.DidLoad")
}

// DidActive ...
func (o *Scene) DidActive() {
	log.Printf("menu.DidActive")
}
