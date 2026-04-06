# Summit — Complete Build Brief

## What this game is

Summit is a turn-based survival management game. The player leads a rope team of 3–5 climbers up a mountain and must bring them back alive. The game is not about action — it's about *reading signals and making decisions under incomplete information*. The core tension is that threats always telegraph before they kill, but the telegraphs are subtle, and by the time something is obvious it may already be too late.

The closest reference points are FTL (team resource management, permadeath, cascading crises) and Subnautica (dread created by ambient information, not jump scares). The feel target for animations and sound is Balatro — every consequence plays out sequentially with micro-feedback on each step, not all at once.

The player never types. Every input is a button click. Maximum 4 choices at any decision point.

---

## What the player sees and does

There are two views of the mountain. The **side view** is the primary screen — a cross-section of the mountain showing camps at different altitudes, with each climber as a small dot that moves up and down between camps. The **top-down view** is a tile grid showing the team's horizontal spread across the terrain.

On the right is a **radio log** — a scrolling panel in monospace font showing timestamped messages from climbers. This is the primary information channel. Players learn to read tone, not just content. A climber who normally says "Morning base. Cold night but slept well. Ready when you are." and today says "Fine." is a whisper-level threat signal.

At the bottom is a **resource panel** showing oxygen cylinders per camp, food, fuel, medicine, and money.

At the top is a **speed control** — Pause / 1x / 4x / 12x. The game runs in real time at the chosen speed. At 12x a full day passes in seconds. The player pauses to make careful decisions, runs at 12x during safe stretches.

When a **crisis fires**, the game pauses automatically. A panel slides in over the dimmed mountain view. It shows the situation, the climber's current stats, and 2–4 action buttons. A real-time countdown of 4 minutes runs. If the player does nothing, the last (worst) option activates automatically. The game unpauses when the choice is made.

When there is **no crisis**, the player clicks any climber dot to open their panel and give orders: advance to next camp, hold position, descend, adjust oxygen flow.

---

## The threat system — the entire game logic lives here

Every threat in the game has exactly three levels. They always arrive in order. They never skip.

**Level 1 — Whisper.** No UI element fires. No pause. No alert. The signal exists only in the radio log or in a stat that ticked slightly faster than normal. A player who isn't paying attention will miss it entirely. Example: Erik's AMS went from 22 to 28 in six hours. Normal rate at that altitude is 2–3 per hour. Something is slightly accelerated. No message appears. The number just ticked faster.

**Level 2 — Warning.** A tagged line appears in the radio log. Still no pause, no event panel. The game keeps running. The player has a 6-hour in-game window to address the underlying issue. If they do, the warning clears and a confirmation line appears in the log: `[11:00] Erik — headache gone. Whatever you did, it worked.` If they don't act within 6 hours, the crisis fires.

**Level 3 — Crisis.** Game pauses. Event panel appears. 4-minute real-time countdown. The choices available depend entirely on whether the player acted at Level 2. If they acted: 4 options including good ones. If they ignored the warning: 3 options, all costly. If they ignored both whisper and warning: 2 options, one of them terrible.

This is the core teaching loop. Players lose their first few runs because they didn't know the warning had a timer. Mid-game players watch for warnings and act fast. Late-game players read whispers and act before warnings appear. Skill expression is entirely in how early you intervene.

**Level 4 — Compound Crisis.** When two separate Level 3 crises fire in the same location or overlapping timeframe, they merge into a single panel. The choices presented make it impossible to fully resolve both situations. The player must sacrifice something. These moments are the game's most memorable — they feel inevitable in retrospect because the player built them through five decisions made across four days.

---

## The catalyst system — attacking information, not stats

Catalysts are environmental events that don't hurt climbers directly. They corrupt the information the player uses to make decisions.

**Comms blackout** — radio static makes check-in messages garbled. A whisper-level signal (terse check-in, slightly elevated AMS) becomes invisible. The crisis can arrive without the whisper ever having been readable.

**Broken barometer** — the pressure reading freezes at its last known value. The player loses the subtle pressure-drop signal that precedes storm arrival. The weather panel shows stale data.

