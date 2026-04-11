package sim

import (
	"fmt"
	"summit/data"
	"summit/world"
)

const (
	// Peak thresholds for threat detection
	AmsWhisper = 30
	AmsWarning = 55
	AmsCrisis  = 80
)

type Sim struct {
	State    *world.WorldState
	Mountain data.Mountain
	Team     []*world.Climber
	Log      []string

	// Phase 4 Clock Systems
	Speed           int     // 0 (Pause), 1, 4, 12
	AccumulatedTime float64 // Tracks fractional seconds until next 3-hr turn
}

func NewSim(w *world.WorldState, m data.Mountain) *Sim {
	return &Sim{
		State:    w,
		Mountain: m,
		Team:     make([]*world.Climber, 0),
		Log:      make([]string, 0),
		Speed:    1, // Default to 1x
	}
}

func (s *Sim) LogLine(msg string) {
	s.Log = append(s.Log, msg)
	if len(s.Log) > 15 { // Slightly more log space for GUI
		s.Log = s.Log[1:]
	}
}

// Tick is the primary simulation entry point for Phase 4
func (s *Sim) Tick(delta float64) {
	if s.Speed == 0 {
		return
	}

	// In 1x speed, a full 24h day takes ~8 minutes (480 seconds)
	// 1 game-hour = 20 real seconds.
	// 3 game-hours (one turn) = 60 real seconds.
	s.AccumulatedTime += delta * float64(s.Speed)

	if s.AccumulatedTime >= 60.0 {
		s.AccumulatedTime = 0
		for _, c := range s.Team {
			s.AdvanceTime(c)
		}
	}
}

func (s *Sim) AdvanceTime(c *world.Climber) {
	// 1. Process Order Impacts (Movement or Rest)
	c.Resting = false
	if c.Order == world.OrderClimb || c.Order == world.OrderDescend {
		s.ProcessMovement(c)
	} else if c.Order == world.OrderRest {
		c.Resting = true
	}

	// 2. CALCULATE PHYSIOLOGICAL CHANGE
	var netFitness int
	var netAms int

	// Base Camp is special: Perfect recovery
	if c.Loc == data.LocBase {
		netFitness = data.BaseRecoveryFit
		netAms = -data.BaseRecoveryAms
	} else {
		// ALititude Decay
		fitLoss := data.AltitudeFitLossMin + s.State.Rng.Intn(data.AltitudeFitLossMax-data.AltitudeFitLossMin+1)
		gainRange := data.AmsBaseGain[c.Loc]
		amsGain := float64(gainRange[0] + s.State.Rng.Intn(gainRange[1]-gainRange[0]+1))

		// Night Penalty
		isNight := s.State.Hour >= 18 || s.State.Hour < 6
		if isNight {
			fitLoss += data.NightFitPenalty
			amsGain += data.NightAmsPenalty
		}

		// RESTING LOGIC: Mitigate loss or turn into recovery
		if c.Resting {
			if c.Altitude >= s.Mountain.DeathZone {
				fitLoss = data.DeathZoneFitLoss // Still losing health in Death Zone
				amsGain -= data.DeathZoneAmsLoss 
			} else {
				// At C1/C2/C3, resting is a net POSITIVE
				fitLoss -= data.RestingFitBonus 
				amsGain -= data.RestingAmsBonus

				// ADDITIONAL BONUS if explicitly at a camp (Tents/Shelter)
				if c.Altitude == s.Mountain.CampAltitudes[c.Loc] {
					fitLoss -= data.CampRestBonus
					amsGain -= 1 // Extra AMS relief from better sleep
				}
			}
		}

		// CLIMBING PENALTY (If moving, it's harder)
		if c.Order == world.OrderClimb || c.Order == world.OrderDescend {
			fitLoss += data.ClimbingFitCost
			amsGain += data.ClimbingAmsCost
		}

		// OXYGEN MITIGATION
		if c.O2Active && c.O2Charges > 0 {
			if fitLoss > 0 { fitLoss /= 2 }
			amsGain /= 2
		}

		netFitness = -fitLoss
		netAms = int(amsGain)
	}

	world.ApplyStatChange(c, "Fitness", netFitness, "Environment")
	world.ApplyStatChange(c, "AMS", netAms, "Environment")

	// 3. RECOVERY & CONSUMPTION
	s.ConsumeO2(c)

	// 4. TIME ADVANCEMENT (Only advances after all team logic is done for the turn)
	// We handle this inside Tick normally, but the math here is per-climber
	// Logic from original Sim: 
	s.State.Hour += 3
	if s.State.Hour >= 24 {
		s.State.Hour = 0
		s.State.Day++
		s.ScanForSummitWindow()
	}

	for _, t := range c.ActiveThreats {
		if t.Level == world.LevelWarning || t.Level == world.LevelCrisis {
			t.TurnsLeft--
		}
	}

	// 5. WIND & EXPOSURE
	baseWind := s.State.WeatherCurve[s.State.Day%31]
	if c.Altitude < 6000 {
		s.State.WindSpeed = baseWind / 2
	} else if c.Altitude < 8000 {
		s.State.WindSpeed = baseWind
	} else {
		s.State.WindSpeed = baseWind + 30
	}

	isAtCamp := c.Altitude == s.Mountain.CampAltitudes[c.Loc]
	if isAtCamp && (c.Order == world.OrderRest || c.Order == world.OrderHold) {
		c.ExposureTurns = 0 // Wind can't hit you inside a tent at camp
	} else if s.State.WindSpeed >= 60 {
		c.ExposureTurns++
	} else if c.Resting {
		c.ExposureTurns = 0
	}
}

