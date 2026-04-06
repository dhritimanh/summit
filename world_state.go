package main

import (
	"fmt"
	"math/rand"
	"time"
)

type WorldState struct {
	Seed         int64
	Rng          *rand.Rand
	Day          int
	Hour         int
	WindSpeed    int
	CampO2       map[int]int // Camp Index -> Bottle count (1 bottle = 3 charges)
	WeatherCurve []int       // Wind speed per day
	Archetype    ClimberArchetype
}

func NewWorldState(seed int64) *WorldState {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	r := rand.New(rand.NewSource(seed))

	// Generate deterministic weather curve (30 days + backup)
	curve := make([]int, 32)
	for i := range curve {
		curve[i] = 20 + r.Intn(40) // Base wind range; scaled by altitude in sim.go
	}

	return &WorldState{
		Seed:         seed,
		Rng:          r,
		Day:          1,
		Hour:         6,
		WindSpeed:    20,
		WeatherCurve: curve,
		CampO2: map[int]int{
			LocBase:     99,
			LocCamp1:    10,
			LocCamp2:    8,
			LocHighCamp: 6,
			LocSummit:   0,
		},
	}
}

func (w *WorldState) GetForecast(day int) []string {
	forecast := make([]string, 4)
	for i := 0; i < 4; i++ {
		targetDay := (day + i) % 31
		if targetDay == 0 {
			targetDay = 1
		}
		wind := w.WeatherCurve[targetDay]

		switch i {
		case 0, 1:
			forecast[i] = fmt.Sprintf("Day %d: %d km/h (Stable)", day+i, wind)
		case 2:
			forecast[i] = fmt.Sprintf("Day %d: %d-%d km/h (Moderate)", day+i, wind-10, wind+10)
		case 3:
			forecast[i] = fmt.Sprintf("Day %d: %d-%d km/h (Low)", day+i, wind-20, wind+20)
		}
	}
	return forecast
}
