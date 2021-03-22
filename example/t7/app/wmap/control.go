package wmap

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/control"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

// NewWorldMap ...
func NewWorldMap(s interfaces.Scene) interfaces.UIControl {
	img := utils.CreateRectImage(1, 1, &color.RGBA{255, 255, 255, 255})

	bound := g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(1200, 1200))
	f := control.NewUIControl(s, nil, enum.ControlTypeFrame, "worldmap-frame", bound, g.DefScale(), g.DefCS(), img)

	bound = g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(3120, 2340))
	l := control.NewUIControl(s, nil, enum.ControlTypeLayer, "worldmap-layer", bound, g.DefScale(), g.DefCS(), g.Images["world"])

	l.EventHandler().AddEventListener(enum.EventTypeScroll, func(ev interfaces.UIControl, params map[string]interface{}) {
		ev.Scroll(ev.Parent().GetEdgeType())
		// log.Printf("callback::scroll:: %d", ev.Parent().GetEdgeType())
	})

	// 親子関係を設定
	f.AppendChild(l)

	return f
}

// NewInfoLayer ...
func NewInfoLayer(s interfaces.Scene) interfaces.UIControl {

	bound := g.NewBoundByPosSize(g.NewPoint(10, 50), g.NewSize(500, 900))
	// 閉じるボタン付きウィンドウ生成
	window := control.NewDefaultClosableWindow(s, bound)

	return window
}

// NewScrollView ...
func NewScrollView(s interfaces.Scene) interfaces.UIScrollView {

	bound := g.NewBoundByPosSize(g.NewPoint(0, 20), g.NewSize(500, 880))
	// スクロールビュー生成
	sv := control.NewDefaultScrollView(s, bound)

	return sv
}
