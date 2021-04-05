package wmap

import (
	"image/color"
	"log"
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
	// o.debug += fmt.Sprintf("scroller.Update: %d/%d: (%0.2f, %0.2f)\n", o.count, o.dest, o.current.X(), o.current.Y())
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
	log.Printf("scroller: from=%#v, to=%#v, startedAt=%v, finishedAt=%v, tps=%v, dest=%v", o.from, o.to, o.startedAt, o.finishedAt, tps, o.dest)
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
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	o.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.UIControl, params map[string]interface{}) {
		log.Printf("callback::click!! %T", ev)
		if t, ok := ev.(*site); ok {
			log.Printf("hoge: %v", t)
			// il.ShowSiteInfo(t.obj)
		}
	})
	return o
}

// GetControls ...
func (o *site) GetControls() []interfaces.UIControl {
	ret := o.UIControl.GetControls()
	ret[0] = o
	return ret
}

func (o *site) updatePosition() {
	if o.Parent() == nil {
		return
	}
	parent := o.Parent()
	_, layerSize := parent.Bound().ToPosSize()
	layerPos := parent.Position(enum.TypeGlobal)
	x := layerPos.X() + float64(layerSize.W())*o.source.Location.X()
	y := layerPos.Y() + float64(layerSize.H())*o.source.Location.Y()

	newPos := g.NewPoint(x, y)
	o.Bound().SetDelta(newPos, nil)
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