**Stolen cache** — the rival expedition takes supplies from a shared camp. This instantly converts a warning-level oxygen situation into a crisis.

Catalysts are strictly an information problem. They never change the underlying stats — they change what the player can see. This is architecturally important: the real world state and the filtered view the player sees are always separate objects.

---

## The domino system — fixes that create new problems

When a player resolves a crisis, the resolution can plant a new whisper in a different climber or system. The player is never told this is happening.

Examples:
- Emergency descent with escort: the escorting climber burns fitness at 3x normal rate. Twelve hours later, a severe exhaustion whisper fires for the escort.
- Administering dexamethasone: suppresses AMS symptoms but causes a 35-point spike in Summit Fever. If Summit Fever crosses 70, the climber may refuse descent orders.
- Ordering full team rest: recovers fitness but burns oxygen and food faster. If supplies were already tight, rest can convert a warning into a crisis.

The domino table is data-driven — a lookup table of `ActionType → consequences` stored separately from the game logic. This means new domino rules can be added without touching the engine.

---

## The information economy

Players start with basic information: radio log text only. They cannot see exact AMS numbers — they infer danger from the tone of check-ins.

Spending money unlocks information tiers:
- **Enhanced** ($400): biometric tags on all climbers. Exact AMS, fitness, and psychology numbers become visible on the climber panel.
- **Expert** ($800 + $400/day retainer): meteorologist on call. The real weather forecast and its actual confidence levels become visible. Without this, the player sees a degraded forecast with artificially reduced confidence.

This creates a meaningful early-game tension: spend money to see clearly, or spend it on oxygen and supplies and navigate blind.

---

## Sound and animation — the Balatro principle

Every consequence plays out sequentially with a micro-feedback event on each step. Never show the full result at once.

When AMS rises, the number ticks up one unit at a time. Each tick plays a ping sound. The pitch of the ping rises as the number rises — so a dangerous AMS level *sounds* dangerous before the player reads it. At the warning threshold, a distinct sound plays and the number briefly changes colour. At the crisis threshold, a heavy thud plays, the screen dims, and the panel slides in.

Climber dots are never completely still. Healthy climbers pulse very slowly. Climbers in trouble pulse faster and shift from white to amber. Critical climbers pulse urgently in red. The player's eye goes to urgency without reading anything.

A low ambient drone runs under the entire game. Its pitch and volume rise as the aggregate danger across all climbers increases. Players begin making decisions based on this sound before they consciously understand why.

Radio log lines for normal check-ins appear instantly. Lines for warnings type out letter by letter. Crisis-level lines type out slowly with a heavier sound per character. The *speed* of text arrival signals severity before the words are read.

When a warning clears because the player acted in time, a confirmation log line appears. This closes the "combo" — the player gets explicit confirmation that they prevented something. Without this, successful prevention feels like nothing happened.

---

## Randomisation and replayability

Every run is generated from a single integer seed. The same seed always produces the same mountain, same climbers, same weather, same event order. The seed is shown at the end of every run so players can share exact runs with each other.

Each run randomises: the climber roster (drawn from an archetype pool — veteran sherpa, rookie journalist, elite alpinist, each with different stat spreads and hidden traits), the weather curve (summit window opens and closes at different points), the event deck (which threats fire and in what order), and the route variant (segment types — ice field, exposed ridge, rock chimney — that change equipment requirements and stamina cost).

---

## The rival expedition

A second team is always on the mountain. Their climbers appear as grey semi-transparent dots on the map. They move independently according to their own schedule. They are competitive intelligence — watching them tells the player things about current conditions. When a grey dot stops moving, something has gone wrong on their team. When they radio on the shared frequency requesting help, the player faces a choice: assist (costs resources, builds relationship, they may help you later) or ignore (free, but they remember).

---

## Build sequence — skateboard to rocket

---

### Phase 1: Skateboard — One climber, terminal only (Tasks 1–6)

**Goal:** Does one decision feel consequential?

