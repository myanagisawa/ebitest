package kitchen

import "github.com/hajimehoshi/ebiten"

type (
	// Scene ...
	Scene interface {
		Update() error
		Draw(r *ebiten.Image)
	}
)
