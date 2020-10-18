package event

import (
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
)

// Handler ...
type Handler struct {
	events map[string]map[*Event]struct{}
}

// AddEventListener ...
func (eh *Handler) AddEventListener(c interfaces.UIControl, name string, callback func(interfaces.UIControl, interfaces.Scene, *ebitest.Point)) {
	ev := &Event{c, callback}
	eh.Set(name, ev)
}

// Firing イベントの発火を行います
func (eh *Handler) Firing(s interfaces.Scene, name string, x, y int) {
	for e := range eh.events[name] {
		if e.target.In(x, y) {
			p := ebitest.NewPoint(float64(x), float64(y))
			e.callback(e.target, s, p)
			break
		}
	}
}

// Set ...
func (eh *Handler) Set(name string, ev interfaces.Event) {
	if eh.events[name] != nil {
		eh.events[name][ev.(*Event)] = struct{}{}
	} else {
		m := map[*Event]struct{}{ev.(*Event): {}}
		eh.events[name] = m
	}
}

// NewEventHandler ...
func NewEventHandler() *Handler {
	eh := &Handler{
		events: map[string]map[*Event]struct{}{},
	}
	return eh
}

// Event ...
type Event struct {
	target   interfaces.UIControl
	callback func(interfaces.UIControl, interfaces.Scene, *ebitest.Point)
}

// // Source ...
// type Source struct {
// 	scene interfaces.Scene
// 	x     int
// 	y     int
// }
