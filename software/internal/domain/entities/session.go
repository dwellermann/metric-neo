package entities

import (
	"fmt"
	"math"
	"metric-neo/internal/domain/valueobjects"
	"time"

	"github.com/google/uuid"
)

// Session ist das Aggregate Root für eine Messreihe.
//
// GO-KONZEPT: Aggregate Root (Domain-Driven Design)
// Session ist der "Container" für alle Schüsse einer Sitzung.
// Sie garantiert Konsistenz und verwaltet Snapshots.
//
// DOMAIN MODEL: Siehe docs/specs/domain-model.md Abschnitt 3.1
// ADR 003: Snapshot-Pattern für Audit Trail
type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"` // Optional: Notizen

	// Optional: Umgebungsbedingungen
	Temperature *valueobjects.Temperature `json:"temperature,omitempty"`

	// GO-KONZEPT: Deep Copy Snapshots (KERN des Patterns!)
	// Diese Felder sind KOPIEN zum Zeitpunkt der Session-Erstellung!
	// Spätere Änderungen an Master-Daten beeinflussen diese Session NICHT.
	//
	// WICHTIG: Wir speichern die GANZEN Objekte, nicht nur IDs!
	ProfileSnapshot    *Profile    `json:"profile_snapshot"`
	ProjectileSnapshot *Projectile `json:"projectile_snapshot"`

	// GO-KONZEPT: Slice von Pointers
	// []*Shot = "Slice (Array) von Pointern zu Shot"
	// Warum Pointer? Shots sind Entities (Identität)
	Shots []*Shot `json:"shots"`
}

// NewSession erstellt eine neue Session mit Snapshots.
//
// GO-KONZEPT: Deep Copy Pattern
// profile und projectile werden KOPIERT, nicht referenziert!
func NewSession(profile *Profile, projectile *Projectile) *Session {
	return &Session{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),

		// GO-KONZEPT: Deep Copy durch Dereferenzierung
		// *profile = "Wert an der Adresse profile"
		// Das erstellt eine KOPIE des Profile-Structs
		//
		// WICHTIG: Bei nested Pointern (z.B. profile.Optic) müssen wir
		// aufpassen - siehe CopyProfile() Funktion unten!
		ProfileSnapshot:    CopyProfile(profile),
		ProjectileSnapshot: CopyProjectile(projectile),

		// GO-KONZEPT: Leerer Slice
		// make([]*Shot, 0) erstellt einen Slice mit Länge 0
		Shots: make([]*Shot, 0),
	}
}

// CopyProfile erstellt eine Deep Copy eines Profile.
//
// GO-KONZEPT: Deep Copy bei nested Pointers
// Problem: profile.Optic ist ein *SightingSystem (Pointer!)
// Wenn wir nur *profile kopieren, zeigt profile.Optic IMMER NOCH
// auf das Original -> keine echte Kopie!
//
// Lösung: Wir müssen ALLE Pointer-Felder manuell kopieren.
func CopyProfile(p *Profile) *Profile {
	if p == nil {
		return nil
	}

	// Shallow Copy (kopiert alle Value-Felder)
	copy := *p

	// Deep Copy für Pointer-Felder
	if p.Optic != nil {
		opticCopy := *p.Optic // Dereferenzieren = Kopieren
		copy.Optic = &opticCopy
	}

	if p.TwistRate != nil {
		twistCopy := *p.TwistRate
		copy.TwistRate = &twistCopy
	}

	if p.DefaultAmmoID != nil {
		idCopy := *p.DefaultAmmoID
		copy.DefaultAmmoID = &idCopy
	}

	return &copy
}

// CopyProjectile erstellt eine Deep Copy eines Projectile.
//
// Einfacher als Profile, da Projectile keine Pointer-Felder hat!
func CopyProjectile(p *Projectile) *Projectile {
	if p == nil {
		return nil
	}

	copy := *p
	return &copy
}

// AddShot fügt einen Shot zur Session hinzu.
//
// GO-KONZEPT: Slice Append
// append() fügt ein Element zu einem Slice hinzu.
// WICHTIG: append() gibt einen NEUEN Slice zurück!
// Deshalb: s.Shots = append(s.Shots, shot)
func (s *Session) AddShot(shot *Shot) {
	s.Shots = append(s.Shots, shot)
}

// RecordShot erstellt einen neuen Shot und fügt ihn hinzu.
//
// Helper-Methode für schnelles Hinzufügen
func (s *Session) RecordShot(velocity valueobjects.Velocity) *Shot {
	shot := NewShot(velocity)
	s.AddShot(shot)
	return shot
}

// ShotCount gibt die Anzahl Schüsse zurück.
func (s *Session) ShotCount() int {
	return len(s.Shots)
}

// ValidShotCount gibt die Anzahl gültiger Schüsse zurück.
func (s *Session) ValidShotCount() int {
	count := 0
	for _, shot := range s.Shots {
		if shot.Valid {
			count++
		}
	}
	return count
}

