package wmap

import (
	"fmt"
	"log"
	"time"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/app/m"
	"github.com/myanagisawa/ebitest/example/t7/lib/game"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/control"
)

type (
	// Scene ...
	Scene struct {
		children []interfaces.UIControl
		manager  *game.Manager
	}
)

var (
	frame interfaces.UIControl
	layer interfaces.UIControl

	window   interfaces.UIControl
	testList *control.UIScrollView

	sites  []*site
	routes []*route

	scrollProg *scroller

	routeWindow interfaces.UIControl
)

// NewScene ...
func NewScene(manager *game.Manager) *Scene {
	s := &Scene{
		children: []interfaces.UIControl{},
		manager:  manager,
	}

	return s
}

// Label ...
func (o *Scene) Label() string {
	return "menu"
}

// TransitionTo ...
func (o *Scene) TransitionTo(t enum.SceneEnum) {
	o.manager.TransitionTo(t)
}

// GetControls ...
func (o *Scene) GetControls() []interfaces.UIControl {
	ret := []interfaces.UIControl{}
	for _, child := range o.children {
		ret = append(ret, child.GetControls()...)
	}
	return ret
}

// Update ...
func (o *Scene) Update() error {

	if scrollProg != nil {
		if scrollProg.Update() {
			// log.Printf("%s", scrollProg.debug)
			layer.SetMoving(scrollProg.GetCurrentPoint())
		}
		if scrollProg.Completed() {
			log.Printf("scroll.completed: \n%s", scrollProg.debug)
			layer.FinishMoving()
			scrollProg = nil
		}
	}
	return nil
}

// DidLoad ...
func (o *Scene) DidLoad() {
	log.Printf("map.DidLoad")

	// シーンフレーム
	frame = NewSceneFrame(o)
	o.children = append(o.children, frame)

	// マップレイヤー
	layer = NewWorldMapLayer(o)
	// 親子関係を設定
	frame.AppendChild(layer)

	// サイト情報作成
	if gamedata, ok := o.manager.GameData(enum.DataTypeSite); ok {
		if ifsites, ok := gamedata.(*m.Sites); ok {

			objsites := *ifsites
			_sites := make([]*site, len(objsites))
			for i := range objsites {
				r := objsites[i]
				site := createSite(o, &r)
				site.updatePosition(layer)
				_sites[i] = site
			}
			sites = _sites
		}
	}

	// 経路情報作成
	if gamedata, ok := o.manager.GameData(enum.DataTypeRoute); ok {
		if ifroutes, ok := gamedata.(*m.Routes); ok {

			objroutes := *ifroutes
			_routes := make([]*route, len(objroutes))
			for i := range objroutes {
				r := objroutes[i]
				route := createRoute(o, &r)
				_routes[i] = route
				// layer.AppendChild(route)
			}
			routes = _routes
		}
	}

	// 重なり順でlayerに追加
	if routes != nil {
		for _, row := range routes {
			layer.AppendChild(row)
		}
	}
	if sites != nil {
		for _, row := range sites {
			layer.AppendChild(row)
		}
	}

	window = NewInfoLayer(o)
	frame.AppendChild(window)

	testList = NewScrollView(o).(*control.UIScrollView)
	window.AppendChild(testList)
}

// DidActive ...
func (o *Scene) DidActive() {
	log.Printf("map.DidActive")

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
	// }

	if sites != nil {
		// ヘッダを設定
		cols := []interface{}{
			"CODE", "種類", "名前", "位置",
		}
		testList.SetHeaderRow(cols)

		data := make([][]interface{}, len(sites))

		for i, site := range sites {
			row := make([]interface{}, 4)
			row[0] = site.source.Code
			row[1] = site.source.Type.Name()
			row[2] = site.source.Name
			x, y := site.source.Location.Get()
			row[3] = fmt.Sprintf("x=%0.2f, y=%0.2f", x, y)
			// row[4] = site.Code
			data[i] = row
		}
		// 行を設定
		testList.AppendRows(data, func(row *control.ScrollViewRow) {
			code := fmt.Sprintf("%s", row.GetDS()[0])
			moveTo(code)
		})

		// scrollView1.SetRowClickFunc(func(row interface{}, pos *g.Point, params map[string]interface{}) {
		// 	log.Printf("かすたむヘッダクリックイベントだよ")
		// }, func(idx int, row interface{}, pos *g.Point, params map[string]interface{}) {
		// 	log.Printf("かすたむ行クリックイベントだよ(%d)", idx)
		// 	if r, ok := row.(interfaces.ListRow); ok {
		// 		log.Printf("  %#v", r.Source())
		// 		// log.Printf("  別フレーム参照: %s", mainf.Label())
		// 		name := r.Source()[2]
		// 		if sites, ok := gm.DataSet(enum.DataTypeSite).(*obj.Sites); ok {
		// 			site := sites.GetByName(name)
		// 			if s, ok := site.(*obj.Site); ok {
		// 				maplayer.MoveTo(s.Code)
		// 				log.Printf("  MoveTo: %s", s.Name)
		// 			}
		// 		}
		// 	}
		// })

		// 初期表示
		moveTo("site-1")
	}

}

// moveTo ...
func moveTo(code string) {
	site := getSiteByCode(code)
	if site == nil {
		return
	}
	// 表示領域サイズ
	_, wsize := frame.Bound().ToPosSize()

	// 表示対象
	pos := site.Position(enum.TypeLocal)
	// Mapサイズ
	_, mapsize := layer.Bound().ToPosSize()

	// 表示対象を中央に表示した場合のposを計算
	ax := pos.X() - (float64(wsize.W()) / 2)
	// log.Printf("ax: %0.2f, ax + float64(wsize.W())=%0.2f, mapsize.W()=%0.2f", ax, (ax + float64(wsize.W())), float64(mapsize.W()))
	if ax < 0 {
		// 対象を中央に配置すると左領域が空く場合
		ax = 0
	} else if (ax + float64(wsize.W())) > float64(mapsize.W()) {
		// 対象を中央に配置すると右領域が空く場合
		ax = float64(mapsize.W()) - float64(wsize.W())
	}

	ay := pos.Y() - (float64(wsize.H()) / 2)
	// log.Printf("ay: %0.2f, ay + float64(wsize.H())=%0.2f, mapsize.H()=%0.2f", ay, (ay + float64(wsize.H())), float64(mapsize.H()))
	if ay < 0 {
		// 対象を中央に配置すると上領域が空く場合
		ay = 0
	} else if (ay + float64(wsize.H())) > float64(mapsize.H()) {
		// 対象を中央に配置すると下領域が空く場合
		ay = float64(mapsize.H()) - float64(wsize.H())
	}

	// 表示位置変更
	// layer.Bound().SetPosSize(g.NewPoint(-ax, -ay), nil)

	// 現在の表示位置
	before := layer.Position(enum.TypeLocal)
	scrollProg = newScroller(before, g.NewPoint(-ax, -ay), 300*time.Millisecond)
}

// getSiteByCode ...
func getSiteByCode(code string) *site {
	if sites == nil || len(sites) == 0 {
		return nil
	}
	for i := range sites {
		site := sites[i]
		if site.source.Code == code {
			return site
		}
	}
	return nil
}
