package wmap

import (
	"image/color"
	"log"

	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/control"
	"github.com/myanagisawa/ebitest/models/frame"
	"github.com/myanagisawa/ebitest/models/layer"
	"github.com/myanagisawa/ebitest/models/scene"
	"github.com/myanagisawa/ebitest/utils"
)

type (
	// Scene ...
	Scene struct {
		scene.Base
	}

	// CustomScrollView ...
	CustomScrollView struct {
		control.UIScrollView
	}
)

var (
	ctrl interfaces.UIControl

	cnt int
)

// NewCustomScrollView ...
func NewCustomScrollView(label string, pos *g.Point, size *g.Size) interfaces.UIScrollView {
	base := control.NewUIScrollView(label, pos, size).(*control.UIScrollView)
	return &CustomScrollView{
		UIScrollView: *base,
	}
}

// NewScene ...
func NewScene(m interfaces.GameManager) *Scene {

	s := &Scene{
		Base: *scene.NewScene("MainMap", m).(*scene.Base),
	}

	// メインフレーム
	mainf := frame.NewFrame("main frame", g.NewPoint(200, 20), g.NewSize(g.Width-200, g.Height-220), &color.RGBA{200, 200, 200, 255}, true)
	s.AddFrame(mainf)

	l := layer.NewLayerBaseByImage("map", g.Images["world"], g.NewPoint(0, 0), false)
	mainf.AddLayer(l)

	// c := control.NewSimpleLabel("test", g.Images["btnBase"], g.NewPoint(100, 100), color.Black)
	c := control.NewSimpleLabel("SIMPLE LABEL", g.NewPoint(100, 100), 48, &color.RGBA{0, 0, 255, 255}, enum.FontStyleGenShinGothicMedium)
	c.EventHandler().AddEventListener(enum.EventTypeClick, func(o interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		log.Printf("callback::click")
	})
	l.AddUIControl(c)
	c = control.NewSimpleLabel("シンプルラベル", g.NewPoint(100, 150), 32, &color.RGBA{0, 255, 0, 255}, enum.FontStyleGenShinGothicNormal)
	l.AddUIControl(c)
	c = control.NewSimpleLabel("文字列試験", g.NewPoint(100, 200), 24, &color.RGBA{255, 0, 0, 255}, enum.FontStyleGenShinGothicRegular)
	l.AddUIControl(c)

	c = control.NewSimpleButton("SIMPLE BUTTON", utils.CopyImage(g.Images["btnBase"]), g.NewPoint(100, 350), 16, &color.RGBA{0, 0, 255, 255})
	c.EventHandler().AddEventListener(enum.EventTypeClick, func(o interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		log.Printf("callback::click")
	})
	l.AddUIControl(c)

	ctrl = c

	cnt = 0

	// サブフレーム1（上）
	topf := frame.NewFrame("top frame", g.NewPoint(0, 0), g.NewSize(g.Width, 20), &color.RGBA{0, 0, 0, 255}, false)
	s.AddFrame(topf)

	// サブフレーム2（横）
	fs := g.NewSize(200, g.Height-220)
	sidef := frame.NewFrame("side frame", g.NewPoint(0, 20), fs, &color.RGBA{127, 127, 200, 255}, false)
	s.AddFrame(sidef)

	l = layer.NewLayerBase("Layer1", g.NewPoint(0, 0), fs, &color.RGBA{0, 0, 0, 128}, false)
	sidef.AddLayer(l)

	// スクロールビュー
	// scrollView := control.NewUIScrollView("scrollView1", g.NewPoint(0, 10.0), g.NewSize(f.Size().W(), f.Size().H()/2))
	scrollView := NewCustomScrollView("scrollView1", g.NewPoint(0, 10.0), g.NewSize(sidef.Size().W(), sidef.Size().H()/2))
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

	scrollView.SetRowClickFunc(func(row interface{}, pos *g.Point, params map[string]interface{}) {
		log.Printf("かすたむヘッダクリックイベントだよ")
	}, func(idx int, row interface{}, pos *g.Point, params map[string]interface{}) {
		log.Printf("かすたむ行クリックイベントだよ(%d)", idx)
		if r, ok := row.(interfaces.ListRow); ok {
			log.Printf("  %#v", r.Source())
			log.Printf("  別フレーム参照: %s", mainf.Label())
		}
	})

	// サブフレーム3（下）
	bottomf := frame.NewFrame("bottom frame", g.NewPoint(0, float64(g.Height-200)), g.NewSize(g.Width, 200), &color.RGBA{127, 200, 127, 255}, false)
	s.AddFrame(bottomf)

	return s
}

// Update ...
func (o *Scene) Update() error {

	// if cnt == 10 {
	// 	a := ctrl.Angle(enum.TypeGlobal)
	// 	a++
	// 	ctrl.SetAngle(a)
	// 	cnt = 0
	// }
	// cnt++

	_ = o.Base.Update()
	return nil
}
