package usecase

import (
	"context"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
)

type listStreamGoalEventUseCase struct {
	subscribeGoalEvent contract.SubscribeGoalEvent
}

func (usecase *listStreamGoalEventUseCase) Execute(ctx context.Context) (<-chan model.GoalEvent, error) {
	return usecase.subscribeGoalEvent.Execute(ctx), nil
}

func NewListStreamGoalEventUseCase(subscribeGoalEvent contract.SubscribeGoalEvent) contract.ListStreamGoalEventUseCase {
	return &listStreamGoalEventUseCase{
		subscribeGoalEvent: subscribeGoalEvent,
	}
}
