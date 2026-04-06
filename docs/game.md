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

**Rule: if the current vehicle isn't fun, fix it before building the next one.**

---

### Skateboard — does one decision feel consequential?

One climber. Terminal output only. No graphics. The question being answered: does reading a log line and picking a number create genuine tension?

**What the terminal output looks like:**
```
────────────────────────────────
 Day 3  |  Zara  |  Camp 1
 Fitness: 71   AMS: 58
────────────────────────────────

 !  [Zara] AMS elevated: 58. Headache, moving slow.

 [1] Rest here 6hrs          fitness -5, AMS -10
 [2] Descend one camp        fitness -10, AMS -20
 [3] Ignore it, keep moving  fitness -15, AMS +12

>
```

The numbers after each choice are shown explicitly at this stage — not because the final game shows them, but because during the skateboard you need to verify the mechanics are working correctly before hiding information.

| # | Task | What to build |
|---|------|---------------|
| 1 | Climber struct | `name string`, `fitness int` (0–100), `ams int` (0–100), `loc int` (0=base, 1=camp1, 2=camp2, 3=summit) |
| 2 | `RunGame()` loop | Each iteration: call `AdvanceTime()`, call `CheckThreats()`, call `PrintStatus()`, read stdin, call `ApplyChoice()` |
| 3 | `AdvanceTime()` | At altitude: fitness -8/turn, AMS +2–8 randomly. At base: fitness +12/turn, AMS -15/turn |
| 4 | `CheckThreats()` | AMS ≥35 and rand: print whisper log line, no pause. AMS ≥55: print warning log line, offer 3 choices. AMS ≥80: print crisis panel, offer 3 choices (worst options only — no warning was acted on) |
| 5 | `ApplyChoice()` | Apply fitness delta, AMS delta, location delta from chosen option |
| 6 | Win/lose checks | Fitness ≤0: print death message, exit. Loc==3: print summit message. Loc==0 after summit: print success, exit |

**Files:** everything in `main.go`. No packages. No subdirectories.

**Pass condition:** Play it. Does the moment where AMS hits 80 feel bad? Does reaching the summit feel good? If yes, continue.

---

### Bicycle — does managing a team feel different from managing one person?

Three climbers. Loop through each one per turn. Add a second threat. Add the warning timer. Add choice gating. The question: does watching Erik deteriorate while you're dealing with Zara create a specific kind of dread?

**What changes in the output:**
```
────────────────────────────────
 Day 5  |  09:00
 Zara     Camp 2   Fit:71  AMS:44  [OK]
 Marcus   Camp 1   Fit:88  AMS:22  [OK]
 Erik     Camp 1   Fit:63  AMS:71  [!!]
────────────────────────────────

 RADIO LOG:
 [07:00] Zara — Camp 2. Wind picked up overnight.
 [09:00] !! Erik — AMS elevated: 71. Recommend assessment.
         Warning expires in: 6 hrs

--- ERIK ---
 Fitness: 63   AMS: 71   Location: Camp 1

 [1] Order rest + oxygen      AMS -15 over 6hrs
 [2] Order descent to base    AMS -30, loses 2 days
 [3] Do nothing               Warning timer continues
```

| # | Task | What to build |
|---|------|---------------|
| 7 | Three climbers | Array of climber structs, loop through each per turn for both time advance and threat check |
| 8 | Team status line | One summary line per climber showing name, location, fitness, AMS, status tag [OK]/[!]/[!!] |
| 9 | Second threat: frostbite | Triggers when wind speed (new field, randomised per turn at altitude) is high and climber has been exposed 3+ turns. Whisper: no output. Warning: log line. Crisis: finger/toe damage, permanent fitness cap reduced |
| 10 | Warning timer | Each active warning stores a `turnsRemaining int`. Decrements each turn. When it hits 0, crisis fires regardless of player action |
| 11 | Choice gating | `warningActedOn bool` field on each active warning. If false when crisis fires, remove the best choice from the options list |
| 12 | Radio log as persistent list | Store last 10 log lines, print them above the current prompt every turn |

**New files:** Extract `climber.go` for the struct. `threat.go` for threat type and level constants. Everything else still in `main.go`.

**Pass condition:** Play it. When two climbers are in trouble simultaneously and you can only act on one this turn, does it feel genuinely bad to choose? If yes, continue.

---

### Scooter — does each run feel different enough to play again?

Add resources. Add the seed system. Add weather. The question: after losing, do you immediately want to try again with a different approach?

**What changes:**
- Oxygen cylinders are now a real constraint. Each climber consumes 1 cylinder per turn at altitude when assigned O2. Running out mid-route triggers an oxygen crisis event.
- The run starts with a seed number printed at the top. Each run, climber names, starting stats, weather curve, and event deck order are generated from this seed.
- A 4-day weather forecast is visible. Confidence degrades for later days. The summit window appears and closes based on the seed — sometimes it's day 18, sometimes day 24.

| # | Task | What to build |
|---|------|---------------|
| 13 | Oxygen resource | `o2PerCamp [4]int` array. Each turn at altitude, if climber has O2 assigned, decrement their camp's supply. If supply hits 0 mid-route, fire oxygen crisis event immediately |
| 14 | Oxygen crisis event | Crisis panel: options are improvise descent (dangerous), share from another camp (requires porter turn), abort attempt. Available options depend on whether the O2 warning was acted on |
| 15 | Run seed | `rand.New(rand.NewSource(seed))` passed to all generation functions. Print seed at game start and game end |
| 16 | Climber roster generation | Draw 3 climbers from archetype pool. Each archetype has stat ranges and one hidden trait (e.g. altitude sickness susceptibility means AMS ticks 1.5x faster) |
| 17 | Weather curve generation | Generate 30-day pressure curve from seed. Summit window = consecutive days where wind < 40km/h. Window position varies per seed |
| 18 | Event deck shuffle | All threat types in a list, shuffle order from seed. This determines which threat fires first in ambiguous situations |

**New files:** `world/resource.go`, `world/weather.go`, `rng/seed.go`, `rng/roster.go`, `rng/weather_seed.go`, `rng/event_deck.go`, `data/climber_archetypes.go`

**Pass condition:** Lose a run. Immediately start a new one with a different seed. Does the new run feel like a genuinely different situation? If yes, continue.

---

### Car — does the mountain feel like a place?

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

---

### Plane — does the game feel alive?

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

---

### Rocket — does every system talk to every other system?

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