One climber. Terminal output only. No graphics. The question being answered: does reading a log line and picking a number create genuine tension?

**What the terminal output looks like:**
```
────────────────────────────────
 Day 3  |  09:00  |  Zara  |  Camp 1
 Altitude: 6065m   Fitness: 71   AMS: 58
────────────────────────────────

 !  [Zara] AMS elevated: 58. Headache. Warning expires in 4 turns.

 [1] Rest here (short)       fitness +8, AMS -5
 [2] Descend one camp        fitness -5, AMS -10
 [3] Ignore it, keep moving  fitness -8, AMS +6

>
```

The numbers after each choice are shown explicitly at this stage — not because the final game shows them, but because during the skateboard you need to verify the mechanics are working correctly before hiding information. Note: We have already implemented the 3-hour tactical turn system (8 turns/day) for better control.

| # | Task | What to build |
|---|------|---------------|
| 1 | Climber & Mountain | `Climber` with `name`, `fitness` (0-100), `ams` (0-100), `altitude` (meters). `Mountain` modular (Everest-scale: Base 5364m to Summit 8848m) |
| 2 | `RunGame()` loop | Action → Decay → Check logic to ensure players can respond to threats before environmental impact. |
| 3 | `AdvanceTime()` | 3-hour turns. At altitude: fitness -3 to -5, AMS scales by height. Resting significantly mitigates decay. |
| 4 | `CheckThreats()` | Whisper (AMS≥30, 50% chance). Warning (AMS≥55, 12h timer). Crisis (AMS≥80 or timeout, restricted choices). |
| 5 | `ApplyChoice()` | Apply fitness/AMS deltas. Advance height (~400m/turn climb). Reaching next camp takes 2-4 turns. |
| 6 | Win/lose checks | Fitness≤0 or AMS=100: death. Loc==summit: victory marker. Loc==base after summit: success + exit. |

**Files:** Modular structure across `main.go`, `climber.go`, `sim.go`, `ui.go`, `events.go`, and `mountain.go`. This is **Better** than the original "Single file" requirement as it prepares for scale.

**Pass condition:** Play it. Does the moment where AMS hits 80 feel bad? Does the 8,000m+ Death Zone feel terrifying? If yes, continue.

---

### Phase 2: Bicycle — Architecture & Environmental Threats (Tasks 7–13)

**Goal:** Can we scale our logic without breaking anything, and does a time-pressured warning feel different? (Still one climber)

Before adding complex systems like Frostbite, Meds, or Team Psychology, we need a rigorous architecture. We will introduce a centralized **Stat Application System** and **Unit Tests** to mathematically verify all core physics instead of relying purely on playtesting. Then, we add Wind and Frostbite.

| # | Task | What to build |
|---|------|---------------|
| 7 | Stat & Modifier Abstraction | Create `ApplyStatChange(climber, stat, amount, source)` pattern. Ensure caps and modifiers (e.g., frostbite max-health limits, meds buffs, psych hits) trigger cleanly here instead of via direct value manipulation. |
| 8 | Unit Testing Suite | Add `sim_test.go`. Write mathematical tests verifying Ambient Math (Base vs Death Zone), Night Penalties, and Resting mitigation to ensure verifiability. |
| 9 | Wind Speed System | Add `WindSpeed int` to World state (sim.go). Varies predictably based on altitude + random gust factor. |
| 10 | Second threat: frostbite | Triggers when wind speed is high and exposure turns > 3. Whisper: no output. Warning: log line. Crisis: permanent fitness cap reduction (handled by Task 7). |
| 11 | Warning timer | Active warnings store `turnsRemaining int` (e.g., 4 turns). Decrements each turn. When it hits 0, crisis fires regardless of player action. |
| 12 | Choice gating | `warningActedOn bool` field. If false when crisis fires, remove the best choice from the crisis panel. |
| 13 | Radio log as persistent | Store last 10 log lines, print them above the current prompt every turn. |

**New files:** Extract `threat.go` for threat type and level constants. `sim_test.go` for verifiability suite.

