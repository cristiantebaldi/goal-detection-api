package events

import (
	"context"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	"github.com/cristiantebaldi/goal-detection-api/pkg/eventbus"
)

type subscribeGoalEvent struct {
	hub *eventbus.Hub
}

func (subscriber *subscribeGoalEvent) Execute(ctx context.Context) <-chan model.GoalEvent {
	rawEvents := subscriber.hub.Subscribe(ctx)
	events := make(chan model.GoalEvent, 16)
	go func() {
		defer close(events)
		for raw := range rawEvents {
			event, ok := raw.(model.GoalEvent)
			if !ok {
				continue
			}
			select {
			case events <- event:
			case <-ctx.Done():
				return
			}
		}
	}()
	return events
}

func NewSubscribeGoalEvent(hub *eventbus.Hub) contract.SubscribeGoalEvent {
	return &subscribeGoalEvent{hub: hub}
}
