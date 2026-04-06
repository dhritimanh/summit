package main

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

type Modifier struct {
	Type          ModifierType
	MaxFitPenalty int
}


// Climber represents a single rope team member
type Climber struct {
	Name          string
	Fitness       int
	AMS           int
	Loc           int
	Summited      bool
	ActiveThreats map[ThreatType]*ActiveThreat
	ExposureTurns int // Turns exposed to high winds
	O2Charges     int // 0-3 (1 bottle = 3 charges)
	Resting       bool
	Altitude      int
	Modifiers     []Modifier
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

// ApplyStatChange safely modifies stats, applying any active modifiers and caps.
// source could be "Environment", "Choice", or "Item"
func ApplyStatChange(c *Climber, stat string, delta int, source string) {
	if stat == "Fitness" {
		maxFit := FitnessMax
		for _, m := range c.Modifiers {
			maxFit -= m.MaxFitPenalty
		}
		maxFit = clamp(maxFit, FitnessMin, FitnessMax)

		c.Fitness = clamp(c.Fitness+delta, FitnessMin, maxFit)
	} else if stat == "AMS" {
		c.AMS = clamp(c.AMS+delta, AmsMin, AmsMax)
	}
}
