package wmap

import (
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"time"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/char"
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

	// ウィンドウ本体
	img := utils.CreateRectImage(1, 1, &color.RGBA{32, 32, 32, 127})
	bound := g.NewBoundByPosSize(g.NewPoint(10, 50), g.NewSize(500, 900))
	l := control.NewUIControl(s, nil, enum.ControlTypeLayer, "info-layer", bound, g.DefScale(), g.DefCS(), img)

	// ヘッダバー
	bound = g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(500, 20))
	h := control.NewUIControl(s, nil, enum.ControlTypeDefault, "header-bar", bound, g.DefScale(), g.DefCS(), img)
	h.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	h.EventHandler().AddEventListener(enum.EventTypeDragging, func(ev interfaces.UIControl, params map[string]interface{}) {
		dp := g.NewPoint(params["dx"].(float64), params["dy"].(float64))
		ev.Parent().SetMoving(dp)
	})
	h.EventHandler().AddEventListener(enum.EventTypeDragDrop, func(ev interfaces.UIControl, params map[string]interface{}) {
		dp := g.NewPoint(params["dx"].(float64), params["dy"].(float64))
		ev.Parent().Bound().SetDelta(dp, nil)
		ev.Parent().SetMoving(nil)
	})

	// 親子関係を設定
	l.AppendChild(h)

	// 閉じるボタン
	fset := char.Res.Get(14, enum.FontStyleGenShinGothicBold)
	ti := fset.GetStringImage("×")
	ti = utils.TextColorTo(ti.(draw.Image), &color.RGBA{192, 192, 192, 255})
	size := ti.Bounds().Size()
	bound = g.NewBoundByPosSize(g.NewPoint(5, 0), g.NewSize(size.X, size.Y))
	btn := control.NewUIControl(s, nil, enum.ControlTypeDefault, "close-btn", bound, g.DefScale(), g.DefCS(), ti)
	btn.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.5, 0.5, 0.5, 1.0)
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	btn.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.UIControl, params map[string]interface{}) {
		window := ev.Parent().Parent()
		window.Remove()
		log.Printf("callback::click")
	})

	// 親子関係を設定
	h.AppendChild(btn)

	return l
}

