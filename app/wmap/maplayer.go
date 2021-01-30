package wmap

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/app/obj"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/models/control"
	"github.com/myanagisawa/ebitest/models/layer"
)

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
	return o
}

func (o *site) draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	// // 描画位置指定
	// op.GeoM.Reset()
	// op.GeoM.Scale(1.0, 1.0)

	size := g.NewSize(o.obj.Image.Size())
	// // 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	// op.GeoM.Translate(-float64(size.W())/2, -float64(size.H())/2)
	// // ユニットの座標に移動
	// op.GeoM.Translate(float64(size.W())/2, float64(size.H())/2)

	// // Site描画座標を計算
	// sitePos := o.parent.getLocation(o.obj.Location)
	// op.GeoM.Translate(sitePos.X(), sitePos.Y())

	// // log.Printf("site: %s, (%0.2f, %0.2f)", o.obj.Code, sitePos.X(), sitePos.Y())
	// screen.DrawImage(o.obj.Image, op)
	o.Base.Draw(screen)

	// site名を描画
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
	return o
}

func (o *route) draw(screen *ebiten.Image) {
	// pos1 := o.parent.getLocation(o.obj.Site1.Location)
	// pos2 := o.parent.getLocation(o.obj.Site2.Location)

	// size1 := g.NewSize(o.obj.Site1.Image.Size())
	// size2 := g.NewSize(o.obj.Site2.Image.Size())

	// // 2点の座標
	// x1, y1 := pos1.X()+float64(size1.W())/2, pos1.Y()+float64(size1.H())/2
	// x2, y2 := pos2.X()+float64(size2.W())/2, pos2.Y()+float64(size2.H())/2

	// // 距離を算出
	// distance := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))

	// // 角度を算出
	// rad := math.Atan2(y2-y1, x2-x1)

	// op := &ebiten.DrawImageOptions{}

	// routeSize := g.NewSize(o.obj.Image.Size())
	// // 描画位置指定
	// op.GeoM.Reset()
	// op.GeoM.Scale(distance/float64(routeSize.W()), 1.0)

	// op.GeoM.Rotate(rad)
	// op.GeoM.Translate(x1, y1)

	// // log.Printf("route: %s, (%0.2f, %0.2f)", o.obj.Code, x1, y1)
	// screen.DrawImage(o.obj.Image, op)

	o.Base.Draw(screen)
}

// MapLayer ...
type MapLayer struct {
	layer.Base
	sites  []site
	routes []route
}

// NewMapLayer ...
func NewMapLayer() *MapLayer {
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

	return ml
}

func (o *MapLayer) getLocation(loc *g.Point) *g.Point {
	layerSize := o.Size()
	layerPos := o.Position(enum.TypeGlobal)
	x := layerPos.X() + float64(layerSize.W())*loc.X()
	y := layerPos.Y() + float64(layerSize.H())*loc.Y()

	return g.NewPoint(x, y)
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

// func (o *MapLayer) drawSite(site *obj.Site, screen *ebiten.Image) {
// 	op := &ebiten.DrawImageOptions{}

// 	// 描画位置指定
// 	op.GeoM.Reset()
// 	op.GeoM.Scale(1.0, 1.0)

// 	size := g.NewSize(site.Image.Size())
// 	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
// 	op.GeoM.Translate(-float64(size.W())/2, -float64(size.H())/2)
// 	// ユニットの座標に移動
// 	op.GeoM.Translate(float64(size.W())/2, float64(size.H())/2)

// 	// Site描画座標を計算
// 	sitePos := o.getSitePosition(site)
// 	op.GeoM.Translate(sitePos.X(), sitePos.Y())

// 	// log.Printf("site: %s, (%0.2f, %0.2f)", site.Code, x, y)
// 	screen.DrawImage(site.Image, op)

// 	// site名を描画
// 	textSize := g.NewSize(site.Text.Size())
// 	op.GeoM.Translate(-float64(textSize.W()-size.W())/2, float64(size.H()))
// 	screen.DrawImage(site.Text, op)
// }

// func (o *MapLayer) getSitePosition(site *obj.Site) *g.Point {
// 	bgSize := o.Size()
// 	x := o.Position(enum.TypeGlobal).X() + float64(bgSize.W())*site.Location.X()
// 	y := o.Position(enum.TypeGlobal).Y() + float64(bgSize.H())*site.Location.Y()

// 	return g.NewPoint(x, y)
// }

// func (o *MapLayer) drawRoute(route *obj.Route, screen *ebiten.Image) {
// 	pos1 := o.getSitePosition(route.Site1)
// 	pos2 := o.getSitePosition(route.Site2)

// 	size1 := g.NewSize(route.Site1.Image.Size())
// 	size2 := g.NewSize(route.Site2.Image.Size())

// 	// 2点の座標
// 	x1, y1 := pos1.X()+float64(size1.W())/2, pos1.Y()+float64(size1.H())/2
// 	x2, y2 := pos2.X()+float64(size2.W())/2, pos2.Y()+float64(size2.H())/2

// 	// 距離を算出
// 	distance := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))

// 	// 角度を算出
// 	rad := math.Atan2(y2-y1, x2-x1)

// 	op := &ebiten.DrawImageOptions{}

// 	routeSize := g.NewSize(route.Image.Size())
// 	// 描画位置指定
// 	op.GeoM.Reset()
// 	op.GeoM.Scale(distance/float64(routeSize.W()), 1.0)

// 	op.GeoM.Rotate(rad)
// 	op.GeoM.Translate(x1, y1)
// 	// // log.Printf("site: %s, (%0.2f, %0.2f)", site.Code, x, y)
// 	screen.DrawImage(route.Image, op)

// }
