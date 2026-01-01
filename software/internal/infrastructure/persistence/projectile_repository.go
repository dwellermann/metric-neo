package persistence

import (
	"encoding/json"
	"fmt"
	"metric-neo/internal/domain/entities"
	"os"
	"path/filepath"
)

// ProjectileRepository verwaltet die Persistierung von Projectile-Entities.
//
// GO-KONZEPT: Repository Pattern (Domain-Driven Design)
// Das Repository abstrahiert die Datenspeicherung. Die Domain-Schicht
// (entities/) weiß NICHTS über JSON oder Files - das ist gut!
//
// Vorteile:
// - Wir könnten später zu SQL wechseln, ohne die Entity zu ändern
// - Tests können ein Mock-Repository nutzen
// - Klare Trennung: Domain-Logik vs. Technik
type ProjectileRepository struct {
	// dataDir ist das Verzeichnis, wo JSON-Dateien gespeichert werden
	// Laut ADR 003: /data/inventory/
	dataDir string
}

// NewProjectileRepository erstellt eine neue Repository-Instanz.
//
// GO-KONZEPT: Dependency Injection
// Wir übergeben den dataDir-Pfad, statt ihn hart zu kodieren.
// Das macht Tests einfacher (wir können ein temp-Verzeichnis nutzen).
func NewProjectileRepository(dataDir string) *ProjectileRepository {
	return &ProjectileRepository{
		dataDir: dataDir,
	}
}

// Save speichert ein Projectile als JSON-Datei.
//
// GO-KONZEPT: Error Wrapping mit fmt.Errorf()
// Wenn ein Fehler auftritt, "wrappen" wir ihn mit Kontext.
// Das hilft beim Debugging: "failed to save projectile: permission denied"
// statt nur: "permission denied"
func (r *ProjectileRepository) Save(p *entities.Projectile) error {
	// GO-KONZEPT: os.MkdirAll() - Rekursive Directory-Erstellung
	// Entspricht `mkdir -p` in der Shell
	// 0755 = Berechtigungen (rwxr-xr-x)
	if err := os.MkdirAll(r.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Dateiname: {UUID}.json
	// filepath.Join() ist platform-agnostic (Windows: \, Linux: /)
	filename := filepath.Join(r.dataDir, p.ID+".json")

	// GO-KONZEPT: json.MarshalIndent() für lesbare JSON-Dateien
	// Parameter:
	// - p: Das zu serialisierende Objekt
	// - "": Kein Prefix
	// - "  ": 2 Spaces Einrückung (macht JSON schön lesbar)
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projectile: %w", err)
	}

	// GO-KONZEPT: os.WriteFile() - Atomic File Write
	// 0644 = Berechtigungen (rw-r--r--)
	// Entspricht: Owner kann lesen+schreiben, Rest nur lesen
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	return nil
}

// Load lädt ein Projectile anhand seiner ID.
//
// GO-KONZEPT: Pointer Return für Entities
// Wir geben *entities.Projectile zurück (nicht eine Kopie).
func (r *ProjectileRepository) Load(id string) (*entities.Projectile, error) {
	filename := filepath.Join(r.dataDir, id+".json")

	// GO-KONZEPT: os.ReadFile() - File lesen
	// Gibt []byte (Byte-Array) zurück
	data, err := os.ReadFile(filename)
	if err != nil {
		// GO-KONZEPT: os.IsNotExist() für spezifische Fehlerprüfung
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("projectile %s not found", id)
		}
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Unmarshal JSON zu Projectile
	var projectile entities.Projectile
	if err := json.Unmarshal(data, &projectile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projectile: %w", err)
	}

	// GO-KONZEPT: Address-Of Operator (&)
	// &projectile macht aus dem Value einen Pointer
	return &projectile, nil
}

// List gibt alle gespeicherten Projectiles zurück.
//
// GO-KONZEPT: Slice Return ([]*entities.Projectile)
// Ein Slice ist wie ein Array, aber dynamisch wachsend.
// []*entities.Projectile = "Slice von Pointern auf Projectile"
func (r *ProjectileRepository) List() ([]*entities.Projectile, error) {
	// GO-KONZEPT: os.ReadDir() - Verzeichnis auflisten
	// Gibt []os.DirEntry zurück
	entries, err := os.ReadDir(r.dataDir)
	if err != nil {
		// Wenn Verzeichnis nicht existiert, geben wir leere Liste zurück
		if os.IsNotExist(err) {
			return []*entities.Projectile{}, nil
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// GO-KONZEPT: Slice mit Capacity pre-allocieren
	// make([]*entities.Projectile, 0, len(entries))
	// - Länge 0 (leer)
	// - Capacity = Anzahl Dateien (vermeidet Reallocations)
	projectiles := make([]*entities.Projectile, 0, len(entries))

	// GO-KONZEPT: Range Loop über Slice
	for _, entry := range entries {
		// Überspringe Verzeichnisse, nur .json Dateien laden
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		// Extrahiere ID aus Dateiname (entferne .json Extension)
		id := entry.Name()[:len(entry.Name())-5] // "abc.json" -> "abc"

		// Lade Projectile
		p, err := r.Load(id)
		if err != nil {
			// Bei Fehler: Überspringe diese Datei, aber fahre fort
			// In Production würdest du das loggen!
			continue
		}

		// GO-KONZEPT: append() - Element zu Slice hinzufügen
		// append() gibt einen NEUEN Slice zurück (Slices sind semi-immutable)
		projectiles = append(projectiles, p)
	}

	return projectiles, nil
}

// Delete löscht ein Projectile.
func (r *ProjectileRepository) Delete(id string) error {
	filename := filepath.Join(r.dataDir, id+".json")

	// GO-KONZEPT: os.Remove() - Datei löschen
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("projectile %s not found", id)
		}
		return fmt.Errorf("failed to delete file %s: %w", filename, err)
	}

	return nil
}
