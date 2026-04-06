# Summit — Implementation Plan

Build a turn-based survival management game where the player leads a rope team up a mountain, reading subtle signals and making decisions under incomplete information. The game follows a 6-phase incremental build (Skateboard → Rocket), each phase playable and testable before the next begins.

## Technology Recommendation

> [!IMPORTANT]
> **Recommendation: Go + Ebitengine**
>
> Both stacks are viable, but Go/Ebitengine is the stronger choice for this specific game:

| Factor | Go / Ebitengine | Lua / LÖVE2D |
|---|---|---|
| **Data-driven architecture** | Go structs + interfaces map directly to the design's `World`, `FilteredWorld`, `ThreatRegistry`, `DominoTable` patterns. Type safety catches mismatches at compile time. | Lua tables are flexible but offer no compile-time safety for the complex nested state this game requires. |
| **Serialization (save/load)** | `encoding/json` works out of the box on typed structs. The design explicitly calls for JSON serialization. | Requires manual serialization or a library. More boilerplate for deeply nested state. |
| **Seed-based determinism** | `math/rand` with `NewSource(seed)` is well-defined and portable. The design doc literally uses this API. | `love.math.newRandomGenerator(seed)` works but Lua's floating point can introduce subtle cross-platform drift. |
| **Audio** | Ebitengine has built-in audio with pitch control, looping, and volume — all needed for the Balatro-style feedback and ambient drone. | LÖVE2D audio is excellent and arguably easier. Slight edge to LÖVE here. |
| **Code organization** | Go packages (`world/`, `threat/`, `catalyst/`, `ui/`) map 1:1 to the design doc's file structure. | Lua modules work but require more discipline to maintain clean boundaries at this scale. |
| **The design doc itself** | Written in Go idioms (`ClimberState`, `ApplyChoice()`, `rand.New(rand.NewSource(seed))`). Following it in Go means zero translation overhead. | Every code reference needs mental translation from Go to Lua. |
| **Cross-platform distribution** | Single binary. `go build` and ship. | Requires LÖVE runtime bundled or installed separately. |

**The design doc was written for Go.** Every struct name, every function signature, every package layout is Go-native. Building in Lua would mean constantly translating the spec. Go/Ebitengine is the path of least friction.

> [!NOTE]
> If you prefer Lua/LÖVE2D for other reasons (faster prototyping, personal familiarity, modding support), the plan below can be adapted. Just say the word.

---

## Proposed Changes

The build follows the design doc's 6 phases exactly. Each phase is self-contained — playable and testable before moving to the next.

---

### Phase 1: Skateboard — One climber, terminal only (Tasks 1–6)

**Goal:** Does one decision feel consequential?

#### [NEW] [main.go](file:///Users/raynwow/m_projects/anti/summit/main.go)
- `RunGame()` loop: advance time → check threats → print status → read stdin → apply choice
- `AdvanceTime()`: altitude-based fitness decay and AMS growth
- `CheckThreats()`: AMS threshold checks → whisper / warning / crisis
- `ApplyChoice()`: apply stat deltas from selected option
- `PrintStatus()`: formatted terminal output block
- Win/lose checks: fitness ≤ 0 = death, loc == summit then back to base = success

**Single file. No packages. Pure terminal I/O.**

---

### Phase 2: Bicycle — Three climbers, warning timers (Tasks 7–12)

**Goal:** Does managing a team feel different from managing one person?

#### [NEW] [climber.go](file:///Users/raynwow/m_projects/anti/summit/climber.go)
- `ClimberState` struct extracted from main
- Array of 3 climbers, loop through each per turn

#### [NEW] [threat.go](file:///Users/raynwow/m_projects/anti/summit/threat.go)
- `ThreatType` and `ThreatLevel` enums
- Second threat: frostbite (wind speed + exposure turns)
- Warning timer: `turnsRemaining int`, decrements each turn, auto-escalates to crisis at 0
- Choice gating: `warningActedOn bool` — if false at crisis time, best option removed

#### [MODIFY] [main.go](file:///Users/raynwow/m_projects/anti/summit/main.go)
- Team status line: one row per climber with `[OK]`/`[!]`/`[!!]` tags
- Radio log as persistent list (last 10 lines, printed above prompt)

---

### Phase 3: Scooter — Resources, seeds, weather (Tasks 13–18)

**Goal:** Does each run feel different enough to play again?

#### [NEW] [world/resource.go](file:///Users/raynwow/m_projects/anti/summit/world/resource.go)
- O2 per camp array, food, fuel, medicine, money
- O2 consumption: 1 cylinder/turn at altitude when assigned

#### [NEW] [world/weather.go](file:///Users/raynwow/m_projects/anti/summit/world/weather.go)
- 30-day pressure curve, wind speed, forecast (real + degraded)

#### [NEW] [rng/seed.go](file:///Users/raynwow/m_projects/anti/summit/rng/seed.go)
- `RunSeed` struct wrapping `rand.Rand`, single source of truth

#### [NEW] [rng/roster.go](file:///Users/raynwow/m_projects/anti/summit/rng/roster.go)
- Generate 3 climbers from archetype pool with stat ranges and hidden traits

