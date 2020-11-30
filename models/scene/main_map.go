package scene

import (
	"image/color"

	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/frame"
	"github.com/myanagisawa/ebitest/models/layer"
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
			label: "MainMap",
		},
	}

	// サブフレーム1（横）
	f := frame.NewFrame("side frame", ebitest.NewPoint(0, 20), ebitest.NewSize(200, ebitest.Height-220), color.RGBA{127, 127, 200, 255}, false)
	s.AddFrame(f)

	img := ebitest.CreateRectImage(150, 300, color.RGBA{0, 0, 0, 128})
	l := layer.NewLayerBase("Layer1", img, ebitest.NewPoint(25, 25), ebitest.NewScale(1, 1), true)
	f.AddLayer(l)

	// サブフレーム2（下）
	f = frame.NewFrame("bottom frame", ebitest.NewPoint(0, float64(ebitest.Height-200)), ebitest.NewSize(ebitest.Width, 200), color.RGBA{127, 200, 127, 255}, false)
	s.AddFrame(f)

	// メインフレーム
	f = frame.NewFrame("main frame", ebitest.NewPoint(200, 20), ebitest.NewSize(ebitest.Width-200, ebitest.Height-220), color.RGBA{200, 200, 200, 255}, true)
	s.AddFrame(f)

	l = layer.NewLayerBase("map", ebitest.Images["world"], ebitest.NewPoint(0, 0), ebitest.NewScale(1, 1), false)
	f.AddLayer(l)

	return s
}
