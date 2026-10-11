package model

import (
	"testing"
	"time"
)

func TestBoundsContains(t *testing.T) {
	bounds := Bounds{X: 40, YMin: 8.5, YMax: 11.5}

	t.Run("deve considerar gol quando x cruza a linha no centro da trave", func(t *testing.T) {
		if !bounds.Contains(Position{X: 40, Y: 10}) {
			t.Fatal("esperava que (40, 10) estivesse dentro do gol")
		}
	})

	t.Run("não deve considerar gol quando a bola para na trave", func(t *testing.T) {
		if bounds.Contains(Position{X: 39.8, Y: 10}) {
			t.Fatal("não esperava que (39.8, 10) estivesse dentro do gol")
		}
	})

	t.Run("não deve considerar gol quando y está fora da trave", func(t *testing.T) {
		if bounds.Contains(Position{X: 41, Y: 8.4}) {
			t.Fatal("não esperava que (41, 8.4) estivesse dentro do gol")
		}
	})
}

func TestGoalEventToStreamEvent(t *testing.T) {
	t.Run("deve formatar o payload SSE com event e detected_date", func(t *testing.T) {
		detectedDate := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
		stream := NewGoalEvent(detectedDate).ToStreamEvent()

		if stream.Event != GoalEventDetected {
			t.Fatalf("esperava event %s, obteve %s", GoalEventDetected, stream.Event)
		}
		if stream.DetectedDate != "2026-09-07T12:00:00.000Z" {
			t.Fatalf("esperava detected_date 2026-09-07T12:00:00.000Z, obteve %s", stream.DetectedDate)
		}
	})
}
