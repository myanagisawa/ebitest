package kitchen

import "github.com/hajimehoshi/ebiten"

type (
	// Scene ...
	Scene interface {
		Update() error
		Draw(r *ebiten.Image)
	}

	// Looper ...
	Looper struct {
		num int
		vec int
		min int
		max int
	}
)

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
