package eventbus

import (
	"context"
	"sync"
)

const defaultSubscriberBufferSize = 16

type Hub struct {
	mu          sync.RWMutex
	subscribers map[chan any]struct{}
	bufferSize  int
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[chan any]struct{}),
		bufferSize:  defaultSubscriberBufferSize,
	}
}

func (hub *Hub) Subscribe(ctx context.Context) <-chan any {
	ch := make(chan any, hub.bufferSize)

	hub.mu.Lock()
	hub.subscribers[ch] = struct{}{}
	hub.mu.Unlock()

	go func() {
		<-ctx.Done()
		hub.unsubscribe(ch)
	}()

	return ch
}

func (hub *Hub) unsubscribe(ch chan any) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if _, ok := hub.subscribers[ch]; !ok {
		return
	}

	delete(hub.subscribers, ch)
	close(ch)
}

func (hub *Hub) Publish(event any) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	for ch := range hub.subscribers {
		select {
		case ch <- event:
		default:
			// Drop if subscriber is slow to avoid blocking writers.
		}
	}
}

func (hub *Hub) subscriberCount() int {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	return len(hub.subscribers)
}
