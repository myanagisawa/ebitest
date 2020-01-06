package ebitest

import (
	"fmt"
	"math/rand"

	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	transitionMaxCount = 20
)

type (
	// Scene ...
	Scene interface {
		Update(state *GameState) error
		Draw(r *ebiten.Image)
	}

	// SceneManager ...
	SceneManager struct {
		current         Scene
		next            Scene
		objects         []Scene
		objectImages    []ebiten.Image
		transitionCount int
		paths           []string
	}

	// GameState ...
	GameState struct {
		SceneManager *SceneManager
		Input        *Input
	}
)

var (
	transitionFrom *ebiten.Image
	transitionTo   *ebiten.Image
)

func init() {
	transitionFrom, _ = ebiten.NewImage(ScreenWidth, ScreenHeight, ebiten.FilterDefault)
	transitionTo, _ = ebiten.NewImage(ScreenWidth, ScreenHeight, ebiten.FilterDefault)
}

// PathToImage ...
func (s *SceneManager) PathToImage(idx int) *ebiten.Image {
	if idx > len(s.paths) {
		panic("idx out of range.")
	}
	path := s.paths[idx]
	img, _, err := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return img
}

// Update ...
func (s *SceneManager) Update(input *Input) error {
	c, b := input.Control()
	if b {
		if c.String() == "A" {
			idx := rand.Intn(len(s.objectImages) - 1)
			s.objects = append(s.objects, NewObject(&s.objectImages[idx], ScreenWidth, ScreenHeight))
			log.Println("type A")
		}
	}

	state := &GameState{
		SceneManager: s,
		Input:        input,
	}

	for _, o := range s.objects {
		o.Update(state)
	}

	if s.transitionCount == 0 {
		return s.current.Update(state)
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

// Draw ...
func (s *SceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
	} else {
		s.DrawTransition(r)
	}
	for _, o := range s.objects {
		o.Draw(r)
	}
}

// GoTo ...
func (s *SceneManager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}

// DrawTransition ...
func (s *SceneManager) DrawTransition(r *ebiten.Image) {
	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transitionTo.Clear()
	s.next.Draw(transitionTo)

	r.DrawImage(transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(transitionTo, op)
	ebitenutil.DebugPrint(r, fmt.Sprintf("count=%d", s.transitionCount))
}
