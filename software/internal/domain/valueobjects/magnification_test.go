package valueobjects

import (
	"testing"
)

// GO-KONZEPT: Table-Driven Tests
// Tests für Magnification Value Object

func TestNewMagnification(t *testing.T) {
	tests := []struct {
		name    string
		factor  float64
		wantErr bool
	}{
		{
			name:    "no magnification (iron sights, red dot)",
			factor:  1.0,
			wantErr: false,
		},
		{
			name:    "typical air rifle scope",
			factor:  4.0,
			wantErr: false,
		},
		{
			name:    "variable scope at 6x",
			factor:  6.0,
			wantErr: false,
		},
		{
			name:    "high magnification scope",
			factor:  12.5,
			wantErr: false,
		},
		{
			name:    "extreme long-range scope",
			factor:  25.0,
			wantErr: false,
		},
		{
			name:    "below 1x is invalid",
			factor:  0.5,
			wantErr: true,
		},
		{
			name:    "zero magnification is invalid",
			factor:  0.0,
			wantErr: true,
		},
		{
			name:    "negative magnification is invalid",
			factor:  -4.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mag, err := NewMagnification(tt.factor)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewMagnification(%v) expected error, got nil", tt.factor)
				}
			} else {
				if err != nil {
					t.Fatalf("NewMagnification(%v) unexpected error: %v", tt.factor, err)
				}
				if mag.Factor() != tt.factor {
					t.Errorf("NewMagnification(%v).Factor() = %v, want %v",
						tt.factor, mag.Factor(), tt.factor)
				}
			}
		})
	}
}

func TestMagnification_String(t *testing.T) {
	tests := []struct {
		name     string
		mag      Magnification
		expected string
	}{
		{
			name:     "no magnification",
			mag:      Magnification(1.0),
			expected: "1.0x",
		},
		{
			name:     "typical air rifle scope",
			mag:      Magnification(4.0),
			expected: "4.0x",
		},
		{
			name:     "variable scope with decimal",
			mag:      Magnification(6.5),
			expected: "6.5x",
		},
		{
			name:     "high magnification",
			mag:      Magnification(12.5),
			expected: "12.5x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mag.String()
			if got != tt.expected {
				t.Errorf("Magnification.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// GO-KONZEPT: Edge Case Testing
// Testet den Grenzwert bei genau 1.0x
func TestMagnification_EdgeCases(t *testing.T) {
	t.Run("exactly 1.0x is valid", func(t *testing.T) {
		mag, err := NewMagnification(1.0)
		if err != nil {
			t.Fatalf("NewMagnification(1.0) unexpected error: %v", err)
		}
		if mag.Factor() != 1.0 {
			t.Errorf("got %.1fx, want 1.0x", mag.Factor())
		}
	})

	t.Run("slightly below 1.0x is invalid", func(t *testing.T) {
		_, err := NewMagnification(0.999)
		if err == nil {
			t.Error("NewMagnification(0.999) expected error, got nil")
		}
	})
}

// GO-KONZEPT: Realistische Testdaten
// Testet mit echten Scope-Spezifikationen
func TestMagnification_RealWorldScopes(t *testing.T) {
	tests := []struct {
		name   string
		factor float64
		scope  string // Welches echte Zielfernrohr
	}{
		{
			name:   "Hawke Vantage 4x32 (fixed)",
			factor: 4.0,
			scope:  "Fixed 4x scope for air rifles",
		},
		{
			name:   "Walther 3-9x44 (set to 6x)",
			factor: 6.0,
			scope:  "Variable scope, mid-range setting",
		},
		{
			name:   "Steyr Challenge E iron sights",
			factor: 1.0,
			scope:  "No magnification (iron sights)",
		},
		{
			name:   "Schmidt & Bender 12.5-50x (max)",
			factor: 50.0,
			scope:  "Extreme long-range precision scope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mag, err := NewMagnification(tt.factor)
			if err != nil {
				t.Fatalf("NewMagnification(%v) for %s failed: %v",
					tt.factor, tt.scope, err)
			}

			if mag.Factor() != tt.factor {
				t.Errorf("Expected %vx for %s, got %vx",
					tt.factor, tt.scope, mag.Factor())
			}
		})
	}
}
