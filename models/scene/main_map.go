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
			label: "MainMap",
		},
	}

	// メインフレーム
	f := frame.NewFrame("main frame", ebitest.NewPoint(200, 20), ebitest.NewSize(ebitest.Width-200, ebitest.Height-220), color.RGBA{200, 200, 200, 255}, true)
	s.AddFrame(f)

	l := layer.NewLayerBase("map", ebitest.Images["world"], ebitest.NewPoint(0, 0), ebitest.NewScale(1, 1), false)
	f.AddLayer(l)

	c := control.NewControlBase("test", ebitest.Images["btnBase"], ebitest.NewPoint(100, 100), color.Black)
	l.AddUIControl(c)

	ctrl = c

	cnt = 0

	// サブフレーム1（横）
	f = frame.NewFrame("side frame", ebitest.NewPoint(0, 20), ebitest.NewSize(200, ebitest.Height-220), color.RGBA{127, 127, 200, 255}, false)
	s.AddFrame(f)

	img := ebitest.CreateRectImage(150, 300, color.RGBA{0, 0, 0, 128})
	l = layer.NewLayerBase("Layer1", img, ebitest.NewPoint(25, 25), ebitest.NewScale(1, 1), true)
	f.AddLayer(l)

	// スクロールビュー
	cols := []interface{}{
		"ID", "カラム名ですColumn", "3番目", "HOGE",
	}
	data := [][]interface{}{
		{1, "１行目", "あいうえお", 12345},
		{2, "２行目", "かきくけこ", 12345},
		{3, "３行目", "さしすせそ", 12345},
		{4, "４行目", "たちつてと", 12345},
		{5, "５行目", "なにぬねの", 12345},
		{6, "６行目", "はひふへほ", 12345},
		{7, "７行目", "まみむめも", 12345},
		{8, "８行目", "やいゆえよ", 12345},
		{9, "９行目", "らりるれろ", 12345},
		{10, "１０行目", "わをん", 12345},
		{11, "１１行目", "あいうえお", 12345},
		{12, "１２行目", "かきくけこ", 12345},
		{13, "１３行目", "さしすせそ", 12345},
		{14, "１４行目", "たちつてと", 12345},
		{15, "１５行目", "なにぬねの", 12345},
		{16, "１６行目", "はひふへほ", 12345},
		{17, "１７行目", "まみむめも", 12345},
		{18, "１８行目", "やいゆえよ", 12345},
		{19, "１９行目", "らりるれろ", 12345},
		{20, "２０行目", "わをん", 12345},
	}
	c = control.NewUIScrollViewByList(l, cols, data, 700, 250, 30, ebitest.NewPoint(20.0, 20.0))
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, point *ebitest.Point) {
		// log.Printf("%s clicked", target.Label())
		t := target.(*control.UIScrollViewImpl)
		index := t.GetIndexOfFocusedRow()
		log.Printf("clicked: index: %d, data: %#v", index, data[index])
	})

	// サブフレーム2（下）
	f = frame.NewFrame("bottom frame", ebitest.NewPoint(0, float64(ebitest.Height-200)), ebitest.NewSize(ebitest.Width, 200), color.RGBA{127, 200, 127, 255}, false)
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
