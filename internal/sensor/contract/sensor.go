package contract

import (
	"context"
	"net/http"

	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/model"
)

type CreateSensorReadingUseCase interface {
	Execute(context.Context, model.CreateSensorReadingInput) (model.CreateSensorReadingOutput, error)
}

type CreateSensorReadingController interface {
	Execute(http.ResponseWriter, *http.Request)
}