**Pass condition:** 1. `go test` passes perfectly. 2. Playtest: When a frostbite warning fires but you have a summit window opening, do you risk the turn to climb, or do you act?


---

### Phase 3: Scooter — Resources, seeds, weather (Tasks 13–18)

Add resources. Add the seed system. Add weather. The question: after losing, do you immediately want to try again with a different approach?

**What changes:**
- Oxygen cylinders are now a real constraint. A standard oxygen bottle holds 3 "Charges" (approx 9 hours). Each climber consumes 1 Charge per turn (3 hours) at altitude when assigned O2. Running out mid-route triggers an oxygen crisis event.
- The run starts with a seed number printed at the top. Each run, climber names, starting stats, weather curve, and event deck order are generated from this seed.
- A 4-day weather forecast is visible. Confidence degrades for later days. The summit window appears and closes based on the seed — sometimes it's day 18, sometimes day 24.

| # | Task | What to build |
|---|------|---------------|
| 13 | Oxygen resource | `CampO2 [4]int` array (measured in Bottles). 1 Bottle = 3 Charges. Climbers track `O2Charges int`. Each turn, if assigned O2, burn 1 charge. If 0, load a new bottle from camp if at camp. If no bottles, O2 Crisis immediately. |
| 14 | Oxygen crisis event | Crisis panel: options are improvise descent (dangerous), push without O2 (massive AMS spike), abort attempt. |
| 15 | Run seed | `rand.New(rand.NewSource(seed))` passed to all generation functions. Print seed at game start and game end |
| 16 | Climber archetype | Instead of a hardcoded "Zara", draw 1 active climber from an archetype pool based on the seed. Each archetype has stat ranges and one hidden trait (e.g. altitude sickness susceptibility means AMS ticks 1.5x faster) |
| 17 | Weather curve generation | Generate 30-day pressure curve from seed. Summit window = consecutive days where wind < 40km/h. Window position varies per seed |
| 18 | Event deck shuffle | All threat types in a list, shuffle order from seed. This determines which threat fires first in ambiguous situations |

**New files:** `world/resource.go`, `world/weather.go`, `rng/seed.go`, `rng/roster.go`, `rng/weather_seed.go`, `rng/event_deck.go`, `data/climber_archetypes.go`

**Pass condition:** Lose a run. Immediately start a new one with a different seed. Does the new run feel like a genuinely different situation? If yes, continue.





**Goal:** Does each run feel different enough to play again?

#### [NEW] [world/resource.go](file:///Users/raynwow/m_projects/anti/summit/world/resource.go)
- O2 bottle arrays per camp, food, fuel, medicine, money
- O2 consumption: 1 charge/turn (3 hours), 3 charges per bottle

#### [NEW] [world/weather.go](file:///Users/raynwow/m_projects/anti/summit/world/weather.go)
- 30-day pressure curve, wind speed, forecast (real + degraded)

#### [NEW] [rng/seed.go](file:///Users/raynwow/m_projects/anti/summit/rng/seed.go)
- `RunSeed` struct wrapping `rand.Rand`, single source of truth

#### [NEW] [rng/roster.go](file:///Users/raynwow/m_projects/anti/summit/rng/roster.go)
- Generate the climber from an archetype pool with stat ranges and hidden traits

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

This is the first graphical stage. Ebitengine (Go) or Löve2D (Lua) — the builder chooses based on language decision. The question: do the dots on the mountain communicate danger before the player reads any numbers?

**What the screen looks like:**

Left 60% of screen: mountain cross-section. A triangle or stepped polygon representing the mountain profile. Camps marked at fixed altitude positions (base, camp 1, camp 2, high camp, summit). Each climber is a circle dot positioned at their current camp. Dot colour: white (AMS < 55), amber (AMS 55–79), red (AMS ≥ 80). Dots pulse — slow pulse when healthy, faster when in warning state, urgent when critical.

Right 20% of screen: radio log panel. Monospace font. Scrolling. Timestamps on every line.

Bottom strip: resource panel. O2 cylinders shown as small icons per camp. Food/fuel/medicine as number counters.

