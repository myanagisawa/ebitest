package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	initTS time.Time

	// initX and initY represents the position when dragging starts.
	initX int
	initY int

	// currentX and currentY represents the current position
	currentX int
	currentY int

	released bool

	dragging bool

	mouseDownTargets []interfaces.StrokeTarget
	// draggingObject represents a object (sprite in this case)
	// that is being dragged.
	targetObject interfaces.StrokeTarget
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

	if !s.dragging {
		if s.initX != s.currentX || s.initY != s.currentY {
			s.dragging = true
		}
	}
	// log.Printf("stroke: x: %d->%d, y: %d->%d", s.initX, s.currentX, s.initY, s.currentY)
}

// IsReleased ...
func (s *Stroke) IsReleased() bool {
	return s.released
}

// IsDragging ...
func (s *Stroke) IsDragging() bool {
	return s.dragging
}

// Position ...
func (s *Stroke) Position() (int, int) {
	return s.currentX, s.currentY
}

// PositionDiff ...
func (s *Stroke) PositionDiff() (float64, float64) {
	dx := float64(s.currentX - s.initX)
	dy := float64(s.currentY - s.initY)
	// log.Printf("dx=%d, dy=%d", dx, dy)
	return dx, dy
}

// DraggingObject ...
func (s *Stroke) DraggingObject() interfaces.StrokeTarget {
	return s.targetObject
}

// SetDraggingObject ...
func (s *Stroke) SetDraggingObject(object interfaces.StrokeTarget) {
	s.targetObject = object
}

// MouseDownTargets ...
func (s *Stroke) MouseDownTargets() []interfaces.StrokeTarget {
	return s.mouseDownTargets
}

// SetMouseDownTargets ...
func (s *Stroke) SetMouseDownTargets(targets []interfaces.StrokeTarget) {
	s.mouseDownTargets = targets
}

// StrokeTime ...
func (s *Stroke) StrokeTime() int {
	now := time.Now()
	duration := now.Sub(s.initTS)
	return int(duration.Milliseconds())
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
		initTS:   time.Now(),
		initX:    cx,
		initY:    cy,
		currentX: cx,
		currentY: cy,
	}
}
