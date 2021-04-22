package control

import (
	"image/color"
	"image/draw"
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/char"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

var (
	WindowEventDraggingCallback = func(ev interfaces.UIControl, params map[string]interface{}) {
		dp := g.NewPoint(params["dx"].(float64), params["dy"].(float64))
		ev.Parent().SetMoving(dp)
	}
	WindowEventDragDropCallback = func(ev interfaces.UIControl, params map[string]interface{}) {
		dp := g.NewPoint(params["dx"].(float64), params["dy"].(float64))
		ev.Parent().Bound().SetDelta(dp, nil)
		ev.Parent().SetMoving(nil)
	}
	WindowEventClickCallback = func(ev interfaces.UIControl, params map[string]interface{}) {
		w := ev.Parent().Parent()
		w.Remove()
		log.Printf("callback::click")
	}
)

// NewDefaultClosableWindow ...
func NewDefaultClosableWindow(s interfaces.Scene, bound *g.Bound) interfaces.UIScrollView {

	return NewClosableWindow(s, bound, WindowEventDraggingCallback, WindowEventDragDropCallback, WindowEventClickCallback)
}

// NewClosableWindow ...
func NewClosableWindow(s interfaces.Scene, bound *g.Bound,
	eventDraggingCallback func(ev interfaces.UIControl, params map[string]interface{}),
	eventDragDropCallback func(ev interfaces.UIControl, params map[string]interface{}),
	eventClickCallback func(ev interfaces.UIControl, params map[string]interface{}),
) interfaces.UIControl {
	headerHeight := 20

	// ウィンドウ本体
	img := utils.CreateRectImage(1, 1, &color.RGBA{32, 32, 32, 127})
	l := NewUIControl(s, nil, enum.ControlTypeLayer, "info-layer", bound, g.DefScale(), g.DefCS(), img)

	// ヘッダバー
	_bound := g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(500, headerHeight))
	h := NewUIControl(s, nil, enum.ControlTypeDefault, "header-bar", _bound, g.DefScale(), g.DefCS(), img)
	h.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	if eventDraggingCallback != nil {
		h.EventHandler().AddEventListener(enum.EventTypeDragging, eventDraggingCallback)
	}
	if eventDragDropCallback != nil {
		h.EventHandler().AddEventListener(enum.EventTypeDragDrop, eventDragDropCallback)
	}

	// 親子関係を設定
	l.AppendChild(h)

	// 閉じるボタン
	fset := char.Res.Get(14, enum.FontStyleGenShinGothicBold)
	ti := fset.GetStringImage("×")
	ti = utils.TextColorTo(ti.(draw.Image), &color.RGBA{192, 192, 192, 255})
	size := ti.Bounds().Size()
	_bound = g.NewBoundByPosSize(g.NewPoint(5, 0), g.NewSize(size.X, size.Y))
	btn := NewUIControl(s, nil, enum.ControlTypeDefault, "close-btn", _bound, g.DefScale(), g.DefCS(), ti)
	btn.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.5, 0.5, 0.5, 1.0)
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	if eventClickCallback != nil {
		btn.EventHandler().AddEventListener(enum.EventTypeClick, eventClickCallback)
	}

	// 親子関係を設定
	h.AppendChild(btn)

	return l
}

// NewImmutableWindow ...
func NewImmutableWindow(s interfaces.Scene, bound *g.Bound) interfaces.UIControl {
	// ウィンドウ本体
	img := utils.CreateRectImage(1, 1, &color.RGBA{32, 32, 32, 127})
	l := NewUIControl(s, nil, enum.ControlTypeLayer, "immutable-window", bound, g.DefScale(), g.DefCS(), img)

	return l
}
