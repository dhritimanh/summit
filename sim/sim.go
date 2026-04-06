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

	// Environmental Systems
	BaseRecoveryFit    = 6
	BaseRecoveryAms    = 8
	AltitudeFitLossMin = 3
	AltitudeFitLossMax = 5
)

var AmsGainByLoc = map[int][2]int{
	data.LocBase:     {-8, -5},
	data.LocCamp1:    {0, 2},
	data.LocCamp2:    {1, 3},
	data.LocHighCamp: {2, 5},
	data.LocSummit:   {3, 6},
}

type Sim struct {
	State    *world.WorldState
	Mountain data.Mountain
	Team     []*world.Climber
	Log      []string
}

func NewSim(w *world.WorldState, m data.Mountain) *Sim {
	return &Sim{
		State:    w,
		Mountain: m,
		Team:     make([]*world.Climber, 0),
		Log:      make([]string, 0),
	}
}

func (s *Sim) LogLine(msg string) {
	s.Log = append(s.Log, msg)
	if len(s.Log) > 10 {
		s.Log = s.Log[1:]
	}
}

func (s *Sim) AdvanceTime(c *world.Climber) {
	if c.Loc > data.LocBase {
		fitLoss := AltitudeFitLossMin + s.State.Rng.Intn(AltitudeFitLossMax-AltitudeFitLossMin+1)
		if c.Resting {
			if c.Altitude >= s.Mountain.DeathZone {
				fitLoss = 4
			} else {
				fitLoss = 1
			}
		}

		isNight := s.State.Hour >= 18 || s.State.Hour < 6
		if isNight {
			fitLoss += 2
		}

		gainRange := AmsGainByLoc[c.Loc]
		amsGain := float64(gainRange[0] + s.State.Rng.Intn(gainRange[1]-gainRange[0]+1))
		
		// Use individual climber's archetype for specialized susceptibility
		if amsGain > 0 {
			amsGain *= c.Archetype.AMSSusPercent
		}
		
		if isNight {
			amsGain += 1
		}

		world.ApplyStatChange(c, "Fitness", -fitLoss, "Environment")
		world.ApplyStatChange(c, "AMS", int(amsGain), "Environment")
	} else {
		world.ApplyStatChange(c, "Fitness", BaseRecoveryFit, "Environment")
		world.ApplyStatChange(c, "AMS", -BaseRecoveryAms, "Environment")
	}

	c.Resting = false
	s.ConsumeO2(c)

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

	baseWind := s.State.WeatherCurve[s.State.Day%31]
	if c.Altitude < 6000 {
		s.State.WindSpeed = baseWind / 2
	} else if c.Altitude < 8000 {
		s.State.WindSpeed = baseWind
	} else {
		s.State.WindSpeed = baseWind + 30
	}

	if s.State.WindSpeed >= 60 {
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
	if c.O2Charges > 0 {
		c.O2Charges--
	}
	if c.O2Charges == 0 && s.State.CampO2[c.Loc] > 0 {
		s.State.CampO2[c.Loc]--
		c.O2Charges = 3
		s.LogLine(fmt.Sprintf("%s loaded a fresh O2 bottle at %s.", c.Name, s.Mountain.CampNames[c.Loc]))
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
