package application

import (
	"testing"
)

func TestProfileService_CreateProfile(t *testing.T) {
	dir := t.TempDir()
	service := NewProfileService(dir)

	// Erstelle Profile
	result := service.CreateProfile(
		"Walther LG400",
		"air_rifle",
		420.0,
		500.0,
		50.0,
	)

	if !result.Success {
		t.Fatalf("CreateProfile failed: %s", result.Error)
	}

	dto := result.Data
	if dto.Name != "Walther LG400" {
		t.Errorf("Name = %s, want Walther LG400", dto.Name)
	}

	if dto.ID == "" {
		t.Error("UUID not generated")
	}

	t.Logf("✓ Created Profile: %s (ID: %s)", dto.Name, dto.ID)
}

func TestProfileService_CreateProfile_Validation(t *testing.T) {
	dir := t.TempDir()
	service := NewProfileService(dir)

	// Test 1: Leerer Name
	result := service.CreateProfile("", "air_rifle", 420.0, 500.0, 50.0)
	if result.Success {
		t.Error("Expected error for empty name")
	}
	t.Logf("✓ Validation: %s", result.Error)

	// Test 2: Negative Länge
	result = service.CreateProfile("Test", "air_rifle", -10.0, 500.0, 50.0)
	if result.Success {
		t.Error("Expected error for negative length")
	}
	t.Logf("✓ Validation: %s", result.Error)
}

func TestProfileService_LoadAndList(t *testing.T) {
	dir := t.TempDir()
	service := NewProfileService(dir)

	// Erstelle 2 Profile
	result1 := service.CreateProfile("Profile 1", "air_rifle", 420.0, 500.0, 50.0)
	result2 := service.CreateProfile("Profile 2", "air_pistol", 230.0, 500.0, 50.0)

	if !result1.Success || !result2.Success {
		t.Fatal("Failed to create profiles")
	}

	// Lade einzelnes Profile
	loadResult := service.LoadProfile(result1.Data.ID)
	if !loadResult.Success {
		t.Fatalf("LoadProfile failed: %s", loadResult.Error)
	}

	if loadResult.Data.Name != "Profile 1" {
		t.Errorf("Loaded wrong profile")
	}

	// Liste alle
	listResult := service.ListProfiles()
	if !listResult.Success {
		t.Fatalf("ListProfiles failed: %s", listResult.Error)
	}

	if len(listResult.Data) != 2 {
		t.Errorf("List count = %d, want 2", len(listResult.Data))
	}

	t.Logf("✓ Listed %d profiles", len(listResult.Data))
}

func TestProfileService_SetOptic(t *testing.T) {
	dir := t.TempDir()
	service := NewProfileService(dir)

	// Erstelle Profile
	createResult := service.CreateProfile("Test Profile", "air_rifle", 420.0, 500.0, 50.0)
	if !createResult.Success {
		t.Fatal("Failed to create profile")
	}

	profileID := createResult.Data.ID

	// Füge Optik hinzu
	result := service.SetOptic(
		profileID,
		"scope",
		"Walther FT 6-24x50",
		450.0,
		6.0,
		24.0,
	)

	if !result.Success {
		t.Fatalf("SetOptic failed: %s", result.Error)
	}

	dto := result.Data
	if dto.Optic == nil {
		t.Fatal("Optic not set")
	}

	if dto.Optic.ModelName != "Walther FT 6-24x50" {
		t.Errorf("Optic name = %s", dto.Optic.ModelName)
	}

	// TotalWeight sollte jetzt höher sein (Waffe + Optik)
	if dto.TotalWeightG <= 450.0 {
		t.Error("TotalWeight should include optic weight")
	}

	t.Logf("✓ Optic added: %s (Total weight: %.0fg)", dto.Optic.ModelName, dto.TotalWeightG)
}

func TestProfileService_DeleteProfile(t *testing.T) {
	dir := t.TempDir()
	service := NewProfileService(dir)

	// Erstelle Profile
	createResult := service.CreateProfile("Delete Me", "air_rifle", 420.0, 500.0, 50.0)
	profileID := createResult.Data.ID

	// Lösche
	deleteResult := service.DeleteProfile(profileID)
	if !deleteResult.Success {
		t.Fatalf("Delete failed: %s", deleteResult.Error)
	}

	// Load sollte fehlschlagen
	loadResult := service.LoadProfile(profileID)
	if loadResult.Success {
		t.Error("Profile should not exist after delete")
	}

	t.Log("✓ Profile deleted successfully")
}

func TestProfileService_Stateless(t *testing.T) {
	dir := t.TempDir()
	service := NewProfileService(dir)

	// Erstelle Profile
	createResult := service.CreateProfile("Stateless Test", "air_rifle", 420.0, 500.0, 50.0)
	profileID := createResult.Data.ID

	// Füge Optik hinzu
	service.SetOptic(profileID, "scope", "Test Scope", 400.0, 6.0, 24.0)

	// Neuer Service (simuliert App-Neustart)
	newService := NewProfileService(dir)

	// Lade Profile
	loadResult := newService.LoadProfile(profileID)
	if !loadResult.Success {
		t.Fatal("Load after 'restart' failed")
	}

	// Optik sollte noch da sein (wurde persistiert)
	if loadResult.Data.Optic == nil {
		t.Error("Optic not persisted")
	}

	t.Log("✓ Stateless pattern works - data persisted across service instances")
}
