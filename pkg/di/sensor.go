package di

import (
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/controller"
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/usecase"
	"github.com/cristiantebaldi/goal-detection-api/pkg/config"
	"github.com/cristiantebaldi/goal-detection-api/pkg/eventbus"
)

func NewCreateSensorReadingController(cfg config.Config, hub *eventbus.Hub) contract.CreateSensorReadingController {
	return controller.NewCreateSensorReadingController(
		usecase.NewCreateSensorReadingUseCase(
			NewCalculatePositionUseCase(cfg.Anchors),
			NewDetectGoalUseCase(cfg.Goal, hub),
		),
	)
}
