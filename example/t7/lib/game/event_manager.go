package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/input"
)

var (
	prevHoverd interfaces.UIControl
	stroke     *input.Stroke

	evparams map[string]interface{}
)

// SetStroke ...
func SetStroke(controls []interfaces.UIControl) {
	if stroke != nil {
		t := stroke.MouseDownTargets()
		log.Printf("別のstrokeがあるため追加できません: targets=%#v", t)
		return
	}
	targets := make([]interfaces.UIControl, len(controls))
	idx := 0
	for _, control := range controls {
		if control.In() {
			targets[idx] = control
			idx++
		}
	}
	targets = targets[:idx]

	if len(targets) > 0 {
		// log.Printf("targets: %#v", targets)
		stroke = input.NewStroke(&input.MouseStrokeSource{})
		stroke.SetMouseDownTargets(targets)
	}
}

// GetEventTarget ...
func GetEventTarget(controls []interfaces.UIControl, et enum.EventTypeEnum) (interfaces.UIControl, bool) {

	for _, control := range controls {
		if control.EventHandler() != nil && control.EventHandler().Has(et) {
			switch et {
			case enum.EventTypeScroll:
				if frame := control.Parent(); frame.Type() == enum.ControlTypeFrame {
					if frame.GetEdgeType() != enum.EdgeTypeNotEdge {
						return control, true
					}
				}
			default:
				if control.In() {
					return control, true
				}
			}
		}
	}
	return nil, false
}

// ExecEvent ...
func ExecEvent(controls []interfaces.UIControl) error {
	// log.Printf("------アプデ-------")
	evparams = make(map[string]interface{})

	// カーソル処理
	{
		// ホバーイベント
		if current, ok := GetEventTarget(controls, enum.EventTypeFocus); ok {
			// log.Printf("0: current=%p, prev=%p", current, prevHoverd)
			if current == prevHoverd {
				// log.Printf("same target")
			} else {
				// フォーカス対象が変わった

				if prevHoverd != nil {
					// 前のフォーカスを外す処理
					evparams["event_type"] = enum.EventTypeBlur
					prevHoverd.EventHandler().Firing(enum.EventTypeFocus, prevHoverd, evparams)
				}

				// 新しいフォーカス処理
				evparams["event_type"] = enum.EventTypeFocus
				current.EventHandler().Firing(enum.EventTypeFocus, current, evparams)
				prevHoverd = current
				// log.Printf("1: prevHoverd=%p", prevHoverd)
			}
		} else {
			// フォーカス対象なし
			if prevHoverd != nil {
				// フォーカス中のコントロールあり
				evparams["event_type"] = enum.EventTypeBlur
				prevHoverd.EventHandler().Firing(enum.EventTypeFocus, prevHoverd, evparams)
			}

			prevHoverd = nil
		}

		// タップイベント
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// マウスタップ
			SetStroke(controls)
		}

		// タップ状態からの状態遷移イベント処理（クリック、D&D、ロングタップ）
		if stroke != nil {
			stroke.Update()

			dx, dy := stroke.PositionDiff()
			evparams["dx"], evparams["dy"] = dx, dy

			// eventCompleted := false
			if target, ok := stroke.Target(); ok {
				// log.Printf("stroke:target: %T", target)
				switch stroke.CurrentEvent() {
				case enum.EventTypeClick:
					target.EventHandler().Firing(enum.EventTypeClick, target, evparams)
					// eventCompleted = true
				case enum.EventTypeDragging:
					target.EventHandler().Firing(enum.EventTypeDragging, target, evparams)
				case enum.EventTypeDragDrop:
					target.EventHandler().Firing(enum.EventTypeDragDrop, target, evparams)
					// eventCompleted = true
				case enum.EventTypeLongPress:
					target.EventHandler().Firing(enum.EventTypeLongPress, target, evparams)
				case enum.EventTypeLongPressReleased:
					target.EventHandler().Firing(enum.EventTypeLongPressReleased, target, evparams)
					// eventCompleted = true
				}
				// tname := fmt.Sprintf("%s", reflect.TypeOf(target))
				// log.Printf("EventType(%d): target: %s", g.stroke.CurrentEvent(), tname)
			}

			if stroke.IsReleased() {
				log.Printf("EventCompleted")
				stroke = nil
			}
		}

		// 他のイベントが発生していない場合
		if prevHoverd == nil && stroke == nil {
			// スクロールイベント
			if target, ok := GetEventTarget(controls, enum.EventTypeScroll); ok {
				target.EventHandler().Firing(enum.EventTypeScroll, target, evparams)
				// log.Printf("EventType(Scroll): target: %s", target.Label())
			}
		}
	}

	// ホイール処理
	{
		xoff, yoff := ebiten.Wheel()
		if xoff != 0 || yoff != 0 {
			if target, ok := GetEventTarget(controls, enum.EventTypeWheel); ok {
				// log.Printf("EventType(Wheel): target: %T", target)

				evparams["xoff"], evparams["yoff"] = xoff, yoff
				target.EventHandler().Firing(enum.EventTypeWheel, target, evparams)
			}
		}
	}

	return nil
}
