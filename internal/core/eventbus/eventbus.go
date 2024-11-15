package eventbus

import "sync"

type EventBus struct {
	subscribers map[string][]func(interface{})
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(interface{})),
	}
}

func (eb *EventBus) Subscribe(event string, callback func(interface{})) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[event] = append(eb.subscribers[event], callback)
}

func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if handlers, found := eb.subscribers[event]; found {
		for _, handler := range handlers {
			go handler(data)
		}
	}
}
