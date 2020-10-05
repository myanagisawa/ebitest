package scenes

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (

	// BattleScene ...
	BattleScene struct {
		manager      *GameManager
		eventHandler *EventHandler
		layers       []Layer
		activeLayer  Layer
	}
)

// NewBattleScene ...
func NewBattleScene(m *GameManager) Scene {

	s := &BattleScene{
		manager: m,
	}

	s.eventHandler = &EventHandler{
		events: map[string]map[*Event]struct{}{},
	}

	s.layers = append(s.layers, NewLayerBase())
	s.layers = append(s.layers, NewBattleMap(s))

	l1 := NewTestWindow()
	l1.parent = s
	s.layers = append(s.layers, l1)

	return s
}

// SetEvent ...
func (s *BattleScene) SetEvent(name string, e *Event) {
	if s.eventHandler.events[name] != nil {
		s.eventHandler.events[name][e] = struct{}{}
	} else {
		m := map[*Event]struct{}{e: {}}
		s.eventHandler.events[name] = m
	}
}

// LayerAt ...
func (s *BattleScene) LayerAt(x, y int) Layer {
	for i := len(s.layers) - 1; i >= 0; i-- {
		l := s.layers[i]
		if l.IsModal() {
			return l
		}
		if l.In(x, y) {
			return l
		}
	}

	return nil
}

// Update ...
func (s *BattleScene) Update(screen *ebiten.Image) error {
	s.activeLayer = s.LayerAt(ebiten.CursorPosition())
	if s.activeLayer != nil {
		// log.Printf("activeLayer: %#v", s.activeLayer.Label())

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// click イベントを発火
			x, y := s.activeLayer.LocalPosition(ebiten.CursorPosition())
			s.eventHandler.Firing(s, "click", x, y)
		}
	}

	for _, layer := range s.layers {
		layer.Update(screen)
	}

	return nil
}

// Draw ...
func (s *BattleScene) Draw(screen *ebiten.Image) {
	for _, layer := range s.layers {
		layer.Draw(screen)
	}
	active := " - "
	if s.activeLayer != nil {
		active = s.activeLayer.Label()
	}

	dbg := fmt.Sprintf("FPS: %0.2f\nactive: %s", ebiten.CurrentFPS(), active)
	ebitenutil.DebugPrint(screen, dbg)
}

// Layout ...
func (s *BattleScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
