package event

import (
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
)

// Handler ...
type Handler struct {
	events map[enum.EventTypeEnum]interfaces.Event
}

// AddEventListener ...
func (o *Handler) AddEventListener(t enum.EventTypeEnum, callback func(ev interfaces.UIControl, params map[string]interface{})) {
	ev := &Event{callback}
	o.events[t] = ev
}

// Firing イベントの発火を行います
func (o *Handler) Firing(t enum.EventTypeEnum, ev interfaces.UIControl, params map[string]interface{}) {
	if e, ok := o.events[t].(*Event); ok {
		params["event"] = t
		e.callback(ev, params)
	}
}

// Has 指定のイベント種別の保持状態を返します
func (o *Handler) Has(t enum.EventTypeEnum) bool {
	_, ok := o.events[t]
	return ok
}

// NewEventHandler ...
func NewEventHandler() interfaces.EventHandler {
	eh := &Handler{
		events: map[enum.EventTypeEnum]interfaces.Event{},
	}
	return eh
}

// Event ...
type Event struct {
	callback func(ev interfaces.UIControl, params map[string]interface{})
}
