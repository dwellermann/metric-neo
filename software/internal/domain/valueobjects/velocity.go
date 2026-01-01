package valueobjects

import (
	"encoding/json"
	"fmt"
)

// Velocity repräsentiert eine Geschwindigkeit in Metern pro Sekunde (m/s).
//
// Wie bei Mass erstellen wir einen Custom Type für Type Safety und Semantik.
type Velocity float64

// NewVelocity erstellt eine neue Velocity mit Validierung.
func NewVelocity(metersPerSecond float64) (Velocity, error) {
	if metersPerSecond < 0 {
		// GO-KONZEPT: fmt.Errorf für formatierte Error-Messages
		// Ähnlich wie console.error() in JS, aber als Rückgabewert
		return 0, fmt.Errorf("velocity cannot be negative, got: %.2f m/s", metersPerSecond)
	}
	return Velocity(metersPerSecond), nil
}

// MetersPerSecond gibt die Geschwindigkeit in m/s zurück.
func (v Velocity) MetersPerSecond() float64 {
	return float64(v)
}

// VelocityFromFPS konvertiert Feet per Second zu m/s.
// 1 fps = 0.3048 m/s
//
// WICHTIG für Metric Neo: Die LMBR-Hardware sendet in fps,
// aber wir speichern IMMER in m/s (siehe domain-model.md Abschnitt 5.4).
func VelocityFromFPS(fps float64) (Velocity, error) {
	const fpsToMps = 0.3048
	return NewVelocity(fps * fpsToMps)
}

// String implementiert fmt.Stringer für schöne Ausgabe.
func (v Velocity) String() string {
	return fmt.Sprintf("%.2f m/s", v.MetersPerSecond())
}

// MarshalJSON implementiert json.Marshaler Interface.
func (v Velocity) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(v))
}

// UnmarshalJSON implementiert json.Unmarshaler Interface.
func (v *Velocity) UnmarshalJSON(data []byte) error {
	var mps float64
	if err := json.Unmarshal(data, &mps); err != nil {
		return err
	}
	
	velocity, err := NewVelocity(mps)
	if err != nil {
		return err
	}
	
	*v = velocity
	return nil
}
