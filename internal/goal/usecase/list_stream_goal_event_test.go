package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
)

type fakeSubscribeGoalEvent struct {
	events chan model.GoalEvent
}

func (f *fakeSubscribeGoalEvent) Execute(ctx context.Context) <-chan model.GoalEvent {
	return f.events
}

var _ contract.SubscribeGoalEvent = (*fakeSubscribeGoalEvent)(nil)

func TestListStreamGoalEvent(t *testing.T) {
	t.Run("deve encaminhar o evento inscrito", func(t *testing.T) {
		events := make(chan model.GoalEvent, 1)
		usecase := NewListStreamGoalEventUseCase(&fakeSubscribeGoalEvent{events: events})
		ctx := t.Context()

		stream, err := usecase.Execute(ctx)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}

		expected := model.NewGoalEvent(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC))
		events <- expected

		select {
		case got := <-stream:
			if got.Type != expected.Type {
				t.Fatalf("esperava tipo %s, obteve %s", expected.Type, got.Type)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout esperando evento")
		}
	})
}
