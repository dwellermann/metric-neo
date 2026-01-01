package entities

import (
	"encoding/json"
	"metric-neo/internal/domain/valueobjects"
	"testing"
)

func TestSightingSystemType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		typ   SightingSystemType
		valid bool
	}{
		{"scope is valid", SightingTypeScope, true},
		{"red dot is valid", SightingTypeRedDot, true},
		{"diopter is valid", SightingTypeDiopter, true},
		{"open sights is valid", SightingTypeOpenSights, true},
		{"invalid type", SightingSystemType("laser"), false},
		{"empty type", SightingSystemType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.typ.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestNewSightingSystem(t *testing.T) {
	weight, _ := valueobjects.NewMass(850.0) // Gramm
	minMag, _ := valueobjects.NewMagnification(3.0)
	maxMag, _ := valueobjects.NewMagnification(9.0)

	tests := []struct {
		name      string
		typ       SightingSystemType
		model     string
		weight    valueobjects.Mass
		minMag    valueobjects.Magnification
		maxMag    valueobjects.Magnification
		wantError bool
	}{
		{
			name:      "valid variable scope",
			typ:       SightingTypeScope,
			model:     "Walther 3-9x44",
			weight:    weight,
			minMag:    minMag,
			maxMag:    maxMag,
			wantError: false,
		},
		{
			name:      "invalid type",
			typ:       SightingSystemType("invalid"),
			model:     "Test",
			weight:    weight,
			minMag:    minMag,
			maxMag:    maxMag,
			wantError: true,
		},
		{
			name:      "empty model name",
			typ:       SightingTypeScope,
			model:     "",
			weight:    weight,
			minMag:    minMag,
			maxMag:    maxMag,
			wantError: true,
		},
		{
			name:      "min > max magnification",
			typ:       SightingTypeScope,
			model:     "Invalid Scope",
			weight:    weight,
			minMag:    maxMag, // 9x
			maxMag:    minMag, // 3x - falsch herum!
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sight, err := NewSightingSystem(tt.typ, tt.model, tt.weight, tt.minMag, tt.maxMag)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if sight == nil {
					t.Fatal("got nil sighting system")
				}
				if sight.ID == "" {
					t.Error("UUID not generated")
				}
			}
		})
	}
}

