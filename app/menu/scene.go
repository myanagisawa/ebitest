package menu

import (
	"image/color"
	"log"

	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/control"
	"github.com/myanagisawa/ebitest/models/frame"
	"github.com/myanagisawa/ebitest/models/layer"
	"github.com/myanagisawa/ebitest/models/scene"
	"github.com/myanagisawa/ebitest/utils"
)

type (
	// Scene ...
	Scene struct {
		scene.Base
	}
)

// NewScene ...
func NewScene(m interfaces.GameManager) *Scene {

	s := &Scene{
		Base: *scene.NewScene("menu scene", m).(*scene.Base),
	}

	// メインフレーム
	mainf := frame.NewFrame("main frame", g.NewPoint(0, 0), g.NewSize(g.Width, g.Height), &color.RGBA{65, 105, 225, 255}, false)
	s.AddFrame(mainf)

	l := layer.NewLayerBase("menu group", g.NewPoint(30, 100), g.NewSize(500, 700), &color.RGBA{0, 0, 0, 127}, true)
	mainf.AddLayer(l)

	c := control.NewSimpleButton("mapへ", utils.CopyImage(g.Images["btnBase"]), g.NewPoint(100, 250), 16, &color.RGBA{0, 0, 255, 255})
	c.EventHandler().AddEventListener(enum.EventTypeClick, func(o interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		log.Printf("callback::click")
		m.TransitionTo(enum.MapEnum)
	})
	l.AddUIControl(c)

	return s
}
