package valueobjects

import (
	"encoding/json"
	"fmt"
	"math"
)

type Energy float64

// Joules gibt die Energie in Joule zurück.
func (e Energy) Joules() float64 {
	return float64(e)
}

// String implementiert fmt.Stringer.
func (e Energy) String() string {
	return fmt.Sprintf("%.2f J", e.Joules())
}

// MarshalJSON implementiert json.Marshaler Interface.
func (e Energy) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(e))
}

func (e *Energy) UnmarshalJSON(data []byte) error {
	var joules float64
	if err := json.Unmarshal(data, &joules); err != nil {
		return err
	}
	*e = Energy(joules)
	return nil
}

func CalculateEnergy(mass Mass, velocity Velocity) Energy {
	// Formel: E = 1/2 * m * v²
	// m in kg, v in m/s -> Ergebnis in Joule

	m := mass.Kilograms()           // Konvertierung zu kg
	v := velocity.MetersPerSecond() // m/s

	joules := 0.5 * m * math.Pow(v, 2)

	return Energy(joules)
}
