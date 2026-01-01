package persistence

import (
	"encoding/json"
	"fmt"
	"metric-neo/internal/domain/entities"
	"os"
	"path/filepath"
)

// SightRepository verwaltet die Persistierung von SightingSystem-Entities.
//
// GO-KONZEPT: Repository Pattern (Domain-Driven Design)
// Analog zu ProfileRepository und ProjectileRepository
type SightRepository struct {
	storageDir string
}

// NewSightRepository erstellt eine neue Repository-Instanz.
func NewSightRepository(storageDir string) *SightRepository {
	return &SightRepository{
		storageDir: storageDir,
	}
}

// Save speichert ein SightingSystem als JSON-Datei.
func (r *SightRepository) Save(sight *entities.SightingSystem) error {
	if sight == nil {
		return fmt.Errorf("sighting system cannot be nil")
	}

	// Erstelle Verzeichnis falls nicht vorhanden
	err := os.MkdirAll(r.storageDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Serialisiere zu JSON (mit Einrückung für Lesbarkeit)
	data, err := json.MarshalIndent(sight, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sighting system: %w", err)
	}

	// Dateiname: <uuid>.json
	filename := filepath.Join(r.storageDir, sight.ID+".json")

	// Schreibe Datei (0644 = rw-r--r--)
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load lädt ein SightingSystem anhand seiner ID.
func (r *SightRepository) Load(id string) (*entities.SightingSystem, error) {
	if id == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	filename := filepath.Join(r.storageDir, id+".json")

	// Prüfe ob Datei existiert
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("sighting system not found: %s", id)
	}

	// Lese Datei
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Deserialisiere JSON
	var sight entities.SightingSystem
	if err := json.Unmarshal(data, &sight); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sighting system: %w", err)
	}

	return &sight, nil
}

// List gibt alle gespeicherten SightingSystems zurück.
func (r *SightRepository) List() ([]*entities.SightingSystem, error) {
	// Prüfe ob Verzeichnis existiert
	if _, err := os.Stat(r.storageDir); os.IsNotExist(err) {
		// Noch keine Sights gespeichert - gebe leeres Array zurück
		return []*entities.SightingSystem{}, nil
	}

	// Lese alle .json Dateien
	entries, err := os.ReadDir(r.storageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var sights []*entities.SightingSystem
	for _, entry := range entries {
		// Überspringe Verzeichnisse und Nicht-JSON-Dateien
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		// Extrahiere ID aus Dateiname
		id := entry.Name()[:len(entry.Name())-5] // Entferne ".json"

		// Lade Sight
		sight, err := r.Load(id)
		if err != nil {
			// Log error aber überspringe fehlerhafte Dateien
			continue
		}

		sights = append(sights, sight)
	}

	return sights, nil
}

// Delete löscht ein SightingSystem.
func (r *SightRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	filename := filepath.Join(r.storageDir, id+".json")

	// Prüfe ob Datei existiert
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("sighting system not found: %s", id)
	}

	// Lösche Datei
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