// CalculateAverageVelocity berechnet die Durchschnittsgeschwindigkeit.
//
// GO-KONZEPT: Aggregat-Berechnung über Slice
// Nur gültige Schüsse werden berücksichtigt!
func (s *Session) CalculateAverageVelocity() (valueobjects.Velocity, error) {
	validShots := s.getValidShots()

	if len(validShots) == 0 {
		return 0, fmt.Errorf("no valid shots in session")
	}

	sum := 0.0
	for _, shot := range validShots {
		sum += shot.Velocity.MetersPerSecond()
	}

	avg := sum / float64(len(validShots))
	return valueobjects.NewVelocity(avg)
}

// CalculateStandardDeviation berechnet die Standardabweichung der Geschwindigkeit.
//
// GO-KONZEPT: Mathematische Berechnungen
// SD = √(Σ(x - μ)² / n)
func (s *Session) CalculateStandardDeviation() (float64, error) {
	validShots := s.getValidShots()

	if len(validShots) < 2 {
		return 0, fmt.Errorf("need at least 2 valid shots for standard deviation")
	}

	// Durchschnitt berechnen
	avg, err := s.CalculateAverageVelocity()
	if err != nil {
		return 0, err
	}
	avgValue := avg.MetersPerSecond()

	// Varianz berechnen
	variance := 0.0
	for _, shot := range validShots {
		diff := shot.Velocity.MetersPerSecond() - avgValue
		variance += diff * diff
	}
	variance /= float64(len(validShots))

	// Standardabweichung = Wurzel der Varianz
	return math.Sqrt(variance), nil
}

// CalculateAverageEnergy berechnet die durchschnittliche Energie.
//
// WICHTIG: Nutzt ProjectileSnapshot für Gewicht!
// Das garantiert historische Korrektheit (ADR 003).
func (s *Session) CalculateAverageEnergy() (valueobjects.Energy, error) {
	validShots := s.getValidShots()

	if len(validShots) == 0 {
		return 0, fmt.Errorf("no valid shots in session")
	}

	sum := 0.0
	for _, shot := range validShots {
		energy := shot.CalculateEnergy(s.ProjectileSnapshot.Weight)
		sum += energy.Joules()
	}

	avg := sum / float64(len(validShots))
	return valueobjects.Energy(avg), nil
}

// MinVelocity gibt die niedrigste gemessene Geschwindigkeit zurück.
func (s *Session) MinVelocity() (valueobjects.Velocity, error) {
	validShots := s.getValidShots()

	if len(validShots) == 0 {
		return 0, fmt.Errorf("no valid shots in session")
	}

	min := validShots[0].Velocity.MetersPerSecond()
	for _, shot := range validShots[1:] {
		if shot.Velocity.MetersPerSecond() < min {
			min = shot.Velocity.MetersPerSecond()
		}
	}

	return valueobjects.NewVelocity(min)
}

// MaxVelocity gibt die höchste gemessene Geschwindigkeit zurück.
func (s *Session) MaxVelocity() (valueobjects.Velocity, error) {
	validShots := s.getValidShots()

	if len(validShots) == 0 {
		return 0, fmt.Errorf("no valid shots in session")
	}

	max := validShots[0].Velocity.MetersPerSecond()
	for _, shot := range validShots[1:] {
		if shot.Velocity.MetersPerSecond() > max {
			max = shot.Velocity.MetersPerSecond()
		}
	}

	return valueobjects.NewVelocity(max)
}

// ExtremeSpread gibt die Spreizung (Max - Min) zurück.
//
// Wichtige Metrik in der Ballistik: Je kleiner, desto konsistenter!
func (s *Session) ExtremeSpread() (float64, error) {
	min, err := s.MinVelocity()
	if err != nil {
		return 0, err
	}

	max, err := s.MaxVelocity()
	if err != nil {
		return 0, err
	}

	return max.MetersPerSecond() - min.MetersPerSecond(), nil
}

// getValidShots ist ein interner Helper für gültige Schüsse.
//
// GO-KONZEPT: Slice Filtering
func (s *Session) getValidShots() []*Shot {
	valid := make([]*Shot, 0)
	for _, shot := range s.Shots {
		if shot.Valid {
			valid = append(valid, shot)
		}
	}
	return valid
}

// SetTemperature setzt die Umgebungstemperatur.
func (s *Session) SetTemperature(temp valueobjects.Temperature) {
	s.Temperature = &temp
}

// SetNote setzt eine Notiz.
func (s *Session) SetNote(note string) {
	s.Note = note
}

// String implementiert fmt.Stringer.
func (s *Session) String() string {
	timeStr := s.CreatedAt.Format("2006-01-02 15:04")
	return fmt.Sprintf("Session %s - %s (%d shots)",
		timeStr,
		s.ProfileSnapshot.Name,
		s.ShotCount(),
	)
}
