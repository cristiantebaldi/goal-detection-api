package di

import (
	"github.com/cristiantebaldi/goal-detection-api/internal/health/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/health/controller"
)

func NewHealthController() contract.HealthController {
	return controller.NewHealthController()
}
