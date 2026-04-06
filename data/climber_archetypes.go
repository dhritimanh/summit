package data

type ClimberArchetype struct {
	NamePool      []string
	BaseFitness   int
	AMSSusPercent float64 // 1.0 = normal, 1.5 = 50% more AMS gain
}

var ArchetypePool = []ClimberArchetype{
	{
		NamePool:      []string{"Zara", "Maya", "Elena"},
		BaseFitness:   90,
		AMSSusPercent: 1.0,
	},
	{
		NamePool:      []string{"Erik", "Lukas", "Johan"},
		BaseFitness:   80,
		AMSSusPercent: 1.2,
	},
	{
		NamePool:      []string{"Karin", "Siri", "Ingrid"},
		BaseFitness:   95,
		AMSSusPercent: 1.5,
	},
}
