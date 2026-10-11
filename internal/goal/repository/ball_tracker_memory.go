package repository

import (
	"context"
	"sync"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
)

type ballTrackerMemoryRepository struct {
	mu      sync.Mutex
	tracker model.BallTracker
}

func (repository *ballTrackerMemoryRepository) Get(ctx context.Context) model.BallTracker {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	return repository.tracker
}

func (repository *ballTrackerMemoryRepository) Save(ctx context.Context, tracker model.BallTracker) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.tracker = tracker
}

func (repository *ballTrackerMemoryRepository) Reset(ctx context.Context) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.tracker = model.BallTracker{State: model.BallStateUnknown}
}

func NewBallTrackerMemoryRepository() contract.BallTrackerRepository {
	return &ballTrackerMemoryRepository{
		tracker: model.BallTracker{State: model.BallStateUnknown},
	}
}
