---

## Phase 4 — Car ✅
**Built:** Ebitengine Graphical Transition, Real-Time Clock, Persistent Orders.
**Status:** Build passes. Interactive GUI is fully functional.

### What was built
- [x] **Graphical UI**: Windowed application using Ebitengine (960x720).
- [x] **Real-Time Simulation**: 60 FPS update loop with `AccumulatedTime` logic.
- [x] **Speed Control**: Manual `Pause`, `1x`, `4x`, and `12x` buttons (and hotkeys 1-3).
- [x] **Persistent Commands**: Climbers now follow `Orders` (Climb/Rest/Hold) instead of turn-by-turn prompts.
- [x] **Manual Oxygen**: O2 is now an assigned toggle; provides 50% mitigation to stat decay.
- [x] **Ambient Visuals**: Day/Night skybox cycle and "Panic Pulse" heart-rate indicators for AMS levels.
- [x] **Interactive HUD**: Detailed side-panel for selected climbers with non-overlapping zone layout.
- [x] **Three-Zone Dashboard**: Split screen into Mountain View (60%), Command Bar (Bottom), and Intel Sidebar (Right).
- [x] **Visual Vitals**: Real-time progress bars for Fitness and AMS, and iconic oxygen tank status.
- [x] **Multi-Climber Ready**: Interior architecture updated to handle team slices.

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
