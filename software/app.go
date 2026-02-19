package main

import (
	"context"
	"fmt"
	"metric-neo/internal/application"
	"metric-neo/internal/infrastructure/chrono"
	"os"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	configService *application.ConfigService
	// Services werden erst initialisiert nach Setup
	profileService    *application.ProfileService
	projectileService *application.ProjectileService
	sessionService    *application.SessionService
	sightService      *application.SightService

	chronoService      chrono.ChronoService
	chronoVelocityChan chan float32
	chronoErrorChan    chan error
	chronoAutoRunning  bool
	chronoMu           sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
//
// Go-Pattern: Lazy Initialization
// Config wird geladen, Services erst nach Setup initialisiert
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialisiere ConfigService
	configService, err := application.NewConfigService()
	if err != nil {
		runtime.LogFatal(ctx, "Failed to initialize config service: "+err.Error())
		return
	}
	a.configService = configService

	// Prüfe ob Setup benötigt wird
	if configService.NeedsSetup() {
		runtime.LogInfo(ctx, "First start detected - setup required")
		// Frontend zeigt Setup-Dialog
		// User wählt Verzeichnis über SelectDirectory()
		return
	}

	// Setup bereits abgeschlossen - initialisiere Services
	if err := a.initializeServices(); err != nil {
		runtime.LogFatal(ctx, "Failed to initialize services: "+err.Error())
	}
}

// initializeServices erstellt alle Services mit Config
//
// Go-Pattern: Factory Method
// Wird nach Setup oder beim Start (wenn Setup bereits erfolgt) aufgerufen
func (a *App) initializeServices() error {
	cfg := a.configService.GetConfig()
	dataDir := cfg.DataDir

	runtime.LogInfo(a.ctx, "Initializing services with data directory: "+dataDir)

	// Erstelle Services
	a.profileService = application.NewProfileService(dataDir)
	a.projectileService = application.NewProjectileService(dataDir)
	a.sessionService = application.NewSessionService(dataDir)
	a.sightService = application.NewSightService(dataDir)

	if a.chronoService == nil {
		a.chronoService = chrono.NewSerialChrono()
	}

	return nil
}

// --- Setup-Flow für Frontend ---

// NeedsSetup prüft ob Initial-Setup benötigt wird
//
// Frontend kann dies beim Start aufrufen um Setup-Dialog zu zeigen
func (a *App) NeedsSetup() bool {
	return a.configService.NeedsSetup()
}

// GetSuggestedDataDir gibt Vorschlag für Daten-Verzeichnis zurück
//
// Frontend kann dies als Default im Verzeichnis-Auswahl-Dialog nutzen
func (a *App) GetSuggestedDataDir() string {
	dir, err := a.configService.GetSuggestedDataDir()
	if err != nil {
		return ""
	}
	return dir
}

// SelectDataDirectory öffnet Verzeichnis-Auswahl-Dialog
//
// Wails Runtime: runtime.OpenDirectoryDialog
// Gibt nur das gewählte Verzeichnis zurück, initialisiert NICHT
// Frontend zeigt Verzeichnis an, User bestätigt mit "Start"
func (a *App) SelectDataDirectory() application.Result[string] {
	// Starte im Home-Verzeichnis (existiert immer auf Linux/Mac/Windows)
	// User kann dort neuen Ordner erstellen
	homeDir, err := os.UserHomeDir()
	if err != nil {
		runtime.LogError(a.ctx, "Failed to get home directory: "+err.Error())
		return application.Fail[string](err)
	}

	runtime.LogDebug(a.ctx, "Opening directory dialog in: "+homeDir)

	// Öffne Directory-Dialog
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Wähle Verzeichnis für Metric Neo Daten",
		DefaultDirectory: homeDir,
	})

	if err != nil {
		runtime.LogError(a.ctx, "OpenDirectoryDialog failed: "+err.Error())
		return application.Fail[string](err)
	}

	// User hat abgebrochen
	if selectedDir == "" {
		runtime.LogInfo(a.ctx, "User cancelled directory selection")
		return application.FailWithMessage[string]("Kein Verzeichnis gewählt")
	}

	runtime.LogInfo(a.ctx, "User selected directory: "+selectedDir)
	return application.OK(selectedDir)
}

