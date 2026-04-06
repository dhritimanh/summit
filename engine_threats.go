package main

type ThreatLevel string

const (
	LevelNone    ThreatLevel = "none"
	LevelWhisper ThreatLevel = "whisper"
	LevelWarning ThreatLevel = "warning"
	LevelCrisis  ThreatLevel = "crisis"
)

type ThreatType string

const (
	ThreatAMS       ThreatType = "AMS"
	ThreatFrostbite ThreatType = "Frostbite"
	ThreatOxygen    ThreatType = "Oxygen"
)

type ActiveThreat struct {
	Type         ThreatType
	Level        ThreatLevel
	TurnsLeft    int
	WarningActed bool
}
