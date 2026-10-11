package usecase

import (
	"math"
	"testing"

	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/model"
)

func defaultAnchors() []model.Anchor {
	return []model.Anchor{
		{ID: "anchor-1", X: 0, Y: 0, Z: 0},
		{ID: "anchor-2", X: 40, Y: 0, Z: 0},
		{ID: "anchor-3", X: 0, Y: 20, Z: 0},
	}
}

func distance2D(ax, ay, bx, by float64) float64 {
	return math.Hypot(ax-bx, ay-by)
}

func TestCalculatePosition(t *testing.T) {
	usecase := NewCalculatePositionUseCase(defaultAnchors())
	ctx := t.Context()

	t.Run("deve calcular a posição no centro da quadra", func(t *testing.T) {
		x, y := 20.0, 10.0
		position, err := usecase.Execute(ctx, model.CalculatePositionInput{
			Distances: []model.DistanceReading{
				{AnchorID: "anchor-1", DistanceMeters: distance2D(0, 0, x, y)},
				{AnchorID: "anchor-2", DistanceMeters: distance2D(40, 0, x, y)},
				{AnchorID: "anchor-3", DistanceMeters: distance2D(0, 20, x, y)},
			},
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if math.Abs(position.X-x) > 0.01 {
			t.Errorf("esperava X %.2f, obteve %.4f", x, position.X)
		}
		if math.Abs(position.Y-y) > 0.01 {
			t.Errorf("esperava Y %.2f, obteve %.4f", y, position.Y)
		}
	})

	t.Run("deve calcular a posição sobre a linha do gol", func(t *testing.T) {
		x, y := 40.0, 10.0
		position, err := usecase.Execute(ctx, model.CalculatePositionInput{
			Distances: []model.DistanceReading{
				{AnchorID: "anchor-1", DistanceMeters: distance2D(0, 0, x, y)},
				{AnchorID: "anchor-2", DistanceMeters: distance2D(40, 0, x, y)},
				{AnchorID: "anchor-3", DistanceMeters: distance2D(0, 20, x, y)},
			},
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if math.Abs(position.X-x) > 0.01 {
			t.Errorf("esperava X %.2f, obteve %.4f", x, position.X)
		}
		if math.Abs(position.Y-y) > 0.01 {
			t.Errorf("esperava Y %.2f, obteve %.4f", y, position.Y)
		}
	})

	t.Run("deve calcular a posição antes da trave", func(t *testing.T) {
		x, y := 39.8, 10.0
		position, err := usecase.Execute(ctx, model.CalculatePositionInput{
			Distances: []model.DistanceReading{
				{AnchorID: "anchor-1", DistanceMeters: distance2D(0, 0, x, y)},
				{AnchorID: "anchor-2", DistanceMeters: distance2D(40, 0, x, y)},
				{AnchorID: "anchor-3", DistanceMeters: distance2D(0, 20, x, y)},
			},
		})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if math.Abs(position.X-x) > 0.01 {
			t.Errorf("esperava X %.2f, obteve %.4f", x, position.X)
		}
		if position.X >= 40 {
			t.Errorf("esperava X abaixo da linha do gol, obteve %.4f", position.X)
		}
	})

	t.Run("deve retornar erro se houver menos de 3 distâncias", func(t *testing.T) {
		_, err := usecase.Execute(ctx, model.CalculatePositionInput{
			Distances: []model.DistanceReading{
				{AnchorID: "anchor-1", DistanceMeters: 10},
				{AnchorID: "anchor-2", DistanceMeters: 10},
			},
		})
		if err == nil {
			t.Fatal("esperava um erro, obteve nil")
		}
		if _, ok := err.(model.ErrInsufficientAnchors); !ok {
			t.Errorf("esperava ErrInsufficientAnchors, obteve %T", err)
		}
	})

	t.Run("deve retornar erro se a âncora configurada estiver ausente", func(t *testing.T) {
		_, err := usecase.Execute(ctx, model.CalculatePositionInput{
			Distances: []model.DistanceReading{
				{AnchorID: "anchor-1", DistanceMeters: 10},
				{AnchorID: "anchor-2", DistanceMeters: 10},
				{AnchorID: "anchor-99", DistanceMeters: 10},
			},
		})
		if err == nil {
			t.Fatal("esperava um erro, obteve nil")
		}
		missing, ok := err.(model.ErrMissingAnchor)
		if !ok {
			t.Fatalf("esperava ErrMissingAnchor, obteve %T", err)
		}
		if missing.ID != "anchor-3" {
			t.Errorf("esperava âncora ausente anchor-3, obteve %s", missing.ID)
		}
	})
}