// CompleteSetup führt Setup mit gewähltem Verzeichnis durch
//
// Erstellt Config, initialisiert Repository-Verzeichnisse, startet Services
// Wird aufgerufen wenn User auf "Start" klickt nach Verzeichnis-Auswahl
func (a *App) CompleteSetup(dataDir string) application.Result[bool] {
	runtime.LogInfo(a.ctx, "Completing setup with directory: "+dataDir)

	// Führe Setup durch
	if err := a.configService.CompleteSetup(dataDir); err != nil {
		runtime.LogError(a.ctx, "Setup failed: "+err.Error())
		return application.Fail[bool](err)
	}

	// Initialisiere Services
	if err := a.initializeServices(); err != nil {
		runtime.LogError(a.ctx, "Service initialization failed: "+err.Error())
		return application.Fail[bool](err)
	}

	runtime.LogInfo(a.ctx, "Setup completed successfully")
	return application.OK(true)
}

// GetCurrentDataDir gibt aktuelles Daten-Verzeichnis zurück
func (a *App) GetCurrentDataDir() string {
	if a.configService.NeedsSetup() {
		return ""
	}
	return a.configService.GetDataDir()
}

// GetChronoConfig gibt die aktuelle Chrono-Konfiguration zurück
func (a *App) GetChronoConfig() application.Result[application.ChronoConfigDTO] {
	if a.configService == nil || a.configService.NeedsSetup() {
		return application.FailWithMessage[application.ChronoConfigDTO]("Config not initialized - setup not completed")
	}
	return application.OK(a.configService.GetChronoConfig())
}

// UpdateChronoConfig speichert Chrono-Konfiguration
func (a *App) UpdateChronoConfig(enabled bool, port string, baudRate int, autoRecord bool) application.Result[application.ChronoConfigDTO] {
	if a.configService == nil || a.configService.NeedsSetup() {
		return application.FailWithMessage[application.ChronoConfigDTO]("Config not initialized - setup not completed")
	}

	cfg := application.ChronoConfigDTO{
		Enabled:    enabled,
		Port:       port,
		BaudRate:   baudRate,
		AutoRecord: autoRecord,
	}

	if err := a.configService.UpdateChronoConfig(cfg); err != nil {
		return application.Fail[application.ChronoConfigDTO](err)
	}

	if !enabled || port == "" || baudRate <= 0 {
		a.stopChrono()
	}

	return application.OK(a.configService.GetChronoConfig())
}

// SessionPollChrono prüft, ob neue Chrono-Messwerte verfügbar sind und zeichnet sie automatisch auf.
// Option A: Frontend pollt diese Methode regelmäßig.
func (a *App) SessionPollChrono(sessionID string) application.Result[application.ChronoPollResultDTO] {
	if a.configService == nil || a.configService.NeedsSetup() {
		return application.FailWithMessage[application.ChronoPollResultDTO]("Config not initialized - setup not completed")
	}
	if a.sessionService == nil {
		return application.FailWithMessage[application.ChronoPollResultDTO]("Services not initialized - setup not completed")
	}
	if sessionID == "" {
		return application.FailWithMessage[application.ChronoPollResultDTO]("Session-ID darf nicht leer sein")
	}

	cfg := a.configService.GetConfig()
	if !cfg.ChronoEnabled || !cfg.ChronoAutoRecord {
		return application.OK(application.ChronoPollResultDTO{Recorded: false})
	}
	if cfg.ChronoPort == "" || cfg.ChronoBaudRate <= 0 {
		return application.FailWithMessage[application.ChronoPollResultDTO]("Chrono ist nicht korrekt konfiguriert")
	}

	if err := a.ensureChronoRunning(cfg.ChronoPort, cfg.ChronoBaudRate); err != nil {
		return application.Fail[application.ChronoPollResultDTO](err)
	}

	select {
	case err := <-a.chronoErrorChan:
		return application.Fail[application.ChronoPollResultDTO](err)
	case velocity := <-a.chronoVelocityChan:
		recordResult := a.sessionService.RecordShot(sessionID, float64(velocity))
		if !recordResult.Success {
			return application.FailWithMessage[application.ChronoPollResultDTO](recordResult.Error)
		}
		v := float64(velocity)
		dto := application.ChronoPollResultDTO{
			Recorded:    true,
			VelocityMPS: &v,
			Session:     &recordResult.Data,
		}
		return application.OK(dto)
	default:
		return application.OK(application.ChronoPollResultDTO{Recorded: false})
	}
}

