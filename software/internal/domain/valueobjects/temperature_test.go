package valueobjects

import (
	"math"
	"testing"
)

// GO-KONZEPT: Table-Driven Tests
// Tests für Temperature Value Object

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		wantErr bool
	}{
		{
			name:    "typical room temperature",
			celsius: 21.5,
			wantErr: false,
		},
		{
			name:    "freezing point of water",
			celsius: 0.0,
			wantErr: false,
		},
		{
			name:    "boiling point of water",
			celsius: 100.0,
			wantErr: false,
		},
		{
			name:    "cold winter day",
			celsius: -15.0,
			wantErr: false,
		},
		{
			name:    "absolute zero is valid",
			celsius: -273.15,
			wantErr: false,
		},
		{
			name:    "below absolute zero is invalid",
			celsius: -273.16,
			wantErr: true,
		},
		{
			name:    "far below absolute zero",
			celsius: -300.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp, err := NewTemperature(tt.celsius)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTemperature(%v) expected error, got nil", tt.celsius)
				}
			} else {
				if err != nil {
					t.Fatalf("NewTemperature(%v) unexpected error: %v", tt.celsius, err)
				}
				if temp.Celsius() != tt.celsius {
					t.Errorf("NewTemperature(%v).Celsius() = %v, want %v",
						tt.celsius, temp.Celsius(), tt.celsius)
				}
			}
		})
	}
}

func TestTemperature_Fahrenheit(t *testing.T) {
	tests := []struct {
		name           string
		celsius        float64
		wantFahrenheit float64
		tolerance      float64
	}{
		{
			name:           "freezing point of water",
			celsius:        0.0,
			wantFahrenheit: 32.0,
			tolerance:      0.01,
		},
		{
			name:           "boiling point of water",
			celsius:        100.0,
			wantFahrenheit: 212.0,
			tolerance:      0.01,
		},
		{
			name:           "room temperature",
			celsius:        21.0,
			wantFahrenheit: 69.8,
			tolerance:      0.1,
		},
		{
			name:           "absolute zero",
			celsius:        -273.15,
			wantFahrenheit: -459.67,
			tolerance:      0.01,
		},
		{
			name:           "negative celsius",
			celsius:        -40.0,
			wantFahrenheit: -40.0, // -40°C = -40°F (Schnittpunkt!)
			tolerance:      0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := Temperature(tt.celsius)
			got := temp.Fahrenheit()

			if math.Abs(got-tt.wantFahrenheit) > tt.tolerance {
				t.Errorf("Temperature(%v).Fahrenheit() = %.2f, want ~%.2f",
					tt.celsius, got, tt.wantFahrenheit)
			}
		})
	}
}

func TestTemperature_String(t *testing.T) {
	tests := []struct {
		name     string
		temp     Temperature
		expected string
	}{
		{
			name:     "room temperature",
			temp:     Temperature(21.5),
			expected: "21.50 °C",
		},
		{
			name:     "freezing point",
			temp:     Temperature(0.0),
			expected: "0.00 °C",
		},
		{
			name:     "negative temperature",
			temp:     Temperature(-15.0),
			expected: "-15.00 °C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.temp.String()
			if got != tt.expected {
				t.Errorf("Temperature.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// GO-KONZEPT: Test für Edge Cases
// Testet Grenzfälle (Boundary Testing)
func TestTemperature_EdgeCases(t *testing.T) {
	t.Run("exactly at absolute zero", func(t *testing.T) {
		temp, err := NewTemperature(-273.15)
		if err != nil {
			t.Fatalf("NewTemperature(-273.15) unexpected error: %v", err)
		}
		if temp.Celsius() != -273.15 {
			t.Errorf("got %.2f, want -273.15", temp.Celsius())
		}
	})

	t.Run("one planck temperature unit below absolute zero", func(t *testing.T) {
		// Ein winziger Wert unter dem absoluten Nullpunkt
		_, err := NewTemperature(-273.15001)
		if err == nil {
			t.Error("NewTemperature(-273.15001) expected error, got nil")
		}
	})
}

// GO-KONZEPT: Beispiel-Test
// Zeigt in der Dokumentation, wie man Temperature benutzt.
// Ausführen mit: go test -v
func ExampleTemperature_Fahrenheit() {
	temp := Temperature(21.0)
	fahrenheit := temp.Fahrenheit()

	// Output zeigt das erwartete Ergebnis
	println(fahrenheit) // Ungefähr 69.8°F
}
