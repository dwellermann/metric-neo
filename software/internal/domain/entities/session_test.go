package entities

import (
	"encoding/json"
	"math"
	"metric-neo/internal/domain/valueobjects"
	"testing"
)

func TestNewSession(t *testing.T) {
	profile := createTestProfile()
	projectile := createTestProjectile()

	session := NewSession(profile, projectile)

	if session == nil {
		t.Fatal("NewSession() returned nil")
	}

	if session.ID == "" {
		t.Error("UUID not generated")
	}

	if session.ShotCount() != 0 {
		t.Errorf("New session should have 0 shots, got %d", session.ShotCount())
	}
}

// TEST: Snapshot Pattern - Kritischer Test!
func TestSession_SnapshotIsolation(t *testing.T) {
	profile := createTestProfile()
	projectile := createTestProjectile()

	session := NewSession(profile, projectile)

	// Merke Original-Werte
	originalProfileName := session.ProfileSnapshot.Name
	originalProjectileBC := session.ProjectileSnapshot.BC

	// ÄNDERE die Original-Entities
	profile.Name = "CHANGED NAME"
	projectile.BC = 999.0

	// Prüfe: Session-Snapshots sind UNVERÄNDERT!
	if session.ProfileSnapshot.Name != originalProfileName {
		t.Error("ProfileSnapshot was modified! Snapshot isolation broken!")
	}

	if session.ProjectileSnapshot.BC != originalProjectileBC {
		t.Error("ProjectileSnapshot was modified! Snapshot isolation broken!")
	}

	t.Log("✓ Snapshot isolation works - changes to originals don't affect session")
}

// TEST: Deep Copy bei nested Pointern
func TestSession_DeepCopyWithOptic(t *testing.T) {
	profile := createTestProfile()

	// Füge Optik hinzu
	opticWeight, _ := valueobjects.NewMass(350.0)
	mag, _ := valueobjects.NewMagnification(4.0)
	optic, _ := NewFixedSightingSystem(SightingTypeScope, "Original Scope", opticWeight, mag)
	profile.SetOptic(optic)

	projectile := createTestProjectile()
	session := NewSession(profile, projectile)

	// Merke Original
	originalOpticName := session.ProfileSnapshot.Optic.ModelName

	// Ändere Original-Optik
	profile.Optic.ModelName = "CHANGED SCOPE"

	// Prüfe: Session-Snapshot ist unverändert
	if session.ProfileSnapshot.Optic.ModelName != originalOpticName {
		t.Error("Optic in snapshot was modified! Deep copy failed!")
	}

	t.Log("✓ Deep copy works for nested pointers")
}

func TestSession_AddShot(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	velocity, _ := valueobjects.NewVelocity(175.5)
	shot := NewShot(velocity)

	session.AddShot(shot)

	if session.ShotCount() != 1 {
		t.Errorf("Shot count = %d, want 1", session.ShotCount())
	}
}

func TestSession_RecordShot(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	velocity, _ := valueobjects.NewVelocity(175.5)
	shot := session.RecordShot(velocity)

	if shot == nil {
		t.Fatal("RecordShot() returned nil")
	}

	if session.ShotCount() != 1 {
		t.Error("Shot was not added to session")
	}
}

func TestSession_CalculateAverageVelocity(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	// Füge 3 Schüsse hinzu
	session.RecordShot(mustVelocity(175.0))
	session.RecordShot(mustVelocity(176.0))
	session.RecordShot(mustVelocity(174.0))

	avg, err := session.CalculateAverageVelocity()
	if err != nil {
		t.Fatalf("CalculateAverageVelocity() failed: %v", err)
	}

	// Durchschnitt: (175 + 176 + 174) / 3 = 175.0
	expected := 175.0
	if math.Abs(avg.MetersPerSecond()-expected) > 0.01 {
		t.Errorf("Average = %.2f m/s, want %.2f m/s", avg.MetersPerSecond(), expected)
	}
}

func TestSession_CalculateAverageVelocity_OnlyValid(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	// 2 gültige, 1 ungültige
	session.RecordShot(mustVelocity(175.0))
	session.RecordShot(mustVelocity(176.0))
	invalidShot := session.RecordShot(mustVelocity(0.0)) // Fehlmessung
	invalidShot.MarkInvalid()

	if session.ValidShotCount() != 2 {
		t.Errorf("ValidShotCount = %d, want 2", session.ValidShotCount())
	}

	avg, err := session.CalculateAverageVelocity()
	if err != nil {
		t.Fatalf("CalculateAverageVelocity() failed: %v", err)
	}

	// Durchschnitt nur von gültigen: (175 + 176) / 2 = 175.5
	expected := 175.5
	if math.Abs(avg.MetersPerSecond()-expected) > 0.01 {
		t.Errorf("Average = %.2f m/s, want %.2f m/s", avg.MetersPerSecond(), expected)
	}
}

