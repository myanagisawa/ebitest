package game

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/input"
	"github.com/myanagisawa/ebitest/models/scene"
)

/*
	- マスタからゲームデータを生成する生成処理: OK
	- レイヤサイズを画面サイズに合わせるオプション
	- ゲームデータの画面上への描画
*/
var (
	ms runtime.MemStats

	cacheGetObjects []interfaces.EbiObject

	hoverdObject interfaces.EbiObject
)

type (
	// Manager ...
	Manager struct {
		background   *ebiten.Image
		currentScene interfaces.Scene
		scenes       map[enum.SceneEnum]interfaces.Scene
		stroke       *Stroke
		// master       *MasterData
	}
)

// func getSiteByCode(code string, sites []*parts.Site) *parts.Site {
// 	for _, site := range sites {
// 		if site.Code == code {
// 			return site
// 		}
// 	}
// 	return nil
// }

// // GetSites master.MSiteからSiteを作成します
// func (g *Manager) GetSites() []*parts.Site {
// 	sites := make([]*parts.Site, len(g.master.Sites))
// 	for i, site := range g.master.Sites {
// 		sites[i] = parts.NewSite(site.Code, site.Type, site.Name, site.Location)
// 	}
// 	return sites
// }

// // GetRoutes master.MSiteからSiteを作成します
// func (g *Manager) GetRoutes(sites []*parts.Site) []*parts.Route {
// 	routes := make([]*parts.Route, len(g.master.Routes))
// 	for i, route := range g.master.Routes {
// 		routes[i] = parts.NewRoute(
// 			route.Code,
// 			route.Type,
// 			route.Name,
// 			getSiteByCode(route.Site1, sites),
// 			getSiteByCode(route.Site2, sites))
// 	}
// 	return routes
// }

// NewManager ...
func NewManager(screenWidth, screenHeight int) *Manager {
	ebitest.Width, ebitest.Height = screenWidth, screenHeight

	gm := &Manager{
		background: ebiten.NewImageFromImage(ebitest.CreateRectImage(screenWidth, screenHeight, &color.RGBA{0, 0, 0, 255})),
		// master: NewMasterData(),
	}

	scenes := map[enum.SceneEnum]interfaces.Scene{}
	scenes[enum.MapEnum] = scene.NewMap(gm)
	gm.scenes = scenes

	// MainMenuを表示
	gm.TransitionTo(enum.MapEnum)

	return gm
}

// TransitionTo ...
func (g *Manager) TransitionTo(t enum.SceneEnum) {
	// var s interfaces.Scene
	// switch t {
	// case enum.MainMenuEnum:
	// 	s = scene.NewMainMenu(g)
	// case enum.MapEnum:
	// 	s = scene.NewMap(g)
	// default:
	// 	panic(fmt.Sprintf("invalid SceneEnum: %d", t))
	// }
	g.currentScene = g.scenes[t]
}

// SetCurrentScene ...
func (g *Manager) SetCurrentScene(s interfaces.Scene) {
	g.currentScene = s
}

// setStroke ...
func (g *Manager) setStroke(x, y int) {
	if g.stroke != nil {
		log.Printf("別のstrokeがあるため追加できません")
		return
	}
	targets := g.GetEventTargetList(x, y, enum.EventTypeClick, enum.EventTypeDragging, enum.EventTypeLongPress)
	if len(targets) > 0 {
		stroke := NewStroke(&input.MouseStrokeSource{})
		stroke.SetMouseDownTargets(targets)
		g.stroke = stroke
	}
}

// GetObjects ...
func (g *Manager) GetObjects(x, y int) []interfaces.EbiObject {
	if cacheGetObjects != nil {
		return cacheGetObjects
	}
	return g.currentScene.GetObjects(x, y)
}

// GetObject ...
func (g *Manager) GetObject(x, y int) interfaces.EbiObject {
	objs := g.GetObjects(x, y)
	if objs != nil && len(objs) > 0 {
		return objs[0]
	}
	return nil
}

// GetEventTarget ...
func (g *Manager) GetEventTarget(x, y int, et enum.EventTypeEnum) (interfaces.EventOwner, bool) {
	objs := g.GetObjects(x, y)
	// log.Printf("Game::GetEventTarget: %#v", objs)
	if objs != nil && len(objs) > 0 {
		for i := range objs {
			obj := objs[i]
			if t, ok := obj.(interfaces.EventOwner); ok {
				// log.Printf("  t: %#v", t)
				if t.EventHandler() != nil && t.EventHandler().Has(et) {
					return t, true
				}
			}
		}
	}
	return nil, false
}

