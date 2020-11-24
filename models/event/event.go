package event

import (
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Handler ...
type Handler struct {
	events map[string]map[interfaces.Event]struct{}
}

// AddEventListener ...
func (o *Handler) AddEventListener(c interfaces.UIControl, name string, callback func(interfaces.UIControl, *ebitest.Point)) {
	ev := &Event{c, callback}
	o.Set(name, ev)
}

// Firing イベントの発火を行います
func (o *Handler) Firing(s interfaces.Scene, name string, x, y int) {
	for e := range o.events[name] {
		event := e.(*Event)
		if event.target.In(x, y) {
			p := ebitest.NewPoint(float64(x), float64(y))
			event.callback(event.target, p)
			break
		}
	}
}

// Set ...
func (o *Handler) Set(name string, ev interfaces.Event) {
	if o.events[name] != nil {
		o.events[name][ev] = struct{}{}
	} else {
		m := map[interfaces.Event]struct{}{ev: {}}
		o.events[name] = m
	}
}

// NewEventHandler ...
func NewEventHandler() interfaces.EventHandler {
	eh := &Handler{
		events: map[string]map[interfaces.Event]struct{}{},
	}
	return eh
}

// Event ...
type Event struct {
	target   interfaces.UIControl
	callback func(interfaces.UIControl, *ebitest.Point)
}
