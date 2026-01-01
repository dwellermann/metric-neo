package persistence

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"testing"
)

func TestSessionRepository_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	repo := NewSessionRepository(dir)

	// Erstelle Test-Session
	profile := createTestProfile()
	projectile := createTestProjectile()
	session := entities.NewSession(profile, projectile)

	// Füge Schüsse hinzu
	v1, _ := valueobjects.NewVelocity(175.0)
	v2, _ := valueobjects.NewVelocity(176.0)
	session.RecordShot(v1)
	session.RecordShot(v2)

	// Save
	err := repo.Save(session)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load
	loaded, err := repo.Load(session.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify
	if loaded.ID != session.ID {
		t.Errorf("ID mismatch")
	}

	if loaded.ShotCount() != 2 {
		t.Errorf("Shot count: got %d, want 2", loaded.ShotCount())
	}

	if loaded.ProfileSnapshot.Name != profile.Name {
		t.Error("Profile snapshot name mismatch")
	}

	if loaded.ProjectileSnapshot.Name != projectile.Name {
		t.Error("Projectile snapshot name mismatch")
	}
}

// KRITISCHER TEST: Snapshot Isolation nach Load
func TestSessionRepository_SnapshotIsolationAfterLoad(t *testing.T) {
	dir := t.TempDir()
	repo := NewSessionRepository(dir)

	// Erstelle Session mit Snapshots
	profile := createTestProfile()
	projectile := createTestProjectile()

	// Merke Original-Werte
	originalProfileName := profile.Name
	originalProjectileBC := projectile.BC

	session := entities.NewSession(profile, projectile)

	// Save
	err := repo.Save(session)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// ÄNDERE die Originals (simuliert Änderung in anderer Session)
	profile.Name = "MODIFIED PROFILE"
	projectile.BC = 999.0

	// Load
	loaded, err := repo.Load(session.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Prüfe: Geladene Session hat ORIGINALE Werte in Snapshots
	if loaded.ProfileSnapshot.Name != originalProfileName {
		t.Error("ProfileSnapshot changed after load! Should preserve original snapshot.")
	}

	if loaded.ProjectileSnapshot.BC != originalProjectileBC {
		t.Error("ProjectileSnapshot changed after load! Should preserve original snapshot.")
	}

	t.Log("✓ Snapshot isolation preserved after save/load cycle")
}

func TestSessionRepository_WithTemperatureAndNote(t *testing.T) {
	dir := t.TempDir()
	repo := NewSessionRepository(dir)

	session := entities.NewSession(createTestProfile(), createTestProjectile())

	// Setze optionale Felder
	temp, _ := valueobjects.NewTemperature(22.5)
	session.SetTemperature(temp)
	session.SetNote("Test session at range")

	// Save
	repo.Save(session)

	// Load
	loaded, err := repo.Load(session.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify
	if loaded.Temperature == nil {
		t.Fatal("Temperature was not saved")
	}

	if loaded.Temperature.Celsius() != 22.5 {
		t.Error("Temperature value mismatch")
	}

	if loaded.Note != "Test session at range" {
		t.Errorf("Note: got %s, want 'Test session at range'", loaded.Note)
	}
}

func TestSessionRepository_List(t *testing.T) {
	dir := t.TempDir()
	repo := NewSessionRepository(dir)

	// Speichere mehrere Sessions
	s1 := entities.NewSession(createTestProfile(), createTestProjectile())
	s2 := entities.NewSession(createTestProfile(), createTestProjectile())

	repo.Save(s1)
	repo.Save(s2)

	// List
	ids, err := repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(ids) != 2 {
		t.Errorf("List count: got %d, want 2", len(ids))
	}
}

func TestSessionRepository_Delete(t *testing.T) {
	dir := t.TempDir()
	repo := NewSessionRepository(dir)

	session := entities.NewSession(createTestProfile(), createTestProjectile())

	// Save
	repo.Save(session)

	// Delete
	err := repo.Delete(session.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Load sollte fehlschlagen
	_, err = repo.Load(session.ID)
	if err == nil {
		t.Error("Load() should fail after delete")
	}
}

func TestSessionRepository_SaveNil(t *testing.T) {
	dir := t.TempDir()
	repo := NewSessionRepository(dir)

	err := repo.Save(nil)
	if err == nil {
		t.Error("Expected error when saving nil session")
	}
}

// Helper functions
func createTestProfile() *entities.Profile {
	barrelLength, _ := valueobjects.NewLength(420.0)
	triggerWeight, _ := valueobjects.NewMass(500.0)
	sightHeight, _ := valueobjects.NewLength(50.0)

	profile, _ := entities.NewProfile(
		"Test Profile",
		entities.CategoryAirRifle,
		barrelLength,
		triggerWeight,
		sightHeight,
	)
	return profile
}

func createTestProjectile() *entities.Projectile {
	weight, _ := valueobjects.NewMass(0.547)
	projectile, _ := entities.NewProjectile("JSB Exact 4.52mm", weight, 0.024)
	return projectile
}
