package valueobjects

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Mass float64

func NewMass(grams float64) (Mass, error) {
	if grams <= 0 {
		// errors.New() erstellt einen einfachen Error
		return 0, errors.New("mass must be greater than zero")
	}
	return Mass(grams), nil
}

// - Mass ist klein (8 Bytes für float64)
// - Wir wollen Mass nicht verändern (immutable Value Object)
func (m Mass) Grams() float64 {
	return float64(m)
}

func (m Mass) Kilograms() float64 {
	return float64(m) / 1000.0
}

func MassFromGrain(grain float64) (Mass, error) {
	const grainToGram = 0.06479891
	return NewMass(grain * grainToGram)
}

func (m Mass) String() string {
	return fmt.Sprintf("%.3f g", m.Grams())
}

func (m Mass) MarshalJSON() ([]byte, error) {
	// json.Marshal() konvertiert den float64-Wert zu JSON
	return json.Marshal(float64(m))
}

func (m *Mass) UnmarshalJSON(data []byte) error {
	var grams float64

	// Dekodiere JSON-Zahl in float64
	if err := json.Unmarshal(data, &grams); err != nil {
		return err
	}

	// Nutze den Constructor für Validierung!
	// Das stellt sicher, dass ungültige Werte (z.B. -5) abgelehnt werden.
	mass, err := NewMass(grams)
	if err != nil {
		return err
	}

	*m = mass
	return nil
}
