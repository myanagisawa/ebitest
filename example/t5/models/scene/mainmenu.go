package scene

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models/control"
	"github.com/myanagisawa/ebitest/example/t5/models/layer"
)

type (
	// MainMenu ...
	MainMenu struct {
		Base
	}
)

// NewMainMenu ...
func NewMainMenu(m interfaces.GameManager) *MainMenu {

	s := &MainMenu{
		Base: Base{
			label: "MainMenu",
		},
	}

	img := ebitest.CreateRectImage(ebitest.Width, ebitest.Height, color.RGBA{150, 150, 150, 255})
	l := layer.NewLayerBase("Layer1", img, s, nil, nil, 0)
	s.SetLayer(l)

	img = ebitest.CreateRectImage(400, 600, color.RGBA{0, 0, 0, 128})
	l = layer.NewLayerBase("Layer2", img, s, nil, ebitest.NewPoint(10.0, 50.0), 0)
	s.SetLayer(l)

	c := control.NewButton("Battle(dev)", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 50, 100)
	l.AddUIControl(c)
	// l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
	// 	log.Printf("btnBattle clicked")
	// 	source.scene.Manager().TransitionToBattleScene()
	// })

	// img = ebitest.CreateRectImage(600, 400, color.RGBA{255, 32, 32, 128})
	// l = layer.NewLayerBase("Layer3", img, s, nil, ebitest.NewPoint(200.0, 100.0), 0)
	// s.SetLayer(l)

	// img = ebitest.CreateRectImage(600, 400, color.RGBA{32, 32, 255, 128})
	// l = layer.NewLayerBase("Layer4", img, s, nil, ebitest.NewPoint(100.0, 200.0), 0)
	// s.SetLayer(l)

	// img = ebitest.CreateRectImage(600, 400, color.RGBA{32, 255, 32, 64})
	// l = layer.NewLayerBase("Layer5", img, s, ebitest.NewScale(0.5, 0.5), ebitest.NewPoint(500.0, 500.0), 0)
	// s.SetLayer(l)

	// img = ebitest.CreateRectImage(300, 200, color.RGBA{255, 32, 32, 64})
	// l = layer.NewLayerBase("Layer6", img, s, ebitest.NewScale(1.0, 0.5), ebitest.NewPoint(700.0, 300.0), 30)
	// s.SetLayer(l)

	return s
}

// Update ...
func (s *MainMenu) Update(screen *ebiten.Image) error {
	s.activeLayer = s.LayerAt(ebiten.CursorPosition())
	if s.activeLayer != nil {
		// log.Printf("activeLayer: %#v", s.activeLayer.Label())
	}

	for _, layer := range s.layers {
		layer.Update(screen)
	}

	return nil
}

// Draw ...
func (s *MainMenu) Draw(screen *ebiten.Image) {

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
