package model

import "time"

const (
	GoalEventDetected    = "goal"
	SSEHeartbeatInterval = 15 * time.Second
)

type BallState string

const (
	BallStateUnknown BallState = "unknown"
	BallStateInField BallState = "in_field"
	BallStateInGoal  BallState = "in_goal"
)

type BallTracker struct {
	State BallState
}

type Position struct {
	X float64
	Y float64
	Z float64
}

type Bounds struct {
	X    float64
	YMin float64
	YMax float64
}

func (bounds Bounds) Contains(position Position) bool {
	return position.X >= bounds.X && position.Y >= bounds.YMin && position.Y <= bounds.YMax
}

type DetectGoalInput struct {
	Position   Position
	Timestamp  time.Time
	SignalLost bool
}

type DetectGoalOutput struct {
	Detected bool
	Event    GoalEvent
}