Top strip: speed controls (Pause / 1x / 4x / 12x), current day and time.

When a climber dot is clicked: a panel opens showing their stats and 3–4 action buttons. Clicking elsewhere closes it.

When a crisis fires: the background dims to 40% opacity, a panel slides in from the top taking up the centre 50% of the screen, the countdown timer displays in the top right of the panel.

| # | Task | What to build |
|---|------|---------------|
| 19 | Render loop | `Update()` advances game state. `Draw()` renders everything. Game clock runs at chosen speed multiplier |
| 20 | Mountain side view | Static polygon for mountain silhouette. Fixed pixel positions for each camp altitude. Climber dots drawn as circles at their camp position |
| 21 | Dot states | Dot colour based on AMS threshold. Pulse animation: `scale = 1.0 + 0.05 * sin(time * pulseSpeed)` where pulseSpeed increases with danger |
| 22 | Crisis panel overlay | On crisis: set `gameClockPaused = true`. Draw dimmed background rect. Slide panel in from top over 0.3 seconds. Show choices as buttons. Show countdown timer |
| 23 | Climber click panel | On dot click: open panel with name, stats, location, action buttons. Clicking a button calls `ApplyChoice()` and closes panel |
| 24 | Radio log panel | Store log lines as strings. Draw last N lines that fit in panel height. New lines appear at bottom, old ones scroll up |
| 25 | Speed controls | Buttons at top. Clicking sets `speedMultiplier` to 0 (pause), 1, 4, or 12. Game clock tick interval divides by multiplier |
| 26 | Resource panel | Draw O2 count per camp as small cylinder icons. Food/fuel/medicine as labelled number counters |

**New packages:** `ui/`, `sim/`, `clock/game_clock.go`

**Pass condition:** Watch the dots for 60 seconds without reading any numbers. Can you tell which climber is in trouble from dot behaviour alone? If yes, continue.



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

Sound and animation. The Balatro principle: every consequence plays out as a sequence of micro-events, each with its own feedback. The question: does the oxygen number dropping feel like something before you consciously read it?

| # | Task | What to build |
|---|------|---------------|
| 27 | Sequential stat reveals | When a stat changes, animate it: current value ticks toward target value at rate of ~5 units per frame. Do not jump to final value instantly |
| 28 | Pitched sound per AMS tick | Each tick of the AMS counter during animation plays a ping. Pitch = `0.8 + (amsValue / 100.0) * 0.6`. High AMS = high pitched sound |
| 29 | Signature sounds | Oxygen drop: hollow click. Warning log line appears: low tone. Crisis panel arrives: heavy thud + brief screen shake (translate screen by 4px for 3 frames). Summit reached: single clear chime |
| 30 | Idle dot animation | All dots pulse gently at all times using `sin(time)`. Speed and amplitude increase with danger state. Nothing is ever completely still |
| 31 | Crisis panel entrance | On crisis: 1 frame flash to white, then dim background over 0.5s, then panel slides from y=-panelHeight to centred position over 0.3s using ease-out curve |
| 32 | Typewriter text | Crisis log lines: reveal one character every 40ms with a soft tick sound per character. Normal log lines: appear instantly |
| 33 | Prevention confirmation | When a warning timer clears because the underlying stat recovered, append to radio log: `[HH:MM] [climber name] — feeling better. Whatever you did worked.` |
| 34 | Ambient drone | Continuous audio loop. Each frame: calculate `dangerLevel = average(AMS of all climbers above 55) / 100`. Set drone pitch to `0.7 + dangerLevel * 0.5`. Set drone volume to `0.1 + dangerLevel * 0.4` |

**New package:** `audio/`

**Pass condition:** Turn off your monitor and listen to the game for 30 seconds. Can you tell when a crisis fires and when things are calm from sound alone? If yes, continue.





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

The full design. Catalysts corrupt information. Dominoes make fixes costly. Compounds make neglect catastrophic. The rival is a silent clock.

