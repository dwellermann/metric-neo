package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// Config enthält die App-Konfiguration
type Config struct {
	DataDir string `json:"dataDir"` // Verzeichnis für Daten-Persistierung

	// Chrono (RS232) Konfiguration
	ChronoEnabled    bool   `json:"chronoEnabled,omitempty"`
	ChronoPort       string `json:"chronoPort,omitempty"`
	ChronoBaudRate   int    `json:"chronoBaudRate,omitempty"`
	ChronoAutoRecord bool   `json:"chronoAutoRecord,omitempty"`
}

// GetConfigPath gibt den Pfad zur config.json zurück
//
// Go-Konzept: Platform-spezifische Pfade
// Linux/Mac: ~/.config/metric-neo/config.json
// Windows: %APPDATA%/metric-neo/config.json
//
// Vorteil: Executable kann verschoben werden, Config bleibt erhalten
func GetConfigPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		// %APPDATA%/metric-neo/
		configDir = filepath.Join(os.Getenv("APPDATA"), "metric-neo")
	case "darwin":
		// ~/Library/Application Support/metric-neo/
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, "Library", "Application Support", "metric-neo")
	default: // linux, freebsd, etc.
		// ~/.config/metric-neo/
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configHome = filepath.Join(home, ".config")
		}
		configDir = filepath.Join(configHome, "metric-neo")
	}

	// Stelle sicher, dass Config-Verzeichnis existiert
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// ConfigExists prüft ob eine Config-Datei existiert
func ConfigExists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// LoadConfig lädt die config.json
//
// Go-Pattern: nil-Return bedeutet "Config existiert nicht"
// Caller muss Setup-Flow starten (Verzeichnis-Auswahl)
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Prüfe ob Config existiert
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil // Keine Config = Setup benötigt
	}

	// Lade existierende Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Validiere dass DataDir existiert
	if _, err := os.Stat(cfg.DataDir); os.IsNotExist(err) {
		// DataDir wurde gelöscht - Config ungültig
		return nil, nil
	}

	return &cfg, nil
}

// SaveConfig speichert config.json
//
// Erstellt automatisch das DataDir falls nicht vorhanden
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Erstelle DataDir falls nicht vorhanden
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// CreateConfig erstellt neue Config mit gewähltem Verzeichnis
//
// Go-Pattern: Constructor-Funktion
// Wird vom Setup-Flow aufgerufen nachdem User Verzeichnis gewählt hat
// Default Chrono BaudRate: 19200
func CreateConfig(dataDir string) (*Config, error) {
	cfg := &Config{
		DataDir:        dataDir,
		ChronoBaudRate: 19200, // Standard BaudRate für Chronographen
	}

	if err := SaveConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
