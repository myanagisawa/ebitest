package scene

import (
	"image/color"
	"log"

	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/control"
	"github.com/myanagisawa/ebitest/models/frame"
	"github.com/myanagisawa/ebitest/models/layer"
)

type (
	// Map ...
	Map struct {
		Base
	}
)

var (
	ctrl interfaces.UIControl

	cnt int
)

// NewMap ...
func NewMap(m interfaces.GameManager) *Map {

	s := &Map{
		Base: Base{
			label:   "MainMap",
			manager: m,
		},
	}

	// メインフレーム
	f := frame.NewFrame("main frame", ebitest.NewPoint(200, 20), ebitest.NewSize(ebitest.Width-200, ebitest.Height-220), &color.RGBA{200, 200, 200, 255}, true)
	s.AddFrame(f)

	l := layer.NewLayerBaseByImage("map", ebitest.Images["world"], ebitest.NewPoint(0, 0), false)
	f.AddLayer(l)

	c := control.NewSimpleLabel("test", ebitest.Images["btnBase"], ebitest.NewPoint(100, 100), color.Black)
	c.EventHandler().AddEventListener(enum.EventTypeClick, func(o interfaces.EventOwner, pos *ebitest.Point, params map[string]interface{}) {
		log.Printf("callback::click")
	})
	l.AddUIControl(c)

	ctrl = c

	cnt = 0

	// サブフレーム1（上）
	f = frame.NewFrame("top frame", ebitest.NewPoint(0, 0), ebitest.NewSize(ebitest.Width, 20), &color.RGBA{0, 0, 0, 255}, false)
	s.AddFrame(f)

	// サブフレーム2（横）
	fs := ebitest.NewSize(200, ebitest.Height-220)
	f = frame.NewFrame("side frame", ebitest.NewPoint(0, 20), fs, &color.RGBA{127, 127, 200, 255}, false)
	s.AddFrame(f)

	l = layer.NewLayerBase("Layer1", ebitest.NewPoint(0, 0), fs, &color.RGBA{0, 0, 0, 128}, false)
	f.AddLayer(l)

	// スクロールビュー
	scrollView := control.NewUIScrollView("scrollView1", ebitest.NewPoint(0, 10.0), ebitest.NewSize(f.Size().W(), f.Size().H()/2))
	l.AddUIControl(scrollView)

	cols := []interface{}{
		"ID", "カラム名", "3番目",
	}
	data := [][]interface{}{
		{1, "１行目", "あいうえお"},
		{2, "２行目", "かきくけこ"},
		{3, "３行目", "さしすせそ"},
		{4, "４行目", "たちつてと"},
		{5, "５行目", "なにぬねの"},
		{6, "６行目", "はひふへほ"},
		{7, "７行目", "まみむめも"},
		{8, "８行目", "やいゆえよ"},
		{9, "９行目", "らりるれろ"},
		{10, "１０行目", "わをん"},
		{11, "１１行目", "あいうえお"},
		{12, "１２行目", "かきくけこ"},
		{13, "１３行目", "さしすせそ"},
		{14, "１４行目", "たちつてと"},
		{15, "１５行目", "なにぬねの"},
		{16, "１６行目", "はひふへほ"},
		{17, "１７行目", "まみむめも"},
		{18, "１８行目", "やいゆえよ"},
		{19, "１９行目", "らりるれろ"},
		{20, "２０行目", "わをん"},
		{1, "１行目", "あいうえお"},
		{2, "２行目", "かきくけこ"},
		{3, "３行目", "さしすせそ"},
		{4, "４行目", "たちつてと"},
		{5, "５行目", "なにぬねの"},
		{6, "６行目", "はひふへほ"},
		{7, "７行目", "まみむめも"},
		{8, "８行目", "やいゆえよ"},
		{9, "９行目", "らりるれろ"},
		{10, "１０行目", "わをん"},
		{11, "１１行目", "あいうえお"},
		{12, "１２行目", "かきくけこ"},
		{13, "１３行目", "さしすせそ"},
		{14, "１４行目", "たちつてと"},
		{15, "１５行目", "なにぬねの"},
		{16, "１６行目", "はひふへほ"},
		{17, "１７行目", "まみむめも"},
		{18, "１８行目", "やいゆえよ"},
		{19, "１９行目", "らりるれろ"},
		{20, "２０行目", "わをん"},
	}
	scrollView.SetDataSource(cols, data)

	// サブフレーム3（下）
	f = frame.NewFrame("bottom frame", ebitest.NewPoint(0, float64(ebitest.Height-200)), ebitest.NewSize(ebitest.Width, 200), &color.RGBA{127, 200, 127, 255}, false)
	s.AddFrame(f)

	return s
}

// Update ...
func (o *Map) Update() error {

	if cnt == 10 {
		a := ctrl.Angle(enum.TypeGlobal)
		a++
		ctrl.SetAngle(a)
		cnt = 0
	}
	cnt++

	_ = o.Base.Update()
	return nil
}
