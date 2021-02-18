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
	s.SetCustomFunc(enum.FuncTypeDidLoad, s.didLoad())
	s.SetCustomFunc(enum.FuncTypeDidActive, s.didActive())

	return s
}

// didLoad ...
func (o *Scene) didLoad() func() {
	return func() {
		// メインフレーム
		mainf := frame.NewFrame(o, "main frame", g.NewPoint(0, 0), g.NewSize(g.Width, g.Height), &color.RGBA{65, 105, 225, 255}, false)

		l := layer.NewLayerBase(mainf, "menu group", g.NewPoint(30, 100), g.NewSize(500, 700), &color.RGBA{0, 0, 0, 127}, true)

		c := control.NewSimpleButton(l, "mapへ", utils.CopyImage(g.Images["btnBase"]), g.NewPoint(100, 250), 16, &color.RGBA{0, 0, 255, 255})
		c.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
			log.Printf("callback::click")
			o.Manager().TransitionTo(enum.MapEnum)
		})
		l.AddUIControl(c)
		log.Printf("menu.DidLoad")
	}
}

// didActive ...
func (o *Scene) didActive() func() {
	return func() {
		log.Printf("menu.DidActive")
	}
}
