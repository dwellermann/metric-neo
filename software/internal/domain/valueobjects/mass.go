package valueobjects

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Mass repräsentiert eine Masse in Gramm.
//
// GO-KONZEPT: Custom Type
// Wir erstellen einen neuen Typ basierend auf float64, um:
// 1. Type Safety zu gewährleisten (Mass kann nicht versehentlich mit Velocity gemischt werden)
// 2. Methoden an den Typ zu binden (OOP-Stil in Go)
// 3. Semantik im Code zu verbessern (Mass statt float64)
type Mass float64

// NewMass erstellt eine neue Mass-Instanz mit Validierung.
//
// GO-KONZEPT: Constructor Pattern
// Go hat keine Konstruktoren wie TypeScript/Java. Stattdessen nutzen wir
// Funktionen mit dem Präfix "New" als Konvention.
//
// GO-KONZEPT: Error Handling
// Go gibt Fehler als zweiten Return-Wert zurück (kein try/catch!).
// Der Aufrufer MUSS den Fehler prüfen.
func NewMass(grams float64) (Mass, error) {
	if grams <= 0 {
		// errors.New() erstellt einen einfachen Error
		return 0, errors.New("mass must be greater than zero")
	}
	return Mass(grams), nil
}

// Grams gibt die Masse in Gramm zurück.
//
// GO-KONZEPT: Method (Value Receiver)
// Die Syntax (m Mass) bindet diese Funktion an den Mass-Typ.
// Wir nutzen einen VALUE receiver (kein Pointer), weil:
// - Mass ist klein (8 Bytes für float64)
// - Wir wollen Mass nicht verändern (immutable Value Object)
func (m Mass) Grams() float64 {
	return float64(m)
}

// Kilograms gibt die Masse in Kilogramm zurück (für Energieberechnung).
func (m Mass) Kilograms() float64 {
	return float64(m) / 1000.0
}

// FromGrain erstellt eine Mass aus Grain (häufig in der Ballistik).
// 1 Grain = 0.06479891 Gramm
//
// GO-KONZEPT: Package-Level Function
// Diese Funktion gehört zum Package, nicht zu einem spezifischen Typ.
func MassFromGrain(grain float64) (Mass, error) {
	const grainToGram = 0.06479891
	return NewMass(grain * grainToGram)
}

// String implementiert das fmt.Stringer Interface.
//
// GO-KONZEPT: Interfaces
// In Go werden Interfaces IMPLIZIT implementiert (kein "implements" Keyword!).
// Wenn ein Typ eine Methode String() string hat, erfüllt er automatisch
// das fmt.Stringer Interface. Das ermöglicht schöne Ausgabe mit fmt.Println().
func (m Mass) String() string {
	return fmt.Sprintf("%.3f g", m.Grams())
}

// MarshalJSON implementiert json.Marshaler Interface.
//
// GO-KONZEPT: Custom JSON Serialization
// Ohne diese Methode würde Go Mass als {"Mass": 0.547} serialisieren.
// Mit dieser Methode wird es als einfache Zahl gespeichert: 0.547
//
// WICHTIG: Der Receiver ist (m Mass), KEIN Pointer!
// Bei Marshaling brauchen wir keinen Pointer, da wir nur lesen.
func (m Mass) MarshalJSON() ([]byte, error) {
	// json.Marshal() konvertiert den float64-Wert zu JSON
	return json.Marshal(float64(m))
}

// UnmarshalJSON implementiert json.Unmarshaler Interface.
//
// GO-KONZEPT: Pointer Receiver bei Unmarshal
// Der Receiver ist (m *Mass), ein POINTER!
// Warum? Weil wir den Wert von m ÄNDERN müssen.
// Ohne Pointer würden wir nur eine Kopie ändern.
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

	// Setze den Wert (daher *Mass Pointer)
	*m = mass
	return nil
}
