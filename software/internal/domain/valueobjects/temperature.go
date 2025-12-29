package valueobjects

import "fmt"

// Temperature repräsentiert eine Temperatur in Grad Celsius (°C).
type Temperature float64

// NewTemperature erstellt eine neue Temperature mit Validierung.
func NewTemperature(celsius float64) (Temperature, error) {
	if celsius < -273.15 {
		return 0, fmt.Errorf("temperature cannot be below absolute zero (-273.15 °C), got: %.2f °C", celsius)
	}
	return Temperature(celsius), nil
}

// Celsius gibt die Temperatur in Grad Celsius zurück.
func (t Temperature) Celsius() float64 {
	return float64(t)
}

// Fahrenheit gibt die Temperatur in Grad Fahrenheit zurück.
// Formel: °F = °C × 9/5 + 32
func (t Temperature) Fahrenheit() float64 {
	return float64(t)*9.0/5.0 + 32.0
}

// String implementiert fmt.Stringer für schöne Ausgabe.
func (t Temperature) String() string {
	return fmt.Sprintf("%.2f °C", t.Celsius())
}
