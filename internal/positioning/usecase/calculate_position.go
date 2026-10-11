package usecase

import (
	"context"
	"math"

	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/positioning/model"
)

type calculatePositionUseCase struct {
	anchors []model.Anchor
}

func (usecase *calculatePositionUseCase) Execute(ctx context.Context, input model.CalculatePositionInput) (model.Position, error) {
	if len(usecase.anchors) < 3 || len(input.Distances) < 3 {
		return model.Position{}, model.ErrInsufficientAnchors{}
	}

	distancesByAnchor := make(map[string]float64, len(input.Distances))
	for _, reading := range input.Distances {
		distancesByAnchor[reading.AnchorID] = reading.DistanceMeters
	}

	resolved := make([]resolvedAnchor, 0, len(usecase.anchors))
	for _, anchor := range usecase.anchors {
		distance, ok := distancesByAnchor[anchor.ID]
		if !ok {
			return model.Position{}, model.ErrMissingAnchor{ID: anchor.ID}
		}
		resolved = append(resolved, resolvedAnchor{
			anchor:   anchor,
			distance: distance,
		})
	}

	x, y, err := trilaterate2D(resolved)
	if err != nil {
		return model.Position{}, err
	}

	return model.Position{X: x, Y: y, Z: 0}, nil
}

type resolvedAnchor struct {
	anchor   model.Anchor
	distance float64
}

func trilaterate2D(anchors []resolvedAnchor) (float64, float64, error) {
	a1 := anchors[0]
	a2 := anchors[1]
	a3 := anchors[2]

	a11 := 2 * (a2.anchor.X - a1.anchor.X)
	a12 := 2 * (a2.anchor.Y - a1.anchor.Y)
	a21 := 2 * (a3.anchor.X - a1.anchor.X)
	a22 := 2 * (a3.anchor.Y - a1.anchor.Y)

	r1sq := a1.distance * a1.distance
	r2sq := a2.distance * a2.distance
	r3sq := a3.distance * a3.distance

	b1 := r1sq - r2sq - a1.anchor.X*a1.anchor.X + a2.anchor.X*a2.anchor.X - a1.anchor.Y*a1.anchor.Y + a2.anchor.Y*a2.anchor.Y
	b2 := r1sq - r3sq - a1.anchor.X*a1.anchor.X + a3.anchor.X*a3.anchor.X - a1.anchor.Y*a1.anchor.Y + a3.anchor.Y*a3.anchor.Y

	det := a11*a22 - a12*a21
	if math.Abs(det) < 1e-9 {
		return 0, 0, model.ErrColinearAnchors{}
	}

	x := (b1*a22 - b2*a12) / det
	y := (a11*b2 - a21*b1) / det
	return x, y, nil
}

func NewCalculatePositionUseCase(anchors []model.Anchor) contract.CalculatePositionUseCase {
	return &calculatePositionUseCase{anchors: anchors}
}