| # | Task | What to build |
|---|------|---------------|
| 35 | `FilteredWorld` struct | A copy of `World` with visibility flags. `catalyst.Apply(world)` returns this copy. All UI rendering reads from `FilteredWorld`, not `World`. `World` is only mutated by `ActionResolver` |
| 36 | Comms blackout catalyst | When active: garble radio log check-in text (replace words with static characters), set `climber.AMSRiskVisible = false` on filtered copy so AMS number shows as `--` in UI |
| 37 | Broken barometer catalyst | When active: freeze `FilteredWorld.Weather.BarometerReading` at last known value. Real pressure continues changing in `World` |
| 38 | Information tier system | `TierBasic`: AMS shows as `--`, forecast confidence artificially reduced. `TierEnhanced` ($400): exact numbers visible. `TierExpert` ($800 + retainer): real forecast confidence visible. Upgrade button in resource panel |
| 39 | Domino seeder | After `ApplyChoice()`, look up `action.Type` in `domino_table`. If a rule exists, plant a `Whisper` struct in `ThreatRegistry` with `TriggerAt = world.Now + rule.Delay`. The whisper fires normally when its time comes |
| 40 | Compound crisis detector | Each tick: get list of active Level 3 crises. For each pair, check if same camp or within 6 game-hours. If yes: build merged panel with fewer total choices, replace both individual panels in event queue |
| 41 | Rival expedition | Second set of climber dots in grey. Move on their own schedule (pre-generated from seed). When their dot stops: add line to radio log. When they radio for help: fire a special event panel with help/ignore choices |
| 42 | All 12 threat types | HACE, HAPE, Summit Fever, Frostbite, Oxygen Crisis, Weather Break, Serac Collapse, Porter Strike, Psych Break, Rival Trouble, Equipment Fail, Sponsor Pressure. Each as a definition file with whisper condition, warning text, crisis panel text, all possible choices, and which choices require prior warning action |
| 43 | `clock/wall_clock.go` | Runs independently of game speed. `StartCountdown(4 * time.Minute, onExpire)`. `onExpire` calls `ApplyChoice(panel.DefaultIx)`. Game clock is paused while panel is visible, wall clock is not |
| 44 | `clock/warning_timer.go` | Each active warning has a `GameDuration` expiry. Decremented by real game-time elapsed each tick (not wall time). When it reaches zero, escalate to crisis |
| 45 | `resolution/choice_scoper.go` | At crisis build time: query `ThreatRegistry.WasActedOnAtWarning(threatID)`. If false: remove choices flagged `RequiresWarningAction`. If resource unavailable: show choice greyed with reason text, disabled |
| 46 | Top-down view | Toggle button switches between side view and top-down tile grid. Grid shows terrain type per tile (ice/rock/snow from `rng/route.go`). Climber dots on grid at their tile position |
| 47 | Porter system | Porters are a separate resource. Sending a load to a camp takes 2 days and costs money. Porter strike threat: if unpaid, loads stop moving. Visible in resource panel as "porter team: available / on strike" |
| 48 | Run seed at end | On any terminal outcome (death/retreat/summit/all-safe), show: `Run seed: 48291. Share this to replay the same mountain.` |
| 49 | Save / load | Serialise full `World` + `ThreatRegistry` + `clock` state to JSON. Load restores exact state. Auto-save every in-game day |

**All packages now live:** `world/`, `catalyst/`, `threat/` + `definitions/`, `resolution/`, `clock/`, `sim/`, `ui/`, `audio/`, `rng/`, `data/`, `save/`




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



---

## Complete file structure at rocket stage

