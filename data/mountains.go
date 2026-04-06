package data

// Loc constants for camp IDs
const (
	LocBase     = 0
	LocCamp1    = 1
	LocCamp2    = 2
	LocHighCamp = 3
	LocSummit   = 4
)

// Mountain represents a climbable peak (Modular for DLC)
type Mountain struct {
	Name          string
	CampNames     map[int]string
	CampAltitudes map[int]int
	DeathZone     int // Altitude where "Death Zone" begins (usually 8000m)
}

// Choice represents a player action
type Choice struct {
	Label      string
	FitDelta   int
	AmsDelta   int
	LocDelta   int
	ClearsWarn bool
}

// Pre-defined mountains
var Everest = Mountain{
	Name: "Everest",
	CampNames: map[int]string{
		LocBase:     "Base Camp",
		LocCamp1:    "Camp 1",
		LocCamp2:    "Camp 2",
		LocHighCamp: "Camp 4 (High)",
		LocSummit:   "Summit",
	},
	CampAltitudes: map[int]int{
		LocBase:     5364,
		LocCamp1:    6065,
		LocCamp2:    6492,
		LocHighCamp: 7950,
		LocSummit:   8848,
	},
	DeathZone: 8000,
}
