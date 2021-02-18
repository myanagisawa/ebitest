package wmap

import (
	"fmt"
	"image/color"
	"log"

	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/app/obj"
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
	gm          interfaces.GameManager
	maplayer    *MapLayer
	il          *infoLayer
	scrollView1 interfaces.UIScrollView
)

// NewCustomScrollView ...
func NewCustomScrollView(l interfaces.Layer, label string, pos *g.Point, size *g.Size) interfaces.UIScrollView {
	base := control.NewUIScrollView(l, label, pos, size).(*control.UIScrollView)
	o := &CustomScrollView{
		UIScrollView: *base,
	}
	l.AddUIControl(o)
	return o
}

// NewScene ...
func NewScene(m interfaces.GameManager) *Scene {
	gm = m
	s := &Scene{
		Base: *scene.NewScene("MainMap", m).(*scene.Base),
	}
	s.SetCustomFunc(enum.FuncTypeDidLoad, s.didLoad())
	s.SetCustomFunc(enum.FuncTypeDidActive, s.didActive())

	return s
}

// didLoad ...
func (o *Scene) didLoad() func() {
	return func() {
		// メインフレーム
		mainf := frame.NewFrame(o, "main frame", g.NewPoint(300, 20), g.NewSize(g.Width-300, g.Height-220), &color.RGBA{200, 200, 200, 255}, true)

		// マップレイヤ
		maplayer = NewMapLayer(mainf)

		// c := control.NewSimpleLabel("test", g.Images["btnBase"], g.NewPoint(100, 100), color.Black)
		c := control.NewSimpleLabel(maplayer, "SIMPLE LABEL", g.NewPoint(100, 100), 48, &color.RGBA{0, 0, 255, 255}, enum.FontStyleGenShinGothicMedium)
		c.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
			log.Printf("callback::click")
		})

		c = control.NewSimpleLabel(maplayer, "シンプルラベル", g.NewPoint(100, 150), 32, &color.RGBA{0, 255, 0, 255}, enum.FontStyleGenShinGothicNormal)
		c = control.NewSimpleLabel(maplayer, "文字列試験", g.NewPoint(100, 200), 24, &color.RGBA{255, 0, 0, 255}, enum.FontStyleGenShinGothicRegular)

		c = control.NewSimpleButton(maplayer, "SIMPLE BUTTON", utils.CopyImage(g.Images["btnBase"]), g.NewPoint(100, 350), 16, &color.RGBA{0, 0, 255, 255})
		c.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
			log.Printf("callback::click")
			o.Manager().TransitionTo(enum.MenuEnum)
		})

		// 情報表示レイヤ
		w, h := mainf.Size().Get()
		il = newInfoLayer(mainf, w, h)

		// サブフレーム1（上）
		_ = frame.NewFrame(o, "top frame", g.NewPoint(0, 0), g.NewSize(g.Width, 20), &color.RGBA{0, 0, 0, 255}, false)

		// サブフレーム2（横）
		fs := g.NewSize(300, g.Height-220)
		sidef := frame.NewFrame(o, "side frame", g.NewPoint(0, 20), fs, &color.RGBA{127, 127, 200, 255}, false)

		l := layer.NewLayerBase(sidef, "Layer1", g.NewPoint(0, 0), fs, &color.RGBA{0, 0, 0, 128}, false)

		// スクロールビュー
		scrollView1 = NewCustomScrollView(l, "scrollView1", g.NewPoint(0, 10.0), g.NewSize(sidef.Size().W(), sidef.Size().H()/2))

		// サブフレーム3（下）
		_ = frame.NewFrame(o, "bottom frame", g.NewPoint(0, float64(g.Height-200)), g.NewSize(g.Width, 200), &color.RGBA{127, 200, 127, 255}, false)

		log.Printf("MainMap.DidLoad")
	}
}

