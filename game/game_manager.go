package game

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/myanagisawa/ebitest/app/appscene"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/utils"
)

/*
	- マスタからゲームデータを生成する生成処理: OK
	- レイヤサイズを画面サイズに合わせるオプション
	- ゲームデータの画面上への描画
*/
var (
	ms runtime.MemStats

	eventManager *EventManager
)

type (
	// Manager ...
	Manager struct {
		background   *ebiten.Image
		currentScene interfaces.Scene
		scenes       map[enum.SceneEnum]interfaces.Scene
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
		background: ebiten.NewImageFromImage(utils.CreateRectImage(screenWidth, screenHeight, &color.RGBA{0, 0, 0, 255})),
		// master: NewMasterData(),
	}

	scenes := map[enum.SceneEnum]interfaces.Scene{}
	scenes[enum.MapEnum] = appscene.NewMap(gm)
	gm.scenes = scenes

	// MainMenuを表示
	gm.TransitionTo(enum.MapEnum)

	eventManager = &EventManager{
		manager: gm,
	}

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

// Update ...
func (g *Manager) Update() error {
	// イベント処理
	eventManager.Update()

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
	focused := eventManager.GetObject(ebiten.CursorPosition())
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
