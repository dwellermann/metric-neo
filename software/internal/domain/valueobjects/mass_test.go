package valueobjects

import (
	"encoding/json"
	"testing"
)

// GO-KONZEPT: Testing in Go
// - Test-Funktionen müssen mit "Test" beginnen
// - Der Parameter *testing.T gibt uns Zugriff auf Test-Utilities (t.Error, t.Fatal, etc.)
// - Tests leben im gleichen Package wie der zu testende Code

func TestNewMass(t *testing.T) {
	// GO-KONZEPT: Table-Driven Tests
	// Das ist Go's idiomatischer Weg, mehrere Test-Cases zu strukturieren.
	// Vorteile:
	// - Übersichtlich
	// - Leicht erweiterbar
	// - Zeigt alle Fehler auf einmal (nicht nur den ersten)

	tests := []struct {
		name      string  // Beschreibung des Test-Cases
		input     float64 // Input-Wert
		wantError bool    // Erwarten wir einen Fehler?
	}{
		{
			name:      "valid mass",
			input:     0.547,
			wantError: false,
		},
		{
			name:      "zero mass is invalid",
			input:     0,
			wantError: true,
		},
		{
			name:      "negative mass is invalid",
			input:     -5.0,
			wantError: true,
		},
	}

	// GO-KONZEPT: Range Loop
	// "range" iteriert über Slices, Maps, Arrays, Channels
	// Syntax: for index, value := range collection
	for _, tt := range tests {
		// GO-KONZEPT: Subtests
		// t.Run() erlaubt hierarchische Tests mit eigenen Namen
		// Das macht die Ausgabe von "go test -v" viel übersichtlicher
		t.Run(tt.name, func(t *testing.T) {
			mass, err := NewMass(tt.input)

			// GO-KONZEPT: Error Checking Pattern
			// In Go prüfen wir immer: "Hat die Funktion einen Error zurückgegeben?"
			if tt.wantError {
				if err == nil {
					t.Errorf("NewMass(%v) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Fatalf("NewMass(%v) unexpected error: %v", tt.input, err)
				}
				if mass.Grams() != tt.input {
					t.Errorf("NewMass(%v).Grams() = %v, want %v", tt.input, mass.Grams(), tt.input)
				}
			}
		})
	}
}

func TestMass_Kilograms(t *testing.T) {
	mass, _ := NewMass(547.0) // 547 Gramm
	want := 0.547             // kg

	got := mass.Kilograms()

	// GO-KONZEPT: Float Comparison
	// Floats sollten nie mit == verglichen werden (Rundungsfehler!)
	// Hier ist die Toleranz klein genug, aber für komplexere Fälle
	// würde man eine Epsilon-Funktion nutzen.
	if got != want {
		t.Errorf("Mass.Kilograms() = %v, want %v", got, want)
	}
}

func TestMassFromGrain(t *testing.T) {
	// Test: 8.44 grain = 0.547 gramm (typisches Diabolo-Gewicht)
	grain := 8.44
	wantGrams := 0.547

	mass, err := MassFromGrain(grain)
	if err != nil {
		t.Fatalf("MassFromGrain(%v) unexpected error: %v", grain, err)
	}

	got := mass.Grams()
	// Toleranz für Rundungsfehler
	tolerance := 0.001
	if diff := got - wantGrams; diff > tolerance || diff < -tolerance {
		t.Errorf("MassFromGrain(%v).Grams() = %v, want ~%v", grain, got, wantGrams)
	}
}

func TestMass_String(t *testing.T) {
	mass, _ := NewMass(0.547)
	want := "0.547 g"

	got := mass.String()
	if got != want {
		t.Errorf("Mass.String() = %q, want %q", got, want)
	}
}

// GO-KONZEPT: JSON Marshaling Test (Round-Trip)
// Ein "Round-Trip" Test prüft: Original → JSON → Decoded → Muss gleich sein
func TestMass_JSONMarshaling(t *testing.T) {
	original, _ := NewMass(0.547)

	// GO-KONZEPT: Marshal (Go → JSON)
	// json.Marshal() konvertiert Go-Daten zu JSON-Bytes
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	// Erwartetes JSON: einfache Zahl, nicht {"Mass": 0.547}
	expected := "0.547"
	got := string(jsonData)
	if got != expected {
		t.Errorf("JSON output = %s, want %s", got, expected)
	}

	// GO-KONZEPT: Unmarshal (JSON → Go)
	// json.Unmarshal() braucht einen POINTER!
	var decoded Mass
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	// Prüfe, ob Original == Decoded
	if decoded.Grams() != original.Grams() {
		t.Errorf("Round-trip failed: got %v, want %v", decoded.Grams(), original.Grams())
	}
}

// Test: Ungültige JSON-Werte werden abgelehnt
func TestMass_UnmarshalJSON_Invalid(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"negative mass", "-5.0"},
		{"zero mass", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Mass
			err := json.Unmarshal([]byte(tt.json), &m)
			if err == nil {
				t.Errorf("Expected error for JSON %s, got nil", tt.json)
			}
		})
	}
}
