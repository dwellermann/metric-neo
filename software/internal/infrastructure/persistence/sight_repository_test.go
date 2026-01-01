package persistence

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"os"
	"path/filepath"
	"testing"
)

func TestSightRepository_SaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSightRepository(tempDir)

	// Erstelle SightingSystem
	weight, _ := valueobjects.NewMass(850.0)
	minMag, _ := valueobjects.NewMagnification(6.0)
	maxMag, _ := valueobjects.NewMagnification(24.0)

	sight, err := entities.NewSightingSystem(
		entities.SightingTypeScope,
		"Schmidt & Bender PM II 6-24x50",
		weight,
		minMag,
		maxMag,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Speichere
	if err := repo.Save(sight); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Lade zurück
	loaded, err := repo.Load(sight.ID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Vergleiche
	if loaded.ID != sight.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, sight.ID)
	}

	if loaded.ModelName != "Schmidt & Bender PM II 6-24x50" {
		t.Errorf("ModelName mismatch: got %s", loaded.ModelName)
	}

	if loaded.Type != entities.SightingTypeScope {
		t.Errorf("Type mismatch: got %s", loaded.Type)
	}

	if loaded.Weight.Grams() != 850.0 {
		t.Errorf("Weight mismatch: got %.1fg", loaded.Weight.Grams())
	}

	if loaded.MinMagnification.Factor() != 6.0 {
		t.Errorf("MinMag mismatch: got %.1fx", loaded.MinMagnification.Factor())
	}

	if loaded.MaxMagnification.Factor() != 24.0 {
		t.Errorf("MaxMag mismatch: got %.1fx", loaded.MaxMagnification.Factor())
	}

	t.Logf("✓ Saved and loaded: %s (%.0fg, %s)",
		loaded.ModelName,
		loaded.Weight.Grams(),
		loaded.Type.String())
}

func TestSightRepository_List(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSightRepository(tempDir)

	// Erstelle mehrere Sights
	sights := []struct {
		name string
		typ  entities.SightingSystemType
	}{
		{"Hawke Vantage 4x32", entities.SightingTypeScope},
		{"Walther FT 8-32x56", entities.SightingTypeScope},
		{"Aimpoint Micro T-2", entities.SightingTypeRedDot},
	}

	for _, s := range sights {
		weight, _ := valueobjects.NewMass(500.0)
		mag, _ := valueobjects.NewMagnification(4.0)
		sight, _ := entities.NewFixedSightingSystem(s.typ, s.name, weight, mag)
		repo.Save(sight)
	}

	// Liste alle
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 sights, got %d", len(list))
	}

	t.Logf("✓ Listed %d sighting systems:", len(list))
	for _, s := range list {
		t.Logf("  - %s (%s)", s.ModelName, s.Type.String())
	}
}

func TestSightRepository_Delete(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSightRepository(tempDir)

	// Erstelle Sight
	weight, _ := valueobjects.NewMass(200.0)
	sight, _ := entities.NewIronSights("Standard Open Sights", weight)

	repo.Save(sight)

	// Prüfe dass es existiert
	list, _ := repo.List()
	if len(list) != 1 {
		t.Fatal("Sight not saved")
	}

	// Lösche
	if err := repo.Delete(sight.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Prüfe dass es weg ist
	list, _ = repo.List()
	if len(list) != 0 {
		t.Error("Sight not deleted")
	}

	// Zweites Löschen sollte fehlschlagen
	err := repo.Delete(sight.ID)
	if err == nil {
		t.Error("Delete should fail for non-existent sight")
	}

	t.Log("✓ Sight deleted successfully")
}

func TestSightRepository_LoadNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSightRepository(tempDir)

	_, err := repo.Load("non-existent-id")
	if err == nil {
		t.Error("Load should fail for non-existent ID")
	}

	t.Logf("✓ Error handling: %v", err)
}

func TestSightRepository_SaveNil(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSightRepository(tempDir)

	err := repo.Save(nil)
	if err == nil {
		t.Error("Save should fail for nil sight")
	}

	t.Logf("✓ Validation: %v", err)
}

func TestSightRepository_JSONFormat(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSightRepository(tempDir)

	// Erstelle Sight
	weight, _ := valueobjects.NewMass(750.0)
	minMag, _ := valueobjects.NewMagnification(3.0)
	maxMag, _ := valueobjects.NewMagnification(12.0)

	sight, _ := entities.NewSightingSystem(
		entities.SightingTypeScope,
		"Leupold VX-3 3-12x40",
		weight,
		minMag,
		maxMag,
	)

	repo.Save(sight)

	// Lese JSON direkt
	filename := filepath.Join(tempDir, sight.ID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("✓ JSON Format:\n%s", string(data))

	// Prüfe ob JSON eingerückt ist (Lesbarkeit)
	if len(data) < 100 {
		t.Error("JSON seems too short - might not be indented")
	}
}