func (a *App) ensureChronoRunning(port string, baudRate int) error {
	a.chronoMu.Lock()
	defer a.chronoMu.Unlock()

	if a.chronoService == nil {
		a.chronoService = chrono.NewSerialChrono()
	}

	if !a.chronoService.IsConnected() {
		if err := a.chronoService.Connect(port, baudRate); err != nil {
			return err
		}
	}

	if a.chronoVelocityChan == nil {
		a.chronoVelocityChan = make(chan float32, 16)
	}
	if a.chronoErrorChan == nil {
		a.chronoErrorChan = make(chan error, 8)
	}

	if !a.chronoAutoRunning {
		a.chronoService.StartAutoRead(a.chronoVelocityChan, a.chronoErrorChan)
		a.chronoAutoRunning = true
	}

	return nil
}

func (a *App) stopChrono() {
	a.chronoMu.Lock()
	defer a.chronoMu.Unlock()

	if a.chronoService != nil && a.chronoAutoRunning {
		a.chronoService.StopAutoRead()
		a.chronoAutoRunning = false
	}

	if a.chronoService != nil && a.chronoService.IsConnected() {
		_ = a.chronoService.Disconnect()
	}
}

// GetSystemTheme gibt das System-Theme zurück (dark/light)
//
// Detektiert automatisch basierend auf OS:
// - Linux: gsettings (GNOME)
// - Windows: Registry
// - macOS: defaults
func (a *App) GetSystemTheme() string {
	return application.GetSystemTheme()
}

// ChangeDataDirectory ändert das Daten-Verzeichnis
//
// WARNUNG: Verschiebt NICHT automatisch die Daten!
// Frontend muss vorher fragen ob Daten kopiert werden sollen
func (a *App) ChangeDataDirectory() application.Result[string] {
	currentDir := a.configService.GetDataDir()

	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Neues Daten-Verzeichnis wählen",
		DefaultDirectory: currentDir,
	})

	if err != nil {
		return application.Fail[string](err)
	}

	if selectedDir == "" {
		return application.FailWithMessage[string]("Kein Verzeichnis gewählt")
	}

	if err := a.configService.ChangeDataDir(selectedDir); err != nil {
		return application.Fail[string](err)
	}

	// Services neu initialisieren mit neuem Verzeichnis
	if err := a.initializeServices(); err != nil {
		return application.Fail[string](err)
	}

	return application.OK(selectedDir)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ==================== PROFILE SERVICE DELEGATION ====================

// ProfileCreateProfile erstellt ein neues Profil
func (a *App) ProfileCreateProfile(name string, category string, barrelLengthMM float64, triggerWeightG float64, sightHeightMM float64) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.CreateProfile(name, category, barrelLengthMM, triggerWeightG, sightHeightMM)
}

// ProfileLoadProfile lädt ein Profil nach ID
func (a *App) ProfileLoadProfile(id string) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.LoadProfile(id)
}

// ProfileListProfiles listet alle Profile auf
func (a *App) ProfileListProfiles() application.Result[[]application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[[]application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.ListProfiles()
}

// ProfileUpdateProfile aktualisiert die Basisdaten eines Profils
func (a *App) ProfileUpdateProfile(
	profileID string,
	name string,
	category string,
	barrelLengthMM float64,
	triggerWeightG float64,
	sightHeightMM float64,
) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.UpdateProfile(profileID, name, category, barrelLengthMM, triggerWeightG, sightHeightMM)
}

