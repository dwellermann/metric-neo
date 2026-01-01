package valueobjects

import "fmt"

// Length repräsentiert eine Länge in Millimetern (mm).
//
// GO-KONZEPT: Custom Type für physikalische Einheiten
// Basis-Einheit laut domain-model.md: Millimeter
type Length float64

// NewLength erstellt eine neue Length aus einem Millimeter-Wert.
// Gibt einen Fehler zurück, wenn der Wert negativ ist.
func NewLength(mm float64) (Length, error) {
	if mm < 0 {
		return 0, fmt.Errorf("length cannot be negative, got: %.2f mm", mm)
	}
	return Length(mm), nil
}

// NewLengthFromInches erstellt eine Length aus einem Inch-Wert.
// Nützlich für US-Kataloge (z.B. "18 inch barrel").
func NewLengthFromInches(inches float64) (Length, error) {
	return NewLength(inches * 25.4) // 1 inch = 25.4 mm
}

// Millimeters gibt die Länge in Millimetern zurück (Basiseinheit).
//
// GO-KONZEPT: Method mit Value Receiver
// Syntax: (l Length) bindet die Funktion an den Length-Typ
func (l Length) Millimeters() float64 {
	return float64(l)
}

// Centimeters gibt die Länge in Zentimetern zurück.
func (l Length) Centimeters() float64 {
	return float64(l) / 10.0
}

// Meters gibt die Länge in Metern zurück.
func (l Length) Meters() float64 {
	return float64(l) / 1000.0
}

// Inches gibt die Länge in Inch zurück.
// 1 inch = 25.4 mm
func (l Length) Inches() float64 {
	return float64(l) / 25.4
}

// String implementiert fmt.Stringer für schöne Ausgabe.
// Zeigt die Basiseinheit (mm) an, nicht Meter!
func (l Length) String() string {
	return fmt.Sprintf("%.2f mm", l.Millimeters())
}