// GetEventTargetList ...
func (g *Manager) GetEventTargetList(x, y int, types ...enum.EventTypeEnum) []interfaces.EventOwner {
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
func (g *Manager) Update() error {
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

			eventCompleted := false
			if target, ok := g.stroke.Target(); ok {
				switch g.stroke.CurrentEvent() {
				case enum.EventTypeClick:
					target.EventHandler().Firing(enum.EventTypeClick, target, cursorpos, evparams)
					eventCompleted = true
				case enum.EventTypeDragging:
					target.EventHandler().Firing(enum.EventTypeDragging, target, cursorpos, evparams)
					// target.UpdateStroke(g.stroke)
				case enum.EventTypeDragDrop:
					target.EventHandler().Firing(enum.EventTypeDragDrop, target, cursorpos, evparams)
					// target.UpdatePositionByDelta()
					eventCompleted = true
				case enum.EventTypeLongPress:
					target.EventHandler().Firing(enum.EventTypeLongPress, target, cursorpos, evparams)
				case enum.EventTypeLongPressReleased:
					target.EventHandler().Firing(enum.EventTypeLongPressReleased, target, cursorpos, evparams)
					eventCompleted = true
				}
				tname := fmt.Sprintf("%s", reflect.TypeOf(target))
				log.Printf("EventType(%d): target: %s", g.stroke.CurrentEvent(), tname)
			}

			if eventCompleted {
				log.Printf("EventCompleted")
				g.stroke = nil
			}
		}
	}

	// ホイール処理
	{
		xoff, yoff := ebiten.Wheel()
		if xoff != 0 || yoff != 0 {
			if target, ok := g.GetEventTarget(x, y, enum.EventTypeWheel); ok {
				evparams := make(map[string]interface{})
				evparams["xoff"], evparams["yoff"] = xoff, yoff

				target.EventHandler().Firing(enum.EventTypeWheel, target, cursorpos, evparams)

				tname := fmt.Sprintf("%s", reflect.TypeOf(target))
				log.Printf("EventType(Wheel): target: %s", tname)
			}
		}
	}

	if g.currentScene != nil {
		return g.currentScene.Update()
	}
	return nil
}

// Draw ...
func (g *Manager) Draw(screen *ebiten.Image) {
	ebitest.DebugText = ""
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(g.background, op)

	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}

	// x, y := ebiten.CursorPosition()
	// dbg := fmt.Sprintf("%s\nTPS: %0.2f\nFPS: %0.2f\npos: (%d, %d)", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS(), x, y)
	dbg := fmt.Sprintf("%s / TPS: %0.2f / FPS: %0.2f", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS())
	if ebitest.DebugText != "" {
		dbg += fmt.Sprintf("%s", ebitest.DebugText)
	}
	focused := g.GetObject(ebiten.CursorPosition())
	if focused != nil {
		dbg += fmt.Sprintf(" / cursor target: %s", focused.Label())
	}
	ebitenutil.DebugPrint(screen, dbg)
}

// Layout ...
func (g *Manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebitest.Width, ebitest.Height
}

func printMemoryStats() string {
	// --------------------------------------------------------
	// runtime.MemoryStats() から、現在の割当メモリ量などが取得できる.
	//
	// まず、データの受け皿となる runtime.MemStats を初期化し
	// runtime.ReadMemStats(*runtime.MemStats) を呼び出して
	// 取得する.
	// --------------------------------------------------------
	runtime.ReadMemStats(&ms)

	// // Alloc は、現在ヒープに割り当てられているメモリ
	// // HeapAlloc と同じ.
	// output.Stdoutl("Alloc", toKb(ms.Alloc))
	// output.Stdoutl("HeapAlloc", toKb(ms.HeapAlloc))

	// // TotalAlloc は、ヒープに割り当てられたメモリ量の累積
	// // Allocと違い、こちらは増えていくが減ることはない
	// output.Stdoutl("TotalAlloc", toKb(ms.TotalAlloc))

	// // HeapObjects は、ヒープに割り当てられているオブジェクトの数
	// output.Stdoutl("HeapObjects", toKb(ms.HeapObjects))

	// // Sys は、OSから割り当てられたメモリの合計量
	// output.Stdoutl("Sys", toKb(ms.Sys))

	// // NumGC は、実施されたGCの回数
	// output.Stdoutl("NumGC", ms.NumGC)
	return fmt.Sprintf("Alloc, Sys, GC: %dMB, %dMB, %d", toMb(ms.Alloc), toMb(ms.Sys), ms.NumGC)
}

func toKb(bytes uint64) uint64 {
	return bytes / 1024
}

//noinspection GoUnusedFunction
func toMb(bytes uint64) uint64 {
	return toKb(bytes) / 1024
}
