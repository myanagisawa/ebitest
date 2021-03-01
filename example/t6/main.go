package main

import (
	"log"
	"sort"
)

type StructA struct {
	hoge int
}

func NewStructA(hoge int) *StructA {
	o := &StructA{hoge: hoge}
	return o
}
func (o *StructA) GetHoge() {
	o.hoge++
}

type StructB struct {
	hoge     int
	hogeFunc func(self *StructB)
}

func NewStructB(hoge int) *StructB {
	o := &StructB{hoge: hoge}
	o.hogeFunc = func(self *StructB) {
		self.hoge++
	}
	return o
}

func (o *StructB) GetHoge() {
	o.hogeFunc(o)
}

type StructC struct {
	hoge     int
	hogeFunc func(self interface{})
}

func NewStructC(hoge int) *StructC {
	o := &StructC{hoge: hoge}
	o.hogeFunc = func(self interface{}) {
		switch t := self.(type) {
		case *StructC:
			t.hoge++
		}
	}
	return o
}

func (o *StructC) GetHoge() {
	o.hogeFunc(o)
}

func GetHoge(o *StructA) {
	o.hoge++
}

func main() {
	t := []string{"a", "d", "c", "e"}
	sort.Slice(t, func(i int, j int) bool {
		return i > j
	})
	for _, r := range t {
		log.Printf("%s", r)
	}
}
