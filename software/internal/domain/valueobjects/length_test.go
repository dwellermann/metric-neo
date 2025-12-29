package valueobjects

import (
	"math"
	"testing"
)

// GO-KONZEPT: Table-Driven Tests
// Tests für Length Value Object

func TestNewLength(t *testing.T) {
	tests := []struct {
		name    string
		mm      float64
		wantErr bool
	}{
		{
			name:    "typical rifle barrel length",
			mm:      450.0,
			wantErr: false,
		},
		{
			name:    "sight height on air rifle",
			mm:      45.0,
			wantErr: false,
		},
		{
			name:    "caliber 4.5mm",
			mm:      4.5,
			wantErr: false,
		},
		{
			name:    "zero length is valid",
			mm:      0.0,
			wantErr: false,
		},
		{
			name:    "negative length is invalid",
			mm:      -10.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			length, err := NewLength(tt.mm)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewLength(%v) expected error, got nil", tt.mm)
				}
			} else {
				if err != nil {
					t.Fatalf("NewLength(%v) unexpected error: %v", tt.mm, err)
				}
				if length.Millimeters() != tt.mm {
					t.Errorf("NewLength(%v).Millimeters() = %v, want %v",
						tt.mm, length.Millimeters(), tt.mm)
				}
			}
		})
	}
}

func TestLengthConversions(t *testing.T) {
	// 450 mm Lauflänge (typisch für Luftgewehr)
	length := Length(450.0)

	// Test Millimeters (Basiseinheit)
	if got := length.Millimeters(); got != 450.0 {
		t.Errorf("Millimeters() = %.2f, want 450.00", got)
	}

	// Test Centimeters
	wantCm := 45.0
	if got := length.Centimeters(); got != wantCm {
		t.Errorf("Centimeters() = %.2f, want %.2f", got, wantCm)
	}

	// Test Meters
	wantM := 0.45
	if got := length.Meters(); math.Abs(got-wantM) > 0.0001 {
		t.Errorf("Meters() = %.3f, want %.3f", got, wantM)
	}

	// Test Inches (450 mm ≈ 17.72 inches)
	wantInches := 17.7165
	tolerance := 0.001
	if got := length.Inches(); math.Abs(got-wantInches) > tolerance {
		t.Errorf("Inches() = %.4f, want ~%.4f", got, wantInches)
	}
}

func TestNewLengthFromInches(t *testing.T) {
	tests := []struct {
		name      string
		inches    float64
		wantMm    float64
		tolerance float64
	}{
		{
			name:      "18 inch barrel (US catalog)",
			inches:    18.0,
			wantMm:    457.2,
			tolerance: 0.01,
		},
		{
			name:      "typical pistol barrel 4 inches",
			inches:    4.0,
			wantMm:    101.6,
			tolerance: 0.01,
		},
		{
			name:      "zero inches",
			inches:    0.0,
			wantMm:    0.0,
			tolerance: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			length, err := NewLengthFromInches(tt.inches)
			if err != nil {
				t.Fatalf("NewLengthFromInches(%v) unexpected error: %v", tt.inches, err)
			}

			got := length.Millimeters()
			if math.Abs(got-tt.wantMm) > tt.tolerance {
				t.Errorf("NewLengthFromInches(%v).Millimeters() = %.2f, want ~%.2f",
					tt.inches, got, tt.wantMm)
			}
		})
	}
}

func TestNewLengthFromInches_Negative(t *testing.T) {
	// Negative Inches sollten Fehler werfen
	_, err := NewLengthFromInches(-5.0)
	if err == nil {
		t.Error("NewLengthFromInches(-5.0) expected error, got nil")
	}
}

func TestLength_String(t *testing.T) {
	tests := []struct {
		name     string
		length   Length
		expected string
	}{
		{
			name:     "caliber 4.5mm",
			length:   Length(4.5),
			expected: "4.50 mm",
		},
		{
			name:     "rifle barrel 450mm",
			length:   Length(450.0),
			expected: "450.00 mm",
		},
		{
			name:     "zero length",
			length:   Length(0.0),
			expected: "0.00 mm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.length.String()
			if got != tt.expected {
				t.Errorf("Length.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// GO-KONZEPT: Benchmark Test
// Zeigt die Performance von Konvertierungen.
// Ausführen mit: go test -bench=.
func BenchmarkLengthConversions(b *testing.B) {
	length := Length(450.0)

	for i := 0; i < b.N; i++ {
		_ = length.Inches()
	}
}