#### [NEW] [rng/weather_seed.go](file:///Users/raynwow/m_projects/anti/summit/rng/weather_seed.go)
- Generate weather curve from seed; summit window position varies per seed

#### [NEW] [rng/event_deck.go](file:///Users/raynwow/m_projects/anti/summit/rng/event_deck.go)
- Shuffle threat firing order from seed

#### [NEW] [data/climber_archetypes.go](file:///Users/raynwow/m_projects/anti/summit/data/climber_archetypes.go)
- Archetype pool: name pools, stat ranges, hidden traits (e.g., AMS susceptibility = 1.5x tick rate)

#### [MODIFY] [main.go](file:///Users/raynwow/m_projects/anti/summit/main.go)
- Oxygen crisis event with gated choices
- Seed printed at start and end of run

---

### Phase 4: Car — First graphical stage with Ebitengine (Tasks 19–26)

**Goal:** Do the dots on the mountain communicate danger before the player reads any numbers?

#### [NEW] [ui/ui.go](file:///Users/raynwow/m_projects/anti/summit/ui/ui.go)
- Implements `ebiten.Game` interface (`Update()`, `Draw()`, `Layout()`)
- Owns all sub-panels, routes input

#### [NEW] [ui/map_view.go](file:///Users/raynwow/m_projects/anti/summit/ui/map_view.go)
- Mountain polygon silhouette, fixed camp altitude positions
- Climber dots: color by AMS threshold (white/amber/red), pulse animation via `sin(time * pulseSpeed)`

#### [NEW] [ui/climber_panel.go](file:///Users/raynwow/m_projects/anti/summit/ui/climber_panel.go)
- Opens on dot click: name, stats, location, action buttons

#### [NEW] [ui/event_panel.go](file:///Users/raynwow/m_projects/anti/summit/ui/event_panel.go)
- Crisis overlay: dim background, slide-in animation, countdown timer, choice buttons

#### [NEW] [ui/radio_log.go](file:///Users/raynwow/m_projects/anti/summit/ui/radio_log.go)
- Scrolling monospace log panel, last N lines visible

#### [NEW] [ui/resource_panel.go](file:///Users/raynwow/m_projects/anti/summit/ui/resource_panel.go)
- O2 cylinder icons per camp, food/fuel/medicine counters

#### [NEW] [ui/hud.go](file:///Users/raynwow/m_projects/anti/summit/ui/hud.go)
- Speed controls (Pause/1x/4x/12x), day/time display

#### [NEW] [sim/sim.go](file:///Users/raynwow/m_projects/anti/summit/sim/sim.go)
- `Sim` struct: owns World, ThreatRegistry, EventQueue

#### [NEW] [sim/tick.go](file:///Users/raynwow/m_projects/anti/summit/sim/tick.go)
- `Tick()`: AdvanceTime → ThreatEval → WarningTimers → EventQueue.Drain

#### [NEW] [clock/game_clock.go](file:///Users/raynwow/m_projects/anti/summit/clock/game_clock.go)
- Tick interval, speed multiplier (0/1/4/12x), pause flag

#### [MODIFY] [main.go](file:///Users/raynwow/m_projects/anti/summit/main.go)
- Switches from terminal loop to `ebiten.RunGame()`

---

### Phase 5: Plane — Sound and animation, Balatro-style feedback (Tasks 27–34)

**Goal:** Does the oxygen number dropping *feel* like something before you consciously read it?

#### [NEW] [audio/audio.go](file:///Users/raynwow/m_projects/anti/summit/audio/audio.go)
- Initialize audio system, master volume control

#### [NEW] [audio/sounds.go](file:///Users/raynwow/m_projects/anti/summit/audio/sounds.go)
- Load sound assets, play functions with pitch parameter
- Pitched AMS tick pings, hollow O2 click, warning tone, crisis thud + screen shake, summit chime

#### [NEW] [audio/drone.go](file:///Users/raynwow/m_projects/anti/summit/audio/drone.go)
- Continuous ambient loop, pitch/volume driven by aggregate danger level each frame

#### [MODIFY] [ui/map_view.go](file:///Users/raynwow/m_projects/anti/summit/ui/map_view.go)
- Idle dot animation: all dots pulse at all times, speed/amplitude increase with danger

#### [MODIFY] [ui/event_panel.go](file:///Users/raynwow/m_projects/anti/summit/ui/event_panel.go)
- Enhanced entrance: 1-frame white flash → dim over 0.5s → slide-in with ease-out

#### [MODIFY] [ui/radio_log.go](file:///Users/raynwow/m_projects/anti/summit/ui/radio_log.go)
- Typewriter text: crisis lines reveal 1 char/40ms with tick sound; normal lines instant
- Prevention confirmation line when warning clears

#### [MODIFY] Various UI files
- Sequential stat reveals: values tick toward target at ~5 units/frame, never jump

---

### Phase 6: Rocket — Full systems integration (Tasks 35–49)

**Goal:** Does every system talk to every other system?

