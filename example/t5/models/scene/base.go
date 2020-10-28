package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
)

// Base ...
type Base struct {
	label       string
	layers      []interfaces.Layer
	activeLayer interfaces.Layer
}

// Label ...
func (s *Base) Label() string {
	return s.label
}

// Draw ...
func (s *Base) Draw(screen *ebiten.Image) {
	return
}

// LayerAt ...
func (s *Base) LayerAt(x, y int) interfaces.Layer {
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
func (s *Base) ActiveLayer() interfaces.Layer {
	return s.activeLayer
}

// GetLayerByLabel ...
func (s *Base) GetLayerByLabel(label string) interfaces.Layer {
	for _, layer := range s.layers {
		log.Printf("GetLayerByLabel: %s == %s, %v", layer.Label(), label, layer.Label() == label)
		if layer.Label() == label {
			return layer
		}
	}
	return nil
}

// SetLayer ...
func (s *Base) SetLayer(l interfaces.Layer) {
	s.layers = append(s.layers, l)
}

// DeleteLayer ...
func (s *Base) DeleteLayer(l interfaces.Layer) {
	var layers []interfaces.Layer
	for _, layer := range s.layers {
		if l != layer {
			layers = append(layers, layer)
		}
	}
	s.layers = layers
}

// Layout ...
func (s *Base) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