func TestSession_CalculateStandardDeviation(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	// Füge Schüsse mit bekannter SD hinzu
	// Werte: 170, 175, 180 -> Mean = 175, SD ≈ 4.08
	session.RecordShot(mustVelocity(170.0))
	session.RecordShot(mustVelocity(175.0))
	session.RecordShot(mustVelocity(180.0))

	sd, err := session.CalculateStandardDeviation()
	if err != nil {
		t.Fatalf("CalculateStandardDeviation() failed: %v", err)
	}

	expected := 4.08
	tolerance := 0.1

	if math.Abs(sd-expected) > tolerance {
		t.Errorf("SD = %.2f, want ~%.2f", sd, expected)
	}

	t.Logf("Standard Deviation: %.2f m/s", sd)
}

func TestSession_MinMaxVelocity(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	session.RecordShot(mustVelocity(170.0))
	session.RecordShot(mustVelocity(175.0))
	session.RecordShot(mustVelocity(180.0))
	session.RecordShot(mustVelocity(172.0))

	min, err := session.MinVelocity()
	if err != nil {
		t.Fatalf("MinVelocity() failed: %v", err)
	}

	if min.MetersPerSecond() != 170.0 {
		t.Errorf("Min = %.1f, want 170.0", min.MetersPerSecond())
	}

	max, err := session.MaxVelocity()
	if err != nil {
		t.Fatalf("MaxVelocity() failed: %v", err)
	}

	if max.MetersPerSecond() != 180.0 {
		t.Errorf("Max = %.1f, want 180.0", max.MetersPerSecond())
	}
}

func TestSession_ExtremeSpread(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	session.RecordShot(mustVelocity(170.0))
	session.RecordShot(mustVelocity(180.0))

	es, err := session.ExtremeSpread()
	if err != nil {
		t.Fatalf("ExtremeSpread() failed: %v", err)
	}

	expected := 10.0 // 180 - 170
	if math.Abs(es-expected) > 0.01 {
		t.Errorf("Extreme Spread = %.1f, want %.1f", es, expected)
	}

	t.Logf("Extreme Spread: %.1f m/s", es)
}

func TestSession_CalculateAverageEnergy(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	// JSB Exact (0.547g) bei verschiedenen Geschwindigkeiten
	session.RecordShot(mustVelocity(175.0))
	session.RecordShot(mustVelocity(176.0))
	session.RecordShot(mustVelocity(174.0))

	avgEnergy, err := session.CalculateAverageEnergy()
	if err != nil {
		t.Fatalf("CalculateAverageEnergy() failed: %v", err)
	}

	// Bei ~175 m/s und 0.547g: ~8.4 J
	if avgEnergy.Joules() < 8.0 || avgEnergy.Joules() > 9.0 {
		t.Errorf("Average energy = %.2f J, expected ~8.4 J", avgEnergy.Joules())
	}

	t.Logf("Average Energy: %.2f J", avgEnergy.Joules())
}

func TestSession_JSONMarshaling(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	// Füge einige Schüsse hinzu
	session.RecordShot(mustVelocity(175.0))
	session.RecordShot(mustVelocity(176.0))
	session.RecordShot(mustVelocity(174.0))

	// Setze Temperatur
	temp, _ := valueobjects.NewTemperature(21.5)
	session.SetTemperature(temp)

	// Setze Notiz
	session.SetNote("Test session for validation")

	// Marshal
	jsonData, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	t.Logf("JSON length: %d bytes", len(jsonData))

	// Unmarshal
	var decoded Session
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	// Verify
	if decoded.ID != session.ID {
		t.Errorf("ID mismatch")
	}

	if decoded.ShotCount() != session.ShotCount() {
		t.Errorf("Shot count: got %d, want %d", decoded.ShotCount(), session.ShotCount())
	}

	if decoded.Note != session.Note {
		t.Errorf("Note mismatch")
	}

	if decoded.ProfileSnapshot.Name != session.ProfileSnapshot.Name {
		t.Errorf("Profile snapshot name mismatch")
	}

	t.Log("✓ JSON round-trip successful")
}

func TestSession_EmptySession_Errors(t *testing.T) {
	session := NewSession(createTestProfile(), createTestProjectile())

	// Keine Schüsse -> alle Berechnungen sollten Fehler geben
	_, err := session.CalculateAverageVelocity()
	if err == nil {
		t.Error("Expected error for empty session")
	}

	_, err = session.CalculateStandardDeviation()
	if err == nil {
		t.Error("Expected error for SD with < 2 shots")
	}

	_, err = session.MinVelocity()
	if err == nil {
		t.Error("Expected error for min velocity")
	}
}

// Helper: Erstellt Velocity oder panic (für Tests)
func mustVelocity(mps float64) valueobjects.Velocity {
	v, err := valueobjects.NewVelocity(mps)
	if err != nil {
		panic(err)
	}
	return v
}