// ProfileDeleteProfile löscht ein Profil
func (a *App) ProfileDeleteProfile(id string) application.Result[bool] {
	if a.profileService == nil {
		return application.FailWithMessage[bool]("Services not initialized - setup not completed")
	}
	return a.profileService.DeleteProfile(id)
}

// ProfileSetOptic setzt eine Optik für ein Profil
func (a *App) ProfileSetOptic(profileID string, opticType string, modelName string, weightG float64, minMagnification float64, maxMagnification float64) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.SetOptic(profileID, opticType, modelName, weightG, minMagnification, maxMagnification)
}

// ProfileRemoveOptic entfernt die Optik von einem Profil
func (a *App) ProfileRemoveOptic(profileID string) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.RemoveOptic(profileID)
}

// ProfileLinkOpticByID verknüpft eine bestehende Optik mit einem Profil
func (a *App) ProfileLinkOpticByID(profileID string, sightID string) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.LinkOpticByID(profileID, sightID)
}

// ProfileUpdateTwistRate aktualisiert die Drallrate eines Profils
func (a *App) ProfileSetTwistRate(profileID string, twistRateMM float64) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.SetTwistRate(profileID, twistRateMM)
}

// ProfileRemoveTwistRate entfernt den Drall eines Profils
func (a *App) ProfileRemoveTwistRate(profileID string) application.Result[application.ProfileDTO] {
	if a.profileService == nil {
		return application.FailWithMessage[application.ProfileDTO]("Services not initialized - setup not completed")
	}
	return a.profileService.RemoveTwistRate(profileID)
}

// SightCreateSight erstellt eine neue Optik
func (a *App) SightCreateSight(typeName string, modelName string, weightG float64, minMagnification float64, maxMagnification float64) application.Result[application.SightDTO] {
	if a.sightService == nil {
		return application.FailWithMessage[application.SightDTO]("Services not initialized - setup not completed")
	}
	return a.sightService.CreateSight(typeName, modelName, weightG, minMagnification, maxMagnification)
}

// SightLoadSight lädt eine Optik
func (a *App) SightLoadSight(id string) application.Result[application.SightDTO] {
	if a.sightService == nil {
		return application.FailWithMessage[application.SightDTO]("Services not initialized - setup not completed")
	}
	return a.sightService.LoadSight(id)
}

// SightListSights listet alle Optiken
func (a *App) SightListSights() application.Result[[]application.SightDTO] {
	if a.sightService == nil {
		return application.FailWithMessage[[]application.SightDTO]("Services not initialized - setup not completed")
	}
	return a.sightService.ListSights()
}

// SightUpdateSight aktualisiert eine Optik
func (a *App) SightUpdateSight(id string, typeName string, modelName string, weightG float64, minMagnification float64, maxMagnification float64) application.Result[application.SightDTO] {
	if a.sightService == nil {
		return application.FailWithMessage[application.SightDTO]("Services not initialized - setup not completed")
	}
	return a.sightService.UpdateSight(id, typeName, modelName, weightG, minMagnification, maxMagnification)
}

// SightDeleteSight löscht eine Optik
func (a *App) SightDeleteSight(id string) application.Result[bool] {
	if a.sightService == nil {
		return application.FailWithMessage[bool]("Services not initialized - setup not completed")
	}
	return a.sightService.DeleteSight(id)
}

// ==================== PROJECTILE SERVICE DELEGATION ====================

// ProjectileCreateProjectile erstellt ein neues Projektil/Munitionstyp
func (a *App) ProjectileCreateProjectile(name string, weightGrams float64, bc float64) application.Result[application.ProjectileDTO] {
	if a.projectileService == nil {
		return application.FailWithMessage[application.ProjectileDTO]("Services not initialized - setup not completed")
	}
	return a.projectileService.CreateProjectile(name, weightGrams, bc)
}

