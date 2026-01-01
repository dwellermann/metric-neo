package entities

import (
	"encoding/json"
	"metric-neo/internal/domain/valueobjects"
	"testing"
)

// GO-KONZEPT: Table-Driven Tests für Entity Constructor
func TestNewProjectile(t *testing.T) {
	weight, _ := valueobjects.NewMass(0.547) // 8.44 grain

	tests := []struct {
		name      string
		projName  string
		weight    valueobjects.Mass
		bc        float64
		wantError bool
	}{
		{
			name:      "valid JSB Exact diabolo",
			projName:  "JSB Exact 4.52",
			weight:    weight,
			bc:        0.022,
			wantError: false,
		},
		{
			name:      "valid with zero BC",
			projName:  "Unknown Diabolo",
			weight:    weight,
			bc:        0.0,
			wantError: false,
		},
		{
			name:      "empty name is invalid",
			projName:  "",
			weight:    weight,
			bc:        0.022,
			wantError: true,
		},
		{
			name:      "negative BC is invalid",
			projName:  "Bad Diabolo",
			weight:    weight,
			bc:        -0.5,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectile, err := NewProjectile(tt.projName, tt.weight, tt.bc)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewProjectile() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("NewProjectile() unexpected error: %v", err)
				}

				// GO-KONZEPT: Nil-Check für Pointer
				// Bei Pointern müssen wir prüfen, ob sie nil sind!
				if projectile == nil {
					t.Fatal("NewProjectile() returned nil without error")
				}

				// Prüfe, dass UUID generiert wurde
				if projectile.ID == "" {
					t.Error("NewProjectile() did not generate UUID")
				}

				// Prüfe, dass Werte korrekt gesetzt wurden
				if projectile.Name != tt.projName {
					t.Errorf("Name = %q, want %q", projectile.Name, tt.projName)
				}

				if projectile.Weight.Grams() != tt.weight.Grams() {
					t.Errorf("Weight = %v, want %v", projectile.Weight.Grams(), tt.weight.Grams())
				}

				if projectile.BC != tt.bc {
					t.Errorf("BC = %v, want %v", projectile.BC, tt.bc)
				}
			}
		})
	}
}

// GO-KONZEPT: Test für UUID-Eindeutigkeit
func TestNewProjectile_UniqueIDs(t *testing.T) {
	weight, _ := valueobjects.NewMass(0.547)

	// Erstelle 100 Projectiles
	// GO-KONZEPT: make() für Slices
	// make([]*Projectile, 0, 100) erstellt einen Slice mit:
	// - Länge 0 (leer)
	// - Capacity 100 (vorreservierter Speicher)
	projectiles := make([]*Projectile, 0, 100)

	for i := 0; i < 100; i++ {
		p, _ := NewProjectile("Test", weight, 0.022)
		projectiles = append(projectiles, p)
	}

	// Prüfe, dass alle IDs unterschiedlich sind
	// GO-KONZEPT: Map als Set
	// In Go gibt es kein Set. Wir nutzen map[string]bool:
	// - Key = ID
	// - Value = true (unwichtig, wir nutzen nur die Keys)
	seen := make(map[string]bool)

	for _, p := range projectiles {
		if seen[p.ID] {
			t.Errorf("Duplicate UUID detected: %s", p.ID)
		}
		seen[p.ID] = true
	}
}

// Test: UpdateBC Methode
func TestProjectile_UpdateBC(t *testing.T) {
	weight, _ := valueobjects.NewMass(0.547)
	p, _ := NewProjectile("JSB Exact", weight, 0.022)

	// Update mit gültigem Wert
	err := p.UpdateBC(0.025)
	if err != nil {
		t.Fatalf("UpdateBC(0.025) failed: %v", err)
	}

	if p.BC != 0.025 {
		t.Errorf("BC after update = %v, want 0.025", p.BC)
	}

	// Update mit ungültigem Wert (negativ)
	err = p.UpdateBC(-0.1)
	if err == nil {
		t.Error("UpdateBC(-0.1) expected error, got nil")
	}

	// BC sollte unverändert sein (0.025)
	if p.BC != 0.025 {
		t.Errorf("BC after failed update = %v, want 0.025 (unchanged)", p.BC)
	}
}

// Test: String() Methode
func TestProjectile_String(t *testing.T) {
	weight, _ := valueobjects.NewMass(0.547)
	p, _ := NewProjectile("JSB Exact 4.52", weight, 0.022)

	got := p.String()
	// String sollte Name, Gewicht und BC enthalten
	// Wir prüfen nicht den exakten Wortlaut, nur dass die Infos da sind
	if got == "" {
		t.Error("String() returned empty string")
	}

	// GO-KONZEPT: t.Logf() für Debug-Ausgabe
	// Wird nur bei go test -v angezeigt
	t.Logf("Projectile.String() = %q", got)
}

// GO-KONZEPT: JSON Round-Trip Test für Entity
func TestProjectile_JSONMarshaling(t *testing.T) {
	weight, _ := valueobjects.NewMass(0.547)
	original, _ := NewProjectile("JSB Exact 4.52", weight, 0.022)

	// Marshal zu JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	// GO-KONZEPT: t.Logf() zum Debuggen
	t.Logf("JSON: %s", string(jsonData))

	// Unmarshal zurück zu Projectile
	var decoded Projectile
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}

	// Prüfe alle Felder
	if decoded.ID != original.ID {
		t.Errorf("ID: got %q, want %q", decoded.ID, original.ID)
	}

	if decoded.Name != original.Name {
		t.Errorf("Name: got %q, want %q", decoded.Name, original.Name)
	}

	if decoded.Weight.Grams() != original.Weight.Grams() {
		t.Errorf("Weight: got %v, want %v", decoded.Weight.Grams(), original.Weight.Grams())
	}

	if decoded.BC != original.BC {
		t.Errorf("BC: got %v, want %v", decoded.BC, original.BC)
	}
}

// GO-KONZEPT: Test mit echten Munitionsdaten
func TestProjectile_RealWorldData(t *testing.T) {
	tests := []struct {
		name         string
		manufacturer string
		grains       float64
		bc           float64
	}{
		{
			name:         "JSB Exact 4.52mm (8.44gr)",
			manufacturer: "JSB Match Diabolo",
			grains:       8.44,
			bc:           0.022,
		},
		{
			name:         "H&N Baracuda Match (10.65gr)",
			manufacturer: "H&N Sport",
			grains:       10.65,
			bc:           0.029,
		},
		{
			name:         "RWS Meisterkugeln (7.0gr)",
			manufacturer: "RWS",
			grains:       7.0,
			bc:           0.018,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Konvertiere Grain zu Mass
			weight, err := valueobjects.MassFromGrain(tt.grains)
			if err != nil {
				t.Fatalf("MassFromGrain() failed: %v", err)
			}

			p, err := NewProjectile(tt.name, weight, tt.bc)
			if err != nil {
				t.Fatalf("NewProjectile() failed: %v", err)
			}

			t.Logf("Created: %s", p.String())
		})
	}
}
