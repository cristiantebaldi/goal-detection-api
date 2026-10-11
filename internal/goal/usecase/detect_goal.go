package usecase

import (
	"context"
	"sync"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
)

type detectGoalUseCase struct {
	mu                    sync.Mutex
	ballTrackerRepository contract.BallTrackerRepository
	bounds                model.Bounds
	publishGoalEvent      contract.PublishGoalEvent
}

func (usecase *detectGoalUseCase) Execute(ctx context.Context, input model.DetectGoalInput) (model.DetectGoalOutput, error) {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if input.SignalLost {
		usecase.ballTrackerRepository.Reset(ctx)
		return model.DetectGoalOutput{Detected: false}, nil
	}

	tracker := usecase.ballTrackerRepository.Get(ctx)
	inGoal := usecase.bounds.Contains(input.Position)
	output := model.DetectGoalOutput{Detected: false}

	switch tracker.State {
	case model.BallStateInField:
		if inGoal {
			event := model.NewGoalEvent(input.Timestamp)
			usecase.publishGoalEvent.Execute(ctx, event)
			tracker.State = model.BallStateInGoal
			output = model.DetectGoalOutput{Detected: true, Event: event}
		}
	case model.BallStateInGoal:
		if !inGoal {
			tracker.State = model.BallStateInField
		}
	default:
		if inGoal {
			tracker.State = model.BallStateInGoal
		} else {
			tracker.State = model.BallStateInField
		}
	}

	usecase.ballTrackerRepository.Save(ctx, tracker)
	return output, nil
}

func NewDetectGoalUseCase(
	ballTrackerRepository contract.BallTrackerRepository,
	bounds model.Bounds,
	publishGoalEvent contract.PublishGoalEvent,
) contract.DetectGoalUseCase {
	return &detectGoalUseCase{
		ballTrackerRepository: ballTrackerRepository,
		bounds:                bounds,
		publishGoalEvent:      publishGoalEvent,
	}
}
