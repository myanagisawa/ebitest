package ex3

import "github.com/hajimehoshi/ebiten"

type (
	// Circle ...
	Circle struct {
		x, y float64
		r    int
	}

	// Scene ...
	Scene interface {
		Update() error
		Draw(r *ebiten.Image)
		GetSize() (int, int)
	}

	// SceneImpl ...
	SceneImpl struct {
	}

	// Size ...
	Size struct {
		Width  int
		Height int
	}

	// Looper ...
	Looper struct {
		num int
		vec int
		min int
		max int
	}
)

// NewSceneImpl ...
func NewSceneImpl() *SceneImpl {
	s := &SceneImpl{}
	return s
}

// Update ...
func (s *SceneImpl) Update() error {
	return nil
}

// Draw ...
func (s *SceneImpl) Draw(r *ebiten.Image) {
}

// GetSize ...
func (s *SceneImpl) GetSize() (int, int) {
	return 0, 0
}

// NewLooper ...
func NewLooper(num, step, min, max int) *Looper {
	l := &Looper{num, step, min, max}
	return l
}

// Get ...
func (l *Looper) Get() int {
	ret := l.num
	if l.num <= l.min || l.num >= l.max {
		l.vec = l.vec * -1
	}
	l.num += l.vec
	return ret
}