func (s *Sim) ConsumeO2(c *world.Climber) {
	if c.Loc == data.LocBase {
		c.O2Charges = 3
		return
	}
	
	// Only consume if the player has manually activated O2
	if !c.O2Active {
		return
	}

	if c.O2Charges > 0 {
		c.O2Charges--
	}
	if c.O2Charges == 0 && s.State.CampO2[c.Loc] > 0 {
		s.State.CampO2[c.Loc]--
		c.O2Charges = 3
		s.LogLine(fmt.Sprintf("%s loaded a fresh O2 bottle at %s.", c.Name, s.Mountain.CampNames[c.Loc]))
	}
}

func (s *Sim) ProcessMovement(c *world.Climber) {
	if c.Order == world.OrderClimb {
		climbGain := 300 + s.State.Rng.Intn(250)
		c.Altitude += climbGain
		target := s.Mountain.CampAltitudes[c.Loc+1]
		if c.Altitude >= target {
			c.Altitude = target
			c.Loc++
			c.Order = world.OrderHold // Stop at the camp automatically
			s.LogLine(fmt.Sprintf("%s reached %s.", c.Name, s.Mountain.CampNames[c.Loc]))
		}
	} else if c.Order == world.OrderDescend {
		c.Loc = world.Clamp(c.Loc-1, data.LocBase, data.LocSummit)
		c.Altitude = s.Mountain.CampAltitudes[c.Loc]
		c.Order = world.OrderHold
		s.LogLine(fmt.Sprintf("%s descended to %s.", c.Name, s.Mountain.CampNames[c.Loc]))
	}
}

func (s *Sim) CheckThreats(c *world.Climber) *world.ActiveThreat {
	if c.ActiveThreats == nil {
		c.ActiveThreats = make(map[world.ThreatType]*world.ActiveThreat)
	}

	for _, t := range c.ActiveThreats {
		if t.Level == world.LevelWarning && t.TurnsLeft <= 0 {
			t.Level = world.LevelCrisis
			return t
		}
	}

	for _, t := range c.ActiveThreats {
		if t.Level == world.LevelWarning || t.Level == world.LevelCrisis {
			return t
		}
	}

	if c.O2Charges <= 0 && c.Loc > data.LocBase {
		return s.AssignThreat(c, world.ThreatOxygen, world.LevelCrisis)
	}

	if c.AMS >= AmsCrisis {
		return s.AssignThreat(c, world.ThreatAMS, world.LevelCrisis)
	}
	if c.AMS >= AmsWarning {
		return s.AssignThreat(c, world.ThreatAMS, world.LevelWarning)
	}

	if s.State.WindSpeed >= 60 && c.ExposureTurns > 3 {
		if c.ExposureTurns > 6 {
			return s.AssignThreat(c, world.ThreatFrostbite, world.LevelCrisis)
		} else if c.ExposureTurns > 4 {
			return s.AssignThreat(c, world.ThreatFrostbite, world.LevelWarning)
		} else {
			return s.AssignThreat(c, world.ThreatFrostbite, world.LevelWhisper)
		}
	}

	if c.AMS >= AmsWhisper && s.State.Rng.Intn(2) == 0 {
		return &world.ActiveThreat{Type: world.ThreatAMS, Level: world.LevelWhisper}
	}

	return nil
}

func (s *Sim) ScanForSummitWindow() {
	consecutiveLowWind := 0
	windowFound := false
	startDay := 0

	for i := 0; i < 5; i++ {
		targetDay := (s.State.Day + i) % 31
		if targetDay == 0 {
			targetDay = 1
		}
		wind := s.State.WeatherCurve[targetDay]
		if wind < 40 {
			if consecutiveLowWind == 0 {
				startDay = s.State.Day + i
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
		s.LogLine(fmt.Sprintf("[METEO] Base reports stable window opening around Day %d.", startDay))
	}
}

func (s *Sim) AssignThreat(c *world.Climber, tType world.ThreatType, level world.ThreatLevel) *world.ActiveThreat {
	t, exists := c.ActiveThreats[tType]
	if !exists {
		t = &world.ActiveThreat{Type: tType}
		c.ActiveThreats[tType] = t
	}

	if t.Level != world.LevelCrisis {
		if level == world.LevelWarning && t.Level != world.LevelWarning {
			t.Level = world.LevelWarning
			t.TurnsLeft = 4
			t.WarningActed = false
		} else {
			t.Level = level
		}
	}

	if tType == world.ThreatFrostbite && t.Level == world.LevelCrisis {
		hasMod := false
		for _, m := range c.Modifiers {
			if m.Type == world.ModFrostbite {
				hasMod = true
			}
		}
		if !hasMod {
			c.Modifiers = append(c.Modifiers, world.Modifier{Type: world.ModFrostbite, MaxFitPenalty: 25})
			world.ApplyStatChange(c, "Fitness", 0, "Environment")
		}
	}

	return t
}

func (s *Sim) CheckTerminal(c *world.Climber) (bool, string) {
	if c.Fitness <= 0 {
		return true, "death"
	}
	if c.AMS >= 100 {
		return true, "ams_death"
	}
	if c.Loc == data.LocSummit && !c.Summited {
		c.Summited = true
		return false, "summit"
	}
	if c.Summited && c.Loc == data.LocBase {
		return true, "success"
	}
	return false, ""
}