// didActive ...
func (o *Scene) didActive() func() {
	return func() {

		if sites, ok := gm.DataSet(enum.DataTypeSite).(*obj.Sites); ok {
			cols := []interface{}{
				"ID", "種類", "名前", "位置",
			}
			data := make([][]interface{}, len(*sites))

			for i, site := range *sites {
				row := make([]interface{}, 4)
				row[0] = i
				row[1] = site.Type.Name()
				row[2] = site.Name
				x, y := site.Location.Get()
				row[3] = fmt.Sprintf("x=%0.2f, y=%0.2f", x, y)
				// row[4] = site.Code
				data[i] = row
			}
			scrollView1.SetDataSource(cols, data)

			scrollView1.SetRowClickFunc(func(row interface{}, pos *g.Point, params map[string]interface{}) {
				log.Printf("かすたむヘッダクリックイベントだよ")
			}, func(idx int, row interface{}, pos *g.Point, params map[string]interface{}) {
				log.Printf("かすたむ行クリックイベントだよ(%d)", idx)
				if r, ok := row.(interfaces.ListRow); ok {
					log.Printf("  %#v", r.Source())
					// log.Printf("  別フレーム参照: %s", mainf.Label())
					name := r.Source()[2]
					if sites, ok := gm.DataSet(enum.DataTypeSite).(*obj.Sites); ok {
						site := sites.GetByName(name)
						if s, ok := site.(*obj.Site); ok {
							maplayer.MoveTo(s.Code)
							log.Printf("  MoveTo: %s", s.Name)
						}
					}
				}
			})

			// 初期表示
			maplayer.MoveTo("site-1")
		}
		// cols := []interface{}{
		// 	"ID", "カラム名", "3番目",
		// }
		// data := [][]interface{}{
		// 	{1, "１行目", "あいうえお"},
		// 	{2, "２行目", "かきくけこ"},
		// 	{3, "３行目", "さしすせそ"},
		// 	{4, "４行目", "たちつてと"},
		// 	{5, "５行目", "なにぬねの"},
		// 	{6, "６行目", "はひふへほ"},
		// 	{7, "７行目", "まみむめも"},
		// 	{8, "８行目", "やいゆえよ"},
		// 	{9, "９行目", "らりるれろ"},
		// 	{10, "１０行目", "わをん"},
		// 	{11, "１１行目", "あいうえお"},
		// 	{12, "１２行目", "かきくけこ"},
		// 	{13, "１３行目", "さしすせそ"},
		// 	{14, "１４行目", "たちつてと"},
		// 	{15, "１５行目", "なにぬねの"},
		// 	{16, "１６行目", "はひふへほ"},
		// 	{17, "１７行目", "まみむめも"},
		// 	{18, "１８行目", "やいゆえよ"},
		// 	{19, "１９行目", "らりるれろ"},
		// 	{20, "２０行目", "わをん"},
		// 	{1, "１行目", "あいうえお"},
		// 	{2, "２行目", "かきくけこ"},
		// 	{3, "３行目", "さしすせそ"},
		// 	{4, "４行目", "たちつてと"},
		// 	{5, "５行目", "なにぬねの"},
		// 	{6, "６行目", "はひふへほ"},
		// 	{7, "７行目", "まみむめも"},
		// 	{8, "８行目", "やいゆえよ"},
		// 	{9, "９行目", "らりるれろ"},
		// 	{10, "１０行目", "わをん"},
		// 	{11, "１１行目", "あいうえお"},
		// 	{12, "１２行目", "かきくけこ"},
		// 	{13, "１３行目", "さしすせそ"},
		// 	{14, "１４行目", "たちつてと"},
		// 	{15, "１５行目", "なにぬねの"},
		// 	{16, "１６行目", "はひふへほ"},
		// 	{17, "１７行目", "まみむめも"},
		// 	{18, "１８行目", "やいゆえよ"},
		// 	{19, "１９行目", "らりるれろ"},
		// 	{20, "２０行目", "わをん"},
		// }
		// scrollView1.SetDataSource(cols, data)

		log.Printf("MainMap.DidActive")
	}
}

// Update ...
func (o *Scene) Update() error {

	_ = o.Base.Update()
	return nil
}
