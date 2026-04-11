package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"summit/data"
	"summit/sim"
	"summit/world"
)

// ─── Shared UI State ─────────────────────────────────────────────────────────

var (
	Scanner = bufio.NewScanner(os.Stdin)
)

// ─── Rendering Methods ────────────────────────────────────────────────────────

func PrintLog(s *sim.Sim) {
	if len(s.Log) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(" RADIO LOG:")
	for _, line := range s.Log {
		fmt.Println(" " + line)
	}
}

func PrintStatus(s *sim.Sim, c *world.Climber) {
	fmt.Println()
	fmt.Println("╔" + strings.Repeat("═", 50) + "╗")
	timeStr := fmt.Sprintf("%02d:00", s.State.Hour)
	if s.State.Hour >= 18 || s.State.Hour < 6 {
		timeStr += " [NIGHT]"
	}

	fmt.Printf("║  %s — Day %d | %-13s | %-16s ║\n", s.Mountain.Name, s.State.Day, timeStr, c.Name)
	fmt.Println("╠" + strings.Repeat("─", 50) + "╢")

	status := s.Mountain.CampNames[c.Loc]
	if c.Altitude > s.Mountain.CampAltitudes[c.Loc] && c.Loc < data.LocSummit {
		status = fmt.Sprintf("Climbing (%dm remaining)", s.Mountain.CampAltitudes[c.Loc+1]-c.Altitude)
	}

	fmt.Printf("║  LOCATION: %-37s ║\n", status)
	altStr := fmt.Sprintf("%dm", c.Altitude)
	if c.Altitude >= s.Mountain.DeathZone {
		altStr += " [DEATH ZONE]"
	}
	fmt.Printf("║  ALTITUDE: %-37s ║\n", altStr)
	fmt.Println("╟" + strings.Repeat("─", 50) + "╢")

	o2Bar := strings.Repeat("█", c.O2Charges) + strings.Repeat("░", 3-c.O2Charges)
	fmt.Printf("║  FITNESS:  %-12d  AMS: %-19d ║\n", c.Fitness, c.AMS)
	fmt.Printf("║  OXYGEN:   %-12s  O2 Charges: %d/3        ║\n", o2Bar, c.O2Charges)

	hasWarning := false
	warnTurns := 99
	for _, t := range c.ActiveThreats {
		if t.Level == world.LevelWarning {
			hasWarning = true
			if t.TurnsLeft < warnTurns {
				warnTurns = t.TurnsLeft
			}
		}
	}
	if hasWarning {
		fmt.Printf("║  ⚠ WARNING: %-2d turns left before CRISIS     ║\n", warnTurns)
	}
	fmt.Println("╚" + strings.Repeat("═", 50) + "╝")

	PrintDebug(s, c)
	PrintForecast(s)
}

func PrintDebug(s *sim.Sim, c *world.Climber) {
	fmt.Println("\n  [ DEV HUD — INTERNAL SYSTEMS ]")
	fmt.Println("  " + strings.Repeat("┈", 48))
	
	envType := "Passive Base Recovery"
	if c.Loc > data.LocBase {
		envType = "Altitude Decay Active"
		if c.Altitude >= s.Mountain.DeathZone {
			envType += " (DEATH ZONE: Rest ineffective)"
		}
	}
	
	// Safety check for name pool to avoid panic if archetype is zero-value
	archName := "Unknown"
	if len(c.Archetype.NamePool) > 0 {
		archName = c.Archetype.NamePool[0]
	}
	
	fmt.Printf("  Expedition Seed: %d | Archetype: %s (Sus: %.1fx)\n", s.State.Seed, archName, c.Archetype.AMSSusPercent)
	fmt.Printf("  Mode:     %s\n", envType)

	if c.Loc > data.LocBase {
		min, max := data.AltitudeFitLossMin, data.AltitudeFitLossMax
		isNight := s.State.Hour >= 18 || s.State.Hour < 6
		if isNight {
			min += 2
			max += 2
			fmt.Printf("  Pressure: Fit -%d to -%d per turn (Night penalty)\n", min, max)
		} else {
			fmt.Printf("  Pressure: Fit -%d to -%d per turn\n", min, max)
		}
		if c.Resting {
			if c.Altitude >= s.Mountain.DeathZone {
				fmt.Println("  Resting:  Environment decay reduced slightly (-4)")
			} else {
				fmt.Println("  Resting:  Environment decay reduced to -1")
			}
		}
	} else {
		fmt.Printf("  Passive recovery: Fit +%d, AMS -%d\n", data.BaseRecoveryFit, data.BaseRecoveryAms)
	}
	
	fmt.Printf("  Weather:  Wind %d km/h | Exposure: %d turns\n", s.State.WindSpeed, c.ExposureTurns)
	
	fmt.Print("  Supplies: ")
	for i := data.LocCamp1; i <= data.LocHighCamp; i++ {
		fmt.Printf("%s: %d O2 | ", s.Mountain.CampNames[i], s.State.CampO2[i])
	}
	fmt.Println()
	fmt.Println("  " + strings.Repeat("┈", 48))
}