// NewScrollView ...
func NewScrollView(s interfaces.Scene) interfaces.UIScrollView {

	// スクロールビューのベース
	img := utils.CreateRectImage(1, 1, &color.RGBA{32, 32, 32, 32})
	bound := g.NewBoundByPosSize(g.NewPoint(0, 20), g.NewSize(500, 880))
	_b := control.NewUIControl(s, nil, enum.ControlTypeDefault, "scroll-view", bound, g.DefScale(), g.DefCS(), img)
	sv := control.NewUIScrollView(_b.(*control.UIControl), nil, nil, nil).(*control.UIScrollView)

	// ヘッダ部分のベース
	img = utils.CreateRectImage(1, 1, &color.RGBA{64, 64, 64, 255})
	bound = g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(492, 25))
	hb := control.NewUIControl(s, nil, enum.ControlTypeDefault, "header-base", bound, g.DefScale(), g.DefCS(), img)
	header := control.NewScrollViewHeader(hb.(*control.UIControl), &color.RGBA{64, 64, 64, 255}).(*control.ScrollViewHeader)

	// 親子関係を設定
	sv.AppendChild(header)
	sv.SetHeader(header)

	// スクロール部分のベース
	img = utils.CreateRectImage(1, 1, &color.RGBA{0, 0, 0, 127})
	bound = g.NewBoundByPosSize(g.NewPoint(0, header.Bound().Max.Y()+control.LineSpacing), g.NewSize(492, 860))
	lb := control.NewUIControl(s, nil, enum.ControlTypeDefault, "list-base", bound, g.DefScale(), g.DefCS(), img).(*control.UIControl)
	listView := control.NewScrollViewList(lb).(*control.ScrollViewList)
	listView.EventHandler().AddEventListener(enum.EventTypeWheel, func(ev interfaces.UIControl, params map[string]interface{}) {
		if children := ev.GetChildren(); children != nil {
			dy := params["yoff"].(float64)
			if lv, ok := ev.(*control.ScrollViewList); ok {

				// スクロール結果がlistのbound外になる場合はスクロールしない
				lb := lv.ScrollBound()
				b := lv.Bound()
				if lb.Min.Y()+dy > 0 {
					// 上に余白ができる
					// log.Printf("上に余白ができる")
					return
				} else if lb.Max.Y()+dy < b.Max.Y()-b.Min.Y() {
					// 下に余白ができる
					// log.Printf("下に余白ができる")
					return
				}

				for _, row := range lv.GetChildren() {
					row.Bound().SetDelta(g.NewPoint(0, dy), nil)
				}
				lv.SetScrollBarPosition()
			}

		} else {
			log.Printf("wheel: ev=%T", ev)
		}
	})

	// 親子関係を設定
	sv.AppendChild(listView)
	sv.SetList(listView)

	// スクロールバーのベース
	img = utils.CreateRectImage(1, 1, &color.RGBA{127, 127, 127, 255})
	bound = g.NewBoundByPosSize(g.NewPoint(listView.Bound().Max.X(), header.Bound().Max.Y()+control.LineSpacing), g.NewSize(8, 860))
	sb := control.NewUIControl(s, nil, enum.ControlTypeDefault, "scroll-bar", bound, g.DefScale(), g.DefCS(), img).(*control.UIControl)

	// 親子関係を設定
	sv.AppendChild(sb)
	sv.SetScrollBar(sb)

	// 関連漬け
	sb.AddRel("list", listView)
	listView.AddRel("bar", sb)

	// スクロールバーのスライダ
	img = utils.CreateRectImage(1, 1, &color.RGBA{192, 192, 192, 255})
	bound = g.NewBoundByPosSize(g.NewPoint(2, control.LineSpacing), g.NewSize(4, 10))
	sbb := control.NewUIControl(s, nil, enum.ControlTypeDefault, "scroll-bar-slider", bound, g.DefScale(), g.DefCS(), img)
	sbb.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	sbb.EventHandler().AddEventListener(enum.EventTypeDragging, func(ev interfaces.UIControl, params map[string]interface{}) {
		if base, ok := ev.Rel("base").(interfaces.UIControl); ok {
			dy := params["dy"].(float64)
			dp := g.NewPoint(0, dy*-1)
			if lv, ok := base.Rel("list").(*control.ScrollViewList); ok {

				// スクロール結果がlistのbound外になる場合はスクロールしない
				lb := lv.ScrollBound()
				b := lv.Bound()
				if lb.Min.Y()+dp.Y() > 0 {
					// 上に余白ができる
					// log.Printf("上に余白ができる")
					return
				} else if lb.Max.Y()+dp.Y() < b.Max.Y()-b.Min.Y() {
					// 下に余白ができる
					// log.Printf("下に余白ができる")
					return
				}

				for _, row := range lv.GetChildren() {
					row.SetMoving(dp)
				}
				lv.SetScrollBarPosition()

			}
		}
	})
	sbb.EventHandler().AddEventListener(enum.EventTypeDragDrop, func(ev interfaces.UIControl, params map[string]interface{}) {
		if base, ok := ev.Rel("base").(interfaces.UIControl); ok {
			if lv, ok := base.Rel("list").(*control.ScrollViewList); ok {

				for _, row := range lv.GetChildren() {
					if c, ok := row.(*control.UIControl); ok {
						_, dy := c.Moving().Get()
						c.Bound().SetDelta(g.NewPoint(0, dy), nil)
						c.SetMoving(nil)
					}
				}
				lv.SetScrollBarPosition()
			}
		}
	})

	// 親子関係を設定
	sb.AppendChild(sbb)

	// 関連漬け
	sb.AddRel("slider", sbb)
	sbb.AddRel("base", sb)

	return sv
}
