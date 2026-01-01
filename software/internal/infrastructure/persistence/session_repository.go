package persistence

import (
	"encoding/json"
	"fmt"
	"metric-neo/internal/domain/entities"
	"os"
	"path/filepath"
)

type SessionRepository struct {
	storageDir string
}

func NewSessionRepository(storageDir string) *SessionRepository {
	return &SessionRepository{
		storageDir: storageDir,
	}
}

// Save speichert eine Session als JSON-Datei
func (r *SessionRepository) Save(session *entities.Session) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	// Erstelle Verzeichnis falls nicht vorhanden
	err := os.MkdirAll(r.storageDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Serialisiere zu JSON (mit Einrückung für Lesbarkeit)
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Dateiname: <uuid>.json
	filename := filepath.Join(r.storageDir, session.ID+".json")

	// Schreibe Datei (0644 = rw-r--r--)
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load lädt eine Session anhand ihrer ID
func (r *SessionRepository) Load(id string) (*entities.Session, error) {
	filename := filepath.Join(r.storageDir, id+".json")

	// Lese Datei
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session not found: %w", err)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Deserialisiere JSON
	var session entities.Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// List gibt alle gespeicherten Session-IDs zurück
func (r *SessionRepository) List() ([]string, error) {
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

// Delete löscht eine Session anhand ihrer ID
func (r *SessionRepository) Delete(id string) error {
	filename := filepath.Join(r.storageDir, id+".json")

	err := os.Remove(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("session not found: %w", err)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
