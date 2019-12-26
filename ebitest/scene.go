package ebitest

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

type (
	// Scene ...
	Scene struct {
		size int
	}
)

// NewScene ...
func NewScene(size int) (*Scene, error) {
	s := &Scene{
		size: size,
	}
	return s, nil
}

// Update ...
func (s *Scene) Update(input *Input) error {

	return nil
}

// Draw draws the scene to the given sceneImage.
func (s *Scene) Draw(sceneImage *ebiten.Image) {
	log.Printf("Scene: Draw: ")
}
