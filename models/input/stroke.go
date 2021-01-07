package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// StrokeSource represents a input device to provide strokes.
type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}

// MouseStrokeSource is a StrokeSource implementation of mouse.
type MouseStrokeSource struct{}

// Position ...
func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

// IsJustReleased ...
func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

// TouchStrokeSource is a StrokeSource implementation of touch.
type TouchStrokeSource struct {
	ID ebiten.TouchID
}

// Position ...
func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

// IsJustReleased ...
func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}
