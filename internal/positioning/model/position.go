package model

type Anchor struct {
	ID string
	X  float64
	Y  float64
	Z  float64
}

type DistanceReading struct {
	AnchorID       string
	DistanceMeters float64
}

type CalculatePositionInput struct {
	Distances []DistanceReading
}

type Position struct {
	X float64
	Y float64
	Z float64
}
