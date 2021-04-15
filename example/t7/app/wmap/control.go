package wmap

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/app/m"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/control"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

var (
	routeInfo interfaces.UIDialog
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

type scroller struct {
	dx    float64
	dy    float64
	count int
	dest  int
	debug string
}

func (o *scroller) Update() bool {
	o.count++
	return true
}

// GetCurrentPoint ...
func (o *scroller) GetCurrentPoint() *g.Point {
	return g.NewPoint(o.dx*float64(o.count), o.dy*float64(o.count))
}

// Completed ...
func (o *scroller) Completed() bool {
	return o.count == o.dest
}

// newScroller ...
func newScroller(from, to *g.Point, d time.Duration) *scroller {
	tps := ebiten.CurrentTPS()
	dest := int(d.Seconds() * tps)
	if dest == 0 {
		dest = 1 // 0のままだと動かないので、最低1
	}
	// 移動量を計算
	moveX, moveY := to.X()-from.X(), to.Y()-from.Y()

	o := &scroller{
		dx:    moveX / float64(dest),
		dy:    moveY / float64(dest),
		dest:  dest,
		count: 0,
	}
	return o
}

type site struct {
	control.UIControl
	source *m.Site
}

func createSite(s interfaces.Scene, source *m.Site) *site {
	size := source.Image.Bounds().Size()
	_bound := g.NewBoundByPosSize(g.DefPoint(), g.NewSize(size.X, size.Y))
	con := control.NewUIControl(s, nil, enum.ControlTypeDefault, source.Code, _bound, g.DefScale(), g.DefCS(), source.Image).(*control.UIControl)
	o := &site{
		UIControl: *con,
		source:    source,
	}
	o.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
			for _, c := range ev.GetChildren() {
				c.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
			}
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
			for _, c := range ev.GetChildren() {
				c.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
			}
		}
	})
	o.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.UIControl, params map[string]interface{}) {
		log.Printf("callback::click!! %T", ev)
		if t, ok := ev.(*site); ok {
			log.Printf("hoge: %v", t)
			// il.ShowSiteInfo(t.obj)
		}
	})

	// 名前ラベルオブジェクト作成
	label := fmt.Sprintf("%s.name", source.Code)
	textSize := g.NewSize(source.Text.Bounds().Size().X, source.Text.Bounds().Size().Y)
	textBound := g.NewBoundByPosSize(g.NewPoint(-float64(textSize.W()-size.X)/2, float64(size.Y)), textSize)
	name := control.NewUIControl(s, nil, enum.ControlTypeDefault, label, textBound, g.DefScale(), g.DefCS(), source.Text).(*control.UIControl)
	o.AppendChild(name)

	return o
}

// GetControls ...
func (o *site) GetControls() []interfaces.UIControl {
	ret := o.UIControl.GetControls()
	ret[0] = o
	return ret
}

func (o *site) updatePosition(layer interfaces.UIControl) {
	// if o.Parent() == nil {
	// 	return
	// }
	// parent := o.Parent()
	// _, layerSize := parent.Bound().ToPosSize()
	// layerPos := parent.Position(enum.TypeGlobal)

	_, layerSize := layer.Bound().ToPosSize()
	layerPos := layer.Position(enum.TypeGlobal)
	x := layerPos.X() + float64(layerSize.W())*o.source.Location.X()
	y := layerPos.Y() + float64(layerSize.H())*o.source.Location.Y()

	newPos := g.NewPoint(x, y)
	o.Bound().SetDelta(newPos, nil)
}

type route struct {
	control.UIControl
	source *m.Route
}