```
summit/
├── main.go
│
├── world/
│   ├── climber.go       ClimberState: name, fitness, ams, loc, summitFever,
│   │                    psychology, o2Assigned, fitnessDecayMult, hidden traits
│   ├── resource.go      O2 per camp, food, fuel, medicine, money, porter count
│   ├── weather.go       30-day pressure curve, wind speed, forecast (real + degraded)
│   ├── rival.go         rival climber positions, schedule, relationship score
│   └── world.go         World struct owning all above + current day/time
│
├── catalyst/
│   ├── catalyst.go      CatalystType enum, CatalystState, Apply() → FilteredWorld
│   ├── comms_filter.go  BlackoutFilter: garble log text, hide AMS visibility flag
│   ├── sensor_filter.go BrokenBaroFilter: freeze pressure reading in filtered copy
│   └── info_tier.go     InformationTier enum, Upgrade(), ApplyTier() to FilteredWorld
│
├── threat/
│   ├── threat.go        ThreatType enum, ThreatLevel enum, Threat struct
│   ├── registry.go      active threats map, history log, WasActedOnAtWarning()
│   ├── whisper.go       WhisperEval(): check conditions, plant whisper, no UI output
│   ├── warning.go       WarningEval(): emit radio log line, start warning timer
│   ├── crisis.go        CrisisEval(): build EventPanel, call ChoiceScoper
│   ├── compound.go      CompoundDetector: pair active L3s, build merged panel
│   └── definitions/
│       ├── hace.go              whisper: AMS ticks fast. warning: "slurred speech".
│       │                        crisis: ataxia confirmed. choices gated by warning action
│       ├── hape.go              whisper: dry cough in log. warning: O2 sat dropping.
│       │                        crisis: pulmonary, immediate descent required
│       ├── summit_fever.go      whisper: clipped answers near summit. warning: psych drop.
│       │                        crisis: standoff at turnaround, may refuse orders
│       ├── frostbite.go         whisper: wind high, gloves not mentioned. warning: numb fingers.
│       │                        crisis: tissue damage, permanent fitness cap reduced
│       ├── oxygen_crisis.go     whisper: consumption rate vs supply. warning: camp stock critical.
│       │                        crisis: mid-route empty, abort or improvise
│       ├── weather_break.go     whisper: baro drop, confidence falls. warning: storm revised earlier.
│       │                        crisis: window closes mid-push
│       ├── serac_collapse.go    whisper: creaking in log. warning: ice movement reported.
│       │                        crisis: route destroyed, camps buried
│       ├── porter_strike.go     whisper: payment complaint. warning: porters request meeting.
│       │                        crisis: no loads moving until resolved
│       ├── psych_break.go       whisper: increasing silence. warning: psychology below 25.
│       │                        crisis: climber refuses all orders
│       ├── rival_trouble.go     whisper: their dot stops. warning: radio silence.
│       │                        crisis: help request on shared frequency
│       ├── equipment_fail.go    whisper: gear complaint in passing. warning: specific failure.
│       │                        crisis: critical equipment out, workaround required
│       └── sponsor_pressure.go  whisper: sponsor email in log. warning: phone call.
│                                crisis: summit or funding pulled
│
├── resolution/
│   ├── action_resolver.go  ApplyAction(): apply stat deltas, call SeedDomino after
│   ├── domino_seeder.go    SeedDomino(): look up action in domino_table, plant whisper
│   ├── choice_scoper.go    ScopeChoices(): gate options by warning history + resources
│   └── outcome.go          TerminalOutcome: ClimberDeath, ForcedRetreat, SummitReached, AllHomeSafe
│
├── clock/
│   ├── game_clock.go    tick interval, speed multiplier (0/1/4/12x), pause flag
│   ├── wall_clock.go    real-time countdown, onExpire callback, independent of game speed
│   └── warning_timer.go per-warning game-time expiry counter
│
├── sim/
│   ├── sim.go           Sim struct: owns World, CatalystState, ThreatRegistry,
│   │                    CompoundDetector, EventQueue, both clocks
│   ├── tick.go          Tick(): AdvanceTime → catalyst.Apply → ThreatEval →
│   │                    CompoundCheck → WarningTimers.Tick → EventQueue.Drain
│   └── event_queue.go   ordered list of pending event panels, Drain() pushes to UI
│
├── ui/
│   ├── ui.go            implements ebiten.Game (or löve2d equivalent). owns all panels
│   ├── map_view.go      mountain polygon, camp positions, dot renderer, pulse animation
│   ├── topdown_view.go  tile grid, terrain types, climber positions on grid
│   ├── climber_panel.go click handler, stat display, action buttons, opens on dot click
│   ├── event_panel.go   crisis overlay: dim background, slide-in animation, countdown display
│   ├── radio_log.go     scrolling log, typewriter reveal for crisis lines, instant for normal
│   ├── resource_panel.go O2 cylinder icons per camp, food/fuel/medicine counters, porter status
│   ├── weather_panel.go  forecast rows, confidence bars, barometer reading + trend arrow
│   └── hud.go           speed buttons, day/time display, information tier indicator
│
├── audio/
│   ├── audio.go         initialise audio system, master volume
│   ├── sounds.go        load all sound assets, play functions with pitch parameter
│   └── drone.go         ambient drone loop, update pitch and volume each frame from danger level
│
├── rng/
│   ├── seed.go          RunSeed struct wrapping rand.Rand, single source of truth per run
│   ├── roster.go        GenerateClimberRoster(): draw archetypes, roll stats, assign hidden traits
│   ├── weather_seed.go  GenerateWeatherCurve(): 30-day pressure + wind array from seed
│   ├── event_deck.go    ShuffleEventDeck(): ordered list of which threat fires first in ties
│   └── route.go         GenerateRouteVariant(): segment types per altitude band
│
├── data/
│   ├── climber_archetypes.go  archetype pool: name pools, stat ranges, possible hidden traits
│   ├── event_templates.go     all whisper/warning/crisis text strings, indexed by ThreatType+Level
│   └── domino_table.go        map[ActionType][]DominoRule — what each resolution spawns
│
└── save/
    ├── save.go   serialise World + ThreatRegistry + clock state to JSON
    └── load.go   deserialise and restore, validate version compatibility
```

