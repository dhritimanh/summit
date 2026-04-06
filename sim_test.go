package main

import (
	"math/rand"
	"testing"
)

func setupTest() {
	session = NewWorldState(1) // Fixed seed for deterministic tests
}

func TestApplyStatChangeCaps(t *testing.T) {
	c := &Climber{Fitness: 50, AMS: 50}
	
	// Test basic
	ApplyStatChange(c, "Fitness", +10, "Test")
	if c.Fitness != 60 {
		t.Errorf("Expected 60, got %d", c.Fitness)
	}

	// Test max cap
	ApplyStatChange(c, "Fitness", +100, "Test")
	if c.Fitness != 100 {
		t.Errorf("Expected 100, got %d", c.Fitness)
	}

	// Test modifier reducing max cap
	c.Modifiers = []Modifier{{Type: ModFrostbite, MaxFitPenalty: 20}}
	ApplyStatChange(c, "Fitness", 0, "Test") // Re-apply to enforce cap
	if c.Fitness != 80 {
		t.Errorf("Expected cap applied to 80, got %d", c.Fitness)
	}

	// Test applying positive delta respects new cap
	ApplyStatChange(c, "Fitness", 10, "Test")
	if c.Fitness != 80 {
		t.Errorf("Expected to stay at 80 cap, got %d", c.Fitness)
	}

	// Test min cap
	ApplyStatChange(c, "Fitness", -100, "Test")
	if c.Fitness != 0 {
		t.Errorf("Expected 0, got %d", c.Fitness)
	}
}

func TestAdvanceTimeBaseRecovery(t *testing.T) {
	setupTest()
	c := &Climber{Fitness: 50, AMS: 50, Loc: LocBase, O2Charges: 3}
	
	advanceTime(c)
	
	// BaseRecoveryFit is 6, BaseRecoveryAms is 8
	if c.Fitness != 56 {
		t.Errorf("Expected 56 Fitness, got %d", c.Fitness)
	}
	if c.AMS != 42 {
		t.Errorf("Expected 42 AMS, got %d", c.AMS)
	}
}

func TestNightTimePenalty(t *testing.T) {
	setupTest()
	c1 := &Climber{Fitness: 50, AMS: 50, Loc: LocCamp1, Altitude: 6000, O2Charges: 3} // Day
	session.Hour = 12
	session.Rng = rand.New(rand.NewSource(1))
	advanceTime(c1)
	
	c2 := &Climber{Fitness: 50, AMS: 50, Loc: LocCamp1, Altitude: 6000, O2Charges: 3} // Night
	session.Hour = 18
	session.Rng = rand.New(rand.NewSource(1))
	advanceTime(c2)

	// Night should have more fitness loss and more AMS gain
	if c2.Fitness >= c1.Fitness {
		t.Errorf("Expected night to reduce fitness more than day: Day:%d, Night:%d", c1.Fitness, c2.Fitness)
	}
	if c2.AMS <= c1.AMS {
		t.Errorf("Expected night to increase AMS more than day: Day:%d, Night:%d", c1.AMS, c2.AMS)
	}
}

func TestDeathZoneResting(t *testing.T) {
	setupTest()
	// Both resting, but one in death zone, one below.
	c1 := &Climber{Fitness: 50, Loc: LocHighCamp, Altitude: 7950, Resting: true, O2Charges: 3}
	c2 := &Climber{Fitness: 50, Loc: LocSummit, Altitude: 8848, Resting: true, O2Charges: 3}
	
	session.Hour = 12 // make it day
	advanceTime(c1)
	advanceTime(c2)

	// In Death Zone, resting gives fitloss = 4. Below, fitloss = 1.
	if c1.Fitness != 49 { 
		t.Errorf("Expected c1 fitness to lose exactly 1 (49), got %d", c1.Fitness)
	}
	if c2.Fitness != 46 { 
		t.Errorf("Expected c2 fitness to lose exactly 4 (46), got %d", c2.Fitness)
	}
}

func TestOxygenConsumption(t *testing.T) {
	setupTest()
	c := &Climber{Fitness: 90, Loc: LocCamp2, Altitude: 6500, O2Charges: 1}
	session.CampO2[LocCamp2] = 1 // 1 bottle available

	// Turn 1: Consume last charge
	advanceTime(c)
	if c.O2Charges != 3 {
		t.Errorf("Expected auto-reload to 3 charges, got %d", c.O2Charges)
	}
	if session.CampO2[LocCamp2] != 0 {
		t.Errorf("Expected camp bottle count to be 0, got %d", session.CampO2[LocCamp2])
	}

	// Turn 2: Consume 1 charge
	advanceTime(c)
	if c.O2Charges != 2 {
		t.Errorf("Expected charges to be 2, got %d", c.O2Charges)
	}

	// Turn 3: No more bottles at camp, charges hit 0
	c.O2Charges = 1
	advanceTime(c)
	if c.O2Charges != 0 {
		t.Errorf("Expected charges to be 0, got %d", c.O2Charges)
	}
}
