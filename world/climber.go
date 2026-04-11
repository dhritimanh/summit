package world

import "summit/data"

const (
	FitnessMin = 0
	FitnessMax = 100
	AmsMin     = 0
	AmsMax     = 100
)

// ModifierType categorizes active effects
type ModifierType string

const (
	ModFrostbite ModifierType = "Frostbite"
	ModMeds      ModifierType = "Meds"
)

// ClimberOrder represents the current objective for a climber
type ClimberOrder string

const (
	OrderHold    ClimberOrder = "Hold"
	OrderClimb   ClimberOrder = "Climb"
	OrderDescend ClimberOrder = "Descend"
	OrderRest    ClimberOrder = "Rest"
)

type Modifier struct {
	Type          ModifierType
	MaxFitPenalty int
}

// Climber represents a single rope team member
type Climber struct {
	Name          string
	Archetype     data.ClimberArchetype
	Fitness       int
	AMS           int
	Loc           int
	Summited      bool
	ActiveThreats map[ThreatType]*ActiveThreat
	ExposureTurns int // Turns exposed to high winds
	O2Charges     int // 0-3 (1 bottle = 3 charges)
	O2Active      bool
	Resting       bool // Internal state for sim calculation
	Order         ClimberOrder
	OrderTarget   int // Altitude or Loc target
	Altitude      int
	Modifiers     []Modifier
	PulseOffset   float64 // For Visual animation
}

// Utility for keeping stats in range
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ApplyStatChange safely modifies stats, applying any active modifiers and caps.
func ApplyStatChange(c *Climber, stat string, delta int, source string) {
	if stat == "Fitness" {
		maxFit := FitnessMax
		for _, m := range c.Modifiers {
			maxFit -= m.MaxFitPenalty
		}
		maxFit = Clamp(maxFit, FitnessMin, FitnessMax)

		c.Fitness = Clamp(c.Fitness+delta, FitnessMin, maxFit)
	} else if stat == "AMS" {
		c.AMS = Clamp(c.AMS+delta, AmsMin, AmsMax)
	}
}
