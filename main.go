package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	LocBase    = 0
	LocCamp1   = 1
	LocCamp2   = 2
	LocHighCamp = 3
	LocSummit  = 4

	AmsWhisper  = 35
	AmsWarning  = 55
	AmsCrisis   = 80

	FitnessMin = 0
	FitnessMax = 100
	AmsMin     = 0
	AmsMax     = 100
)

var locName = map[int]string{
	LocBase:     "Base Camp",
	LocCamp1:    "Camp 1",
	LocCamp2:    "Camp 2",
	LocHighCamp: "High Camp",
	LocSummit:   "Summit",
}

// ─── Climber ──────────────────────────────────────────────────────────────────

type Climber struct {
	Name          string
	Fitness       int
	AMS           int
	Loc           int
	Summited      bool
	WarningActive bool
	WarningTurns  int // turns remaining before crisis auto-fires
	WarningActed  bool
	Resting       bool
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ─── Radio log ────────────────────────────────────────────────────────────────

var radioLog []string

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

// ─── Time / turn ──────────────────────────────────────────────────────────────

var (
	day  = 1
	hour = 6
	rng  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// amsGainByLoc: AMS grows faster the higher you are
var amsGainByLoc = map[int][2]int{
	LocBase:     {-8, -5}, // recovering (daily return ~40)
	LocCamp1:    {0, 2},
	LocCamp2:    {1, 3},
	LocHighCamp: {2, 5},
	LocSummit:   {3, 6}, // summit is brutal
}

func advanceTime(c *Climber) {
	if c.Loc > LocBase {
		// At altitude: lose fitness, gain AMS scaled by altitude
		fitLoss := 3 + rng.Intn(3) // 3–5 per 3hr turn
		if c.Resting {
			fitLoss = 1 // Resting significantly reduces decay
		}
		gainRange := amsGainByLoc[c.Loc]
		amsGain := gainRange[0] + rng.Intn(gainRange[1]-gainRange[0]+1)
		c.Fitness = clamp(c.Fitness-fitLoss, FitnessMin, FitnessMax)
		c.AMS = clamp(c.AMS+amsGain, AmsMin, AmsMax)
	} else {
		// At base: recover
		c.Fitness = clamp(c.Fitness+6, FitnessMin, FitnessMax)
		c.AMS = clamp(c.AMS-8, AmsMin, AmsMax)
	}

	c.Resting = false // Reset resting flag after time passes

	// Advance clock by 3 hours per turn
	hour += 3
	if hour >= 24 {
		hour = 0
		day++
	}

	// Tick warning timer if active
	if c.WarningActive {
		c.WarningTurns--
	}
}

// ─── Threat detection ─────────────────────────────────────────────────────────

// Returns: "none", "whisper", "warning", "crisis"
func checkThreats(c *Climber) string {
	// Crisis: auto-fire if warning timer expired
	if c.WarningActive && c.WarningTurns <= 0 {
		c.WarningActive = false
		return "crisis"
	}

	// Warning already active — don't re-trigger
	if c.WarningActive {
		return "none"
	}

	if c.AMS >= AmsCrisis {
		return "crisis"
	}
	if c.AMS >= AmsWarning {
		c.WarningActive = true
		c.WarningTurns = 4 // 4 turns (12 in-game hours) to act
		c.WarningActed = false
		return "warning"
	}
	if c.AMS >= 30 && rng.Intn(2) == 0 { // whisper: AMS≥30, 50% chance
		return "whisper"
	}
	return "none"
}

// ─── Status print ─────────────────────────────────────────────────────────────

func printStatus(c *Climber) {
	fmt.Println()
	fmt.Println("────────────────────────────────────────")
	fmt.Printf("  Day %d | %02d:00  |  %s  |  %s\n", day, hour, c.Name, locName[c.Loc])
	fmt.Printf("  Fitness: %-4d  AMS: %d\n", c.Fitness, c.AMS)
	if c.WarningActive {
		fmt.Printf("  ⚠  Warning active — %d turn(s) before crisis\n", c.WarningTurns)
	}
	fmt.Println("────────────────────────────────────────")
}

// ─── Choice display ───────────────────────────────────────────────────────────

type Choice struct {
	Label      string
	FitDelta   int
	AmsDelta   int
	LocDelta   int
	ClearsWarn bool
}

func printChoices(choices []Choice) {
	fmt.Println()
	for i, ch := range choices {
		fitStr := fmt.Sprintf("fitness %+d", ch.FitDelta)
		amsStr := fmt.Sprintf("AMS %+d", ch.AmsDelta)
		fmt.Printf("  [%d] %-32s  %s, %s\n", i+1, ch.Label, fitStr, amsStr)
	}
	fmt.Print("\n> ")
}

// ─── Apply choice ─────────────────────────────────────────────────────────────

func applyChoice(c *Climber, ch Choice) {
	c.Fitness = clamp(c.Fitness+ch.FitDelta, FitnessMin, FitnessMax)
	c.AMS = clamp(c.AMS+ch.AmsDelta, AmsMin, AmsMax)
	c.Loc = clamp(c.Loc+ch.LocDelta, LocBase, LocSummit)
	if ch.ClearsWarn {
		c.WarningActive = false
		c.WarningActed = true
	}
	if strings.Contains(ch.Label, "Rest") {
		c.Resting = true
	}
}

// ─── Event panels ─────────────────────────────────────────────────────────────

func whisperEvent(c *Climber) {
	// Whisper: no pause, no panel — just a subtle radio log blip
	lines := []string{
		fmt.Sprintf("%s — moving slower than usual today.", c.Name),
		fmt.Sprintf("%s — check-in at %02d:00. Short.", c.Name, hour),
		fmt.Sprintf("%s — fine.", c.Name),
	}
	logLine(lines[rng.Intn(len(lines))])
}

func warningEvent(c *Climber) string {
	logLine(fmt.Sprintf("!! %s — AMS elevated: %d. Headache, moving slow. Warning expires in %d turns.", c.Name, c.AMS, c.WarningTurns))

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
	// Gated by whether warning was acted on
	var choices []Choice
	if c.WarningActed {
		// Had warning, acted: still has all options
		choices = []Choice{
			{Label: "Emergency dex injection",  FitDelta: +5,  AmsDelta: -15, LocDelta: 0, ClearsWarn: false},
			{Label: "Immediate descent",        FitDelta: -15, AmsDelta: -20, LocDelta: -2, ClearsWarn: false},
			{Label: "Oxygen boost + rest",      FitDelta: +4,  AmsDelta: -12, LocDelta: 0, ClearsWarn: false},
			{Label: "Push through (desperate)", FitDelta: -20, AmsDelta: +8,  LocDelta: 0, ClearsWarn: false},
		}
	} else {
		// Ignored warning: worst options only
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

// ─── Normal turn ──────────────────────────────────────────────────────────────

func normalTurn(c *Climber) string {
	var choices []Choice

	if c.Loc == LocSummit {
		// At summit: no point advancing — only hold, rest, or descend
		choices = []Choice{
			{Label: "Begin descent",    FitDelta: -4,  AmsDelta: -4,  LocDelta: -1, ClearsWarn: false},
			{Label: "Rest before descending", FitDelta: +6, AmsDelta: +2, LocDelta: 0, ClearsWarn: false},
			{Label: "Hold position",   FitDelta: -2,  AmsDelta: +3,  LocDelta: 0, ClearsWarn: false},
		}
		printLog()
		printStatus(c)
		fmt.Println()
		fmt.Println("  ↓ You are at the summit. Fitness and AMS are working against you — descend.")
		fmt.Println("  Orders for " + c.Name + ":")
		printChoices(choices)
	} else {
		choices = []Choice{
			{Label: "Advance to next camp", FitDelta: -6,  AmsDelta: +2,  LocDelta: +1, ClearsWarn: false},
			{Label: "Hold position",        FitDelta: -2,  AmsDelta: +1,  LocDelta: 0, ClearsWarn: false},
			{Label: "Rest (3 hrs)",         FitDelta: +5,  AmsDelta: -3,  LocDelta: 0, ClearsWarn: false},
			{Label: "Descend one camp",     FitDelta: -4,  AmsDelta: -6,  LocDelta: -1, ClearsWarn: false},
		}
		printLog()
		printStatus(c)
		fmt.Println()
		fmt.Println("  Orders for " + c.Name + ":")
		printChoices(choices)
	}

	return readChoice(c, choices)
}

// ─── Input ────────────────────────────────────────────────────────────────────

var scanner = bufio.NewScanner(os.Stdin)

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

// ─── Win / lose ───────────────────────────────────────────────────────────────

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

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║            S U M M I T               ║")
	fmt.Println("║     A mountain survival game         ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Enter 1–4 to choose. 'q' to quit.")
	fmt.Println()

	c := &Climber{
		Name:    "Zara",
		Fitness: 90,
		AMS:     0,
		Loc:     LocBase,
	}

	logLine("Zara — Base Camp. Ready to move. Clear skies.")

	for {
		// 1. CHECK TERMINAL (Win/Loss)
		// We check this at the START of the loop so the player sees their 
		// final state before the game ends.
		done, outcome := checkTerminal(c)

		switch outcome {
		case "summit":
			fmt.Println()
			fmt.Println("  ★  SUMMIT REACHED — now get her home alive.")
			logLine("Zara — Summit. We made it. Descending now.")
		case "death":
			fmt.Println()
			fmt.Println("████████████████████████████████████████")
			fmt.Println("  ZARA IS DEAD")
			fmt.Println("  Fitness collapsed. She didn't make it down.")
			fmt.Println("████████████████████████████████████████")
		case "ams_death":
			fmt.Println()
			fmt.Println("████████████████████████████████████████")
			fmt.Println("  ZARA IS DEAD")
			fmt.Println("  Severe AMS. HACE. No recovery possible.")
			fmt.Println("████████████████████████████████████████")
		case "success":
			fmt.Println()
			fmt.Println("╔══════════════════════════════════════╗")
			fmt.Println("║  SUMMIT AND SAFE — all home alive.   ║")
			fmt.Println("║  Zara made it. Well led.             ║")
			fmt.Println("╚══════════════════════════════════════╝")
		}

		if done {
			fmt.Println()
			os.Exit(0)
		}

		// 2. CHECK THREATS (Environmental status)
		// This happens before the player acts so they can respond to warnings.
		threat := checkThreats(c)

		// 3. PLAYER INPUT (Get choice and apply it)
		var result string
		switch threat {
		case "whisper":
			whisperEvent(c)
			result = normalTurn(c)
		case "warning":
			result = warningEvent(c)
		case "crisis":
			result = crisisEvent(c)
		default:
			result = normalTurn(c)
		}

		if result == "quit" {
			fmt.Println("\n  Expedition abandoned.")
			os.Exit(0)
		}

		// 4. ADVANCE TIME (Environmental pressure)
		// This happens AFTER the action has been applied.
		advanceTime(c)
	}
}
