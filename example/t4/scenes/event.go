package scenes

// EventHandler ...
type EventHandler struct {
	events map[string]map[*Event]struct{}
}

// Firing イベントの発火を行います
func (e *EventHandler) Firing(s Scene, name string, x, y int) {
	for e := range e.events[name] {
		if e.target.In(x, y) {
			e.callback(e.target, &EventSource{
				scene: s,
				x:     x,
				y:     y,
			})
			break
		}
	}
}

// Set ...
func (e *EventHandler) Set(name string, ev *Event) {
	if e.events[name] != nil {
		e.events[name][ev] = struct{}{}
	} else {
		m := map[*Event]struct{}{ev: {}}
		e.events[name] = m
	}
}

// Event ...
type Event struct {
	target   UIControl
	callback func(UIControl, *EventSource)
}

// EventSource ...
type EventSource struct {
	scene Scene
	x     int
	y     int
}
