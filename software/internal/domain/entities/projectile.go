package entities

import (
	"fmt"
	"metric-neo/internal/domain/valueobjects"

	"github.com/google/uuid"
)

// Projectile repräsentiert ein Munitions-Profil (Stammdaten).
//
// GO-KONZEPT: Struct (vergleichbar mit TypeScript Interface/Class)
// - Structs sind Value Types (werden kopiert bei Zuweisung)
// - Felder mit Großbuchstaben sind EXPORTIERT (public)
// - Felder mit Kleinbuchstaben sind PRIVAT (nur im gleichen Package sichtbar)
//
// DOMAIN MODEL: Siehe docs/specs/domain-model.md Abschnitt 3.4
type Projectile struct {
	// GO-KONZEPT: Struct Tags
	// `json:"id"` sagt dem JSON-Encoder: "Nutze 'id' als Feldname, nicht 'ID'"
	// Ohne Tag würde es als "ID" (Großbuchstaben) im JSON erscheinen.
	ID string `json:"id"`

	// Name des Munitionstyps (z.B. "JSB Exact 4.52")
	Name string `json:"name"`

	// GO-KONZEPT: Embedding von Custom Types
	// Weight ist KEIN float64, sondern unser Value Object Mass!
	// Das erzwingt Type Safety: Man kann nicht versehentlich eine
	// Velocity als Weight setzen.
	Weight valueobjects.Mass `json:"weight"`

	// BC (Ballistic Coefficient) - Optional für zukünftige Features
	// Primitive Typen (float64, int, string) können direkt genutzt werden
	BC float64 `json:"bc"`
}

// NewProjectile ist der Constructor für Projectile.
//
// GO-KONZEPT: Pointer Return (*Projectile)
// Warum geben wir einen POINTER zurück, nicht den Wert?
//
//  1. PERFORMANCE: Projectile könnte später groß werden (viele Felder).
//     Einen Pointer zu kopieren (8 Bytes) ist günstiger als das ganze Struct.
//
//  2. IDENTITÄT: Entities haben eine Identität (UUID). Wir wollen, dass
//     verschiedene Variablen auf DAS GLEICHE Objekt zeigen können.
//     Bei Value-Kopien hätten wir mehrere "Klone".
//
//  3. MUTATION: Wenn wir später Methoden haben, die Projectile ändern
//     (z.B. UpdateBC()), müssen wir einen Pointer nutzen.
//
// GO-BEST-PRACTICE: Entities -> Pointer, Value Objects -> Value
func NewProjectile(name string, weight valueobjects.Mass, bc float64) (*Projectile, error) {
	// Validierung: Name darf nicht leer sein
	if name == "" {
		return nil, fmt.Errorf("projectile name cannot be empty")
	}

	// Validierung: BC muss >= 0 sein (physikalisch sinnvoll)
	if bc < 0 {
		return nil, fmt.Errorf("ballistic coefficient cannot be negative, got: %.3f", bc)
	}

	// GO-KONZEPT: UUID Generierung
	// uuid.New() erstellt eine zufällige UUIDv4
	// .String() konvertiert sie zu "550e8400-e29b-41d4-a716-446655440000"
	id := uuid.New().String()

	// GO-KONZEPT: Struct Literal
	// Wir erstellen eine Instanz und geben ihre ADRESSE zurück (&)
	// Das & (Address-Of Operator) macht aus Projectile einen *Projectile
	return &Projectile{
		ID:     id,
		Name:   name,
		Weight: weight,
		BC:     bc,
	}, nil
}

// String implementiert fmt.Stringer für schöne Ausgabe.
//
// GO-KONZEPT: Method mit Pointer Receiver
// Die Syntax (p *Projectile) bedeutet: Diese Methode gehört zu *Projectile.
//
// WARUM POINTER-RECEIVER?
//   - Consistency: Wenn EINE Methode einen Pointer braucht (z.B. Update()),
//     sollten ALLE Methoden Pointer nutzen (Go Best Practice).
//   - Performance: Bei großen Structs vermeiden wir Kopien.
//
// REGEL: Entities -> Pointer Receiver, Value Objects -> Value Receiver
func (p *Projectile) String() string {
	return fmt.Sprintf("%s (%.3f g, BC: %.3f)", p.Name, p.Weight.Grams(), p.BC)
}

// UpdateBC aktualisiert den ballistischen Koeffizienten.
//
// GO-KONZEPT: Mutation mit Pointer Receiver
// OHNE Pointer (*Projectile) würden wir nur eine KOPIE ändern!
// Der Original-Wert bliebe unverändert.
//
// WICHTIG: Diese Methode ist nur für MASTER-Daten gedacht.
// Snapshots in Sessions sind IMMUTABLE und haben diese Methode nicht!
func (p *Projectile) UpdateBC(newBC float64) error {
	if newBC < 0 {
		return fmt.Errorf("ballistic coefficient cannot be negative")
	}
	p.BC = newBC
	return nil
}