#### [NEW] [catalyst/catalyst.go](file:///Users/raynwow/m_projects/anti/summit/catalyst/catalyst.go)
- `FilteredWorld` struct: copy of World with visibility flags
- `Apply(world) → FilteredWorld` — all UI reads from filtered copy only

#### [NEW] [catalyst/comms_filter.go](file:///Users/raynwow/m_projects/anti/summit/catalyst/comms_filter.go)
- Garble radio log text, hide AMS numbers (`--`)

#### [NEW] [catalyst/sensor_filter.go](file:///Users/raynwow/m_projects/anti/summit/catalyst/sensor_filter.go)
- Freeze barometer reading in filtered copy

#### [NEW] [catalyst/info_tier.go](file:///Users/raynwow/m_projects/anti/summit/catalyst/info_tier.go)
- Basic/Enhanced/Expert tiers, upgrade button, tier-based visibility

#### [NEW] [resolution/action_resolver.go](file:///Users/raynwow/m_projects/anti/summit/resolution/action_resolver.go)
- Apply stat deltas, call domino seeder after

#### [NEW] [resolution/domino_seeder.go](file:///Users/raynwow/m_projects/anti/summit/resolution/domino_seeder.go)
- Lookup action type in domino table, plant delayed whisper

#### [NEW] [resolution/choice_scoper.go](file:///Users/raynwow/m_projects/anti/summit/resolution/choice_scoper.go)
- Gate choices by warning history + resource availability

#### [NEW] [threat/compound.go](file:///Users/raynwow/m_projects/anti/summit/threat/compound.go)
- Detect overlapping L3 crises, build merged panel with fewer choices

#### [NEW] [threat/definitions/*.go](file:///Users/raynwow/m_projects/anti/summit/threat/definitions/) (12 files)
- All 12 threat types: HACE, HAPE, Summit Fever, Frostbite, Oxygen Crisis, Weather Break, Serac Collapse, Porter Strike, Psych Break, Rival Trouble, Equipment Fail, Sponsor Pressure

#### [NEW] [world/rival.go](file:///Users/raynwow/m_projects/anti/summit/world/rival.go)
- Grey dots, independent schedule from seed, help/ignore events

#### [NEW] [ui/topdown_view.go](file:///Users/raynwow/m_projects/anti/summit/ui/topdown_view.go)
- Toggle between side view and top-down tile grid

#### [NEW] [ui/weather_panel.go](file:///Users/raynwow/m_projects/anti/summit/ui/weather_panel.go)
- Forecast rows, confidence bars, barometer + trend arrow

#### [NEW] [clock/wall_clock.go](file:///Users/raynwow/m_projects/anti/summit/clock/wall_clock.go)
- Real-time 4-minute countdown for crisis panels, independent of game speed

#### [NEW] [clock/warning_timer.go](file:///Users/raynwow/m_projects/anti/summit/clock/warning_timer.go)
- Per-warning game-time expiry counter

#### [NEW] [rng/route.go](file:///Users/raynwow/m_projects/anti/summit/rng/route.go)
- Generate route variant: segment types per altitude band

#### [NEW] [data/event_templates.go](file:///Users/raynwow/m_projects/anti/summit/data/event_templates.go)
- All whisper/warning/crisis text strings indexed by ThreatType + Level

#### [NEW] [data/domino_table.go](file:///Users/raynwow/m_projects/anti/summit/data/domino_table.go)
- `map[ActionType][]DominoRule` — what each resolution spawns

#### [NEW] [save/save.go](file:///Users/raynwow/m_projects/anti/summit/save/save.go)
- Serialize full state to JSON

#### [NEW] [save/load.go](file:///Users/raynwow/m_projects/anti/summit/save/load.go)
- Deserialize and restore, validate version compatibility

---

## Open Questions

> [!IMPORTANT]
> **Technology choice:** Do you approve **Go + Ebitengine**, or do you prefer **Lua + LÖVE2D**?

> [!NOTE]
> **Build order:** The plan follows the design doc's Skateboard → Rocket sequence exactly. Should I start with Phase 1 (Skateboard — terminal-only, single climber) and we'll iterate from there?

---

## Verification Plan

### Each Phase — Play Test
The design doc defines explicit pass conditions per phase:

| Phase | Pass Condition |
|---|---|
| Skateboard | AMS hitting 80 feels bad. Reaching summit feels good. |
| Bicycle | Choosing between two climbers in trouble feels genuinely bad. |
| Scooter | Different seed = meaningfully different run. Immediate retry desire after losing. |
| Car | Can tell which climber is in trouble from dot behavior alone (60s, no numbers). |
| Plane | Can tell crisis vs calm from sound alone (30s, monitor off). |
| Rocket | Comms blackout + frostbite warning + oxygen shortage = specific kind of hell built by player decisions. |

### Automated Tests
- `go build ./...` compiles cleanly at every phase
- `go vet ./...` passes at every phase
- Unit tests for core mechanics: `AdvanceTime()`, `CheckThreats()`, `ApplyChoice()`, seed determinism (same seed → same run)

### Manual Verification
- Play each phase before proceeding to the next — **the design doc's rule: if the current vehicle isn't fun, fix it before building the next one**
