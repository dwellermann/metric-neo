package entities

import (
	"encoding/json"
	"math"
	"metric-neo/internal/domain/valueobjects"
	"testing"
	"time"
)

func TestNewShot(t *testing.T) {
	velocity, _ := valueobjects.NewVelocity(175.5)

	shot := NewShot(velocity)

	if shot == nil {
		t.Fatal("NewShot() returned nil")
	}

	if shot.Velocity.MetersPerSecond() != 175.5 {
		t.Errorf("Velocity = %.2f, want 175.5", shot.Velocity.MetersPerSecond())
	}

	if !shot.Valid {
		t.Error("New shot should be valid by default")
	}

	// Timestamp sollte "jetzt" sein (innerhalb 1 Sekunde)
	if time.Since(shot.Timestamp) > time.Second {
		t.Error("Timestamp is not recent")
	}
}

func TestNewShotAt(t *testing.T) {
	velocity, _ := valueobjects.NewVelocity(175.5)
	timestamp := time.Date(2025, 12, 1, 10, 30, 0, 0, time.UTC)

	shot := NewShotAt(velocity, timestamp)

	if !shot.Timestamp.Equal(timestamp) {
		t.Errorf("Timestamp = %v, want %v", shot.Timestamp, timestamp)
	}
}

func TestShot_CalculateEnergy(t *testing.T) {
	velocity, _ := valueobjects.NewVelocity(175.0)
	shot := NewShot(velocity)

	mass, _ := valueobjects.NewMass(0.547)

	energy := shot.CalculateEnergy(mass)

	// E = 0.5 * 0.000547 * 175² ≈ 8.38 J
	expected := 8.38
	tolerance := 0.05

	if math.Abs(energy.Joules()-expected) > tolerance {
		t.Errorf("Energy = %.2f J, want ~%.2f J", energy.Joules(), expected)
	}
}

func TestShot_MarkInvalid(t *testing.T) {
	velocity, _ := valueobjects.NewVelocity(175.5)
	shot := NewShot(velocity)

	if !shot.Valid {
		t.Fatal("Shot should start valid")
	}

	shot.MarkInvalid()

	if shot.Valid {
		t.Error("Shot should be invalid after MarkInvalid()")
	}

	// Kann wieder gültig markiert werden
	shot.MarkValid()

	if !shot.Valid {
		t.Error("Shot should be valid after MarkValid()")
	}
}

func TestShot_ElapsedSince(t *testing.T) {
	velocity, _ := valueobjects.NewVelocity(175.0)

	t1 := time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 12, 1, 10, 0, 5, 0, time.UTC) // 5 Sekunden später

	shot1 := NewShotAt(velocity, t1)
	shot2 := NewShotAt(velocity, t2)

	elapsed := shot2.ElapsedSince(shot1)

	if elapsed != 5*time.Second {
		t.Errorf("Elapsed = %v, want 5s", elapsed)
	}
}

func TestShot_JSONMarshaling(t *testing.T) {
	velocity, _ := valueobjects.NewVelocity(175.5)
	timestamp := time.Date(2025, 12, 1, 10, 30, 0, 0, time.UTC)

	original := NewShotAt(velocity, timestamp)

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	t.Logf("JSON: %s", string(jsonData))

	var decoded Shot
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	// Timestamps in Go vergleichen mit Equal()
	if !decoded.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Timestamp: got %v, want %v", decoded.Timestamp, original.Timestamp)
	}

	if decoded.Velocity.MetersPerSecond() != original.Velocity.MetersPerSecond() {
		t.Errorf("Velocity mismatch")
	}

	if decoded.Valid != original.Valid {
		t.Errorf("Valid mismatch")
	}
}

// Helper - wird in session_test.go wiederverwendet
func createTestProfile() *Profile {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(2500.0)
	sight, _ := valueobjects.NewLength(45.0)
	profile, _ := NewProfile("Test Rifle", CategoryAirRifle, barrel, trigger, sight)
	return profile
}

func createTestProjectile() *Projectile {
	weight, _ := valueobjects.NewMass(0.547)
	projectile, _ := NewProjectile("JSB Exact", weight, 0.022)
	return projectile
}