func TestNewFixedSightingSystem(t *testing.T) {
	weight, _ := valueobjects.NewMass(350.0)
	mag, _ := valueobjects.NewMagnification(4.0)

	sight, err := NewFixedSightingSystem(SightingTypeScope, "Hawke Vantage 4x32", weight, mag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Bei Festbrennweite: Min == Max
	if sight.MinMagnification.Factor() != sight.MaxMagnification.Factor() {
		t.Errorf("Fixed scope should have Min == Max, got %.1fx-%.1fx",
			sight.MinMagnification.Factor(), sight.MaxMagnification.Factor())
	}

	if sight.MinMagnification.Factor() != 4.0 {
		t.Errorf("Magnification = %.1fx, want 4.0x", sight.MinMagnification.Factor())
	}

	if sight.IsVariable() {
		t.Error("Fixed scope should not be variable")
	}
}

func TestNewIronSights(t *testing.T) {
	weight, _ := valueobjects.NewMass(50.0) // Kimme & Korn sind sehr leicht

	sight, err := NewIronSights("Steyr Challenge E Diopter", weight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Kimme & Korn hat immer 1x
	if sight.MinMagnification.Factor() != 1.0 {
		t.Errorf("Iron sights magnification = %.1fx, want 1.0x", sight.MinMagnification.Factor())
	}

	if sight.Type != SightingTypeOpenSights {
		t.Errorf("Type = %v, want %v", sight.Type, SightingTypeOpenSights)
	}
}

func TestSightingSystem_IsVariable(t *testing.T) {
	weight, _ := valueobjects.NewMass(350.0)

	// Variable Optik
	min, _ := valueobjects.NewMagnification(3.0)
	max, _ := valueobjects.NewMagnification(9.0)
	variable, _ := NewSightingSystem(SightingTypeScope, "Walther 3-9x44", weight, min, max)

	if !variable.IsVariable() {
		t.Error("3-9x scope should be variable")
	}

	// Festbrennweite
	fixed, _ := valueobjects.NewMagnification(4.0)
	fixedScope, _ := NewFixedSightingSystem(SightingTypeScope, "Hawke 4x32", weight, fixed)

	if fixedScope.IsVariable() {
		t.Error("4x scope should not be variable")
	}
}

func TestSightingSystem_MagnificationRange(t *testing.T) {
	weight, _ := valueobjects.NewMass(350.0)

	tests := []struct {
		name     string
		minMag   float64
		maxMag   float64
		expected string
	}{
		{"variable 3-9x", 3.0, 9.0, "3.0-9.0x"},
		{"variable 6-24x", 6.0, 24.0, "6.0-24.0x"},
		{"fixed 4x", 4.0, 4.0, "4.0x"},
		{"fixed 1x (iron sights)", 1.0, 1.0, "1.0x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, _ := valueobjects.NewMagnification(tt.minMag)
			max, _ := valueobjects.NewMagnification(tt.maxMag)
			sight, _ := NewSightingSystem(SightingTypeScope, "Test", weight, min, max)

			got := sight.MagnificationRange()
			if got != tt.expected {
				t.Errorf("MagnificationRange() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSightingSystem_JSONMarshaling(t *testing.T) {
	weight, _ := valueobjects.NewMass(850.0)
	min, _ := valueobjects.NewMagnification(3.0)
	max, _ := valueobjects.NewMagnification(9.0)

	original, _ := NewSightingSystem(
		SightingTypeScope,
		"Walther 3-9x44",
		weight,
		min,
		max,
	)

	// Marshal
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	t.Logf("JSON: %s", string(jsonData))

	// Unmarshal
	var decoded SightingSystem
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	// Verify
	if decoded.ID != original.ID {
		t.Errorf("ID: got %q, want %q", decoded.ID, original.ID)
	}

	if decoded.Type != original.Type {
		t.Errorf("Type: got %v, want %v", decoded.Type, original.Type)
	}

	if decoded.ModelName != original.ModelName {
		t.Errorf("ModelName: got %q, want %q", decoded.ModelName, original.ModelName)
	}
}

func TestSightingSystem_RealWorldData(t *testing.T) {
	tests := []struct {
		name    string
		typ     SightingSystemType
		model   string
		weightG float64
		minMag  float64
		maxMag  float64
	}{
		{
			name:    "Hawke Vantage 4x32 (fixed air rifle)",
			typ:     SightingTypeScope,
			model:   "Hawke Vantage 4x32 AO",
			weightG: 350.0,
			minMag:  4.0,
			maxMag:  4.0,
		},
		{
			name:    "Walther 3-9x44 (variable)",
			typ:     SightingTypeScope,
			model:   "Walther 3-9x44 Sniper",
			weightG: 850.0,
			minMag:  3.0,
			maxMag:  9.0,
		},
		{
			name:    "Steyr Diopter (match grade)",
			typ:     SightingTypeDiopter,
			model:   "Steyr Challenge E Diopter",
			weightG: 120.0,
			minMag:  1.0,
			maxMag:  1.0,
		},
		{
			name:    "Aimpoint Micro H-2 (red dot)",
			typ:     SightingTypeRedDot,
			model:   "Aimpoint Micro H-2",
			weightG: 86.0,
			minMag:  1.0,
			maxMag:  1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight, _ := valueobjects.NewMass(tt.weightG)
			minMag, _ := valueobjects.NewMagnification(tt.minMag)
			maxMag, _ := valueobjects.NewMagnification(tt.maxMag)

			sight, err := NewSightingSystem(tt.typ, tt.model, weight, minMag, maxMag)
			if err != nil {
				t.Fatalf("NewSightingSystem() failed: %v", err)
			}

			t.Logf("Created: %s", sight.String())
		})
	}
}
