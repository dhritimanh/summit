package main

import (
	"math/rand"
	"time"
)

// ─── Simulation Constants ─────────────────────────────────────────────────────

const (
	// Peak thresholds for threat detection
	AmsWhisper = 30 // High chance logic
	AmsWarning = 55 // User must respond or face a crisis
	AmsCrisis  = 80 // Potential HACE / death if ignored

	// Environmental Systems
	BaseRecoveryFit    = 6 // +48 per day at Base
	BaseRecoveryAms    = 8 // -64 per day at Base
	AltitudeFitLossMin = 3 // -24 per day at altitude
	AltitudeFitLossMax = 5 // -40 per day at altitude
)

// Global simulation state
var (
	day  = 1
	hour = 6
	rng  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// AMS Growth rates per camp (faster/punishing the higher you are)
var amsGainByLoc = map[int][2]int{
	LocBase:     {-8, -5},  // negative gain = recovery
	LocCamp1:    {0, 2},
	LocCamp2:    {1, 3},
	LocHighCamp: {2, 5},
	LocSummit:   {3, 6},
}

// ─── Simulation Methods ───────────────────────────────────────────────────────

func advanceTime(c *Climber) {
	// SYSTEM: Environmental Pressure (runs after user action)
	if c.Loc > LocBase {
		// ALTITUDE SYSTEM: Passive decay and AMS accumulation
		fitLoss := AltitudeFitLossMin + rng.Intn(AltitudeFitLossMax-AltitudeFitLossMin+1)
		
		isNight := hour >= 18 || hour < 6
		if isNight {
			fitLoss += 2
		}

		if c.Resting {
			if c.Altitude >= currentMountain.DeathZone {
				fitLoss = 4 // Resting barely helps in death zone
			} else {
				fitLoss = 1 // Resting significantly mitigates decay
			}
		}

		gainRange := amsGainByLoc[c.Loc]
		amsGain := gainRange[0] + rng.Intn(gainRange[1]-gainRange[0]+1)
		if isNight {
			amsGain += 1
		}

		
		c.Fitness = clamp(c.Fitness-fitLoss, FitnessMin, FitnessMax)
		c.AMS = clamp(c.AMS+amsGain, AmsMin, AmsMax)
	} else {
		// BASE CAMP SYSTEM: Automatic passive recovery
		c.Fitness = clamp(c.Fitness+BaseRecoveryFit, FitnessMin, FitnessMax)
		c.AMS = clamp(c.AMS-BaseRecoveryAms, AmsMin, AmsMax)
	}

	c.Resting = false // Reset modifier after time passes

	// Advance clock (3hr turns)
	hour += 3
	if hour >= 24 {
		hour = 0
		day++
	}

	// Warning timer countdown (if active)
	if c.WarningActive {
		c.WarningTurns--
	}
}

func checkThreats(c *Climber) string {
	// Crisis: auto-fire if warning timer expired
	if c.WarningActive && c.WarningTurns <= 0 {
		c.WarningActive = false
		return "crisis"
	}

	// Warning already active — wait for timer or action
	if c.WarningActive {
		return "none"
	}

	if c.AMS >= AmsCrisis {
		return "crisis"
	}
	if c.AMS >= AmsWarning {
		c.WarningActive = true
		c.WarningTurns = 4 // 12 in-game hours to act
		c.WarningActed = false
		return "warning"
	}
	if c.AMS >= AmsWhisper && rng.Intn(2) == 0 {
		return "whisper"
	}
	return "none"
}

func checkTerminal(c *Climber) (bool, string) {
	if c.Fitness <= 0 {
		return true, "death"
	}
	if c.AMS >= 100 {
		return true, "ams_death"
	}
	if c.Loc == LocSummit && !c.Summited {
		c.Summited = true
		return false, "summit"
	}
	if c.Summited && c.Loc == LocBase {
		return true, "success"
	}
	return false, ""
}
