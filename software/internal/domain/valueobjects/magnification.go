package valueobjects

import "fmt"

type Magnification float64

func NewMagnification(factor float64) (Magnification, error) {
	if factor < 1 {
		return 0, fmt.Errorf("magnification must be at least 1x, got: %.2fx", factor)
	}
	return Magnification(factor), nil
}

func (m Magnification) Factor() float64 {
	return float64(m)
}

func (m Magnification) String() string {
	return fmt.Sprintf("%.1fx", m.Factor())
}