func PrintForecast(s *sim.Sim) {
	fmt.Println(" FORECAST (Next 4 days):")
	forecast := s.State.GetForecast(s.State.Day)
	for _, line := range forecast {
		fmt.Println("  → " + line)
	}
	fmt.Println()
}

func PrintChoices(choices []data.Choice) {
	fmt.Println()
	for i, ch := range choices {
		fitStr := fmt.Sprintf("fitness %+d", ch.FitDelta)
		amsStr := fmt.Sprintf("AMS %+d", ch.AmsDelta)
		fmt.Printf("  [%d] %-36s  %s, %s\n", i+1, ch.Label, fitStr, amsStr)
	}
	fmt.Print("\n> ")
}

// ─── Interaction Logic (Moved here to break sim <-> ui cycle) ───────────────

func RenderWhisper(s *sim.Sim, c *world.Climber, t *world.ActiveThreat) {
	var lines []string
	if t.Type == world.ThreatFrostbite {
		lines = []string{
			fmt.Sprintf("%s — hands getting numb.", c.Name),
			fmt.Sprintf("%s — wind is brutal up here.", c.Name),
			fmt.Sprintf("%s — shivering heavily.", c.Name),
		}
	} else {
		lines = []string{
			fmt.Sprintf("%s — moving slower than usual today.", c.Name),
			fmt.Sprintf("%s — check-in at %02d:00. Short.", c.Name, s.State.Hour),
			fmt.Sprintf("%s — fine.", c.Name),
		}
	}
	s.LogLine(lines[s.State.Rng.Intn(len(lines))])
}

func PromptWarning(s *sim.Sim, c *world.Climber, t *world.ActiveThreat) string {
	if t.Type == world.ThreatFrostbite {
		s.LogLine(fmt.Sprintf("!! %s — FROSTBITE RISK. High winds. Warning expires in %d turns.", c.Name, t.TurnsLeft))
	} else {
		s.LogLine(fmt.Sprintf("!! %s — AMS elevated: %d. Headache. Warning expires in %d turns.", c.Name, c.AMS, t.TurnsLeft))
	}

	restFitDelta, restAmsDelta := +8, -5
	if c.Altitude >= s.Mountain.DeathZone {
		restFitDelta, restAmsDelta = +2, -1
	}

	choices := []data.Choice{
		{Label: "Rest here (short)", FitDelta: restFitDelta, AmsDelta: restAmsDelta, LocDelta: 0, ClearsWarn: true},
		{Label: "Descend one camp", FitDelta: -5, AmsDelta: -10, LocDelta: -1, ClearsWarn: true},
		{Label: "Ignore it, keep moving", FitDelta: -8, AmsDelta: +6, LocDelta: 0, ClearsWarn: false},
	}

	PrintLog(s)
	PrintStatus(s, c)
	fmt.Println()
	if t.Type == world.ThreatFrostbite {
		fmt.Println("  !  FROSTBITE WARNING — prolonged exposure. Act soon.")
	} else {
		fmt.Println("  !  AMS WARNING — act within the timer or face a harder crisis")
	}
	PrintChoices(choices)
	return ReadChoice(s, c, choices)
}

func PromptCrisis(s *sim.Sim, c *world.Climber, t *world.ActiveThreat) string {
	var choices []data.Choice
	if t.Type == world.ThreatOxygen {
		choices = []data.Choice{
			{Label: "Improvise emergency descent", FitDelta: -20, AmsDelta: -10, LocDelta: -1, ClearsWarn: false},
			{Label: "Push without O2 (Massive AMS)", FitDelta: -10, AmsDelta: +30, LocDelta: 0, ClearsWarn: false},
			{Label: "Abort and stay (Dangerous)", FitDelta: -5, AmsDelta: +10, LocDelta: 0, ClearsWarn: false},
		}
	} else if t.WarningActed {
		choices = []data.Choice{
			{Label: "Emergency dex injection", FitDelta: +5, AmsDelta: -15, LocDelta: 0, ClearsWarn: false},
			{Label: "Immediate descent", FitDelta: -15, AmsDelta: -20, LocDelta: -2, ClearsWarn: false},
			{Label: "Oxygen boost + rest", FitDelta: +4, AmsDelta: -12, LocDelta: 0, ClearsWarn: false},
			{Label: "Push through (desperate)", FitDelta: -20, AmsDelta: +8, LocDelta: 0, ClearsWarn: false},
		}
	} else {
		choices = []data.Choice{
			{Label: "Immediate descent", FitDelta: -25, AmsDelta: -30, LocDelta: -2, ClearsWarn: false},
			{Label: "Push through (desperate)", FitDelta: -35, AmsDelta: +5, LocDelta: 0, ClearsWarn: false},
		}
	}

	PrintLog(s)
	PrintStatus(s, c)
	fmt.Println()
	if t.Type == world.ThreatFrostbite {
		fmt.Println("  !!  FROSTBITE CRISIS — permanent damage sustained (-25 Max Fitness).")
	} else if t.Type == world.ThreatOxygen {
		fmt.Println("  !!  OXYGEN CRISIS — you are out of supplemental air!")
	} else {
		fmt.Println("  !!  AMS CRISIS — situation critical")
	}

	if !t.WarningActed && t.Type != world.ThreatOxygen {
		fmt.Println("      (you ignored the warning — options are limited)")
	}

	PrintChoices(choices)
	return ReadChoice(s, c, choices)
}

