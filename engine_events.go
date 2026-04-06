package main

import "fmt"

func whisperEvent(c *Climber, t *ActiveThreat) {
	// Whisper: Subtle log blips only. We can vary text based on type.
	var lines []string
	if t.Type == ThreatFrostbite {
		lines = []string{
			fmt.Sprintf("%s — hands getting numb.", c.Name),
			fmt.Sprintf("%s — wind is brutal up here.", c.Name),
			fmt.Sprintf("%s — shivering heavily.", c.Name),
		}
	} else {
		lines = []string{
			fmt.Sprintf("%s — moving slower than usual today.", c.Name),
			fmt.Sprintf("%s — check-in at %02d:00. Short.", c.Name, session.Hour),
			fmt.Sprintf("%s — fine.", c.Name),
		}
	}

	logLine(lines[session.Rng.Intn(len(lines))])
}


func warningEvent(c *Climber, t *ActiveThreat) string {
	if t.Type == ThreatFrostbite {
		logLine(fmt.Sprintf("!! %s — FROSTBITE RISK. High winds. Warning expires in %d turns.", c.Name, t.TurnsLeft))
	} else {
		logLine(fmt.Sprintf("!! %s — AMS elevated: %d. Headache. Warning expires in %d turns.", c.Name, c.AMS, t.TurnsLeft))
	}

	restFitDelta := +8
	restAmsDelta := -5
	if c.Altitude >= currentMountain.DeathZone {
		restFitDelta = +2
		restAmsDelta = -1
	}

	choices := []Choice{
		{Label: "Rest here (short)",      FitDelta: restFitDelta,  AmsDelta: restAmsDelta,  LocDelta: 0, ClearsWarn: true},
		{Label: "Descend one camp",       FitDelta: -5,  AmsDelta: -10, LocDelta: -1, ClearsWarn: true},
		{Label: "Ignore it, keep moving", FitDelta: -8,  AmsDelta: +6,  LocDelta: 0, ClearsWarn: false},
	}


	printLog()
	printStatus(c)
	fmt.Println()
	if t.Type == ThreatFrostbite {
		fmt.Println("  !  FROSTBITE WARNING — prolonged exposure. Act soon.")
	} else {
		fmt.Println("  !  AMS WARNING — act within the timer or face a harder crisis")
	}
	printChoices(choices)

	return readChoice(c, choices)
}

func crisisEvent(c *Climber, t *ActiveThreat) string {
	var choices []Choice
	if t.Type == ThreatOxygen {
		// Oxygen depletion crisis
		choices = []Choice{
			{Label: "Improvise emergency descent", FitDelta: -20, AmsDelta: -10, LocDelta: -1, ClearsWarn: false},
			{Label: "Push without O2 (Massive AMS)", FitDelta: -10, AmsDelta: +30,  LocDelta: 0,  ClearsWarn: false},
			{Label: "Abort and stay (Dangerous)", FitDelta: -5,  AmsDelta: +10,  LocDelta: 0,  ClearsWarn: false},
		}
	} else if t.WarningActed {
		// Player was responsible: better recovery options

		choices = []Choice{
			{Label: "Emergency dex injection",  FitDelta: +5,  AmsDelta: -15, LocDelta: 0, ClearsWarn: false},
			{Label: "Immediate descent",        FitDelta: -15, AmsDelta: -20, LocDelta: -2, ClearsWarn: false},
			{Label: "Oxygen boost + rest",      FitDelta: +4,  AmsDelta: -12, LocDelta: 0, ClearsWarn: false},
			{Label: "Push through (desperate)", FitDelta: -20, AmsDelta: +8,  LocDelta: 0, ClearsWarn: false},
		}
	} else {
		// Player ignored warning: punishment scenario
		choices = []Choice{
			{Label: "Immediate descent",        FitDelta: -25, AmsDelta: -30, LocDelta: -2, ClearsWarn: false},
			{Label: "Push through (desperate)", FitDelta: -35, AmsDelta: +5,  LocDelta: 0, ClearsWarn: false},
		}
	}

	printLog()
	printStatus(c)
	fmt.Println()
	if t.Type == ThreatFrostbite {
		fmt.Println("  !!  FROSTBITE CRISIS — permanent damage sustained (-25 Max Fitness).")
	} else if t.Type == ThreatOxygen {
		fmt.Println("  !!  OXYGEN CRISIS — you are out of supplemental air!")
	} else {
		fmt.Println("  !!  AMS CRISIS — situation critical")
	}
	
	if !t.WarningActed && t.Type != ThreatOxygen {
		fmt.Println("      (you ignored the warning — options are limited)")
	}

	printChoices(choices)

	return readChoice(c, choices)
}

func normalTurn(c *Climber) string {
	var choices []Choice

	isNight := session.Hour >= 18 || session.Hour < 6
	inDeathZone := c.Altitude >= currentMountain.DeathZone

	restFitDelta := +5
	restAmsDelta := -3
	if inDeathZone {
		restFitDelta = +1
		restAmsDelta = 0
	}


	climbFitDelta := -10
	if isNight {
		climbFitDelta = -15
	}

	if c.Loc == LocSummit {
		choices = []Choice{
			{Label: "Begin descent",    FitDelta: -4,  AmsDelta: -4,  LocDelta: -1, ClearsWarn: false},
			{Label: "Rest before descending", FitDelta: restFitDelta, AmsDelta: restAmsDelta + 2, LocDelta: 0, ClearsWarn: false},
			{Label: "Hold position",   FitDelta: -2,  AmsDelta: +3,  LocDelta: 0, ClearsWarn: false},
		}
		printLog()
		printStatus(c)
		fmt.Println()
		fmt.Println("  ↓ You are at the summit. Fitness and AMS are working against you — descend.")
	} else {
		advanceLabel := "Advance to next camp"
		if c.Altitude < currentMountain.CampAltitudes[c.Loc+1] && c.Loc < LocSummit {
			advanceLabel = fmt.Sprintf("Continue climbing (%dm total)", currentMountain.CampAltitudes[c.Loc+1]-c.Altitude)
		}
		choices = []Choice{
			{Label: advanceLabel,           FitDelta: climbFitDelta, AmsDelta: +3,  LocDelta: +1, ClearsWarn: false},
			{Label: "Hold position",        FitDelta: -2,  AmsDelta: +1,  LocDelta: 0, ClearsWarn: false},
			{Label: "Rest (3 hrs)",         FitDelta: restFitDelta,  AmsDelta: restAmsDelta,  LocDelta: 0, ClearsWarn: false},
			{Label: "Descend one camp",     FitDelta: -4,  AmsDelta: -6,  LocDelta: -1, ClearsWarn: false},
		}

		printLog()
		printStatus(c)
	}

	fmt.Println()
	fmt.Println("  Orders for " + c.Name + ":")
	printChoices(choices)

	return readChoice(c, choices)
}
