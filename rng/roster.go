package rng

import (
	"summit/world"
	"summit/data"
)

func GetRandomArchetype(w *world.WorldState) data.ClimberArchetype {
	idx := w.Rng.Intn(len(data.ArchetypePool))
	return data.ArchetypePool[idx]
}

func GenerateClimber(w *world.WorldState) *world.Climber {
	arch := GetRandomArchetype(w)
	name := arch.NamePool[w.Rng.Intn(len(arch.NamePool))]

	// Populating World-level archetype for session tracking in single-player
	w.Archetype = arch

	return &world.Climber{
		Name:          name,
		Archetype:     arch,
		Fitness:       arch.BaseFitness,
		AMS:           0,
		Loc:           data.LocBase,
		ActiveThreats: make(map[world.ThreatType]*world.ActiveThreat),
		O2Charges:     3, // Standard start
	}
}
