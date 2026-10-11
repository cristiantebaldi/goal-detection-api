package di

import (
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/controller"
	goalevents "github.com/cristiantebaldi/goal-detection-api/internal/goal/events"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/repository"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/usecase"
	"github.com/cristiantebaldi/goal-detection-api/pkg/config"
	"github.com/cristiantebaldi/goal-detection-api/pkg/eventbus"
)

func NewGoalEventHub() *eventbus.Hub {
	return eventbus.NewHub()
}

func NewDetectGoalUseCase(goal config.Goal, hub *eventbus.Hub) contract.DetectGoalUseCase {
	return usecase.NewDetectGoalUseCase(
		repository.NewBallTrackerMemoryRepository(),
		model.Bounds{X: goal.X, YMin: goal.YMin, YMax: goal.YMax},
		goalevents.NewPublishGoalEvent(hub),
	)
}

func NewListStreamGoalEventController(hub *eventbus.Hub) contract.ListStreamGoalEventController {
	return controller.NewListStreamGoalEventController(
		usecase.NewListStreamGoalEventUseCase(
			goalevents.NewSubscribeGoalEvent(hub),
		),
	)
}
