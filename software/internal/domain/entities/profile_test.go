package entities

import (
	"encoding/json"
	"metric-neo/internal/domain/valueobjects"
	"testing"
)

func TestNewProfile(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0) // 450mm
	trigger, _ := valueobjects.NewMass(500.0)  // 500g Abzugsgewicht
	sight, _ := valueobjects.NewLength(45.0)   // 45mm Visierhöhe

	tests := []struct {
		name      string
		profName  string
		category  ProfileCategory
		wantError bool
	}{
		{
			name:      "valid air rifle",
			profName:  "Steyr Challenge E",
			category:  CategoryAirRifle,
			wantError: false,
		},
		{
			name:      "empty name invalid",
			profName:  "",
			category:  CategoryAirRifle,
			wantError: true,
		},
		{
			name:      "invalid category",
			profName:  "Test",
			category:  ProfileCategory("invalid"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := NewProfile(tt.profName, tt.category, barrel, trigger, sight)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if profile.ID == "" {
					t.Error("UUID not generated")
				}
				if profile.Name != tt.profName {
					t.Errorf("Name = %q, want %q", profile.Name, tt.profName)
				}
			}
		})
	}
}

func TestProfile_OptionalFields(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(500.0)
	sight, _ := valueobjects.NewLength(45.0)

	profile, _ := NewProfile("Test Rifle", CategoryAirRifle, barrel, trigger, sight)

	// Test: Optic anfangs nicht gesetzt
	if profile.HasOptic() {
		t.Error("New profile should not have optic")
	}

	if profile.Optic != nil {
		t.Error("Optic should be nil")
	}

	// Test: TwistRate anfangs nicht gesetzt
	if profile.HasTwistRate() {
		t.Error("New profile should not have twist rate")
	}

	// Test: DefaultAmmo anfangs nicht gesetzt
	if profile.HasDefaultAmmo() {
		t.Error("New profile should not have default ammo")
	}
}

func TestProfile_SetOptic(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(500.0)
	sight, _ := valueobjects.NewLength(45.0)

	profile, _ := NewProfile("Test Rifle", CategoryAirRifle, barrel, trigger, sight)

	// Erstelle Optik
	opticWeight, _ := valueobjects.NewMass(350.0)
	mag, _ := valueobjects.NewMagnification(4.0)
	optic, _ := NewFixedSightingSystem(SightingTypeScope, "Hawke 4x32", opticWeight, mag)

	// Setze Optik
	profile.SetOptic(optic)

	// Prüfe
	if !profile.HasOptic() {
		t.Error("Profile should have optic after SetOptic()")
	}

	if profile.Optic == nil {
		t.Fatal("Optic is nil after SetOptic()")
	}

	if profile.Optic.ModelName != "Hawke 4x32" {
		t.Errorf("Optic model = %q, want %q", profile.Optic.ModelName, "Hawke 4x32")
	}

	// Entferne Optik
	profile.RemoveOptic()

	if profile.HasOptic() {
		t.Error("Profile should not have optic after RemoveOptic()")
	}
}

func TestProfile_SetTwistRate(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(500.0)
	sight, _ := valueobjects.NewLength(45.0)

	profile, _ := NewProfile("Test Rifle", CategoryAirRifle, barrel, trigger, sight)

	// Setze TwistRate
	twist, _ := valueobjects.NewLength(254.0) // 1:10" twist = 254mm
	profile.SetTwistRate(twist)

	if !profile.HasTwistRate() {
		t.Error("Profile should have twist rate after SetTwistRate()")
	}

	if profile.TwistRate == nil {
		t.Fatal("TwistRate is nil")
	}

	if profile.TwistRate.Millimeters() != 254.0 {
		t.Errorf("TwistRate = %.1f mm, want 254.0 mm", profile.TwistRate.Millimeters())
	}
}

