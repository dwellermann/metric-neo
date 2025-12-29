package valueobjects

import "testing"

func TestNewVelocity(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		wantError bool
	}{
		{
			name:      "typical air rifle velocity",
			input:     175.5,
			wantError: false,
		},
		{
			name:      "zero velocity is valid (shot didn't fire)",
			input:     0,
			wantError: false,
		},
		{
			name:      "negative velocity is invalid",
			input:     -10.0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			velocity, err := NewVelocity(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewVelocity(%v) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Fatalf("NewVelocity(%v) unexpected error: %v", tt.input, err)
				}
				if velocity.MetersPerSecond() != tt.input {
					t.Errorf("NewVelocity(%v).MetersPerSecond() = %v, want %v",
						tt.input, velocity.MetersPerSecond(), tt.input)
				}
			}
		})
	}
}

// Test: 575 fps ≈ 175.26 m/s (typisch für Luftgewehr)
func TestVelocityFromFPS(t *testing.T) {
	fps := 575.0
	wantMps := 175.26

	velocity, err := VelocityFromFPS(fps)
	if err != nil {
		t.Fatalf("VelocityFromFPS(%v) unexpected error: %v", fps, err)
	}

	got := velocity.MetersPerSecond()
	tolerance := 0.01 // 1 cm/s Toleranz
	if diff := got - wantMps; diff > tolerance || diff < -tolerance {
		t.Errorf("VelocityFromFPS(%v).MetersPerSecond() = %v, want ~%v", fps, got, wantMps)
	}
}

func TestVelocity_String(t *testing.T) {
	velocity, _ := NewVelocity(175.5)
	want := "175.50 m/s"
	got := velocity.String()
	if got != want {
		t.Errorf("Velocity.String() = %q, want %q", got, want)
	}
}
