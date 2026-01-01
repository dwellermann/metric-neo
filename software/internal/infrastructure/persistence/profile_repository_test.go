package persistence

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"testing"
)

func TestProfileRepository_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	repo := NewProfileRepository(dir)

	// Erstelle Test-Profile
	profile := createTestProfileForRepo()

	// Save
	err := repo.Save(profile)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load
	loaded, err := repo.Load(profile.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify
	if loaded.ID != profile.ID {
		t.Errorf("ID: got %s, want %s", loaded.ID, profile.ID)
	}

	if loaded.Name != profile.Name {
		t.Errorf("Name: got %s, want %s", loaded.Name, profile.Name)
	}

	if loaded.TriggerWeight.Grams() != profile.TriggerWeight.Grams() {
		t.Errorf("TriggerWeight mismatch")
	}
}

func TestProfileRepository_WithOptic(t *testing.T) {
	dir := t.TempDir()
	repo := NewProfileRepository(dir)

	// Erstelle Profile mit Optik
	profile := createTestProfileForRepo()

	opticWeight, _ := valueobjects.NewMass(450.0)
	mag, _ := valueobjects.NewMagnification(6.0)
	optic, _ := entities.NewFixedSightingSystem(entities.SightingTypeScope, "Walther FT 8x56", opticWeight, mag)
	profile.SetOptic(optic)

	// Save
	err := repo.Save(profile)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load
	loaded, err := repo.Load(profile.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify Optic
	if !loaded.HasOptic() {
		t.Fatal("Optic was not saved")
	}

	if loaded.Optic.ModelName != "Walther FT 8x56" {
		t.Errorf("Optic name: got %s, want Walther FT 8x56", loaded.Optic.ModelName)
	}

	if loaded.Optic.MaxMagnification.Factor() != 6.0 {
		t.Error("Magnification mismatch")
	}
}

func TestProfileRepository_List(t *testing.T) {
	dir := t.TempDir()
	repo := NewProfileRepository(dir)

	// Speichere mehrere Profiles
	p1 := createTestProfileForRepo()
	p2 := createTestProfileForRepo()

	repo.Save(p1)
	repo.Save(p2)

	// List
	ids, err := repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(ids) != 2 {
		t.Errorf("List count: got %d, want 2", len(ids))
	}
}

func TestProfileRepository_Delete(t *testing.T) {
	dir := t.TempDir()
	repo := NewProfileRepository(dir)

	profile := createTestProfileForRepo()

	// Save
	repo.Save(profile)

	// Delete
	err := repo.Delete(profile.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Load sollte fehlschlagen
	_, err = repo.Load(profile.ID)
	if err == nil {
		t.Error("Load() should fail after delete")
	}
}

func TestProfileRepository_LoadNonExistent(t *testing.T) {
	dir := t.TempDir()
	repo := NewProfileRepository(dir)

	_, err := repo.Load("non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent profile")
	}
}

func TestProfileRepository_SaveNil(t *testing.T) {
	dir := t.TempDir()
	repo := NewProfileRepository(dir)

	err := repo.Save(nil)
	if err == nil {
		t.Error("Expected error when saving nil profile")
	}
}

// Helper: Erstelle Test-Profile
func createTestProfileForRepo() *entities.Profile {
	barrelLength, _ := valueobjects.NewLength(420.0)
	triggerWeight, _ := valueobjects.NewMass(500.0)
	sightHeight, _ := valueobjects.NewLength(50.0)

	profile, _ := entities.NewProfile(
		"Walther LG400",
		entities.CategoryAirRifle,
		barrelLength,
		triggerWeight,
		sightHeight,
	)
	return profile
}
