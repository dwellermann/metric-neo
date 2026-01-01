package entities

import (
	"fmt"
	"metric-neo/internal/domain/valueobjects"

	"github.com/google/uuid"
)

// SightingSystemType ist ein Enum für verschiedene Zielvorrichtungen.
//
// GO-KONZEPT: Enums in Go
// Go hat KEINE native Enum-Unterstützung wie TypeScript/Java!
// Stattdessen nutzt man einen Custom Type + Konstanten.
//
// Vorteile dieses Patterns:
// - Type Safety: Man kann nicht versehentlich einen String zuweisen
// - Autocomplete in IDEs funktioniert
// - Validierung ist einfach
type SightingSystemType string

// GO-KONZEPT: Konstanten-Block mit iota
// Wir definieren alle erlaubten Werte als Konstanten.
// Der "const (...)" Block gruppiert zusammengehörige Konstanten.
const (
	SightingTypeScope      SightingSystemType = "scope"       // Zielfernrohr
	SightingTypeRedDot     SightingSystemType = "red_dot"     // Red Dot Visier
	SightingTypeDiopter    SightingSystemType = "diopter"     // Diopter (Präzisionsvisier)
	SightingTypeOpenSights SightingSystemType = "open_sights" // Kimme & Korn
)

// IsValid prüft, ob der SightingSystemType gültig ist.
//
// GO-KONZEPT: Validation Method für Enums
// Da Go keine Enum-Validierung hat, schreiben wir sie selbst.
func (t SightingSystemType) IsValid() bool {
	switch t {
	case SightingTypeScope, SightingTypeRedDot, SightingTypeDiopter, SightingTypeOpenSights:
		return true
	}
	return false
}

// String implementiert fmt.Stringer für schöne Ausgabe.
func (t SightingSystemType) String() string {
	// Map für lesbare Namen
	names := map[SightingSystemType]string{
		SightingTypeScope:      "Zielfernrohr",
		SightingTypeRedDot:     "Red Dot",
		SightingTypeDiopter:    "Diopter",
		SightingTypeOpenSights: "Kimme & Korn",
	}

	if name, ok := names[t]; ok {
		return name
	}
	return string(t) // Fallback
}

// SightingSystem repräsentiert eine Zielvorrichtung (Optik/Visierung).
//
// DOMAIN MODEL: Siehe docs/specs/domain-model.md Abschnitt 3.3
type SightingSystem struct {
	ID   string             `json:"id"`
	Type SightingSystemType `json:"type"`

	// ModelName: z.B. "Schmidt & Bender PM II 5-25x56"
	ModelName string `json:"model_name"`

	// Weight: Gewicht der Optik (beeinflusst Balance des Gewehrs)
	Weight valueobjects.Mass `json:"weight"`

	// Magnification: Min/Max Vergrößerung
	// Bei Festbrennweite (z.B. 4x32) ist Min == Max
	MinMagnification valueobjects.Magnification `json:"min_magnification"`
	MaxMagnification valueobjects.Magnification `json:"max_magnification"`
}

// NewSightingSystem erstellt eine neue SightingSystem-Instanz.
//
// GO-KONZEPT: Constructor mit vielen Parametern
// In Go gibt es keine Named Parameters oder Optional Parameters wie in TypeScript!
// Für komplexe Objekte gibt es mehrere Patterns (siehe unten).
func NewSightingSystem(
	sightingType SightingSystemType,
	modelName string,
	weight valueobjects.Mass,
	minMag valueobjects.Magnification,
	maxMag valueobjects.Magnification,
) (*SightingSystem, error) {
	// Validierung: Type muss gültig sein
	if !sightingType.IsValid() {
		return nil, fmt.Errorf("invalid sighting system type: %s", sightingType)
	}

	// Validierung: ModelName darf nicht leer sein
	if modelName == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	// Validierung: Min darf nicht größer als Max sein
	if minMag.Factor() > maxMag.Factor() {
		return nil, fmt.Errorf("min magnification (%.1fx) cannot be greater than max (%.1fx)",
			minMag.Factor(), maxMag.Factor())
	}

	return &SightingSystem{
		ID:               uuid.New().String(),
		Type:             sightingType,
		ModelName:        modelName,
		Weight:           weight,
		MinMagnification: minMag,
		MaxMagnification: maxMag,
	}, nil
}

// NewFixedSightingSystem ist ein Helper für Festbrennweiten-Optiken.
//
// GO-KONZEPT: Constructor Variants (Alternative Constructors)
// Wenn ein Objekt auf verschiedene Arten erstellt werden kann,
// bieten wir mehrere "New*"-Funktionen an.
//
// Beispiel: "Hawke Vantage 4x32" hat feste 4x Vergrößerung
func NewFixedSightingSystem(
	sightingType SightingSystemType,
	modelName string,
	weight valueobjects.Mass,
	magnification valueobjects.Magnification,
) (*SightingSystem, error) {
	// Bei Festbrennweite: Min == Max
	return NewSightingSystem(sightingType, modelName, weight, magnification, magnification)
}

// NewIronSights ist ein Helper für Kimme & Korn (keine Vergrößerung).
//
// GO-KONZEPT: Domain-Specific Constructor
// Kimme & Korn hat IMMER 1x Vergrößerung.
// Dieser Constructor kodiert diese Regel.
func NewIronSights(modelName string, weight valueobjects.Mass) (*SightingSystem, error) {
	oneMag, _ := valueobjects.NewMagnification(1.0)
	return NewSightingSystem(
		SightingTypeOpenSights,
		modelName,
		weight,
		oneMag,
		oneMag,
	)
}

// IsVariable gibt an, ob die Optik variable Vergrößerung hat.
func (s *SightingSystem) IsVariable() bool {
	return s.MinMagnification.Factor() != s.MaxMagnification.Factor()
}

// MagnificationRange gibt den Vergrößerungsbereich als String zurück.
// Beispiel: "3-9x" oder "4x" (bei Festbrennweite)
func (s *SightingSystem) MagnificationRange() string {
	if s.IsVariable() {
		return fmt.Sprintf("%.1f-%.1fx", s.MinMagnification.Factor(), s.MaxMagnification.Factor())
	}
	return fmt.Sprintf("%.1fx", s.MinMagnification.Factor())
}

// String implementiert fmt.Stringer.
func (s *SightingSystem) String() string {
	return fmt.Sprintf("%s %s (%s, %s)",
		s.Type.String(),
		s.ModelName,
		s.MagnificationRange(),
		s.Weight.String(),
	)
}