func PromptNormalTurn(s *sim.Sim, c *world.Climber) string {
	restFitDelta, restAmsDelta := +5, -3
	if c.Altitude >= s.Mountain.DeathZone {
		restFitDelta, restAmsDelta = +1, 0
	}
	climbFitDelta := -10
	if s.State.Hour >= 18 || s.State.Hour < 6 {
		climbFitDelta = -15
	}

	var choices []data.Choice
	if c.Loc == data.LocSummit {
		choices = []data.Choice{
			{Label: "Begin descent", FitDelta: -4, AmsDelta: -4, LocDelta: -1, ClearsWarn: false},
			{Label: "Rest before descending", FitDelta: restFitDelta, AmsDelta: restAmsDelta + 2, LocDelta: 0, ClearsWarn: false},
			{Label: "Hold position", FitDelta: -2, AmsDelta: +3, LocDelta: 0, ClearsWarn: false},
		}
		PrintLog(s)
		PrintStatus(s, c)
		fmt.Println("\n  ↓ You are at the summit. Fitness and AMS are working against you — descend.")
	} else {
		advanceLabel := "Advance to next camp"
		if c.Altitude < s.Mountain.CampAltitudes[c.Loc+1] && c.Loc < data.LocSummit {
			advanceLabel = fmt.Sprintf("Continue climbing (%dm total)", s.Mountain.CampAltitudes[c.Loc+1]-c.Altitude)
		}
		choices = []data.Choice{
			{Label: advanceLabel, FitDelta: climbFitDelta, AmsDelta: +3, LocDelta: +1, ClearsWarn: false},
			{Label: "Hold position", FitDelta: -2, AmsDelta: +1, LocDelta: 0, ClearsWarn: false},
			{Label: "Rest (3 hrs)", FitDelta: restFitDelta, AmsDelta: restAmsDelta, LocDelta: 0, ClearsWarn: false},
			{Label: "Descend one camp", FitDelta: -4, AmsDelta: -6, LocDelta: -1, ClearsWarn: false},
		}
		PrintLog(s)
		PrintStatus(s, c)
	}

	fmt.Println("\n  Orders for " + c.Name + ":")
	PrintChoices(choices)
	return ReadChoice(s, c, choices)
}

func ReadChoice(s *sim.Sim, c *world.Climber, choices []data.Choice) string {
	for {
		Scanner.Scan()
		input := strings.TrimSpace(Scanner.Text())
		switch input {
		case "1", "2", "3", "4":
			idx := int(input[0] - '1')
			if idx < len(choices) {
				ApplyChoice(s, c, choices[idx])
				return "ok"
			}
		case "q", "quit":
			return "quit"
		}
		fmt.Print("  Enter a valid number > ")
	}
}

func ApplyChoice(s *sim.Sim, c *world.Climber, ch data.Choice) {
	world.ApplyStatChange(c, "Fitness", ch.FitDelta, "Choice")
	world.ApplyStatChange(c, "AMS", ch.AmsDelta, "Choice")
	
	if ch.LocDelta > 0 {
		climbGain := 300 + s.State.Rng.Intn(250)
		c.Altitude += climbGain
		target := s.Mountain.CampAltitudes[c.Loc+1]
		if c.Altitude >= target {
			c.Altitude = target
			c.Loc++
			s.LogLine(fmt.Sprintf("%s reached %s.", c.Name, s.Mountain.CampNames[c.Loc]))
		}
	} else if ch.LocDelta < 0 {
		c.Loc = world.Clamp(c.Loc+ch.LocDelta, data.LocBase, data.LocSummit)
		c.Altitude = s.Mountain.CampAltitudes[c.Loc]
		s.LogLine(fmt.Sprintf("%s descended to %s.", c.Name, s.Mountain.CampNames[c.Loc]))
	}

	if ch.ClearsWarn {
		for k, t := range c.ActiveThreats {
			if t.Level == world.LevelWarning {
				t.WarningActed = true
				delete(c.ActiveThreats, k)
			}
		}
	}
	if strings.Contains(ch.Label, "Rest") {
		c.Resting = true
	}
}
