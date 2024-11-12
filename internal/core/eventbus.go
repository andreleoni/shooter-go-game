package core

import "sync"

type EventType string

const (
	CollisionEvent EventType = "collision"
	DestroyEvent   EventType = "destroy"
)

type Event struct {
	Type    EventType
	Payload interface{}
}

type EventListener func(event Event)

type EventBus struct {
	listeners map[EventType][]EventListener
	mu        sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[EventType][]EventListener),
	}
}

func (eb *EventBus) Subscribe(eventType EventType, listener EventListener) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners[eventType] = append(eb.listeners[eventType], listener)
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if listeners, exists := eb.listeners[event.Type]; exists {
		for _, listener := range listeners {
			go listener(event)
		}
	}
}
