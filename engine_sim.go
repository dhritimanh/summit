package main

import (
	"fmt"
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

// Global simulation session (initialized in main)
var session *WorldState


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
		fitLoss := AltitudeFitLossMin + session.Rng.Intn(AltitudeFitLossMax-AltitudeFitLossMin+1)
		
		if c.Resting {
			if c.Altitude >= currentMountain.DeathZone {
				fitLoss = 4 // Resting barely helps in death zone
			} else {
				fitLoss = 1 // Resting significantly mitigates decay
			}
		}

		isNight := session.Hour >= 18 || session.Hour < 6
		if isNight {
			fitLoss += 2
		}

		gainRange := amsGainByLoc[c.Loc]
		amsGain := float64(gainRange[0] + session.Rng.Intn(gainRange[1]-gainRange[0]+1))
		
		// Apply Archetype susceptibility
		if amsGain > 0 {
			amsGain *= session.Archetype.AMSSusPercent
		}

		if isNight {
			amsGain += 1
		}



		
		ApplyStatChange(c, "Fitness", -fitLoss, "Environment")
		ApplyStatChange(c, "AMS", int(amsGain), "Environment")
	} else {

		// BASE CAMP SYSTEM: Automatic passive recovery
		ApplyStatChange(c, "Fitness", BaseRecoveryFit, "Environment")
		ApplyStatChange(c, "AMS", -BaseRecoveryAms, "Environment")
	}

	c.Resting = false // Reset modifier after time passes

	// ─── Phase 3: Oxygen Consumption ───
	consumeO2(c)

	// Advance clock (3hr turns)
	session.Hour += 3
	if session.Hour >= 24 {
		session.Hour = 0
		session.Day++
		// NEW: Scan for window at start of day
		ScanForSummitWindow()
	}



	// Warning timer countdown (if active)
	for _, t := range c.ActiveThreats {
		if t.Level == LevelWarning || t.Level == LevelCrisis {
			t.TurnsLeft--
		}
	}

	// Update wind from curve + altitude scaling
	baseWind := session.WeatherCurve[session.Day % 31]
	if c.Altitude < 6000 {
		session.WindSpeed = baseWind / 2
	} else if c.Altitude < 8000 {
		session.WindSpeed = baseWind
	} else {
		session.WindSpeed = baseWind + 30 // Death zone boost
	}


	if session.WindSpeed >= 60 {
		c.ExposureTurns++
	} else if c.Resting {
		c.ExposureTurns = 0
	}
}

func consumeO2(c *Climber) {
	// Base camp doesn't need O2
	if c.Loc == LocBase {
		c.O2Charges = 3 // Always full at base
		return
	}

	if c.O2Charges > 0 {
		c.O2Charges--
	}

	// Automatic reload if at 0 but bottles are available at camp
	if c.O2Charges == 0 && session.CampO2[c.Loc] > 0 {
		session.CampO2[c.Loc]--
		c.O2Charges = 3
		logLine(fmt.Sprintf("%s loaded a fresh O2 bottle at %s.", c.Name, currentMountain.CampNames[c.Loc]))
	}
}


func checkThreats(c *Climber) *ActiveThreat {
	if c.ActiveThreats == nil {
		c.ActiveThreats = make(map[ThreatType]*ActiveThreat)
	}

	// 1. Check for expired warnings converting to crisis
	for _, t := range c.ActiveThreats {
		if t.Level == LevelWarning && t.TurnsLeft <= 0 {
			t.Level = LevelCrisis
			return t
		}
	}

	// 2. Return an existing ticking warning or crisis
	for _, t := range c.ActiveThreats {
		if t.Level == LevelWarning || t.Level == LevelCrisis {
			return t
		}
	}

	// 3. Oxygen Logic (Highest Priority New Threat)
	if c.O2Charges <= 0 && c.Loc > LocBase {
		return assignThreat(c, ThreatOxygen, LevelCrisis)
	}

	// 4. Process new AMS threats
	if c.AMS >= AmsCrisis {
		return assignThreat(c, ThreatAMS, LevelCrisis)
	}
	if c.AMS >= AmsWarning {
		return assignThreat(c, ThreatAMS, LevelWarning)
	}

	// 5. Frostbite logic
	if session.WindSpeed >= 60 && c.ExposureTurns > 3 {
		if c.ExposureTurns > 6 {
			return assignThreat(c, ThreatFrostbite, LevelCrisis)
		} else if c.ExposureTurns > 4 {
			return assignThreat(c, ThreatFrostbite, LevelWarning)
		} else {
			return assignThreat(c, ThreatFrostbite, LevelWhisper)
		}
	}

	if c.AMS >= AmsWhisper && session.Rng.Intn(2) == 0 {
		return &ActiveThreat{Type: ThreatAMS, Level: LevelWhisper}
	}

	return nil
}

func ScanForSummitWindow() {
	// Look at next 5 days
	consecutiveLowWind := 0
	windowFound := false
	startDay := 0

	for i := 0; i < 5; i++ {
		targetDay := (session.Day + i) % 31
		if targetDay == 0 {
			targetDay = 1
		}
		wind := session.WeatherCurve[targetDay]
		if wind < 40 {
			if consecutiveLowWind == 0 {
				startDay = session.Day + i
			}
			consecutiveLowWind++
		} else {
			consecutiveLowWind = 0
		}

		if consecutiveLowWind >= 3 {
			windowFound = true
			break
		}
	}

	if windowFound {
		logLine(fmt.Sprintf("[METEO] Base reports stable window opening around Day %d.", startDay))
	}
}



func assignThreat(c *Climber, tType ThreatType, level ThreatLevel) *ActiveThreat {
	t, exists := c.ActiveThreats[tType]
	if !exists {
		t = &ActiveThreat{Type: tType}
		c.ActiveThreats[tType] = t
	}

	if t.Level != LevelCrisis {
		if level == LevelWarning && t.Level != LevelWarning {
			t.Level = LevelWarning
			t.TurnsLeft = 4
			t.WarningActed = false
		} else {
			t.Level = level
		}
	}
	
	// Once Frostbite crisis hits, add permanent modifier if it doesn't already exist
	if tType == ThreatFrostbite && t.Level == LevelCrisis {
		hasMod := false
		for _, m := range c.Modifiers {
			if m.Type == ModFrostbite {
				hasMod = true
			}
		}
		if !hasMod {
			c.Modifiers = append(c.Modifiers, Modifier{Type: ModFrostbite, MaxFitPenalty: 25})
			ApplyStatChange(c, "Fitness", 0, "Environment") // Enforce cap immediately
		}
	}

	return t
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
