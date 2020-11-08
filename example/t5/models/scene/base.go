package scene

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
)

var (
	ms runtime.MemStats
)

// Base ...
type Base struct {
	label       string
	layers      []interfaces.Layer
	activeLayer interfaces.Layer
}

// Label ...
func (s *Base) Label() string {
	return s.label
}

// Draw ...
func (s *Base) Draw(screen *ebiten.Image) {
	return
}

// LayerAt ...
func (s *Base) LayerAt(x, y int) interfaces.Layer {
	for i := len(s.layers) - 1; i >= 0; i-- {
		l := s.layers[i]
		if l.IsModal() {
			return l
		}
		if l.In(x, y) {
			return l
		}
	}

	return nil
}

// ActiveLayer ...
func (s *Base) ActiveLayer() interfaces.Layer {
	return s.activeLayer
}

// GetLayerByLabel ...
func (s *Base) GetLayerByLabel(label string) interfaces.Layer {
	for _, layer := range s.layers {
		log.Printf("GetLayerByLabel: %s == %s, %v", layer.Label(), label, layer.Label() == label)
		if layer.Label() == label {
			return layer
		}
	}
	return nil
}

// SetLayer ...
func (s *Base) SetLayer(l interfaces.Layer) {
	s.layers = append(s.layers, l)
}

// DeleteLayer ...
func (s *Base) DeleteLayer(l interfaces.Layer) {
	var layers []interfaces.Layer
	for _, layer := range s.layers {
		if l != layer {
			layers = append(layers, layer)
		}
	}
	s.layers = layers
}

// Layout ...
func (s *Base) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// GetEdgeType ...
func GetEdgeType(x, y int) enum.EdgeTypeEnum {
	minX, maxX := ebitest.EdgeSize, ebitest.Width-ebitest.EdgeSize
	minY, maxY := ebitest.EdgeSize, ebitest.Height-ebitest.EdgeSize

	// 範囲外判定
	if x < -ebitest.EdgeSizeOuter || x > ebitest.Width+ebitest.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	} else if y < -ebitest.EdgeSizeOuter || y > ebitest.Height+ebitest.EdgeSizeOuter {
		return enum.EdgeTypeNotEdge
	}

	// 判定
	if x <= minX && y <= minY {
		return enum.EdgeTypeTopLeft
	} else if x > minX && x < maxX && y <= minY {
		return enum.EdgeTypeTop
	} else if x >= maxX && y <= minY {
		return enum.EdgeTypeTopRight
	} else if x >= maxX && y > minY && y < maxY {
		return enum.EdgeTypeRight
	} else if x >= maxX && y >= maxY {
		return enum.EdgeTypeBottomRight
	} else if x > minX && x < maxX && y >= maxY {
		return enum.EdgeTypeBottom
	} else if x <= minX && y >= maxY {
		return enum.EdgeTypeBottomLeft
	} else if x <= minX && y > minY && y < maxY {
		return enum.EdgeTypeLeft
	}
	return enum.EdgeTypeNotEdge
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
