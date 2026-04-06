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

## Next phase — Bicycle (In Progress)

### What will be added
- [ ] **Architecture & Verifiability**: Extract `ApplyStatChange` loop to safely handle modifiers (like Meds, Psych status, Frostbite debuffs) without breaking core logic.
- [ ] **Automated Testing Suite**: `sim_test.go` to mathematically verify the simulation boundaries without playing an 8-hour session.
- [ ] Wind Speed System + Frostbite threat (exposure).
- [ ] Warning timers & Choice gating (removing options if warnings are ignored).
- [ ] Radio log persistence.


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
