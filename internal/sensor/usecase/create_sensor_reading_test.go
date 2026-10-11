package usecase

import (
	"context"
	"testing"
	"time"

	goalcontract "github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	goalmodel "github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	positioningcontract "github.com/cristiantebaldi/goal-detection-api/internal/positioning/contract"
	positioningmodel "github.com/cristiantebaldi/goal-detection-api/internal/positioning/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/model"
)

type fakeCalculatePositionUseCase struct {
	position positioningmodel.Position
	err      error
	calls    int
}

func (f *fakeCalculatePositionUseCase) Execute(ctx context.Context, input positioningmodel.CalculatePositionInput) (positioningmodel.Position, error) {
	f.calls++
	return f.position, f.err
}

type fakeDetectGoalUseCase struct {
	output        goalmodel.DetectGoalOutput
	err           error
	lastInput     goalmodel.DetectGoalInput
	signalLostN   int
	positionCalls int
}

func (f *fakeDetectGoalUseCase) Execute(ctx context.Context, input goalmodel.DetectGoalInput) (goalmodel.DetectGoalOutput, error) {
	f.lastInput = input
	if input.SignalLost {
		f.signalLostN++
	} else {
		f.positionCalls++
	}
	return f.output, f.err
}

var _ positioningcontract.CalculatePositionUseCase = (*fakeCalculatePositionUseCase)(nil)
var _ goalcontract.DetectGoalUseCase = (*fakeDetectGoalUseCase)(nil)

func TestCreateSensorReading(t *testing.T) {
	timestamp := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	ctx := t.Context()

	t.Run("deve calcular posição, detectar gol e retornar goal_detected", func(t *testing.T) {
		positioning := &fakeCalculatePositionUseCase{position: positioningmodel.Position{X: 41, Y: 10}}
		detector := &fakeDetectGoalUseCase{output: goalmodel.DetectGoalOutput{Detected: true}}
		usecase := NewCreateSensorReadingUseCase(positioning, detector)

		output, err := usecase.Execute(ctx, model.CreateSensorReadingInput{
			Timestamp: timestamp,
			Positions: []model.AnchorDistance{
				{AnchorID: "anchor-1", DistanceMeters: 22},
				{AnchorID: "anchor-2", DistanceMeters: 10},
				{AnchorID: "anchor-3", DistanceMeters: 22},
			},
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !output.GoalDetected {
			t.Fatal("esperava goal_detected true")
		}
		if positioning.calls != 1 {
			t.Fatalf("esperava 1 chamada de positioning, obteve %d", positioning.calls)
		}
		if detector.lastInput.Position.X != 41 {
			t.Errorf("esperava posição X 41, obteve %.2f", detector.lastInput.Position.X)
		}
	})

	t.Run("deve tratar perda de sinal sem calcular posição", func(t *testing.T) {
		positioning := &fakeCalculatePositionUseCase{}
		detector := &fakeDetectGoalUseCase{}
		usecase := NewCreateSensorReadingUseCase(positioning, detector)

		output, err := usecase.Execute(ctx, model.CreateSensorReadingInput{
			Timestamp:  timestamp,
			SignalLost: true,
			Positions:  []model.AnchorDistance{},
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if output.GoalDetected {
			t.Fatal("não esperava gol na perda de sinal")
		}
		if positioning.calls != 0 {
			t.Fatalf("não esperava chamada de positioning, obteve %d", positioning.calls)
		}
		if detector.signalLostN != 1 {
			t.Fatalf("esperava 1 chamada com signal_lost, obteve %d", detector.signalLostN)
		}
	})

	t.Run("deve tratar positions vazio como perda de sinal", func(t *testing.T) {
		positioning := &fakeCalculatePositionUseCase{}
		detector := &fakeDetectGoalUseCase{}
		usecase := NewCreateSensorReadingUseCase(positioning, detector)

		_, err := usecase.Execute(ctx, model.CreateSensorReadingInput{
			Timestamp: timestamp,
			Positions: []model.AnchorDistance{},
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if detector.signalLostN != 1 {
			t.Fatalf("esperava 1 chamada com signal_lost, obteve %d", detector.signalLostN)
		}
	})

	t.Run("deve propagar erro do positioning", func(t *testing.T) {
		positioning := &fakeCalculatePositionUseCase{err: positioningmodel.ErrInsufficientAnchors{}}
		detector := &fakeDetectGoalUseCase{}
		usecase := NewCreateSensorReadingUseCase(positioning, detector)

		_, err := usecase.Execute(ctx, model.CreateSensorReadingInput{
			Timestamp: timestamp,
			Positions: []model.AnchorDistance{
				{AnchorID: "anchor-1", DistanceMeters: 1},
				{AnchorID: "anchor-2", DistanceMeters: 1},
				{AnchorID: "anchor-3", DistanceMeters: 1},
			},
		})
		if err == nil {
			t.Fatal("esperava um erro, obteve nil")
		}
		if _, ok := err.(positioningmodel.ErrInsufficientAnchors); !ok {
			t.Errorf("esperava ErrInsufficientAnchors, obteve %T", err)
		}
	})
}