func TestProfile_TotalWeight(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(500.0)
	sight, _ := valueobjects.NewLength(45.0)

	profile, _ := NewProfile("Test Rifle", CategoryAirRifle, barrel, trigger, sight)

	// Ohne Optik
	weight1 := profile.TotalWeight()
	if weight1.Grams() != 500.0 {
		t.Errorf("Weight without optic = %.1fg, want 500.0g", weight1.Grams())
	}

	// Mit Optik
	opticWeight, _ := valueobjects.NewMass(350.0)
	mag, _ := valueobjects.NewMagnification(4.0)
	optic, _ := NewFixedSightingSystem(SightingTypeScope, "Hawke 4x32", opticWeight, mag)
	profile.SetOptic(optic)

	weight2 := profile.TotalWeight()
	expected := 500.0 + 350.0 // 850g
	if weight2.Grams() != expected {
		t.Errorf("Weight with optic = %.1fg, want %.1fg", weight2.Grams(), expected)
	}
}

func TestProfile_JSONMarshaling(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(2500.0) // Realistisches Gewehrgewicht
	sight, _ := valueobjects.NewLength(45.0)

	profile, _ := NewProfile("Steyr Challenge E", CategoryAirRifle, barrel, trigger, sight)

	// Füge Optik hinzu
	opticWeight, _ := valueobjects.NewMass(350.0)
	mag, _ := valueobjects.NewMagnification(4.0)
	optic, _ := NewFixedSightingSystem(SightingTypeScope, "Hawke 4x32", opticWeight, mag)
	profile.SetOptic(optic)

	// Marshal
	jsonData, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	t.Logf("JSON: %s", string(jsonData))

	// Unmarshal
	var decoded Profile
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	// Verify
	if decoded.Name != profile.Name {
		t.Errorf("Name: got %q, want %q", decoded.Name, profile.Name)
	}

	if !decoded.HasOptic() {
		t.Error("Decoded profile should have optic")
	}

	if decoded.Optic.ModelName != optic.ModelName {
		t.Errorf("Optic model: got %q, want %q", decoded.Optic.ModelName, optic.ModelName)
	}
}

func TestProfile_JSONMarshaling_WithoutOptic(t *testing.T) {
	barrel, _ := valueobjects.NewLength(450.0)
	trigger, _ := valueobjects.NewMass(2500.0)
	sight, _ := valueobjects.NewLength(45.0)

	profile, _ := NewProfile("Walther LG400", CategoryAirRifle, barrel, trigger, sight)

	// OHNE Optik (Kimme & Korn)
	jsonData, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	t.Logf("JSON without optic: %s", string(jsonData))

	// Prüfe, dass "optic" NICHT im JSON ist (wegen omitempty)
	// GO-KONZEPT: String contains check
	jsonStr := string(jsonData)
	if contains(jsonStr, "optic") {
		t.Error("JSON should not contain 'optic' field when nil (omitempty)")
	}

	// Unmarshal sollte trotzdem funktionieren
	var decoded Profile
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	if decoded.HasOptic() {
		t.Error("Decoded profile should not have optic")
	}
}

// Helper
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestProfile_RealWorldData(t *testing.T) {
	tests := []struct {
		name     string
		category ProfileCategory
		barrel   float64 // mm
		sight    float64 // mm
		hasOptic bool
	}{
		{
			name:     "Steyr Challenge E (match rifle)",
			category: CategoryAirRifle,
			barrel:   420.0,
			sight:    45.0,
			hasOptic: false, // Diopter
		},
		{
			name:     "Walther LG400 (scoped)",
			category: CategoryAirRifle,
			barrel:   450.0,
			sight:    50.0,
			hasOptic: true,
		},
		{
			name:     "Feinwerkbau P8X (match pistol)",
			category: CategoryAirPistol,
			barrel:   230.0,
			sight:    30.0,
			hasOptic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			barrel, _ := valueobjects.NewLength(tt.barrel)
			trigger, _ := valueobjects.NewMass(500.0)
			sight, _ := valueobjects.NewLength(tt.sight)

			profile, _ := NewProfile(tt.name, tt.category, barrel, trigger, sight)

			if tt.hasOptic {
				opticWeight, _ := valueobjects.NewMass(350.0)
				mag, _ := valueobjects.NewMagnification(4.0)
				optic, _ := NewFixedSightingSystem(SightingTypeScope, "Test Scope", opticWeight, mag)
				profile.SetOptic(optic)
			}

			t.Logf("Created: %s", profile.String())
		})
	}
}
