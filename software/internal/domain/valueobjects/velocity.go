package valueobjects

import (
	"encoding/json"
	"fmt"
)

type Velocity float64

// NewVelocity erstellt eine neue Velocity mit Validierung.
func NewVelocity(metersPerSecond float64) (Velocity, error) {
	if metersPerSecond < 0 {
		return 0, fmt.Errorf("velocity cannot be negative, got: %.2f m/s", metersPerSecond)
	}
	return Velocity(metersPerSecond), nil
}

// MetersPerSecond gibt die Geschwindigkeit in m/s zurück.
func (v Velocity) MetersPerSecond() float64 {
	return float64(v)
}

func VelocityFromFPS(fps float64) (Velocity, error) {
	const fpsToMps = 0.3048
	return NewVelocity(fps * fpsToMps)
}

func (v Velocity) String() string {
	return fmt.Sprintf("%.2f m/s", v.MetersPerSecond())
}

func (v Velocity) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(v))
}

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
