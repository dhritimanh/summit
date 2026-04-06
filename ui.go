package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ─── Shared UI State ─────────────────────────────────────────────────────────

var (
	radioLog []string
	scanner  = bufio.NewScanner(os.Stdin)
)

// ─── Methods ──────────────────────────────────────────────────────────────────

func logLine(msg string) {
	t := time.Now()
	entry := fmt.Sprintf("[%02d:%02d] %s", t.Hour(), t.Minute(), msg)
	radioLog = append(radioLog, entry)
	if len(radioLog) > 12 {
		radioLog = radioLog[len(radioLog)-12:]
	}
}

func printLog() {
	if len(radioLog) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(" RADIO LOG:")
	for _, line := range radioLog {
		fmt.Println(" " + line)
	}
}

func printStatus(c *Climber) {
	fmt.Println()
	fmt.Println("────────────────────────────────────────")
	fmt.Printf("  Peak: %s\n", currentMountain.Name)
	status := currentMountain.CampNames[c.Loc]
	if c.Altitude > currentMountain.CampAltitudes[c.Loc] && c.Loc < LocSummit {
		status = fmt.Sprintf("Climbing towards %s", currentMountain.CampNames[c.Loc+1])
	}
	fmt.Printf("  Day %d | %02d:00  |  %s  |  %s\n", day, hour, c.Name, status)
	altStr := fmt.Sprintf("%dm", c.Altitude)
	if c.Altitude >= currentMountain.DeathZone {
		altStr += " [DEATH ZONE]"
	}
	fmt.Printf("  Altitude: %-12s Fitness: %-4d  AMS: %d\n", altStr, c.Fitness, c.AMS)
	if c.WarningActive {
		fmt.Printf("  ⚠  Warning active — %d turn(s) before crisis\n", c.WarningTurns)
	}
	fmt.Println("────────────────────────────────────────")
}

func printChoices(choices []Choice) {
	fmt.Println()
	for i, ch := range choices {
		fitStr := fmt.Sprintf("fitness %+d", ch.FitDelta)
		amsStr := fmt.Sprintf("AMS %+d", ch.AmsDelta)
		fmt.Printf("  [%d] %-36s  %s, %s\n", i+1, ch.Label, fitStr, amsStr)
	}
	fmt.Print("\n> ")
}

func readChoice(c *Climber, choices []Choice) string {
	for {
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "1":
			if len(choices) >= 1 {
				applyChoice(c, choices[0])
				return "ok"
			}
		case "2":
			if len(choices) >= 2 {
				applyChoice(c, choices[1])
				return "ok"
			}
		case "3":
			if len(choices) >= 3 {
				applyChoice(c, choices[2])
				return "ok"
			}
		case "4":
			if len(choices) >= 4 {
				applyChoice(c, choices[3])
				return "ok"
			}
		case "q", "quit":
			return "quit"
		}
		fmt.Print("  Enter a valid number > ")
	}
}

func applyChoice(c *Climber, ch Choice) {
	c.Fitness = clamp(c.Fitness+ch.FitDelta, FitnessMin, FitnessMax)
	c.AMS = clamp(c.AMS+ch.AmsDelta, AmsMin, AmsMax)
	
	// Complex Location/Height Logic
	if ch.LocDelta > 0 {
		// Climb Logic (~300m - 550m in 3 hours)
		climbGain := 300 + rng.Intn(250)
		c.Altitude += climbGain
		target := currentMountain.CampAltitudes[c.Loc+1]
		if c.Altitude >= target {
			c.Altitude = target
			c.Loc++
			logLine(fmt.Sprintf("%s reached %s.", c.Name, currentMountain.CampNames[c.Loc]))
		}
	} else if ch.LocDelta < 0 {
		// Descent Logic (Much faster gravity-assisted movement)
		c.Loc = clamp(c.Loc+ch.LocDelta, LocBase, LocSummit)
		c.Altitude = currentMountain.CampAltitudes[c.Loc]
		logLine(fmt.Sprintf("%s descended to %s.", c.Name, currentMountain.CampNames[c.Loc]))
	}

	if ch.ClearsWarn {
		c.WarningActive = false
		c.WarningActed = true
	}
	if strings.Contains(ch.Label, "Rest") {
		c.Resting = true
	}
}
