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

### What is NOT in this phase (by design)
- No graphics
- No multiple climbers
- No resources (O2, food, fuel)
- No seed system
- No weather
- Stats shown explicitly (fitness delta on every choice) — intentional for mechanical verification

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

## Phase 3 — Scooter (Complete)
**Built:** Deterministic Seed-based RNG, Oxygen "Charges" system, Auto-reloading resources, Climber Archetypes, and a 4-day Weather Forecast engine.
**Status:** 100% Complete. Modular structure improved for Phase 4 scaling. Verified by unit tests.

### What has been built
- [x] **Deterministic Seed**: Centralized `WorldState` initialized via `NewWorldState(seed)`.
- [x] **Oxygen Mechanics**: "Charges" system (3 per bottle) with auto-reload from camp supplies.
- [x] **Archetypes**: `ClimberArchetype` system affecting base fitness and AMS susceptibility.
- [x] **Weather Forecast UI**: 4-day interactive panel showing upcoming wind speeds and probability ranges.
- [x] **O2 Crisis**: High-stakes crisis event when oxygen is depleted at altitude.
- [x] **Intelligent Modularization**: Logical file grouping (`engine_`, `data_`, `entity_`) to support Phase 4 expansion.
- [x] **Summit Window**: Automated meteorology scan for stable low-wind periods.

---

## Feedback log

| # | Phase | Note | Status |
|---|---|---|---|
| 1 | Skateboard | **BUG**: "Advance" at Summit wasted turns — loc clamped at 4, fitness drained, player died stuck at top | ✅ Fixed — Summit now shows Begin Descent / Rest / Hold only |
| 2 | Skateboard | **BUG**: AMS gain was flat 2–8/turn regardless of altitude — threat system barely fired in a fast ascent | ✅ Fixed — AMS now scales by altitude, Summit is most punishing |
| 3 | Skateboard | **BUG**: Whisper threshold (AMS≥35, 33% chance) too conservative — never appeared in playtest | ✅ Fixed — threshold lowered to AMS≥30, probability raised to 50% |
| 4 | Skateboard | **BUG**: Resting at High Camp/Summit resulted in net fitness loss because environmental decay (6–10) exceeded rest bonus (+6/8) | ✅ Fixed — 'Rest' actions now halve environmental decay for that turn, and base bonus was buffed. |
| 5 | Skateboard | **CRITICAL LOGIC BUG**: Environmental decay happened *before* player input. If fitness was low, player died "at the start of the turn" without a chance to Rest. | ✅ Fixed — Reordered loop: State Check -> Player Input -> Decay. Action results are processed *before* the environment takes its toll. |
| 6 | Skateboard | **DEATH ZONE BUG**: Resting at the summit yielded a net positive fitness gain, allowing players to sleep infinitely at 8800m. | ✅ Fixed — Added Death Zone (8000m+) constraint to `sim.go` where resting recovery is neutralized and ambient pressure remains lethal. |
| 7 | Skateboard | **LOGIC/REALISM BUG**: Climbing at night (or 24/7 uninterrupted) carried no penalty. | ✅ Fixed — Added Nighttime Penalty (18:00 to 06:00). Ambient decay increases by 2, and climbing costs -15 fitness instead of -10. |
