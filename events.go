package main

import "fmt"

func whisperEvent(c *Climber) {
	// Whisper: Subtle log blips only
	lines := []string{
		fmt.Sprintf("%s — moving slower than usual today.", c.Name),
		fmt.Sprintf("%s — check-in at %02d:00. Short.", c.Name, hour),
		fmt.Sprintf("%s — fine.", c.Name),
	}
	logLine(lines[rng.Intn(len(lines))])
}

func warningEvent(c *Climber) string {
	logLine(fmt.Sprintf("!! %s — AMS elevated: %d. Headache. Warning expires in %d turns.", c.Name, c.AMS, c.WarningTurns))

	choices := []Choice{
		{Label: "Rest here (short)",      FitDelta: +8,  AmsDelta: -5,  LocDelta: 0, ClearsWarn: true},
		{Label: "Descend one camp",       FitDelta: -5,  AmsDelta: -10, LocDelta: -1, ClearsWarn: true},
		{Label: "Ignore it, keep moving", FitDelta: -8,  AmsDelta: +6,  LocDelta: 0, ClearsWarn: false},
	}

	printLog()
	printStatus(c)
	fmt.Println()
	fmt.Println("  !  AMS WARNING — act within the timer or face a harder crisis")
	printChoices(choices)

	return readChoice(c, choices)
}

func crisisEvent(c *Climber) string {
	var choices []Choice
	if c.WarningActed {
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
	fmt.Println("  !!  AMS CRISIS — situation critical")
	if !c.WarningActed {
		fmt.Println("      (you ignored the warning — options are limited)")
	}
	printChoices(choices)

	return readChoice(c, choices)
}

func normalTurn(c *Climber) string {
	var choices []Choice

	if c.Loc == LocSummit {
		choices = []Choice{
			{Label: "Begin descent",    FitDelta: -4,  AmsDelta: -4,  LocDelta: -1, ClearsWarn: false},
			{Label: "Rest before descending", FitDelta: +6, AmsDelta: +2, LocDelta: 0, ClearsWarn: false},
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
			{Label: advanceLabel,           FitDelta: -10, AmsDelta: +3,  LocDelta: +1, ClearsWarn: false},
			{Label: "Hold position",        FitDelta: -2,  AmsDelta: +1,  LocDelta: 0, ClearsWarn: false},
			{Label: "Rest (3 hrs)",         FitDelta: +5,  AmsDelta: -3,  LocDelta: 0, ClearsWarn: false},
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
