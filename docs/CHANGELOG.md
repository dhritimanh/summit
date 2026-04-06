# Summit — Changelog

---

## Phase 1 — Skateboard ✅
**Built:** `main.go` (single file, ~280 LOC)
**Status:** Build passes. Ready for playtesting.

### What was built
- `Climber` struct: `Name`, `Fitness` (0–100), `AMS` (0–100), `Loc` (Base to Summit), `Altitude` (meters)
- **Altitude System**: Everest-scale heights (Base: 5364m, Summit: 8848m)
  - `Advance` now gains ~400–600m per turn. 
  - Reaching a camp takes multiple turns of climbing (e.g., C2 to C4 is 1500m / 3–4 turns).
  - **DEATH ZONE**: Visual marker at 8000m+ where recovery is minimal.
- **Turn Granularity**: 1 turn = 3 hours (8 turns/day).
- `advanceTime()`: Fitness decay and AMS gain scale by both Turn Duration and Altitude.
- `checkThreats()`: Three-level threat system (Whisper / Warning / Crisis) with timed responses.
- `applyChoice()`: Processes climbing/descent/rest logic.
- Radio log: Shared channel with timestamped check-ins.
- Win/Loss: Summit success must return to Base alive. Fitness 0 = death.

---

## Phase 2 — Bicycle ✅
**Built:** Extracted robust simulation architecture, Multi-Threat map parsing, Global Weather State
**Status:** Build passes perfectly on `sim_test.go` suite. Ready for playtesting.

### What was built
- [x] **Architecture & Verifiability**: Abstracted stat updates into `ApplyStatChange` loop to safely handle modifiers (Caps, Buffs, Debuffs) without breaking core state math.
- [x] **Automated Testing Suite**: `sim_test.go` math framework added. Verified Base Recovery, Night Penalty offsets, and Death Zone math.
- [x] **New Thread Models**: `threat.go` created to track `ThreatType` and `ThreatLevel` dynamically.
- [x] **Wind Speed & Frostbite**: `windSpeed` implemented in `sim.go`. Exceeding 60km/h adds `ExposureTurns`. Over 6 turns results in a Frostbite Crisis which permanently cuts `-25 Fitness`.
- [x] **Choice Gating Scale**: Updated Warning timers and actions to tie to a dictionary map of threats `c.ActiveThreats` instead of purely hardcoded string flags.
- [x] **Radio log persistence**: Completed via `ui.go` rendering logic.

---

## Phase 3 — Scooter ✅
**Built:** Deterministic Seed-based RNG, Oxygen "Charges" system, Auto-reloading resources, Climber Archetypes, and a 4-day Weather Forecast engine.
**Status:** 100% Complete. Modular structure improved for Phase 4 scaling. Verified by unit tests.

### What has been built
- [x] **Deterministic Seed**: Centralized `WorldState` initialized via `NewWorldState(seed)`. Every run is repeatable.
- [x] **Oxygen Mechanics**: "Charges" system (3 per bottle) with auto-reload from camp supplies. Mid-route depletion triggers a Crisis.
- [x] **Archetypes**: `ClimberArchetype` system affecting base fitness and AMS susceptibility (e.g., Veteran vs Journalist).
- [x] **Weather Forecast Engine**: 31-day atmospheric curve. 4-day forecast with confidence degradation.
- [x] **Summit Window**: Automated meteorology scan detection.
- [x] **Rocket Architecture**: Codebase refactored into modular packages (`sim`, `ui`, `world`, `rng`, `data`) to prevent circular dependencies and prepare for Ebitengine (Phase 4).

---

## Architecture Milestone: Rocket Alignment 🚀
**Built:** Full package-level separation to support graphical UI integration.
**Status:** Stable and verified.

### Key Changes
- **Cycle-Free Packages**: Resolved the `sim` <-> `ui` circular dependency by moving all interaction logic into `ui/ui.go`. 
- **Orchestrator Pattern**: Introduced the `Sim` struct as the central authority for simulation physics, making the engine "Headless" and ready for GUI attach.
- **Improved Data Safety**: Added safety checks for zero-value archetypes and internalized stat calculation on the climber level.

---

## Feedback log

| # | Phase | Note | Status |
|---|---|---|---|
| 5 | Skateboard | **CRITICAL LOGIC BUG**: Environmental decay happened *before* player input. If fitness was low, player died "at the start of the turn" without a chance to Rest. | ✅ Fixed |
| 6 | Skateboard | **DEATH ZONE BUG**: Resting at the summit yielded a net positive fitness gain. | ✅ Fixed |
| 7 | Skateboard | **LOGIC/REALISM BUG**: Climbing at night carried no penalty. | ✅ Fixed |
| 8 | Refactor | **PANIC**: index out of range [0] in PrintDebug. | ✅ Fixed (Archetype initialization safety) |