func createRoute(s interfaces.Scene, source *m.Route) *route {

	getSiteByMaster := func(_source *m.Site) *site {
		for _, s := range sites {
			if _source.Code == s.source.Code {
				return s
			}
		}
		return nil
	}

	site1 := getSiteByMaster(source.Site1)
	site2 := getSiteByMaster(source.Site2)
	if site1 == nil || site2 == nil {
		return nil
	}

	pos1, size1 := site1.Bound().ToPosSize()
	pos2, size2 := site2.Bound().ToPosSize()

	// 2点の座標
	x1, y1 := pos1.X()+float64(size1.W())/2, pos1.Y()+float64(size1.H())/2
	x2, y2 := pos2.X()+float64(size2.W())/2, pos2.Y()+float64(size2.H())/2

	// 距離を算出
	distance := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	routeSize := source.Image.Bounds().Size()

	// 角度を算出
	rad := math.Atan2(y2-y1, x2-x1)

	// scale := g.NewScale(distance/float64(routeSize.X), 1.0)
	// _bound := g.NewBoundByPosSize(g.NewPoint(x1, y1), g.NewSize(routeSize.X, routeSize.Y))
	_bound := g.NewBoundByPosSize(g.NewPoint(x1, y1), g.NewSize(int(distance), routeSize.Y))
	con := control.NewUIControl(s, nil, enum.ControlTypeDefault, source.Code, _bound, g.DefScale(), g.DefCS(), source.Image).(*control.UIControl)
	con.SetAngle(g.NewAngle(rad))
	o := &route{
		UIControl: *con,
		source:    source,
	}
	o.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
			for _, c := range ev.GetChildren() {
				c.ColorScale().Set(0.75, 0.75, 0.75, 1.0)
			}
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
			for _, c := range ev.GetChildren() {
				c.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
			}
		}
	})
	o.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.UIControl, params map[string]interface{}) {
		log.Printf("callback::click!! %T", ev)
		if t, ok := ev.(*site); ok {
			log.Printf("hoge: %v", t)
			// il.ShowSiteInfo(t.obj)
		}
	})

	// o.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
	// 	log.Printf("callback::focus")
	// 	if t, ok := ev.(interfaces.Focusable); ok {
	// 		t.ToggleHover()
	// 	}
	// 	if t, ok := ev.(*route); ok {
	// 		name := t.Label()
	// 		lpos := t.Layer().Position(enum.TypeGlobal)
	// 		// log.Printf("---------------------")
	// 		// log.Printf("pos: x=%0.1f, y=%0.1f", pos.X(), pos.Y())
	// 		// log.Printf("lpos: x=%0.1f, y=%0.1f", lpos.X(), lpos.Y())
	// 		// log.Printf("tpos: x=%0.1f, y=%0.1f", t.Position(enum.TypeLocal).X(), t.Position(enum.TypeLocal).Y())
	// 		infopos := g.NewPoint(pos.X()-lpos.X()+5, pos.Y()-lpos.Y()-25)
	// 		routeInfo.SetPosition(infopos.Get())

	// 		if et, ok := params["event_type"]; ok {
	// 			switch et.(enum.EventTypeEnum) {
	// 			case enum.EventTypeFocus:
	// 				text := t.obj.Text
	// 				ts := text.Bounds().Size()
	// 				ds := routeInfo.Size(enum.TypeScaled)
	// 				p := g.NewPoint(float64((ds.W()-ts.X)/2), float64((ds.H()-ts.Y)/2))
	// 				label := control.NewUIControl(t.Layer(), "", text, p, g.DefScale(), false)
	// 				routeInfo.SetItems([]interfaces.UIControl{label})
	// 				routeInfo.SetVisible(true)
	// 			case enum.EventTypeBlur:
	// 				routeInfo.SetItems([]interfaces.UIControl{})
	// 				routeInfo.SetVisible(false)
	// 			}
	// 			log.Printf("name: %s, event_type: %t", name, et)
	// 		}
	// 	}
	// })
	// parent.AddUIControl(o)
	return o
}

// GetControls ...
func (o *route) GetControls() []interfaces.UIControl {
	ret := o.UIControl.GetControls()
	ret[0] = o
	return ret
}

// NewSceneFrame ...
func NewSceneFrame(s interfaces.Scene) interfaces.UIControl {
	img := utils.CreateRectImage(1, 1, &color.RGBA{255, 255, 255, 255})

	bound := g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(g.Width, g.Height))
	f := control.NewUIControl(s, nil, enum.ControlTypeFrame, "worldmap-frame", bound, g.DefScale(), g.DefCS(), img)

	return f
}

// NewWorldMap ...
func NewWorldMapLayer(s interfaces.Scene) interfaces.UIControl {
	bound := g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(3120, 2340))
	l := control.NewUIControl(s, nil, enum.ControlTypeLayer, "worldmap-layer", bound, g.DefScale(), g.DefCS(), g.Images["world"])

	l.EventHandler().AddEventListener(enum.EventTypeScroll, func(ev interfaces.UIControl, params map[string]interface{}) {
		ev.Scroll(ev.Parent().GetEdgeType())
		// log.Printf("callback::scroll:: %d", ev.Parent().GetEdgeType())
	})

	return l
}

// NewInfoLayer ...
func NewInfoLayer(s interfaces.Scene) interfaces.UIControl {

	bound := g.NewBoundByPosSize(g.NewPoint(10, 50), g.NewSize(500, 900))
	// 閉じるボタン付きウィンドウ生成
	window := control.NewDefaultClosableWindow(s, bound)

	return window
}

// NewScrollView ...
func NewScrollView(s interfaces.Scene) interfaces.UIScrollView {

	bound := g.NewBoundByPosSize(g.NewPoint(0, 20), g.NewSize(500, 880))
	// スクロールビュー生成
	sv := control.NewDefaultScrollView(s, bound)

	return sv
}
