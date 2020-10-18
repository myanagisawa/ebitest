package scene

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
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

	l := layer.NewLayerBase("Layer1", ebitest.Images["bgFlower"], s, ebitest.NewScale(1.0, 1.0), nil, 0, false)
	s.SetLayer(l)

	c := control.NewButton("一緒にスクロール", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 500, 500)
	l.AddUIControl(c)

	img := ebitest.CreateRectImage(400, 600, color.RGBA{0, 0, 0, 128})
	l = layer.NewLayerBase("Layer2", img, s, ebitest.NewScale(0.7, 0.7), ebitest.NewPoint(10.0, 50.0), 0, true)
	s.SetLayer(l)

	c = control.NewButton("Battle(dev)", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 50, 100)
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
	// l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
	// 	log.Printf("btnBattle clicked")
	// 	source.scene.Manager().TransitionToBattleScene()
	// })

	img = ebitest.CreateRectImage(600, 400, color.RGBA{255, 32, 32, 128})
	l = layer.NewLayerBase("Layer3", img, s, nil, ebitest.NewPoint(200.0, 100.0), 0, false)
	s.SetLayer(l)

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
	control := " - "
	if s.activeLayer != nil {
		active = s.activeLayer.Label()
		c := s.activeLayer.UIControlAt(ebiten.CursorPosition())
		if c != nil {
			control = c.Label()
		}
	}
	dbg := fmt.Sprintf("FPS: %0.2f\nactive: %s\ncontrol: %s", ebiten.CurrentFPS(), active, control)
	ebitenutil.DebugPrint(screen, dbg)

}

// GetEdgeType ...
func GetEdgeType(x, y int) enum.EdgeTypeEnum {
	minX, maxX := ebitest.EdgeSize, ebitest.Width-ebitest.EdgeSize
	minY, maxY := ebitest.EdgeSize, ebitest.Height-ebitest.EdgeSize

	// 範囲外判定
	if x < -ebitest.EdgeSizeOuter || x > ebitest.Width+ebitest.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	} else if y < -ebitest.EdgeSizeOuter || y > ebitest.Height+ebitest.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	}

	// 判定
	if x <= minX && y <= minY {
		return enum.EdgeTypeTopLeft
	} else if x > minX && x < maxX && y <= minY {
		return enum.EdgeTypeTop
	} else if x >= maxX && y <= minY {
		return enum.EdgeTypeTopRight
	} else if x >= maxX && y > minY && y < maxY {
		return enum.EdgeTypeRight
	} else if x >= maxX && y >= maxY {
		return enum.EdgeTypeBottomRight
	} else if x > minX && x < maxX && y >= maxY {
		return enum.EdgeTypeBottom
	} else if x <= minX && y >= maxY {
		return enum.EdgeTypeBottomLeft
	} else if x <= minX && y > minY && y < maxY {
		return enum.EdgeTypeLeft
	}
	return enum.EdgeTypeNotEdge
}
