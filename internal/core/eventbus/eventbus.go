package eventbus

type EventBus struct {
	subscribers map[string][]func(interface{})
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(interface{})),
	}
}

func (eb *EventBus) Subscribe(event string, callback func(interface{})) {
	eb.subscribers[event] = append(eb.subscribers[event], callback)
}

func (eb *EventBus) Publish(event string, data interface{}) {
	if handlers, found := eb.subscribers[event]; found {
		for _, handler := range handlers {
			handler(data)
		}
	}
}
