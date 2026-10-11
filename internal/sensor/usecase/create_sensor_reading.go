package usecase

import (
	"context"

	goalcontract "github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	goalmodel "github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	positioningcontract "github.com/cristiantebaldi/goal-detection-api/internal/positioning/contract"
	positioningmodel "github.com/cristiantebaldi/goal-detection-api/internal/positioning/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/model"
)

type createSensorReadingUseCase struct {
	calculatePositionUseCase positioningcontract.CalculatePositionUseCase
	detectGoalUseCase        goalcontract.DetectGoalUseCase
}

func (usecase *createSensorReadingUseCase) Execute(ctx context.Context, input model.CreateSensorReadingInput) (model.CreateSensorReadingOutput, error) {
	if input.SignalLost || len(input.Positions) == 0 {
		output, err := usecase.detectGoalUseCase.Execute(ctx, goalmodel.DetectGoalInput{
			Timestamp:  input.Timestamp,
			SignalLost: true,
		})
		if err != nil {
			return model.CreateSensorReadingOutput{}, err
		}
		return model.CreateSensorReadingOutput{GoalDetected: output.Detected}, nil
	}

	distances := make([]positioningmodel.DistanceReading, 0, len(input.Positions))
	for _, position := range input.Positions {
		distances = append(distances, positioningmodel.DistanceReading{
			AnchorID:       position.AnchorID,
			DistanceMeters: position.DistanceMeters,
		})
	}

	position, err := usecase.calculatePositionUseCase.Execute(ctx, positioningmodel.CalculatePositionInput{
		Distances: distances,
	})
	if err != nil {
		return model.CreateSensorReadingOutput{}, err
	}

	output, err := usecase.detectGoalUseCase.Execute(ctx, goalmodel.DetectGoalInput{
		Position: goalmodel.Position{
			X: position.X,
			Y: position.Y,
			Z: position.Z,
		},
		Timestamp: input.Timestamp,
	})
	if err != nil {
		return model.CreateSensorReadingOutput{}, err
	}

	return model.CreateSensorReadingOutput{GoalDetected: output.Detected}, nil
}

func NewCreateSensorReadingUseCase(
	calculatePositionUseCase positioningcontract.CalculatePositionUseCase,
	detectGoalUseCase goalcontract.DetectGoalUseCase,
) contract.CreateSensorReadingUseCase {
	return &createSensorReadingUseCase{
		calculatePositionUseCase: calculatePositionUseCase,
		detectGoalUseCase:        detectGoalUseCase,
	}
}
