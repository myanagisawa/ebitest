package menu

import (
	"fmt"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/char"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/control"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

// var (
// 	img = utils.CreateRectImage(1, 1, &color.RGBA{255, 255, 255, 255})
// )

func init() {
	rand.Seed(time.Now().UnixNano())

}

// NewFrame ...
func NewFrame(s interfaces.Scene, label string, bound *g.Bound) interfaces.UIControl {
	img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})

	f := control.NewUIControl(s, nil, enum.ControlTypeFrame, label, bound, g.DefScale(), g.DefCS(), img)

	pos := g.NewPoint(0, 0)
	size := g.NewSize(100, 1200)

	f.SetChildren([]interfaces.UIControl{
		NewLayer(s, f, fmt.Sprintf("%s.layer-1", label), g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewLayer(s, f, fmt.Sprintf("%s.layer-2", label), g.NewBoundByPosSize(pos.SetDelta(100, 0), size)),
		NewLayer(s, f, fmt.Sprintf("%s.layer-3", label), g.NewBoundByPosSize(pos.SetDelta(100, 0), size)),
		NewLayer(s, f, fmt.Sprintf("%s.layer-4", label), g.NewBoundByPosSize(pos.SetDelta(100, 0), size)),
	})

	return f
}

// NewLayer ...
func NewLayer(s interfaces.Scene, parent interfaces.UIControl, label string, bound *g.Bound) interfaces.UIControl {
	img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 192})

	l := control.NewUIControl(s, parent, enum.ControlTypeLayer, label, bound, g.DefScale(), g.DefCS(), img)

	pos := g.NewPoint(0, 0)
	size := g.NewSize(100, 100)

	l.SetChildren([]interfaces.UIControl{
		NewUIControl(s, l, fmt.Sprintf("%s.control-1", label), g.NewBoundByPosSize(pos.SetDelta(0, 0), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-2", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-3", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-4", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-5", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-6", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-7", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-8", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-9", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-10", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-11", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
		NewUIControl(s, l, fmt.Sprintf("%s.control-12", label), g.NewBoundByPosSize(pos.SetDelta(0, 100), size)),
	})

	return l
}

// NewUIControl ...
func NewUIControl(s interfaces.Scene, parent interfaces.UIControl, label string, bound *g.Bound) interfaces.UIControl {
	img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 160})

	c := control.NewUIControl(s, parent, enum.ControlTypeDefault, label, bound, g.DefScale(), g.DefCS(), img)
	c.SetVector(g.NewVector(float64(rand.Intn(3)+2), g.NewAngleByDeg(rand.Intn(360))))

	size := g.NewSize(50, 50)

	c.SetChildren([]interfaces.UIControl{
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-1", label), g.NewBoundByPosSize(g.NewPoint(0, 0), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-2", label), g.NewBoundByPosSize(g.NewPoint(50, 0), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-3", label), g.NewBoundByPosSize(g.NewPoint(0, 50), size)),
		NewChildControl(s, c, fmt.Sprintf("%s.child-control-4", label), g.NewBoundByPosSize(g.NewPoint(50, 50), size)),
	})

	return c
}

// NewChildControl ...
func NewChildControl(s interfaces.Scene, parent interfaces.UIControl, label string, bound *g.Bound) interfaces.UIControl {
	img := utils.CreateRectImage(1, 1, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 127})

	c := control.NewUIControl(s, parent, enum.ControlTypeDefault, label, bound, g.DefScale(), g.DefCS(), img)
	c.SetAngle(g.NewAngleByDeg(rand.Intn(360)))

	f := func(self interfaces.UIControl) {
		l := self.Label()
		deg, _ := strconv.Atoi(l[len(l)-1:])
		deg += rand.Intn(2) - 3
		// log.Printf("updateProc: %s: %d: %T", self.label, deg, self)

		self.Angle(enum.TypeRaw).SetDelta(deg)
	}
	c.SetUpdateProc(f)

	return c
}

// NewButtonControl ...
func NewButtonControl(s interfaces.Scene, label string, bound *g.Bound) interfaces.UIControl {
	img := utils.CopyImage(g.Images["btnBase"])
	fset := char.Res.Get(16, enum.FontStyleGenShinGothicBold)
	ti := fset.GetStringImage(label)
	ti = utils.TextColorTo(ti.(draw.Image), &color.RGBA{0, 0, 0, 255})
	ti = utils.ImageOnTextToCenter(img.(draw.Image), ti)

	c := control.NewUIControl(s, nil, enum.ControlTypeDefault, label, bound, g.DefScale(), g.DefCS(), ti)
	c.EventHandler().AddEventListener(enum.EventTypeFocus, func(ev interfaces.UIControl, params map[string]interface{}) {
		et := params["event_type"].(enum.EventTypeEnum)
		switch et {
		case enum.EventTypeFocus:
			ev.ColorScale().Set(0.5, 0.5, 0.5, 1.0)
		case enum.EventTypeBlur:
			ev.ColorScale().Set(1.0, 1.0, 1.0, 1.0)
		}
	})
	c.EventHandler().AddEventListener(enum.EventTypeClick, func(ev interfaces.UIControl, params map[string]interface{}) {
		log.Printf("callback::click")
		ev.Scene().TransitionTo(enum.MapEnum)
	})

	return c
}
