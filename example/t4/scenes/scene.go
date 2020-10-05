package scenes

import "github.com/hajimehoshi/ebiten"

// Scene ...
type Scene interface {
	ebiten.Game
	Draw(screen *ebiten.Image)
	SetEvent(name string, e *Event)
	SetLayer(l Layer)
	DeleteLayer(l Layer)
	Manager() *GameManager
	LayerAt(x, y int) Layer
	ActiveLayer() Layer
}

// SceneBase ...
type SceneBase struct {
	manager      *GameManager
	eventHandler *EventHandler
	layers       []Layer
	activeLayer  Layer
}

// Draw ...
func (s *SceneBase) Draw(screen *ebiten.Image) {
	return
}

// SetEvent ...
func (s *SceneBase) SetEvent(name string, e *Event) {
	if s.eventHandler.events[name] != nil {
		s.eventHandler.events[name][e] = struct{}{}
	} else {
		m := map[*Event]struct{}{e: {}}
		s.eventHandler.events[name] = m
	}
}

// Manager ...
func (s *SceneBase) Manager() *GameManager {
	return s.manager
}

// LayerAt ...
func (s *SceneBase) LayerAt(x, y int) Layer {
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

// ActiveLayer ...
func (s *SceneBase) ActiveLayer() Layer {
	return s.activeLayer
}

// SetLayer ...
func (s *SceneBase) SetLayer(l Layer) {
	s.layers = append(s.layers, l)
}

// DeleteLayer ...
func (s *SceneBase) DeleteLayer(l Layer) {
	var layers []Layer
	for _, layer := range s.layers {
		if l != layer {
			layers = append(layers, layer)
		}
	}
	s.layers = layers
}

// Layout ...
func (s *SceneBase) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
