package contract

import (
	"context"
	"net/http"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
)

type PublishGoalEvent interface {
	Execute(context.Context, model.GoalEvent)
}

type SubscribeGoalEvent interface {
	Execute(context.Context) <-chan model.GoalEvent
}

type BallTrackerRepository interface {
	Get(context.Context) model.BallTracker
	Save(context.Context, model.BallTracker)
	Reset(context.Context)
}

type DetectGoalUseCase interface {
	Execute(context.Context, model.DetectGoalInput) (model.DetectGoalOutput, error)
}

type ListStreamGoalEventUseCase interface {
	Execute(context.Context) (<-chan model.GoalEvent, error)
}

type ListStreamGoalEventController interface {
	Execute(http.ResponseWriter, *http.Request)
}
