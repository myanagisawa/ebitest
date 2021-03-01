package main

import (
	"testing"
)

// func BenchmarkAppend_AllocateEveryTime(b *testing.B) {
// 	base := []string{}
// 	b.ResetTimer()
// 	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
// 	for i := 0; i < b.N; i++ {
// 		// 都度append
// 		base = append(base, fmt.Sprintf("no%d", i))
// 	}
// }

// func BenchmarkAppend_AllocateOnce(b *testing.B) {
// 	//最初に長さを決める
// 	base := make([]string, b.N)
// 	b.ResetTimer()
// 	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
// 	for i := 0; i < b.N; i++ {
// 		base[i] = fmt.Sprintf("no%d", i)
// 	}
// }
func BenchmarkExecStructC(b *testing.B) {
	o := NewStructC(3)
	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		o.GetHoge()
	}
}

func BenchmarkExecStructB(b *testing.B) {
	o := NewStructB(2)
	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		o.GetHoge()
	}
}

func BenchmarkExecStructA2(b *testing.B) {
	o := NewStructA(1)
	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		GetHoge(o)
	}
}

func BenchmarkExecStructA(b *testing.B) {
	o := NewStructA(1)
	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		o.GetHoge()
	}
}