// ProjectileLoadProjectile lädt ein Projektil nach ID
func (a *App) ProjectileLoadProjectile(id string) application.Result[application.ProjectileDTO] {
	if a.projectileService == nil {
		return application.FailWithMessage[application.ProjectileDTO]("Services not initialized - setup not completed")
	}
	return a.projectileService.LoadProjectile(id)
}

// ProjectileListProjectiles listet alle Projektile auf
func (a *App) ProjectileListProjectiles() application.Result[[]application.ProjectileDTO] {
	if a.projectileService == nil {
		return application.FailWithMessage[[]application.ProjectileDTO]("Services not initialized - setup not completed")
	}
	return a.projectileService.ListProjectiles()
}

// ProjectileDeleteProjectile löscht ein Projektil
func (a *App) ProjectileDeleteProjectile(id string) application.Result[bool] {
	if a.projectileService == nil {
		return application.FailWithMessage[bool]("Services not initialized - setup not completed")
	}
	return a.projectileService.DeleteProjectile(id)
}

// ProjectileUpdateBC aktualisiert den BC-Wert eines Projektils
func (a *App) ProjectileUpdateBC(projectileID string, newBC float64) application.Result[application.ProjectileDTO] {
	if a.projectileService == nil {
		return application.FailWithMessage[application.ProjectileDTO]("Services not initialized - setup not completed")
	}
	return a.projectileService.UpdateBC(projectileID, newBC)
}

// ProjectileUpdateProjectile aktualisiert die Basisdaten eines Projektils
func (a *App) ProjectileUpdateProjectile(projectileID string, name string, weightGrams float64, bc float64) application.Result[application.ProjectileDTO] {
	if a.projectileService == nil {
		return application.FailWithMessage[application.ProjectileDTO]("Services not initialized - setup not completed")
	}
	return a.projectileService.UpdateProjectile(projectileID, name, weightGrams, bc)
}

// ==================== SESSION SERVICE DELEGATION ====================

// SessionCreateSession erstellt eine neue Session
func (a *App) SessionCreateSession(profileID string, projectileID string, temperatureCelsius *float64, note string) application.Result[application.SessionDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[application.SessionDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.CreateSession(profileID, projectileID, temperatureCelsius, note)
}

// SessionRecordShot zeichnet einen neuen Schuss in einer Session auf
func (a *App) SessionRecordShot(sessionID string, velocityMPS float64) application.Result[application.SessionDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[application.SessionDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.RecordShot(sessionID, velocityMPS)
}

// SessionMarkShotInvalid markiert einen Schuss als ungültig
func (a *App) SessionMarkShotInvalid(sessionID string, shotIndex int) application.Result[application.SessionDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[application.SessionDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.MarkShotInvalid(sessionID, shotIndex)
}

// SessionGetStatistics gibt die Statistiken einer Session zurück
func (a *App) SessionGetStatistics(sessionID string) application.Result[application.StatisticsDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[application.StatisticsDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.GetStatistics(sessionID)
}

// SessionLoadSession lädt eine komplette Session mit allen Schüssen
func (a *App) SessionLoadSession(id string) application.Result[application.SessionDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[application.SessionDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.LoadSession(id)
}

// SessionListSessions listet alle Sessions auf (Metadaten ohne Schüsse)
func (a *App) SessionListSessions() application.Result[[]application.SessionMetaDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[[]application.SessionMetaDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.ListSessions()
}

// SessionDeleteSession löscht eine Session
func (a *App) SessionDeleteSession(id string) application.Result[bool] {
	if a.sessionService == nil {
		return application.FailWithMessage[bool]("Services not initialized - setup not completed")
	}
	return a.sessionService.DeleteSession(id)
}

// SessionUpdateNote aktualisiert die Notiz einer Session
func (a *App) SessionUpdateNote(sessionID string, note string) application.Result[application.SessionDTO] {
	if a.sessionService == nil {
		return application.FailWithMessage[application.SessionDTO]("Services not initialized - setup not completed")
	}
	return a.sessionService.UpdateNote(sessionID, note)
}
