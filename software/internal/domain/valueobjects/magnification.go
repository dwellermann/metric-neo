package valueobjects

import "fmt"

// Magnification repräsentiert die Vergrößerung eines Zielfernrohrs.
//
// GO-KONZEPT: Custom Type für physikalische Einheiten
// Basis-Einheit laut domain-model.md: Vergrößerungsfaktor (z.B. 4x, 10x)
type Magnification float64

// NewMagnification erstellt eine neue Magnification mit Validierung.
func NewMagnification(factor float64) (Magnification, error) {
	if factor < 1 {
		return 0, fmt.Errorf("magnification must be at least 1x, got: %.2fx", factor)
	}
	return Magnification(factor), nil
}

// Factor gibt den Vergrößerungsfaktor zurück.
func (m Magnification) Factor() float64 {
	return float64(m)
}

// String implementiert fmt.Stringer für schöne Ausgabe.
func (m Magnification) String() string {
	return fmt.Sprintf("%.1fx", m.Factor())
}
