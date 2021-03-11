package game

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/myanagisawa/ebitest/app/g"
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
		data         map[enum.DataTypeEnum]interfaces.DataSet
	}
)

// NewManager ...
func NewManager(screenWidth, screenHeight int) *Manager {
	g.Width, g.Height = screenWidth, screenHeight

	gm := &Manager{
		background: ebiten.NewImageFromImage(utils.CreateRectImage(screenWidth, screenHeight, &color.RGBA{0, 0, 0, 255})),
		scenes:     make(map[enum.SceneEnum]interfaces.Scene),
		data:       make(map[enum.DataTypeEnum]interfaces.DataSet),
	}

	eventManager = &EventManager{
		manager: gm,
	}

	return gm
}

// AddDataSet ...
func (o *Manager) AddDataSet(key enum.DataTypeEnum, set interfaces.DataSet) {
	o.data[key] = set
}

// DataSet ...
func (o *Manager) DataSet(key enum.DataTypeEnum) interfaces.DataSet {
	return o.data[key]
}

// SetScene ...
func (o *Manager) SetScene(key enum.SceneEnum, scene interfaces.Scene) {
	o.scenes[key] = scene
	scene.ExecDidLoad()
}

// TransitionTo ...
func (o *Manager) TransitionTo(t enum.SceneEnum) {
	s := o.scenes[t]
	o.currentScene = s
	s.ExecDidActive()
	log.Printf("TransitionTo: %#v", o.currentScene)
}

// Update ...
func (o *Manager) Update() error {
	// log.Printf("game.Manager.Update")
	// イベント処理
	eventManager.Update()

	if o.currentScene != nil {
		return o.currentScene.Update()
	}
	return nil
}

// Draw ...
func (o *Manager) Draw(screen *ebiten.Image) {
	g.DebugText = ""
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(o.background, op)

	if o.currentScene != nil {
		o.currentScene.Draw(screen)

		// test
		// log.Printf("---- test")
		count := 0
		scene := o.currentScene
		for _, frame := range scene.Frames() {
			// frame描画
			count++
			frame.Draw(screen)

			for _, layer := range frame.Layers() {
				// layer描画
				count++
				layer.Draw(screen)

				for _, control := range layer.UIControls() {

					// control描画
					count++
					control.Draw(screen)

					// 子要素描画
					for _, child := range control.Children() {
						count++
						child.Draw(screen)
					}
				}
			}
		}
		// log.Printf("---- test: %d items", count)

		// test
	}

	// x, y := ebiten.CursorPosition()
	// dbg := fmt.Sprintf("%s\nTPS: %0.2f\nFPS: %0.2f\npos: (%d, %d)", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS(), x, y)
	dbg := fmt.Sprintf("%s / TPS: %0.2f / FPS: %0.2f", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS())
	if g.DebugText != "" {
		dbg += fmt.Sprintf("%s", g.DebugText)
	}
	focused := eventManager.GetObject(ebiten.CursorPosition())
	if focused != nil {
		dbg += fmt.Sprintf(" / cursor target: %s", focused.Label())
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
	// return fmt.Sprintf("Alloc, Sys, GC: %dMB, %dMB, %d", toMb(ms.Alloc), toMb(ms.Sys), ms.NumGC)
}

func toKb(bytes uint64) uint64 {
	return bytes / 1024
}

//noinspection GoUnusedFunction
func toMb(bytes uint64) uint64 {
	return toKb(bytes) / 1024
}
