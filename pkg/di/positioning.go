package di

import (
	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/usecase"
	"github.com/cristiantebaldi/goal-detection-api/pkg/config"
)

func NewCalculatePositionUseCase(anchors []config.Anchor) contract.CalculatePositionUseCase {
	converted := make([]model.Anchor, 0, len(anchors))
	for _, anchor := range anchors {
		converted = append(converted, model.Anchor{
			ID: anchor.ID,
			X:  anchor.X,
			Y:  anchor.Y,
			Z:  anchor.Z,
		})
	}
	return usecase.NewCalculatePositionUseCase(converted)
}
