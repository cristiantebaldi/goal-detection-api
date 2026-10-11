package eventbus

import (
	"context"
	"testing"
	"time"
)

type testEvent struct {
	ID int
}

func TestHub(t *testing.T) {
	t.Run("deve entregar o evento publicado para o inscrito", func(t *testing.T) {
		hub := NewHub()
		ctx := t.Context()

		ch := hub.Subscribe(ctx)
		hub.Publish(testEvent{ID: 42})

		select {
		case event := <-ch:
			got, ok := event.(testEvent)
			if !ok {
				t.Fatalf("esperava testEvent, obteve %T", event)
			}
			if got.ID != 42 {
				t.Fatalf("esperava id 42, obteve %d", got.ID)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout esperando evento")
		}
	})

	t.Run("deve cancelar a inscrição quando o contexto é cancelado", func(t *testing.T) {
		hub := NewHub()
		ctx, cancel := context.WithCancel(context.Background())

		_ = hub.Subscribe(ctx)
		if hub.subscriberCount() != 1 {
			t.Fatalf("esperava 1 inscrito, obteve %d", hub.subscriberCount())
		}

		cancel()
		deadline := time.Now().Add(time.Second)
		for time.Now().Before(deadline) {
			if hub.subscriberCount() == 0 {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}

		t.Fatalf("esperava 0 inscritos após cancel, obteve %d", hub.subscriberCount())
	})

	t.Run("não deve bloquear o publish quando o inscrito é lento", func(t *testing.T) {
		hub := NewHub()
		hub.bufferSize = 1

		ctx := t.Context()

		ch := hub.Subscribe(ctx)

		hub.Publish(testEvent{ID: 1})
		hub.Publish(testEvent{ID: 2})
		hub.Publish(testEvent{ID: 3})

		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatal("esperava pelo menos um evento no buffer")
		}
	})
}
