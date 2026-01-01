package valueobjects

import (
	"encoding/json"
	"fmt"
	"math"
)

// Energy repräsentiert kinetische Energie in Joule.
//
// WICHTIG: Energy wird nie direkt erstellt, sondern IMMER berechnet!
// Siehe CalculateEnergy() Funktion unten.
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

// UnmarshalJSON implementiert json.Unmarshaler Interface.
// HINWEIS: Energy sollte normalerweise nie direkt deserialisiert werden,
// da es immer berechnet wird. Wir implementieren es trotzdem für Vollständigkeit.
func (e *Energy) UnmarshalJSON(data []byte) error {
	var joules float64
	if err := json.Unmarshal(data, &joules); err != nil {
		return err
	}
	*e = Energy(joules)
	return nil
}

// CalculateEnergy berechnet kinetische Energie nach E = 1/2 * m * v²
//
// GO-KONZEPT: Package-Level Function (kein Method)
// Diese Funktion gehört zum Package, nicht zu einem Typ, weil sie
// ZWEI verschiedene Value Objects kombiniert (Mass + Velocity).
//
// GO-KONZEPT: Named Return Value
// Wir könnten auch einfach "Energy" schreiben, aber "energy Energy" ist
// expliziter und dokumentiert den Return-Wert.
func CalculateEnergy(mass Mass, velocity Velocity) Energy {
	// Formel: E = 1/2 * m * v²
	// m in kg, v in m/s -> Ergebnis in Joule

	m := mass.Kilograms()           // Konvertierung zu kg
	v := velocity.MetersPerSecond() // m/s

	// GO-KONZEPT: math.Pow() für Exponenten
	// Go hat keinen ** Operator wie Python. Stattdessen math.Pow(base, exponent)
	joules := 0.5 * m * math.Pow(v, 2)

	return Energy(joules)
}
