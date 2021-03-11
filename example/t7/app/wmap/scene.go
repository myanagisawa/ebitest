package wmap

import (
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
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
		children: []interfaces.UIControl{},
		manager:  manager,
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

// DidLoad ...
func (o *Scene) DidLoad() {
	log.Printf("map.DidLoad")
	f := NewWorldMap(o)
	o.children = append(o.children, f)

	l := NewInfoLayer(o, f)
	f.children = append(f.children, l)
}

// DidActive ...
func (o *Scene) DidActive() {
	log.Printf("map.DidActive")
}
