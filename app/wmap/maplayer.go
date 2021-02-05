package wmap

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/app/obj"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/functions"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/control"
	"github.com/myanagisawa/ebitest/models/layer"
)

var (
	scrollProg *scroller

	routeInfo interfaces.UIControl
)

type scroller struct {
	from       *g.Point
	to         *g.Point
	current    *g.Point
	startedAt  *time.Time
	finishedAt *time.Time
	dest       int
	count      int
	debug      string
}

func (o *scroller) Update() bool {
	// log.Printf("o.count: %d / %d", o.count, o.dest)
	if o.dest == 0 {
		return false
	}
	o.count++
	dx := (o.to.X() - o.from.X()) / float64(o.dest)
	dy := (o.to.Y() - o.from.Y()) / float64(o.dest)

	o.current = g.NewPoint(dx*float64(o.count)+o.from.X(), dy*float64(o.count)+o.from.Y())
	o.debug += fmt.Sprintf("scrolle.Update: %d/%d: (%0.2f, %0.2f)\n", o.count, o.dest, o.current.X(), o.current.Y())
	return true
}

// GetCurrentPoint ...
func (o *scroller) GetCurrentPoint() *g.Point {
	return o.current
}

// Completed ...
func (o *scroller) Completed() bool {
	return o.count == o.dest
}

// newScroller ...
func newScroller(from, to *g.Point, d time.Duration) *scroller {
	st := time.Now()
	fn := st.Add(d)
	tps := ebiten.CurrentTPS()
	o := &scroller{
		from:       from,
		to:         to,
		startedAt:  &st,
		finishedAt: &fn,
		dest:       int(d.Seconds() * tps),
		count:      0,
	}
	// log.Printf("scroller: from=%#v, to=%#v, startedAt=%v, finishedAt=%v, tps=%v, dest=%v", o.from, o.to, o.startedAt, o.finishedAt, tps, o.dest)
	return o
}

type site struct {
	control.Base
	parent *MapLayer
	obj    *obj.Site
}

func createSite(obj *obj.Site, parent *MapLayer) *site {
	pos := parent.getLocation(obj.Location)
	base := control.NewControlBase(obj.Code, obj.Image, pos).(*control.Base)
	base.SetLayer(parent)
	o := &site{
		Base:   *base,
		parent: parent,
		obj:    obj,
	}
	o.EventHandler().AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)
	o.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		log.Printf("callback::click")
	})

	return o
}

func (o *site) draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op = o.Base.DrawWithOptions(screen, op)

	// site名を描画
	size := g.NewSize(o.obj.Image.Size())
	textSize := g.NewSize(o.obj.Text.Size())
	op.GeoM.Translate(-float64(textSize.W()-size.W())/2, float64(size.H()))
	screen.DrawImage(o.obj.Text, op)
}

type route struct {
	control.Base
	parent *MapLayer
	obj    *obj.Route
}

func createRoute(obj *obj.Route, parent *MapLayer) *route {
	pos1 := parent.getLocation(obj.Site1.Location)
	pos2 := parent.getLocation(obj.Site2.Location)

	size1 := g.NewSize(obj.Site1.Image.Size())
	size2 := g.NewSize(obj.Site2.Image.Size())

	// 2点の座標
	x1, y1 := pos1.X()+float64(size1.W())/2, pos1.Y()+float64(size1.H())/2
	x2, y2 := pos2.X()+float64(size2.W())/2, pos2.Y()+float64(size2.H())/2

	// 距離を算出
	distance := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	routeSize := g.NewSize(obj.Image.Size())

	// 角度を算出
	rad := math.Atan2(y2-y1, x2-x1)

	base := control.NewControlBase(obj.Code, obj.Image, g.NewPoint(x1, y1)).(*control.Base)
	base.SetLayer(parent)
	base.SetScale(distance/float64(routeSize.W()), 1.0)
	base.SetAngle(rad)
	o := &route{
		Base:   *base,
		parent: parent,
		obj:    obj,
	}
	o.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		log.Printf("callback::focus")
		if t, ok := ev.(interfaces.Focusable); ok {
			t.ToggleHover()
		}
		if t, ok := ev.(*route); ok {
			name := t.Label()
			routeInfo.SetPosition(t.Position(enum.TypeLocal).Get())
			log.Printf("name: %s", name)
		}
	})
	return o
}

