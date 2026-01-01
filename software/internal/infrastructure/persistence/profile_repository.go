package persistence

import (
	"encoding/json"
	"fmt"
	"metric-neo/internal/domain/entities"
	"os"
	"path/filepath"
)

type ProfileRepository struct {
	storageDir string
}

func NewProfileRepository(storageDir string) *ProfileRepository {
	return &ProfileRepository{
		storageDir: storageDir,
	}
}

// Save speichert ein Profile als JSON-Datei
func (r *ProfileRepository) Save(profile *entities.Profile) error {
	if profile == nil {
		return fmt.Errorf("profile cannot be nil")
	}

	// Erstelle Verzeichnis falls nicht vorhanden
	err := os.MkdirAll(r.storageDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Serialisiere zu JSON (mit Einrückung für Lesbarkeit)
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	// Dateiname: <uuid>.json
	filename := filepath.Join(r.storageDir, profile.ID+".json")

	// Schreibe Datei (0644 = rw-r--r--)
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load lädt ein Profile anhand seiner ID
func (r *ProfileRepository) Load(id string) (*entities.Profile, error) {
	filename := filepath.Join(r.storageDir, id+".json")

	// Lese Datei
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile not found: %w", err)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Deserialisiere JSON
	var profile entities.Profile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &profile, nil
}

// List gibt alle gespeicherten Profile-IDs zurück
func (r *ProfileRepository) List() ([]string, error) {
	entries, err := os.ReadDir(r.storageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // Leeres Array statt Fehler
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			// Entferne .json Extension um ID zu erhalten
			id := entry.Name()[:len(entry.Name())-5]
			ids = append(ids, id)
		}
	}

	return ids, nil
}

// Delete löscht ein Profile anhand seiner ID
func (r *ProfileRepository) Delete(id string) error {
	filename := filepath.Join(r.storageDir, id+".json")

	err := os.Remove(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("profile not found: %w", err)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
