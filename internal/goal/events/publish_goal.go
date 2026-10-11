package events

import (
	"context"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	"github.com/cristiantebaldi/goal-detection-api/pkg/eventbus"
)

type publishGoalEvent struct {
	hub *eventbus.Hub
}

func (publisher *publishGoalEvent) Execute(ctx context.Context, event model.GoalEvent) {
	publisher.hub.Publish(event)
}

func NewPublishGoalEvent(hub *eventbus.Hub) contract.PublishGoalEvent {
	return &publishGoalEvent{hub: hub}
}
