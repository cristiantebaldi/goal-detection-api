package model

import "time"

type GoalEvent struct {
	Type         string
	DetectedDate time.Time
}

func NewGoalEvent(detectedDate time.Time) GoalEvent {
	return GoalEvent{
		Type:         GoalEventDetected,
		DetectedDate: detectedDate,
	}
}

type GoalStreamEvent struct {
	Event        string `json:"event"`
	DetectedDate string `json:"detected_date"`
}

func (event GoalEvent) ToStreamEvent() GoalStreamEvent {
	return GoalStreamEvent{
		Event:        event.Type,
		DetectedDate: event.DetectedDate.UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}
