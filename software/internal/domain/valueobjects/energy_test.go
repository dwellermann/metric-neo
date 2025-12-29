package valueobjects

import (
	"math"
	"testing"
)

func TestCalculateEnergy(t *testing.T) {
	tests := []struct {
		name        string
		massGrams   float64
		velocityMps float64
		wantJoules  float64
	}{
		{
			name:        "JSB Exact 4.52mm at 175 m/s",
			massGrams:   0.547, // 8.44 grain
			velocityMps: 175.0,
			wantJoules:  8.38, // E = 0.5 * 0.000547 * 175² ≈ 8.38 J
		},
		{
			name:        "heavier diabolo at same velocity",
			massGrams:   0.685, // ~10.5 grain
			velocityMps: 175.0,
			wantJoules:  10.50, // E = 0.5 * 0.000685 * 175² ≈ 10.50 J
		},
		{
			name:        "zero velocity gives zero energy",
			massGrams:   0.547,
			velocityMps: 0,
			wantJoules:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mass, _ := NewMass(tt.massGrams)
			velocity, _ := NewVelocity(tt.velocityMps)

			energy := CalculateEnergy(mass, velocity)
			got := energy.Joules()

			// GO-KONZEPT: Floating-Point Vergleiche
			// Niemals float64-Werte direkt mit == vergleichen!
			// Nutze stattdessen eine Toleranz (epsilon).
			tolerance := 0.02 // 20 Millijoule Toleranz (für Rundungsfehler)

			if math.Abs(got-tt.wantJoules) > tolerance {
				t.Errorf("CalculateEnergy(%v g, %v m/s) = %v J, want ~%v J",
					tt.massGrams, tt.velocityMps, got, tt.wantJoules)
			}
		})
	}
}

func TestEnergy_String(t *testing.T) {
	energy := Energy(7.5)
	want := "7.50 J"
	got := energy.String()
	if got != want {
		t.Errorf("Energy.String() = %q, want %q", got, want)
	}
}