func (o *route) In(x, y int) bool {
	// オブジェクトの位置を取得する
	pos := o.Base.Position(enum.TypeGlobal)

	// 無回転当たり判定オブジェクトを取得
	bsize := o.Base.Size(enum.TypeScaled)
	bsize = g.NewSize(bsize.W(), bsize.H()*3)

	// あたり判定objの原点を取得(だいたい中心を原点としているが描画時の回転軸に合わせてlefttopを原点にする)
	// height側を実線より太くしている関係で、Y側の原点は
	rectCenter := g.NewPoint(pos.X(), pos.Y()+float64(bsize.H())/2)

	// ポインタ座標を当たり判定objとの相対座標に変換
	relativeX := float64(x) - rectCenter.X()
	relativeY := float64(y) - rectCenter.Y()

	// ポインタ座標を座標変換する（回転を打ち消す）
	rad := o.Angle(enum.TypeLocal)
	transformPos := g.NewPoint(
		math.Cos(rad)*relativeX+math.Sin(rad)*relativeY,
		-math.Sin(rad)*relativeX+math.Cos(rad)*relativeY)

	// 当たり判定objと変換したポインタ座標で当たり判定を行う
	objWidth := float64(bsize.W())
	objHeight := float64(bsize.H())
	if 0 <= transformPos.X() && objWidth >= transformPos.X() &&
		-objHeight/2 <= transformPos.Y() && objHeight/2 >= transformPos.Y() {
		// log.Printf("%s にフォーカス！", o.Label())
		return true
	}

	// // オブジェクトの位置を取得する
	// pos := o.Base.Position(enum.TypeGlobal)

	// // オブジェクトの見た目サイズを取得する
	// bsize := o.Base.Size(enum.TypeScaled)
	// bsize = g.NewSize(bsize.W(), bsize.H()*2)
	// // log.Printf("%s サイズ: %#v", o.Label(), bsize)

	// // 見た目サイズの原点を取得する
	// rectCenter := g.NewPoint(pos.X()+float64(bsize.W())/2, pos.Y()+float64(bsize.H())/2)
	// // rectCenter := pos

	// // ポインタ座標を矩形との相対座標に変換する
	// relativeX := float64(x) - rectCenter.X()
	// relativeY := float64(y) - rectCenter.Y()

	// // ポインタ座標を座標変換する
	// /*
	// 	// 相対座標に対して矩形の回転を打ち消す逆行列を掛ける
	// 	Vec2 transform_pos = Vec2(
	// 		cosf(rad) * relative_position.X + sinf(rad) * relative_position.Y,
	// 		-sinf(rad) * relative_position.X + cosf(rad) * relative_position.Y
	// 	);
	// */
	// rad := o.Angle(enum.TypeLocal)
	// transformPos := g.NewPoint(
	// 	math.Cos(rad)*relativeX+math.Sin(rad)*relativeY,
	// 	-math.Sin(rad)*relativeX+math.Cos(rad)*relativeY)

	// // transformPos = g.NewPoint(transformPos.X()-float64(bsize.W())/2, transformPos.Y()-float64(bsize.H())/2)

	// // ポインタ座標と矩形の当たり判定を行う
	// /*
	// 	// 矩形と点の当たり判定を行う
	// 	if (-rect.Width / 2.0f <= transform_pos.X && rect.Width / 2.0f >= transform_pos.X &&
	// 		-rect.Height / 2.0f <= transform_pos.Y && rect.Height / 2.0f >= transform_pos.Y)
	// 	{
	// */
	// objWidth := float64(bsize.W())
	// objHeight := float64(bsize.H())
	// if o.Label() == "route-1" {
	// 	log.Printf("%s: ポインタ: x=%d, y=%d", o.Label(), x, y)
	// 	log.Printf("  オブジェクト位置: x=%0.1f, y=%0.1f", pos.X(), pos.Y())
	// 	log.Printf("  オブジェクトのサイズ: w=%d, h=%d", bsize.W(), bsize.H())
	// 	log.Printf("  ポインタの相対座標: x=%0.1f, y=%0.1f", relativeX, relativeY)
	// 	log.Printf("  角度: rad=%0.1f (%0.1f)", rad, rad*(180/math.Pi))
	// 	log.Printf("  座標変換: x=%0.1f, y=%0.1f", transformPos.X(), transformPos.Y())
	// 	log.Printf("    判定: left=%0.1f, right=%0.1f, top=%0.1f, bottom=%0.1f  -> %v", -objWidth/2, objWidth/2, -objHeight/2, objHeight/2, -objWidth/2 <= transformPos.X() && objWidth/2 >= transformPos.X() &&
	// 		-objHeight/2 <= transformPos.Y() && objHeight/2 >= transformPos.Y())
	// }
	// if -objWidth/2 <= transformPos.X() && objWidth/2 >= transformPos.X() &&
	// 	-objHeight/2 <= transformPos.Y() && objHeight/2 >= transformPos.Y() {
	// 	log.Printf("%s にフォーカス！", o.Label())
	// 	return true
	// }

	return false
}

