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
	return o
}

func (o *route) draw(screen *ebiten.Image) {
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
	if ax < 0 {
		// 対象を中央に配置すると左領域が空く場合
		ax = 0
	} else if (ax + float64(wsize.W())) > float64(mapsize.W()) {
		// 対象を中央に配置すると右領域が空く場合
		ax = float64(mapsize.W()) - float64(wsize.W())
	}

	ay := pos.Y() - (float64(wsize.H()) / 2)
	if ay < 0 {
		// 対象を中央に配置すると上領域が空く場合
		ay = 0
	} else if (ay + float64(wsize.H())) > float64(mapsize.H()) {
		// 対象を中央に配置すると下領域が空く場合
		ax = float64(mapsize.H()) - float64(wsize.H())
	}

	// 現在の表示位置
	//before := o.Position(enum.TypeLocal)

	// 表示位置変更
	o.SetPosition(-ax, -ay)
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
