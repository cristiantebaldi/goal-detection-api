package events

import (
	"testing"
	"time"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	"github.com/cristiantebaldi/goal-detection-api/pkg/eventbus"
)

func TestPublishGoalEvent(t *testing.T) {
	t.Run("deve publicar model.GoalEvent no hub", func(t *testing.T) {
		hub := eventbus.NewHub()
		ctx := t.Context()
		subscribed := hub.Subscribe(ctx)
		publisher := NewPublishGoalEvent(hub)
		detectedDate := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
		event := model.NewGoalEvent(detectedDate)

		publisher.Execute(ctx, event)

		got, ok := (<-subscribed).(model.GoalEvent)
		if !ok {
			t.Fatal("esperava model.GoalEvent")
		}
		if got.Type != model.GoalEventDetected {
			t.Fatalf("esperava tipo %s, obteve %s", model.GoalEventDetected, got.Type)
		}
	})
}

func TestSubscribeGoalEvent(t *testing.T) {
	t.Run("deve encaminhar model.GoalEvent publicado no hub", func(t *testing.T) {
		hub := eventbus.NewHub()
		ctx := t.Context()
		subscribed := NewSubscribeGoalEvent(hub).Execute(ctx)
		event := model.NewGoalEvent(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC))

		hub.Publish(event)

		select {
		case got := <-subscribed:
			if got.Type != model.GoalEventDetected {
				t.Fatalf("esperava tipo %s, obteve %s", model.GoalEventDetected, got.Type)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout esperando evento")
		}
	})

	t.Run("deve ignorar valor que não é model.GoalEvent", func(t *testing.T) {
		hub := eventbus.NewHub()
		ctx := t.Context()
		subscribed := NewSubscribeGoalEvent(hub).Execute(ctx)
		event := model.NewGoalEvent(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC))

		hub.Publish("nao-e-gol")
		hub.Publish(event)

		select {
		case got := <-subscribed:
			if got.Type != model.GoalEventDetected {
				t.Fatalf("esperava tipo %s, obteve %s", model.GoalEventDetected, got.Type)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout esperando evento")
		}
	})
}
