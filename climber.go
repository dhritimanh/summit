package main

const (
	FitnessMin = 0
	FitnessMax = 100
	AmsMin     = 0
	AmsMax     = 100
)

// Climber represents a single rope team member
type Climber struct {
	Name          string
	Fitness       int
	AMS           int
	Loc           int
	Summited      bool
	WarningActive bool
	WarningTurns  int // turns remaining before crisis auto-fires
	WarningActed  bool
	Resting       bool
	Altitude      int
}

// Utility for keeping stats in range
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
