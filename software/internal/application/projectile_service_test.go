package application

import (
	"testing"
)

func TestProjectileService_CreateProjectile(t *testing.T) {
	dir := t.TempDir()
	service := NewProjectileService(dir)

	result := service.CreateProjectile("JSB Exact 4.52mm", 0.547, 0.024)

	if !result.Success {
		t.Fatalf("CreateProjectile failed: %s", result.Error)
	}

	dto := result.Data
	if dto.Name != "JSB Exact 4.52mm" {
		t.Errorf("Name = %s", dto.Name)
	}

	if dto.WeightGrams != 0.547 {
		t.Errorf("Weight = %.3f, want 0.547", dto.WeightGrams)
	}

	t.Logf("✓ Created: %s (%.3fg, BC: %.3f)", dto.Name, dto.WeightGrams, dto.BC)
}

func TestProjectileService_Validation(t *testing.T) {
	dir := t.TempDir()
	service := NewProjectileService(dir)

	// Test: Leerer Name
	result := service.CreateProjectile("", 0.547, 0.024)
	if result.Success {
		t.Error("Expected error for empty name")
	}

	// Test: Negativer BC
	result = service.CreateProjectile("Test", 0.547, -0.1)
	if result.Success {
		t.Error("Expected error for negative BC")
	}

	// Test: BC > 1
	result = service.CreateProjectile("Test", 0.547, 1.5)
	if result.Success {
		t.Error("Expected error for BC > 1")
	}

	t.Log("✓ Validation working")
}

func TestProjectileService_UpdateBC(t *testing.T) {
	dir := t.TempDir()
	service := NewProjectileService(dir)

	// Erstelle Projectile
	createResult := service.CreateProjectile("Test Pellet", 0.547, 0.020)
	projectileID := createResult.Data.ID

	// Update BC
	updateResult := service.UpdateBC(projectileID, 0.025)
	if !updateResult.Success {
		t.Fatalf("UpdateBC failed: %s", updateResult.Error)
	}

	if updateResult.Data.BC != 0.025 {
		t.Errorf("BC = %.3f, want 0.025", updateResult.Data.BC)
	}

	t.Logf("✓ BC updated: %.3f → %.3f", 0.020, updateResult.Data.BC)
}

func TestProjectileService_ListProjectiles(t *testing.T) {
	dir := t.TempDir()
	service := NewProjectileService(dir)

	// Erstelle mehrere Projectiles
	service.CreateProjectile("JSB Exact", 0.547, 0.024)
	service.CreateProjectile("H&N Baracuda", 0.690, 0.029)
	service.CreateProjectile("RWS Meisterkugeln", 0.454, 0.018)

	// Liste
	listResult := service.ListProjectiles()
	if !listResult.Success {
		t.Fatalf("ListProjectiles failed: %s", listResult.Error)
	}

	if len(listResult.Data) != 3 {
		t.Errorf("Count = %d, want 3", len(listResult.Data))
	}

	for _, p := range listResult.Data {
		t.Logf("  - %s (%.3fg, BC: %.3f)", p.Name, p.WeightGrams, p.BC)
	}
}
