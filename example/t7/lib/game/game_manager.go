package game

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

/*
	- マスタからゲームデータを生成する生成処理: OK
	- レイヤサイズを画面サイズに合わせるオプション
	- ゲームデータの画面上への描画
*/
var (
	ms          runtime.MemStats
	withoutDraw bool

	dbg string
)

type (
	// Manager ...
	Manager struct {
		background   *ebiten.Image
		currentScene interfaces.Scene
		scenes       map[enum.SceneEnum]interfaces.Scene
	}
)

// NewManager ...
func NewManager(screenWidth, screenHeight int) *Manager {
	g.Width, g.Height = screenWidth, screenHeight

	gm := &Manager{
		background: ebiten.NewImageFromImage(utils.CreateRectImage(screenWidth, screenHeight, &color.RGBA{0, 0, 0, 255})),
		scenes:     make(map[enum.SceneEnum]interfaces.Scene),
	}
	withoutDraw = false

	return gm
}

// SetScene ...
func (o *Manager) SetScene(key enum.SceneEnum, scene interfaces.Scene) {
	o.scenes[key] = scene
	scene.DidLoad()
}

// TransitionTo ...
func (o *Manager) TransitionTo(t enum.SceneEnum) {
	s := o.scenes[t]
	o.currentScene = s
	s.DidActive()
	log.Printf("TransitionTo: %#v", o.currentScene)
}

// Update ...
func (o *Manager) Update() error {

	if o.currentScene != nil {
		// i := 0
		for _, child := range o.currentScene.GetControls() {
			child.Update()
			// i++
		}
		// log.Printf("-- update: %d controls", i)
	}

	dbg = fmt.Sprintf("%s / TPS: %0.2f / FPS: %0.2f", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS())
	// log.Printf("%s", dbg)

	return nil
}

// Draw ...
func (o *Manager) Draw(screen *ebiten.Image) {
	if withoutDraw {
		return
	}

	if o.currentScene != nil {
		// i := 0
		for _, child := range o.currentScene.GetControls() {
			child.Draw(screen)
			// i++
		}
		// log.Printf("-- draw: %d controls", i)
	}

	ebitenutil.DebugPrint(screen, dbg)
}

// Layout ...
func (o *Manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
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
	return fmt.Sprintf("Alloc, Obj, Sys, GC: %dMB, %d, %dMB, %d", toMb(ms.Alloc), toKb(ms.HeapObjects), toMb(ms.Sys), ms.NumGC)
}

func toKb(bytes uint64) uint64 {
	return bytes / 1024
}

//noinspection GoUnusedFunction
func toMb(bytes uint64) uint64 {
	return toKb(bytes) / 1024
}
