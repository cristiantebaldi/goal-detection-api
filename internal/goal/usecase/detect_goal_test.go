package usecase

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/repository"
)

type fakePublishGoalEvent struct {
	mu     sync.Mutex
	events []model.GoalEvent
}

func (f *fakePublishGoalEvent) Execute(ctx context.Context, event model.GoalEvent) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, event)
}

func (f *fakePublishGoalEvent) published() []model.GoalEvent {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]model.GoalEvent, len(f.events))
	copy(out, f.events)
	return out
}

var _ contract.PublishGoalEvent = (*fakePublishGoalEvent)(nil)

func newDetectGoalUseCase(publisher contract.PublishGoalEvent) contract.DetectGoalUseCase {
	return NewDetectGoalUseCase(
		repository.NewBallTrackerMemoryRepository(),
		model.Bounds{X: 40, YMin: 8.5, YMax: 11.5},
		publisher,
	)
}

func TestDetectGoal(t *testing.T) {
	timestamp := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	ctx := t.Context()

	t.Run("deve publicar gol quando a bola cruza a linha", func(t *testing.T) {
		publisher := &fakePublishGoalEvent{}
		usecase := newDetectGoalUseCase(publisher)

		output, err := usecase.Execute(ctx, model.DetectGoalInput{
			Position:  model.Position{X: 20, Y: 10},
			Timestamp: timestamp,
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.Detected {
			t.Fatal("não esperava gol fora da área")
		}

		output, err = usecase.Execute(ctx, model.DetectGoalInput{
			Position:  model.Position{X: 41, Y: 10},
			Timestamp: timestamp,
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !output.Detected {
			t.Fatal("esperava gol após cruzar a linha")
		}
		published := publisher.published()
		if len(published) != 1 {
			t.Fatalf("esperava 1 evento publicado, obteve %d", len(published))
		}
		if published[0].Type != model.GoalEventDetected {
			t.Errorf("esperava tipo %s, obteve %s", model.GoalEventDetected, published[0].Type)
		}
	})

	t.Run("não deve publicar quando a bola para na trave", func(t *testing.T) {
		publisher := &fakePublishGoalEvent{}
		usecase := newDetectGoalUseCase(publisher)

		_, err := usecase.Execute(ctx, model.DetectGoalInput{
			Position:  model.Position{X: 20, Y: 10},
			Timestamp: timestamp,
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}

		output, err := usecase.Execute(ctx, model.DetectGoalInput{
			Position:  model.Position{X: 39.8, Y: 10},
			Timestamp: timestamp,
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.Detected {
			t.Fatal("não esperava gol na trave")
		}
		if len(publisher.published()) != 0 {
			t.Fatalf("não esperava evento publicado, obteve %d", len(publisher.published()))
		}
	})

	t.Run("não deve republicar enquanto a bola permanece no gol", func(t *testing.T) {
		publisher := &fakePublishGoalEvent{}
		usecase := newDetectGoalUseCase(publisher)

		_, _ = usecase.Execute(ctx, model.DetectGoalInput{Position: model.Position{X: 20, Y: 10}, Timestamp: timestamp})
		_, _ = usecase.Execute(ctx, model.DetectGoalInput{Position: model.Position{X: 41, Y: 10}, Timestamp: timestamp})
		output, err := usecase.Execute(ctx, model.DetectGoalInput{Position: model.Position{X: 41.2, Y: 10}, Timestamp: timestamp})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.Detected {
			t.Fatal("não esperava republicar gol")
		}
		if len(publisher.published()) != 1 {
			t.Fatalf("esperava 1 evento publicado, obteve %d", len(publisher.published()))
		}
	})

	t.Run("deve zerar o estado quando o sinal é perdido", func(t *testing.T) {
		publisher := &fakePublishGoalEvent{}
		usecase := newDetectGoalUseCase(publisher)

		_, _ = usecase.Execute(ctx, model.DetectGoalInput{Position: model.Position{X: 20, Y: 10}, Timestamp: timestamp})
		output, err := usecase.Execute(ctx, model.DetectGoalInput{SignalLost: true, Timestamp: timestamp})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.Detected {
			t.Fatal("não esperava gol na perda de sinal")
		}

		output, err = usecase.Execute(ctx, model.DetectGoalInput{Position: model.Position{X: 22, Y: 10}, Timestamp: timestamp})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.Detected {
			t.Fatal("não esperava gol fantasma após recuperação do sinal")
		}
		if len(publisher.published()) != 0 {
			t.Fatalf("não esperava evento publicado, obteve %d", len(publisher.published()))
		}
	})

	t.Run("não deve publicar se a primeira posição válida já está no gol", func(t *testing.T) {
		publisher := &fakePublishGoalEvent{}
		usecase := newDetectGoalUseCase(publisher)

		output, err := usecase.Execute(ctx, model.DetectGoalInput{
			Position:  model.Position{X: 41, Y: 10},
			Timestamp: timestamp,
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.Detected {
			t.Fatal("não esperava gol sem cruzamento confirmado")
		}
		if len(publisher.published()) != 0 {
			t.Fatalf("não esperava evento publicado, obteve %d", len(publisher.published()))
		}
	})
}
