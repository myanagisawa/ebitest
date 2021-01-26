package functions

import (
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
)

var (

	// CommonEventCallback 共通イベントコールバック
	CommonEventCallback = func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		et := params["event"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeClick:
		case enum.EventTypeFocus:
			if t, ok := ev.(interfaces.Focusable); ok {
				t.ToggleHover()
			}
		case enum.EventTypeBlur:
		case enum.EventTypeDragging:
			if t, ok := ev.(interfaces.Draggable); ok {
				t.DidStroke(params["dx"].(float64), params["dy"].(float64))
			}
		case enum.EventTypeDragDrop:
			if t, ok := ev.(interfaces.Draggable); ok {
				t.FinishStroke()
			}
		case enum.EventTypeLongPress:
		case enum.EventTypeLongPressReleased:
		case enum.EventTypeWheel:
			if t, ok := ev.(interfaces.Wheelable); ok {
				t.DidWheel(params["xoff"].(float64), params["yoff"].(float64))
			}
		}
		// if t, ok := o.(interfaces.EbiObject); ok {
		// 	log.Printf("Event: %s: %s", et.Name(), t.Label())
		// }
	}
)
