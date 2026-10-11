package contract

import (
	"context"

	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/model"
)

type CalculatePositionUseCase interface {
	Execute(context.Context, model.CalculatePositionInput) (model.Position, error)
}
