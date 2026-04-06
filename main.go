package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║            S U M M I T               ║")
	fmt.Println("║     A mountain survival game         ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Enter 1–4 to choose. 'q' to quit.")
	fmt.Println()

	// Climber state
	c := &Climber{
		Name:     "Zara",
		Fitness:  90,
		AMS:      0,
		Loc:      LocBase,
		Altitude: currentMountain.CampAltitudes[LocBase],
	}

	logLine(fmt.Sprintf("%s — %s. Ready to move. Clear skies.", c.Name, currentMountain.Name))

	for {
		// 1. CHECK TERMINAL (Win/Loss)
		// We check this at the START of the loop so the player sees the 
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
		// This happens before player input so they can respond to warnings.
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
