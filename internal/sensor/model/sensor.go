package model

import "time"

type AnchorDistance struct {
	AnchorID       string  `json:"anchor_id" validate:"required"`
	DistanceMeters float64 `json:"distance_meters" validate:"gte=0"`
}

type CreateSensorReadingInput struct {
	Timestamp  time.Time        `json:"timestamp" validate:"required"`
	SignalLost bool             `json:"signal_lost"`
	Positions  []AnchorDistance `json:"positions" validate:"dive"`
}

type CreateSensorReadingOutput struct {
	GoalDetected bool `json:"goal_detected"`
}
