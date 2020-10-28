package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	// initX and initY represents the position when dragging starts.
	initX int
	initY int

	// currentX and currentY represents the current position
	currentX int
	currentY int

	released bool

	// draggingObject represents a object (sprite in this case)
	// that is being dragged.
	draggingObject interface{}
}

// Update ...
func (s *Stroke) Update() {
	if s.released {
		return
	}
	if s.source.IsJustReleased() {
		s.released = true
		return
	}
	x, y := s.source.Position()
	s.currentX = x
	s.currentY = y
	// log.Printf("stroke: x: %d->%d, y: %d->%d", s.initX, s.currentX, s.initY, s.currentY)
}

// IsReleased ...
func (s *Stroke) IsReleased() bool {
	return s.released
}

// Position ...
func (s *Stroke) Position() (int, int) {
	return s.currentX, s.currentY
}

// PositionDiff ...
func (s *Stroke) PositionDiff() (int, int) {
	dx := s.currentX - s.initX
	dy := s.currentY - s.initY
	// log.Printf("dx=%d, dy=%d", dx, dy)
	return dx, dy
}

// DraggingObject ...
func (s *Stroke) DraggingObject() interface{} {
	return s.draggingObject
}

// SetDraggingObject ...
func (s *Stroke) SetDraggingObject(object interface{}) {
	s.draggingObject = object
}

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

// NewStroke ...
func NewStroke(source StrokeSource) *Stroke {
	cx, cy := source.Position()
	return &Stroke{
		source:   source,
		initX:    cx,
		initY:    cy,
		currentX: cx,
		currentY: cy,
	}
}
