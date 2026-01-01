package persistence

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"os"
	"path/filepath"
	"testing"
)

// GO-KONZEPT: Setup/Teardown Pattern mit Helper-Funktionen
// Diese Funktion erstellt ein temporäres Verzeichnis für Tests.
// Sie gibt auch eine Cleanup-Funktion zurück (Modern Go Pattern!).
func setupTestRepository(t *testing.T) (*ProjectileRepository, func()) {
	// GO-KONZEPT: t.TempDir() - Automatisches Temp-Verzeichnis
	// - Erstellt ein temp-Verzeichnis
	// - Löscht es automatisch nach Test-Ende
	// - Sehr praktisch für File-I/O Tests!
	tempDir := t.TempDir()

	repo := NewProjectileRepository(tempDir)

	// GO-KONZEPT: Cleanup Function Return
	// Der Aufrufer kann defer cleanup() nutzen
	cleanup := func() {
		// t.TempDir() räumt automatisch auf, aber wir könnten hier
		// zusätzliche Cleanup-Logik einfügen (z.B. Connections schließen)
	}

	return repo, cleanup
}

// Test: Save und Load Round-Trip
func TestProjectileRepository_SaveAndLoad(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	// Erstelle Test-Projectile
	weight, _ := valueobjects.NewMass(0.547)
	original, _ := entities.NewProjectile("JSB Exact 4.52", weight, 0.022)

	// Speichere
	err := repo.Save(original)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Lade zurück
	loaded, err := repo.Load(original.ID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Vergleiche alle Felder
	if loaded.ID != original.ID {
		t.Errorf("ID: got %q, want %q", loaded.ID, original.ID)
	}

	if loaded.Name != original.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, original.Name)
	}

	if loaded.Weight.Grams() != original.Weight.Grams() {
		t.Errorf("Weight: got %v, want %v", loaded.Weight.Grams(), original.Weight.Grams())
	}

	if loaded.BC != original.BC {
		t.Errorf("BC: got %v, want %v", loaded.BC, original.BC)
	}
}

// Test: Load nicht-existierende ID
func TestProjectileRepository_LoadNotFound(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	// Versuche, nicht-existierende ID zu laden
	_, err := repo.Load("non-existent-id")
	if err == nil {
		t.Error("Load() expected error for non-existent ID, got nil")
	}
}

// Test: List gibt alle Projectiles zurück
func TestProjectileRepository_List(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	// Leere Liste am Anfang
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("List() on empty repo: got %d items, want 0", len(list))
	}

	// Füge 3 Projectiles hinzu
	weight1, _ := valueobjects.NewMass(0.547)
	p1, _ := entities.NewProjectile("JSB Exact", weight1, 0.022)
	repo.Save(p1)

	weight2, _ := valueobjects.NewMass(0.690)
	p2, _ := entities.NewProjectile("H&N Baracuda", weight2, 0.029)
	repo.Save(p2)

	weight3, _ := valueobjects.NewMass(0.454)
	p3, _ := entities.NewProjectile("RWS Meisterkugeln", weight3, 0.018)
	repo.Save(p3)

	// Liste erneut
	list, err = repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	// GO-KONZEPT: Slice-Länge prüfen
	if len(list) != 3 {
		t.Errorf("List() returned %d items, want 3", len(list))
	}

	// Prüfe, dass alle IDs vorhanden sind
	// GO-KONZEPT: Map als Set für schnelle Lookup
	ids := make(map[string]bool)
	for _, p := range list {
		ids[p.ID] = true
	}

	if !ids[p1.ID] {
		t.Errorf("List() missing projectile %s", p1.ID)
	}
	if !ids[p2.ID] {
		t.Errorf("List() missing projectile %s", p2.ID)
	}
	if !ids[p3.ID] {
		t.Errorf("List() missing projectile %s", p3.ID)
	}
}

// Test: Delete löscht Projectile
func TestProjectileRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	// Erstelle und speichere
	weight, _ := valueobjects.NewMass(0.547)
	p, _ := entities.NewProjectile("Test Diabolo", weight, 0.022)
	repo.Save(p)

	// Prüfe, dass Datei existiert
	filename := filepath.Join(repo.dataDir, p.ID+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	// Lösche
	err := repo.Delete(p.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Prüfe, dass Datei nicht mehr existiert
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Error("File still exists after Delete()")
	}

	// Load sollte jetzt fehlschlagen
	_, err = repo.Load(p.ID)
	if err == nil {
		t.Error("Load() after Delete() should fail")
	}
}

// Test: Delete nicht-existierende ID
func TestProjectileRepository_DeleteNotFound(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	err := repo.Delete("non-existent-id")
	if err == nil {
		t.Error("Delete() expected error for non-existent ID, got nil")
	}
}

// GO-KONZEPT: Integration Test mit echten Daten
func TestProjectileRepository_Integration(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	t.Log("=== Integration Test: Full Workflow ===")

	// 1. Erstelle mehrere Projectiles
	projectiles := []*entities.Projectile{}

	weight1, _ := valueobjects.MassFromGrain(8.44)
	p1, _ := entities.NewProjectile("JSB Exact 4.52mm", weight1, 0.022)
	projectiles = append(projectiles, p1)

	weight2, _ := valueobjects.MassFromGrain(10.65)
	p2, _ := entities.NewProjectile("H&N Baracuda Match", weight2, 0.029)
	projectiles = append(projectiles, p2)

	weight3, _ := valueobjects.MassFromGrain(7.0)
	p3, _ := entities.NewProjectile("RWS Meisterkugeln", weight3, 0.018)
	projectiles = append(projectiles, p3)

	// 2. Speichere alle
	for _, p := range projectiles {
		if err := repo.Save(p); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
		t.Logf("Saved: %s", p.String())
	}

	// 3. Liste alle auf
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	t.Logf("Found %d projectiles in repository", len(list))

	// 4. Lösche eines
	if err := repo.Delete(p2.ID); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	t.Logf("Deleted: %s", p2.Name)

	// 5. Liste erneut - sollte nur noch 2 sein
	list, err = repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("After delete: got %d projectiles, want 2", len(list))
	}

	t.Log("✓ Integration test passed")
}
