package application

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigService verwaltet Setup und Config-Zugriff
//
// Go-Pattern: Service kapselt Setup-Logik
// Verantwortlichkeiten:
// - Config laden/erstellen
// - Setup-Status prüfen
// - Default-Verzeichnisse vorschlagen
type ConfigService struct {
	config *Config
}

// ChronoConfigDTO kapselt Chrono-Einstellungen für UI/Bindings.
type ChronoConfigDTO struct {
	Enabled    bool   `json:"enabled"`
	Port       string `json:"port"`
	BaudRate   int    `json:"baudRate"`
	AutoRecord bool   `json:"autoRecord"`
}

// NewConfigService erstellt neuen ConfigService
//
// Go-Pattern: Constructor lädt Config falls vorhanden
// Gibt nil config zurück wenn Setup benötigt wird
func NewConfigService() (*ConfigService, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return &ConfigService{
		config: cfg,
	}, nil
}

// NeedsSetup prüft ob Initial-Setup benötigt wird
func (s *ConfigService) NeedsSetup() bool {
	return s.config == nil
}

// GetSuggestedDataDir schlägt ein Standard-Datenverzeichnis vor
//
// Go-Konzept: Platform-agnostische User-Verzeichnisse
// Vorschlag: ~/Documents/MetricNeo bzw. equivalent
// User kann im Dialog aber anderes Verzeichnis wählen
func (s *ConfigService) GetSuggestedDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Platform-spezifische Vorschläge
	return filepath.Join(home, "Documents", "MetricNeo"), nil
}

// CompleteSetup erstellt Config mit gewähltem Verzeichnis
//
// Wird von app.go aufgerufen nachdem User Verzeichnis gewählt hat
// Initialisiert die Verzeichnisstruktur gemäß ADR 003:
// - /inventory/profiles/   (Stammdaten: Waffen/Geräte)
// - /inventory/projectiles/ (Stammdaten: Munition)
// - /inventory/sights/      (Stammdaten: Zielvorrichtungen)
// - /sessions/              (Transaktionsdaten: Messprotokolle)
func (s *ConfigService) CompleteSetup(dataDir string) error {
	// Erstelle Config
	cfg, err := CreateConfig(dataDir)
	if err != nil {
		return err
	}

	// Initialisiere Verzeichnisstruktur gemäß ADR 003
	// Stammdaten unter /inventory/, Sessions direkt im Root
	inventorySubdirs := []string{
		filepath.Join("inventory", "profiles"),
		filepath.Join("inventory", "projectiles"),
		filepath.Join("inventory", "sights"),
	}
	transactionDirs := []string{"sessions"}

	// Erstelle alle Verzeichnisse
	allDirs := append(inventorySubdirs, transactionDirs...)
	for _, dir := range allDirs {
		path := filepath.Join(dataDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	s.config = cfg
	return nil
}

// GetConfig gibt die geladene Config zurück
//
// Go-Pattern: Getter mit Validierung
// Panics wenn Setup nicht abgeschlossen (Caller-Fehler)
func (s *ConfigService) GetConfig() *Config {
	if s.config == nil {
		panic("config not initialized - setup must be completed first")
	}
	return s.config
}

// GetDataDir gibt das konfigurierte Daten-Verzeichnis zurück
func (s *ConfigService) GetDataDir() string {
	return s.GetConfig().DataDir
}

// ChangeDataDir ändert das Daten-Verzeichnis
//
// Use Case: User möchte Daten in anderes Verzeichnis verschieben
// WICHTIG: Verschiebt NICHT die Daten automatisch!
// App muss vorher fragen ob Daten kopiert werden sollen
func (s *ConfigService) ChangeDataDir(newDataDir string) error {
	s.config.DataDir = newDataDir
	return SaveConfig(s.config)
}

// GetChronoConfig gibt die aktuelle Chrono-Konfiguration zurück.
func (s *ConfigService) GetChronoConfig() ChronoConfigDTO {
	if s.config == nil {
		return ChronoConfigDTO{}
	}

	return ChronoConfigDTO{
		Enabled:    s.config.ChronoEnabled,
		Port:       s.config.ChronoPort,
		BaudRate:   s.config.ChronoBaudRate,
		AutoRecord: s.config.ChronoAutoRecord,
	}
}

// UpdateChronoConfig speichert die Chrono-Konfiguration.
func (s *ConfigService) UpdateChronoConfig(cfg ChronoConfigDTO) error {
	if s.config == nil {
		return fmt.Errorf("config not initialized")
	}

	s.config.ChronoEnabled = cfg.Enabled
	s.config.ChronoPort = cfg.Port
	s.config.ChronoBaudRate = cfg.BaudRate
	s.config.ChronoAutoRecord = cfg.AutoRecord

	return SaveConfig(s.config)
}
