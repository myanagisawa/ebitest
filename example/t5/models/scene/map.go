package scene

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models/control"
	"github.com/myanagisawa/ebitest/example/t5/models/layer"
)

var (
	dbg string
)

type (
	// Map ...
	Map struct {
		Base
	}
)

// NewMap ...
func NewMap(m interfaces.GameManager) *Map {

	s := &Map{
		Base: Base{
			label: "Map",
		},
	}

	l := layer.NewMapLayer("Layer1", ebitest.Images["world"], s, ebitest.NewScale(1.0, 1.0), nil, 0, false)
	l.Sites = m.GetSites()
	l.Routes = m.GetRoutes(l.Sites)
	s.SetLayer(l)
	c := control.NewButton("menu", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 500, 500)
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, scene interfaces.Scene, point *ebitest.Point) {
		log.Printf("%s clicked", target.Label())
		m.TransitionTo(enum.MainMenuEnum)
	})
	c = control.NewButton("拡大/縮小", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 50, 30)
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, scene interfaces.Scene, point *ebitest.Point) {
		log.Printf("%s clicked", target.Label())
		layer := scene.GetLayerByLabel("Layer1")
		s := layer.EbiObjects()[0].Scale()
		if s.X() >= 1.0 {
			s.Set(0.5, 0.5)
		} else {
			s.Set(1.0, 1.0)
		}
	})

	return s
}

// Update ...
func (s *Map) Update() error {
	et := GetEdgeType(ebiten.CursorPosition())
	if et != enum.EdgeTypeNotEdge {
		s.layers[0].Scroll(et)
	}

	s.activeLayer = s.LayerAt(ebiten.CursorPosition())
	if s.activeLayer != nil {
		// log.Printf("activeLayer: %#v", s.activeLayer.Label())
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// click イベントを発火
			s.activeLayer.EventHandler().Firing(s, "click", x, y)
		}
	}

	for _, layer := range s.layers {
		layer.Update()
	}

	return nil
}

// Draw ...
func (s *Map) Draw(screen *ebiten.Image) {

	for _, layer := range s.layers {
		layer.Draw(screen)
	}

	active := " - "
	control := " - "
	if s.activeLayer != nil {
		eo := s.activeLayer.EbiObjects()[0]
		px, py := eo.GlobalPosition().Get()
		active = fmt.Sprintf("%s: (%d, %d)", s.activeLayer.LabelFull(), int(px), int(py))
		c := s.activeLayer.UIControlAt(ebiten.CursorPosition())
		if c != nil {
			px, py = c.EbiObjects()[0].GlobalPosition().Get()
			control = fmt.Sprintf("%s: (%d, %d)", c.Label(), int(px), int(py))
		}
	}

	x, y := ebiten.CursorPosition()
	dbg := fmt.Sprintf("%s\nTPS: %0.2f\nFPS: %0.2f\npos: (%d, %d)\nactive:\n - layer: %s\n - control: %s", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS(), x, y, active, control)
	// dbg = fmt.Sprintf("%s", printMemoryStats())
	ebitenutil.DebugPrint(screen, dbg)
}
