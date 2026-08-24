package events

import "sync"

type Handler func(Event)

type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{handlers: map[string][]Handler{}}
}

func (b *InMemoryBus) Subscribe(eventType string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], h)
}

func (b *InMemoryBus) Publish(e Event) {
	b.mu.RLock()
	h := append([]Handler{}, b.handlers[e.EventType()]...)
	b.mu.RUnlock()
	for _, handler := range h {
		handler(e)
	}
}