---

## What done looks like at each stage

**Skateboard:** Terminal prints a status block each turn. AMS climbing creates visible dread. Crisis panel appears with degraded options if warning was ignored. Summit and death both feel distinct.

**Bicycle:** Three climbers shown simultaneously. Being unable to act on all three this turn feels like a real sacrifice. Warning timers create urgency without pausing the game.

**Scooter:** Different seed = meaningfully different run. Oxygen running out mid-route feels like a genuine crisis not a scripted event. Player wants to immediately retry after losing.

**Car:** Dots communicate danger without text. Mountain feels like a place. Game can be watched at 4x speed and the dots tell the story.

**Plane:** Stat changes feel weighty. Sound tells you how dangerous things are before you read numbers. Prevention feels rewarding because the game confirms it.

**Rocket:** Every system affects every other system. A comms blackout during a frostbite warning during an oxygen shortage feels like a specific kind of hell that the player created through their own decisions. The seed at the end makes you want to share the run.

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

#### [MODIFY] Various UI files
- Sequential stat reveals: values tick toward target at ~5 units/frame, never jump

---

### Phase 7: Interstellar — The Rope Team (Tasks 50–55)

**Goal:** Does managing a team of 3–5 climbers feel like a different game?

Now that the environmental simulation, UI, and sound are perfected for one person, we add the final layer: Team Complexity.

| # | Task | What to build |
|---|------|---------------|
| 50 | Multi-Climber Array | Change `Sim.Zara` to `Sim.Team []Climber`. Update all loops to iterate per member. |
| 51 | Support Logic | Being in the same camp as a teammate provides a +2 Fitness recovery bonus per turn. |
| 52 | Interdynamics | New "Shared Warning" logic: if one climber is in Crisis, a teammate at the same location can use their action to help, improving the first climber's options. |
| 53 | Personality Traits | Each climber archetype has a "Psychology" stat that affects how they react to others' failures. |
| 54 | Team UI | Status summary line for all team members at the top of every screen. |
| 55 | Stress Escalation | If more than 2 climbers are in warning state, AMS growth for everyone else is increased by 1.2x. |

**Pass condition:** Play with 3 climbers. When two are in trouble simultaneously and you can only act on one this turn, does it create the intended "Sophie's Choice" dread?

---

## Complete file structure at interstellar stage
