package ebitest

import "github.com/hajimehoshi/ebiten"

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
		transitionCount int
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

// Update ...
func (s *SceneManager) Update(input *Input) error {
	if s.transitionCount == 0 {
		return s.current.Update(&GameState{
			SceneManager: s,
			Input:        input,
		})
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
		return
	}

	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transitionTo.Clear()
	s.next.Draw(transitionTo)

	r.DrawImage(transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(transitionTo, op)
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
