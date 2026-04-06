package main

import (
	"fmt"
	"os"
	"summit/data"
	"summit/rng"
	"summit/sim"
	"summit/ui"
	"summit/world"
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

	// Initialize Sim with new WorldState and default Mountain
	w := world.NewWorldState(0)
	s := sim.NewSim(w, data.Everest)

	// Climber generation via RNG package
	c := rng.GenerateClimber(w)
	c.Altitude = s.Mountain.CampAltitudes[data.LocBase]
	s.Team = append(s.Team, c)

	fmt.Printf(" [ SESSION SEED: %d ]\n", w.Seed)
	s.LogLine(fmt.Sprintf("%s — %s. Ready to move. Clear skies.", c.Name, s.Mountain.Name))

	for {
		// 1. CHECK TERMINAL (Win/Loss)
		done, outcome := s.CheckTerminal(c)

		switch outcome {
		case "summit":
			fmt.Println("\n  ★  SUMMIT REACHED — now get her home alive.")
			s.LogLine(fmt.Sprintf("%s — Summit. We made it. Descending now.", c.Name))
		case "death":
			fmt.Println("\n████████████████████████████████████████")
			fmt.Printf("  %s IS DEAD\n", c.Name)
			fmt.Println("  Fitness collapsed. She didn't make it down.")
			fmt.Println("████████████████████████████████████████")
		case "ams_death":
			fmt.Println("\n██████████████████════██████████████████")
			fmt.Printf("  %s IS DEAD\n", c.Name)
			fmt.Println("  Severe AMS. HACE. No recovery possible.")
			fmt.Println("████████████████████████████████████████")
		case "success":
			fmt.Println("\n╔══════════════════════════════════════╗")
			fmt.Println("║  SUMMIT AND SAFE — all home alive.   ║")
			fmt.Printf("║  %s made it. Well led.             ║\n", c.Name)
			fmt.Println("╚══════════════════════════════════════╝")
		}

		if done {
			fmt.Println()
			os.Exit(0)
		}

		// 2. CHECK THREATS
		threat := s.CheckThreats(c)

		// 3. PLAYER INPUT
		var result string
		if threat == nil || threat.Level == world.LevelNone {
			result = ui.PromptNormalTurn(s, c)
		} else {
			switch threat.Level {
			case world.LevelWhisper:
				ui.RenderWhisper(s, c, threat)
				result = ui.PromptNormalTurn(s, c)
			case world.LevelWarning:
				result = ui.PromptWarning(s, c, threat)
			case world.LevelCrisis:
				result = ui.PromptCrisis(s, c, threat)
			}
		}

		if result == "quit" {
			fmt.Println("\n  Expedition abandoned.")
			os.Exit(0)
		}

		// 4. ADVANCE TIME
		s.AdvanceTime(c)
	}
}