// GetObjects ...
func (o *route) GetObjects(x, y int) []interfaces.EbiObject {
	if o.In(x, y) {
		return []interfaces.EbiObject{o}
	}
	return nil
}

func (o *route) draw(screen *ebiten.Image) {
	// o.Base.Draw(screen)

	op := &ebiten.DrawImageOptions{}
	_ = o.Base.DrawWithOptions(screen, op)

	// site名を描画
	// op.GeoM.Scale(1.0, 1.0)
	// size := g.NewSize(o.obj.Image.Size())
	// textSize := g.NewSize(o.obj.Text.Size())
	// op.GeoM.Translate(-float64(textSize.W()-size.W())/2, float64(size.H()))
	// screen.DrawImage(o.obj.Text, op)
}

// MapLayer ...
type MapLayer struct {
	layer.Base
	sites  []site
	routes []route
}

// NewMapLayer ...
func NewMapLayer() *MapLayer {
	scrollProg = nil

	l := layer.NewLayerBaseByImage("map", g.Images["world"], g.NewPoint(0, 0), false).(*layer.Base)
	ml := &MapLayer{
		Base:   *l,
		sites:  []site{},
		routes: []route{},
	}

	ifsites, ok := gm.DataSet(enum.DataTypeSite).(*obj.Sites)
	if ok {
		objsites := *ifsites
		sites := make([]site, len(objsites))
		for i := range objsites {
			r := objsites[i]
			site := createSite(&r, ml)
			sites[i] = *site
		}
		ml.sites = sites
	}

	ifroutes, ok := gm.DataSet(enum.DataTypeRoute).(*obj.Routes)
	if ok {
		objroutes := *ifroutes
		routes := make([]route, len(objroutes))
		for i := range objroutes {
			r := objroutes[i]
			route := createRoute(&r, ml)
			routes[i] = *route
		}
		ml.routes = routes
	}

	routeInfo = control.NewSimpleLabel("route info", g.NewPoint(0, 0), 14, &color.RGBA{0, 0, 0, 255}, enum.FontStyleGenShinGothicRegular)
	ml.AddUIControl(routeInfo)

	return ml
}

// GetObjects ...
func (o *MapLayer) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	for i := len(o.sites) - 1; i >= 0; i-- {
		c := &o.sites[i]
		objs = append(objs, c.GetObjects(x, y)...)
	}
	for i := len(o.routes) - 1; i >= 0; i-- {
		c := &o.routes[i]
		objs = append(objs, c.GetObjects(x, y)...)
	}

	objs = append(objs, o.Base.GetObjects(x, y)...)
	// log.Printf("LayerBase::GetObjects: %#v", objs)
	return objs
}

func (o *MapLayer) getLocation(loc *g.Point) *g.Point {
	layerSize := o.Size()
	layerPos := o.Position(enum.TypeGlobal)
	x := layerPos.X() + float64(layerSize.W())*loc.X()
	y := layerPos.Y() + float64(layerSize.H())*loc.Y()

	return g.NewPoint(x, y)
}

// getSiteByCode ...
func (o *MapLayer) getSiteByCode(code string) *site {
	if o.sites == nil || len(o.sites) == 0 {
		return nil
	}
	for i := range o.sites {
		site := o.sites[i]
		if site.obj.Code == code {
			return &site
		}
	}
	return nil
}

// MoveTo ...
func (o *MapLayer) MoveTo(code string) {
	site := o.getSiteByCode(code)
	if site == nil {
		return
	}
	// MAP表示フレーム
	frame := o.Frame()
	// 表示領域サイズ
	wsize := frame.Size()

	// 表示対象
	pos := site.Position(enum.TypeLocal)
	// Mapサイズ
	mapsize := o.Size()

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
	// o.SetPosition(-ax, -ay)

	// 現在の表示位置
	before := o.Position(enum.TypeLocal)
	scrollProg = newScroller(before, g.NewPoint(-ax, -ay), 300*time.Millisecond)
}

// Update ...
func (o *MapLayer) Update() error {

	if scrollProg != nil {
		if scrollProg.Update() {
			o.SetPosition(scrollProg.GetCurrentPoint().Get())
		}
		if scrollProg.Completed() {
			// log.Printf("%s", scrollProg.debug)
			scrollProg = nil
		}
	} else {
		o.Base.Update()
	}

	return nil
}

// Draw ...
func (o *MapLayer) Draw(screen *ebiten.Image) {
	// 共通Draw処理
	o.Base.Draw(screen)

	// log.Printf("---------------- route draw [%d] --------------", len(o.routes))
	for _, route := range o.routes {
		route.draw(screen)
	}

	// log.Printf("---------------- site draw [%d] --------------", len(o.sites))
	for _, site := range o.sites {
		site.draw(screen)
	}

}
