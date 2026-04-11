package data

// PhysiologyConfig stores the master balancing constants for the simulation
// Adjusting these changes the overall difficulty/realism of the game.
const (
	// Base Recovery (Perfect safety)
	BaseRecoveryFit = 6
	BaseRecoveryAms = 8

	// Ambient Altitude Decay (Passive fit loss per 3-hour turn)
	AltitudeFitLossMin = 3
	AltitudeFitLossMax = 5

	// Action Modifiers
	// BALANCING NOTE (Fitness): At Camp 2, avg AltitudeFitLoss is 4. Night penalty is 2 (Total Night Loss = 6).
	// To prevent "infinite healing" at high altitude without O2, the max camp recovery must be balanced.
	// Resting anywhere provides +4. Resting AT CAMP provides an extra +2 (Total Camp Rest = 6).
	// Net fitness at Camp 2: Day = +2 gain. Night = 0 (perfectly stable, cannot heal).
	RestingFitBonus  = 4 // Standard rest anywhere (Reduced from 6 to stall super-recovery)
	CampRestBonus    = 2 // ADDITIONAL bonus if at an established camp
	NightFitPenalty  = 2 // Added to fitLoss between 18:00 - 06:00
	ClimbingFitCost  = 2 // Added to fitLoss (makes movement harder)

	// AMS Modifiers
	// BALANCING NOTE (AMS): Camp 2 avg AMS gain is 2. Night penalty adds 1 (Total Night Gain = 3).
	// Resting AMS bonus is 2. Extra Camp rest bonus adds 1 (Total Camp Relief = -3 AMS).
	// Net AMS change at Camp 2: Day = (2 - 3) = -1 drain. Night = (3 - 3) = 0 (stable).
	RestingAmsBonus  = 2 // Subtracted from amsGain (Reduced from 3 to keep high altitudes risky)
	NightAmsPenalty  = 1 // Added to amsGain
	ClimbingAmsCost  = 1 // Added to amsGain

	// The Death Zone (8,000m+)
	DeathZoneFitLoss = 4 // Forced fitness loss even while resting
	DeathZoneAmsLoss = 1 // Minimum ams gain even while resting
)

// AmsBaseGain tracks the raw AMS pressure per altitude band
var AmsBaseGain = map[int][2]int{
	LocBase:     {-8, -5}, // Negative = recovery
	LocCamp1:    {0, 2},
	LocCamp2:    {1, 3},
	LocHighCamp: {2, 5},
	LocSummit:   {3, 6},
}
