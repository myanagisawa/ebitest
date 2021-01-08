package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/input"
)

var (
	cacheGetObjects []interfaces.EbiObject

	hoverdObject interfaces.EbiObject
)

type (
	// EventManager ...
	EventManager struct {
		manager *Manager
		stroke  *input.Stroke
	}
)

// setStroke ...
func (g *EventManager) setStroke(x, y int) {
	if g.stroke != nil {
		t := g.stroke.MouseDownTargets()
		log.Printf("別のstrokeがあるため追加できません: targets=%#v", t)
		return
	}
	targets := g.GetEventTargetList(x, y, enum.EventTypeClick, enum.EventTypeDragging, enum.EventTypeLongPress)
	if len(targets) > 0 {
		stroke := input.NewStroke(&input.MouseStrokeSource{})
		stroke.SetMouseDownTargets(targets)
		g.stroke = stroke
	}
}

// GetObjects ...
func (g *EventManager) GetObjects(x, y int) []interfaces.EbiObject {
	if cacheGetObjects != nil {
		return cacheGetObjects
	}
	return g.manager.currentScene.GetObjects(x, y)
}

// GetObject ...
func (g *EventManager) GetObject(x, y int) interfaces.EbiObject {
	objs := g.GetObjects(x, y)
	if objs != nil && len(objs) > 0 {
		return objs[0]
	}
	return nil
}

// GetEventTarget ...
func (g *EventManager) GetEventTarget(x, y int, et enum.EventTypeEnum) (interfaces.EventOwner, bool) {
	objs := g.GetObjects(x, y)
	// log.Printf("Game::GetEventTarget: %#v", objs)
	if objs != nil && len(objs) > 0 {
		for i := range objs {
			obj := objs[i]
			if t, ok := obj.(interfaces.EventOwner); ok {
				// log.Printf("  t: %#v", t)
				if t.EventHandler() != nil && t.EventHandler().Has(et) {
					switch et {
					case enum.EventTypeScroll:
						if o, ok := t.(interfaces.Scrollable); ok {
							if o.GetEdgeType(x, y) != enum.EdgeTypeNotEdge {
								return o.(interfaces.EventOwner), true
							}
						}
					default:
						if o, ok := t.(interfaces.EbiObject); ok {
							if o.In(x, y) {
								return o.(interfaces.EventOwner), true
							}
						}
					}
				}
			}
		}
	}
	return nil, false
}

// GetEventTargetList ...
func (g *EventManager) GetEventTargetList(x, y int, types ...enum.EventTypeEnum) []interfaces.EventOwner {
	targets := []interfaces.EventOwner{}
	objs := g.GetObjects(ebiten.CursorPosition())
	for i := range objs {
		obj := objs[i]
		if t, ok := obj.(interfaces.EventOwner); ok {
			// log.Printf("t: %#v", t)
			for j := range types {
				et := types[j]
				if t.EventHandler().Has(et) {
					targets = append(targets, t)
				}
			}
		}
	}
	return targets
}

// Update ...
func (g *EventManager) Update() error {
	// --- キャッシュクリア ---
	cacheGetObjects = nil
	// --- キャッシュクリア ---

	x, y := ebiten.CursorPosition()
	cursorpos := ebitest.NewPoint(float64(x), float64(y))

	evparams := make(map[string]interface{})

	// カーソル処理
	{
		// ホバーイベント
		if hoverd, ok := g.GetEventTarget(x, y, enum.EventTypeFocus); ok {
			newHoverdObject := hoverd.(interfaces.EbiObject)
			if newHoverdObject == hoverdObject {
				// log.Printf("same target")
			} else {
				// フォーカス対象が変わった
				if t, ok := hoverdObject.(interfaces.EventOwner); ok {
					// 前のフォーカスを外す処理
					t.EventHandler().Firing(enum.EventTypeFocus, t, cursorpos, evparams)
				}
				// 新しいフォーカス処理
				hoverd.EventHandler().Firing(enum.EventTypeFocus, hoverd, cursorpos, evparams)
				hoverdObject = newHoverdObject
			}
		} else {
			// フォーカス対象なし
			if t, ok := hoverdObject.(interfaces.EventOwner); ok {
				t.EventHandler().Firing(enum.EventTypeFocus, t, cursorpos, evparams)
			}
			hoverdObject = nil
		}

		// タップイベント
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// マウスタップ
			g.setStroke(x, y)
		}

		// タップ状態からの状態遷移イベント処理（クリック、D&D、ロングタップ）
		if g.stroke != nil {
			g.stroke.Update()

			dx, dy := g.stroke.PositionDiff()
			evparams["dx"], evparams["dy"] = dx, dy

			// eventCompleted := false
			if target, ok := g.stroke.Target(); ok {
				switch g.stroke.CurrentEvent() {
				case enum.EventTypeClick:
					target.EventHandler().Firing(enum.EventTypeClick, target, cursorpos, evparams)
					// eventCompleted = true
				case enum.EventTypeDragging:
					target.EventHandler().Firing(enum.EventTypeDragging, target, cursorpos, evparams)
				case enum.EventTypeDragDrop:
					target.EventHandler().Firing(enum.EventTypeDragDrop, target, cursorpos, evparams)
					// eventCompleted = true
				case enum.EventTypeLongPress:
					target.EventHandler().Firing(enum.EventTypeLongPress, target, cursorpos, evparams)
				case enum.EventTypeLongPressReleased:
					target.EventHandler().Firing(enum.EventTypeLongPressReleased, target, cursorpos, evparams)
					// eventCompleted = true
				}
				// tname := fmt.Sprintf("%s", reflect.TypeOf(target))
				// log.Printf("EventType(%d): target: %s", g.stroke.CurrentEvent(), tname)
			}

			if g.stroke.IsReleased() {
				log.Printf("EventCompleted")
				g.stroke = nil
			}
		}

		// 他のイベントが発生していない場合
		if hoverdObject == nil && g.stroke == nil {
			// スクロールイベント
			if target, ok := g.GetEventTarget(x, y, enum.EventTypeScroll); ok {
				target.EventHandler().Firing(enum.EventTypeScroll, target, cursorpos, evparams)
				// tname := fmt.Sprintf("%s", reflect.TypeOf(target))
				// log.Printf("EventType(Wheel): target: %s", tname)
			}
		}
	}

	// ホイール処理
	{
		xoff, yoff := ebiten.Wheel()
		if xoff != 0 || yoff != 0 {
			if target, ok := g.GetEventTarget(x, y, enum.EventTypeWheel); ok {
				evparams["xoff"], evparams["yoff"] = xoff, yoff
				target.EventHandler().Firing(enum.EventTypeWheel, target, cursorpos, evparams)

				// tname := fmt.Sprintf("%s", reflect.TypeOf(target))
				// log.Printf("EventType(Wheel): target: %s", tname)
			}
		}
	}

	return nil
}
