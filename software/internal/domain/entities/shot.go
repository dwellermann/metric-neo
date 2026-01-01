package entities

import (
	"fmt"
	"metric-neo/internal/domain/valueobjects"
	"time"
)

// Shot repräsentiert eine einzelne Geschwindigkeitsmessung.
//
// GO-KONZEPT: Immutable Entity
// Ein Shot ist UNVERÄNDERLICH nach Erstellung!
// Es gibt keine Setter-Methoden (außer SetValid für Fehlerkorrektur).
//
// DOMAIN MODEL: Siehe docs/specs/domain-model.md Abschnitt 3.5
type Shot struct {
	// Timestamp: Exakter Zeitpunkt der Messung
	// GO-KONZEPT: time.Time Type
	// Go's Standard-Typ für Zeitstempel (Teil der stdlib)
	Timestamp time.Time `json:"timestamp"`

	// Velocity: Gemessene Geschwindigkeit (vom Chronographen)
	Velocity valueobjects.Velocity `json:"velocity"`

	// Valid: Markierung für Fehlmessungen
	// Beispiel: Diabolo verfehlt Sensor -> ungültig
	Valid bool `json:"valid"`
}

// NewShot erstellt einen neuen Shot mit aktuellem Zeitstempel.
//
// GO-KONZEPT: time.Now()
// time.Now() gibt den aktuellen Zeitpunkt zurück (UTC oder Local je nach System)
func NewShot(velocity valueobjects.Velocity) *Shot {
	return &Shot{
		Timestamp: time.Now(), // Automatisch aktueller Zeitpunkt
		Velocity:  velocity,
		Valid:     true, // Defaultmäßig gültig
	}
}

// NewShotAt erstellt einen Shot mit spezifischem Zeitstempel.
//
// Nützlich für:
// - Tests (kontrollierbare Timestamps)
// - Import von historischen Daten
// - Replay von Sessions
func NewShotAt(velocity valueobjects.Velocity, timestamp time.Time) *Shot {
	return &Shot{
		Timestamp: timestamp,
		Velocity:  velocity,
		Valid:     true,
	}
}

// CalculateEnergy berechnet die kinetische Energie dieses Schusses.
//
// GO-KONZEPT: Calculated Field (nicht persistiert!)
// Energy wird NICHT gespeichert, sondern on-the-fly berechnet.
// Warum? Siehe ADR 003: Audit Trail - wenn sich Projektilgewicht
// ändert, darf sich die historische Energie NICHT ändern.
//
// Deshalb: Session speichert ProjectileSnapshot, Shot nutzt diesen für Berechnung!
func (s *Shot) CalculateEnergy(projectileWeight valueobjects.Mass) valueobjects.Energy {
	return valueobjects.CalculateEnergy(projectileWeight, s.Velocity)
}

// MarkInvalid markiert den Shot als ungültig.
//
// GO-KONZEPT: Einzige Mutation erlaubt
// Wir erlauben EINE Änderung: Valid-Status.
// Begründung: Nutzer soll Fehlmessungen nachträglich markieren können.
func (s *Shot) MarkInvalid() {
	s.Valid = false
}

// MarkValid markiert den Shot als gültig.
func (s *Shot) MarkValid() {
	s.Valid = true
}

// String implementiert fmt.Stringer.
func (s *Shot) String() string {
	validStr := "✓"
	if !s.Valid {
		validStr = "✗"
	}

	// GO-KONZEPT: time.Time Formatting
	// Format-String: "15:04:05" = Stunden:Minuten:Sekunden
	// Go nutzt ein MERKWÜRDIGES System: Reference Time "Mon Jan 2 15:04:05 MST 2006"
	// Jedes Element dieser Zeit ist die Zahl selbst!
	timeStr := s.Timestamp.Format("15:04:05.000")

	return fmt.Sprintf("[%s] %s %s", timeStr, s.Velocity, validStr)
}

// TimeSince gibt die Zeit seit diesem Shot zurück.
//
// GO-KONZEPT: time.Duration
// Die Differenz zwischen zwei Zeitpunkten ist eine Duration
func (s *Shot) TimeSince() time.Duration {
	return time.Since(s.Timestamp)
}

// ElapsedSince gibt die verstrichene Zeit seit anderem Shot zurück.
//
// Nützlich für Analyse: Wie viel Zeit zwischen Schüssen?
func (s *Shot) ElapsedSince(other *Shot) time.Duration {
	return s.Timestamp.Sub(other.Timestamp)
}
