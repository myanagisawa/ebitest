package event

import (
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
)

// Handler ...
type Handler struct {
	owner  interfaces.EventOwner
	events map[enum.EventTypeEnum]interfaces.Event
}

// AddEventListener ...
func (o *Handler) AddEventListener(t enum.EventTypeEnum, callback func(interfaces.EventOwner, *ebitest.Point)) {
	ev := &Event{o.owner, callback}
	o.events[t] = ev
}

// Firing イベントの発火を行います
func (o *Handler) Firing(t enum.EventTypeEnum, x, y int) {
	event := o.events[t].(*Event)
	event.callback(event.target, ebitest.NewPoint(float64(x), float64(y)))
}

// Has 指定のイベント種別の保持状態を返します
func (o *Handler) Has(t enum.EventTypeEnum) bool {
	_, ok := o.events[t]
	return ok
}

// NewEventHandler ...
func NewEventHandler(o interfaces.EventOwner) interfaces.EventHandler {
	eh := &Handler{
		owner:  o,
		events: map[enum.EventTypeEnum]interfaces.Event{},
	}
	return eh
}

// Event ...
type Event struct {
	target   interfaces.EventOwner
	callback func(interfaces.EventOwner, *ebitest.Point)
}
