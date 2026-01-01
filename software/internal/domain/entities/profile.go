package entities

import (
	"fmt"
	"metric-neo/internal/domain/valueobjects"

	"github.com/google/uuid"
)

// ProfileCategory ist ein Enum für Sportgeräte-Kategorien.
type ProfileCategory string

const (
	CategoryAirRifle  ProfileCategory = "air_rifle"
	CategoryAirPistol ProfileCategory = "air_pistol"
	CategoryBow       ProfileCategory = "bow"
	CategoryFirearm   ProfileCategory = "firearm"
)

// IsValid prüft, ob die ProfileCategory gültig ist.
func (c ProfileCategory) IsValid() bool {
	switch c {
	case CategoryAirRifle, CategoryAirPistol, CategoryBow, CategoryFirearm:
		return true
	}
	return false
}

// String implementiert fmt.Stringer.
func (c ProfileCategory) String() string {
	names := map[ProfileCategory]string{
		CategoryAirRifle:  "Luftgewehr",
		CategoryAirPistol: "Luftpistole",
		CategoryBow:       "Bogen",
		CategoryFirearm:   "Feuerwaffe",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return string(c)
}

// DOMAIN MODEL: Siehe docs/specs/domain-model.md Abschnitt 3.2
type Profile struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"` // z.B. "Steyr Challenge E"
	Category ProfileCategory `json:"category"`

	// Hardware-Specs
	BarrelLength valueobjects.Length `json:"barrel_length"`

	// GO-KONZEPT: Optional Field mit Pointer
	// TwistRate ist OPTIONAL (nicht alle Waffen haben Drall)
	// nil = kein Drall (z.B. Glattrohr)
	// *Length = Pointer, kann nil sein
	TwistRate *valueobjects.Length `json:"twist_rate,omitempty"`

	TriggerWeight valueobjects.Mass   `json:"trigger_weight"`
	SightHeight   valueobjects.Length `json:"sight_height"`

	// GO-KONZEPT: Nested Struct als Pointer (Optional)
	// Optic ist OPTIONAL - manche Profile haben keine Optik (Kimme & Korn)
	// nil = keine Optik
	// *SightingSystem = Pointer zu SightingSystem
	//
	// WICHTIG: Das "omitempty" JSON-Tag bedeutet:
	// "Wenn nil, dann nicht im JSON ausgeben"
	Optic *SightingSystem `json:"optic,omitempty"`

	// GO-KONZEPT: Pointer zu anderem Entity
	// DefaultAmmo ist eine REFERENZ zu einem Projectile
	// Wird als ID gespeichert, nicht als ganzes Objekt!
	DefaultAmmoID *string `json:"default_ammo_id,omitempty"`
}

// NewProfile erstellt ein neues Profile.
//
// GO-KONZEPT: Optional Parameters Pattern
// Problem: Go hat keine Optional Parameters!
// Lösung: Wir nutzen Pointer für optionale Werte.
//   - nil = nicht gesetzt
//   - &value = gesetzt
func NewProfile(
	name string,
	category ProfileCategory,
	barrelLength valueobjects.Length,
	triggerWeight valueobjects.Mass,
	sightHeight valueobjects.Length,
) (*Profile, error) {
	// Validierung
	if name == "" {
		return nil, fmt.Errorf("profile name cannot be empty")
	}

	if !category.IsValid() {
		return nil, fmt.Errorf("invalid profile category: %s", category)
	}

	return &Profile{
		ID:            uuid.New().String(),
		Name:          name,
		Category:      category,
		BarrelLength:  barrelLength,
		TriggerWeight: triggerWeight,
		SightHeight:   sightHeight,
		// Optic bleibt nil (optional)
		// TwistRate bleibt nil (optional)
		// DefaultAmmoID bleibt nil (optional)
	}, nil
}

// SetOptic setzt die Optik.
//
// GO-KONZEPT: Setter für Optional Fields
// Wir nutzen eine Methode statt direktem Feldzugriff,
// um Validierung zu ermöglichen.
func (p *Profile) SetOptic(optic *SightingSystem) {
	p.Optic = optic
}

// RemoveOptic entfernt die Optik.
func (p *Profile) RemoveOptic() {
	p.Optic = nil
}

// HasOptic prüft, ob eine Optik montiert ist.
//
// GO-KONZEPT: Nil-Check für Pointer
// In Go muss man IMMER prüfen, ob ein Pointer nil ist,
// bevor man ihn dereferenziert!
func (p *Profile) HasOptic() bool {
	return p.Optic != nil
}

// SetTwistRate setzt die Drallänge (optional).
func (p *Profile) SetTwistRate(twistRate valueobjects.Length) {
	p.TwistRate = &twistRate
}

// HasTwistRate prüft, ob eine Drallänge gesetzt ist.
func (p *Profile) HasTwistRate() bool {
	return p.TwistRate != nil
}

// SetDefaultAmmo setzt die Standard-Munition.
func (p *Profile) SetDefaultAmmo(projectileID string) {
	p.DefaultAmmoID = &projectileID
}

// HasDefaultAmmo prüft, ob Standard-Munition gesetzt ist.
func (p *Profile) HasDefaultAmmo() bool {
	return p.DefaultAmmoID != nil
}

// TotalWeight berechnet das Gesamtgewicht (Waffe + Optik).
//
// GO-KONZEPT: Calculated Property mit Nil-Check
func (p *Profile) TotalWeight() valueobjects.Mass {
	// Basisgewicht: TriggerWeight ist ein Proxy für Waffengewicht
	// (in Realität würde man ein separates Feld haben)
	baseWeight := p.TriggerWeight

	// Addiere Optik-Gewicht, wenn vorhanden
	if p.HasOptic() {
		// GO-KONZEPT: Value Object Addition
		// Wir können nicht einfach Mass + Mass schreiben!
		// Wir müssen die Werte extrahieren und neu erstellen.
		totalGrams := baseWeight.Grams() + p.Optic.Weight.Grams()
		total, _ := valueobjects.NewMass(totalGrams)
		return total
	}

	return baseWeight
}

// String implementiert fmt.Stringer.
func (p *Profile) String() string {
	opticInfo := "ohne Optik"
	if p.HasOptic() {
		opticInfo = p.Optic.ModelName
	}

	return fmt.Sprintf("%s (%s) - %s",
		p.Name,
		p.Category.String(),
		opticInfo,
	)
}
