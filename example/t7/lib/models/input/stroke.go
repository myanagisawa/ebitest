package input

import (
	"time"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
)

var (
	longtapms = 1000
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

	tapping      bool
	dragging     bool
	longpressing bool

	mouseDownTargets []interfaces.UIControl

	targetObject interfaces.UIControl
}

// Update ...
func (s *Stroke) Update() {
	if s.released {
		return
	}
	if s.source.IsJustReleased() {
		s.released = true
		s.setTarget(s.CurrentEvent())
		return
	}
	x, y := s.source.Position()
	s.currentX = x
	s.currentY = y

	if !s.dragging {
		if s.initX != s.currentX || s.initY != s.currentY {
			s.tapping = false
			s.dragging = true
			s.longpressing = false
		}
	}

	if !s.dragging && s.StrokeTime() > longtapms {
		s.tapping = false
		s.dragging = false
		s.longpressing = true
	}
	// イベントターゲットを設定
	s.setTarget(s.CurrentEvent())
	// log.Printf("stroke: x: %d->%d, y: %d->%d", s.initX, s.currentX, s.initY, s.currentY)
}

// CurrentEvent ...
func (s *Stroke) CurrentEvent() enum.EventTypeEnum {
	if !s.released {
		if s.longpressing {
			return enum.EventTypeLongPress
		}
		if s.dragging {
			return enum.EventTypeDragging
		}
		return enum.EventTypeNone
	}
	if s.longpressing {
		return enum.EventTypeLongPressReleased
	}
	if s.dragging {
		return enum.EventTypeDragDrop
	}
	return enum.EventTypeClick
}

// IsReleased ...
func (s *Stroke) IsReleased() bool {
	return s.released
}

// IsDragging ...
func (s *Stroke) IsDragging() bool {
	return s.dragging
}

// IsLongPress ...
func (s *Stroke) IsLongPress() bool {
	return s.longpressing
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

// Target ...
func (s *Stroke) Target() (interfaces.UIControl, bool) {
	t := s.targetObject
	return t, t != nil
}

// setTarget ...
func (s *Stroke) setTarget(et enum.EventTypeEnum) {
	s.getTargetByEventType(et)
}

// setTarget ...
func (s *Stroke) getTargetByEventType(et enum.EventTypeEnum) {
	if s.mouseDownTargets == nil || len(s.mouseDownTargets) == 0 {
		return
	}
	if s.targetObject != nil {
		return
	}
	for i := range s.mouseDownTargets {
		t := s.mouseDownTargets[i]
		if t.EventHandler().Has(et) {
			s.targetObject = t
		}
	}
}

// MouseDownTargets ...
func (s *Stroke) MouseDownTargets() []interfaces.UIControl {
	return s.mouseDownTargets
}

// SetMouseDownTargets ...
func (s *Stroke) SetMouseDownTargets(targets []interfaces.UIControl) {
	s.mouseDownTargets = targets
}

// StrokeTime ...
func (s *Stroke) StrokeTime() int {
	now := time.Now()
	duration := now.Sub(s.initTS)
	return int(duration.Milliseconds())
}

// NewStroke ...
func NewStroke(source StrokeSource) *Stroke {
	cx, cy := source.Position()
	return &Stroke{
		source:       source,
		initTS:       time.Now(),
		initX:        cx,
		initY:        cy,
		currentX:     cx,
		currentY:     cy,
		tapping:      true,
		dragging:     false,
		longpressing: false,
	}
}
