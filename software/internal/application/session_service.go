package application

import (
	"fmt"
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"metric-neo/internal/infrastructure/persistence"
	"path/filepath"
)

// SessionService orchestriert Session-Use-Cases.
//
// KOMPLEXESTER SERVICE:
// - Lädt Profile UND Projectile aus verschiedenen Repositories
// - Erstellt Session mit SNAPSHOT PATTERN (frozen copies)
// - Verwaltet Shots (RecordShot fügt hinzu und speichert)
// - Berechnet Statistiken
type SessionService struct {
	sessionRepo    *persistence.SessionRepository
	profileRepo    *persistence.ProfileRepository
	projectileRepo *persistence.ProjectileRepository
}

// NewSessionService erstellt einen neuen SessionService.
// dataDir ist das Root-Verzeichnis, Repositories nutzen:
// - sessions/ (direkt im Root)
// - inventory/profiles/ (Stammdaten)
// - inventory/projectiles/ (Stammdaten)
func NewSessionService(dataDir string) *SessionService {
	return &SessionService{
		sessionRepo:    persistence.NewSessionRepository(filepath.Join(dataDir, "sessions")),
		profileRepo:    persistence.NewProfileRepository(filepath.Join(dataDir, "inventory", "profiles")),
		projectileRepo: persistence.NewProjectileRepository(filepath.Join(dataDir, "inventory", "projectiles")),
	}
}

// CreateSession erstellt eine neue Session mit Snapshots.
//
// WICHTIG: Snapshot Pattern!
// ProfileSnapshot und ProjectileSnapshot werden beim Erstellen "eingefroren".
// Spätere Änderungen an Profile/Projectile beeinflussen diese Session NICHT.
func (s *SessionService) CreateSession(
	profileID string,
	projectileID string,
	temperatureCelsius *float64,
	note string,
) Result[SessionDTO] {
	// Validierung
	if profileID == "" {
		return FailWithMessage[SessionDTO]("Profile-ID darf nicht leer sein")
	}
	if projectileID == "" {
		return FailWithMessage[SessionDTO]("Projectile-ID darf nicht leer sein")
	}

	// Lade Profile
	profile, err := s.profileRepo.Load(profileID)
	if err != nil {
		return FailWithMessage[SessionDTO](fmt.Sprintf("Profile nicht gefunden: %s", profileID))
	}

	// Lade Projectile
	projectile, err := s.projectileRepo.Load(projectileID)
	if err != nil {
		return FailWithMessage[SessionDTO](fmt.Sprintf("Projectile nicht gefunden: %s", projectileID))
	}

	// Erstelle Session (mit Snapshot Pattern!)
	// GO-KONZEPT: Die NewSession-Funktion macht Deep Copies von Profile und Projectile
	session := entities.NewSession(profile, projectile)

	// Setze optionale Felder
	if temperatureCelsius != nil {
		temp, err := valueobjects.NewTemperature(*temperatureCelsius)
		if err != nil {
			return Fail[SessionDTO](err)
		}
		session.SetTemperature(temp)
	}

	if note != "" {
		session.SetNote(note)
	}

	// Speichere Session
	if err := s.sessionRepo.Save(session); err != nil {
		return Fail[SessionDTO](err)
	}

	return OK(SessionToDTO(session))
}

// RecordShot fügt einen neuen Schuss zu einer Session hinzu.
//
// WORKFLOW:
// 1. Lade Session aus Repository
// 2. Füge Shot hinzu (in-memory)
// 3. Speichere Session zurück
// 4. Return aktualisierte Session als DTO
func (s *SessionService) RecordShot(
	sessionID string,
	velocityMPS float64,
) Result[SessionDTO] {
	if sessionID == "" {
		return FailWithMessage[SessionDTO]("Session-ID darf nicht leer sein")
	}

	// Validierung erfolgt in NewVelocity() (inkl. Bounds-Check).

	// Lade Session
	session, err := s.sessionRepo.Load(sessionID)
	if err != nil {
		return FailWithMessage[SessionDTO]("Session nicht gefunden")
	}

	// Erstelle Velocity Value Object
	velocity, err := valueobjects.NewVelocity(velocityMPS)
	if err != nil {
		return Fail[SessionDTO](err)
	}

	// Füge Shot hinzu
	session.RecordShot(velocity)

	// Speichere aktualisierte Session
	if err := s.sessionRepo.Save(session); err != nil {
		return Fail[SessionDTO](err)
	}

	return OK(SessionToDTO(session))
}

// MarkShotInvalid markiert einen Schuss als ungültig (Fehlmessung).
//
// GO-KONZEPT: Index-basierter Zugriff
// shotIndex ist 0-basiert (wie in JavaScript Arrays)
func (s *SessionService) MarkShotInvalid(sessionID string, shotIndex int) Result[SessionDTO] {
	if sessionID == "" {
		return FailWithMessage[SessionDTO]("Session-ID darf nicht leer sein")
	}

	session, err := s.sessionRepo.Load(sessionID)
	if err != nil {
		return FailWithMessage[SessionDTO]("Session nicht gefunden")
	}

	// Validiere Index
	if shotIndex < 0 || shotIndex >= len(session.Shots) {
		return FailWithMessage[SessionDTO](fmt.Sprintf("Ungültiger Shot-Index: %d", shotIndex))
	}

	// Markiere als invalid
	session.Shots[shotIndex].MarkInvalid()

	// Speichere
	if err := s.sessionRepo.Save(session); err != nil {
		return Fail[SessionDTO](err)
	}

	return OK(SessionToDTO(session))
}

// GetStatistics berechnet Statistiken für eine Session.
func (s *SessionService) GetStatistics(sessionID string) Result[StatisticsDTO] {
	if sessionID == "" {
		return FailWithMessage[StatisticsDTO]("Session-ID darf nicht leer sein")
	}

	session, err := s.sessionRepo.Load(sessionID)
	if err != nil {
		return FailWithMessage[StatisticsDTO]("Session nicht gefunden")
	}

	// Berechne Statistiken (aus session_dto.go)
	stats, err := GetStatistics(session)
	if err != nil {
		return Fail[StatisticsDTO](err)
	}

	return OK(stats)
}

// LoadSession lädt eine Session als vollständiges DTO (mit allen Shots).
func (s *SessionService) LoadSession(id string) Result[SessionDTO] {
	if id == "" {
		return FailWithMessage[SessionDTO]("ID darf nicht leer sein")
	}

	session, err := s.sessionRepo.Load(id)
	if err != nil {
		return FailWithMessage[SessionDTO](fmt.Sprintf("Session nicht gefunden: %s", id))
	}

	return OK(SessionToDTO(session))
}

// ListSessions gibt alle Sessions als Metadata-DTOs zurück (OHNE Shots).
//
// PERFORMANCE: Lädt nur Metadata, nicht alle Shots!
// Für Details: LoadSession() verwenden.
func (s *SessionService) ListSessions() Result[[]SessionMetaDTO] {
	ids, err := s.sessionRepo.List()
	if err != nil {
		return Fail[[]SessionMetaDTO](err)
	}

	// Lade alle Sessions und konvertiere zu Meta-DTOs
	dtos := make([]SessionMetaDTO, 0, len(ids))
	for _, id := range ids {
		session, err := s.sessionRepo.Load(id)
		if err != nil {
			continue // Skip fehlerhafte
		}
		dtos = append(dtos, SessionToMetaDTO(session))
	}

	return OK(dtos)
}

// DeleteSession löscht eine Session.
func (s *SessionService) DeleteSession(id string) Result[bool] {
	if id == "" {
		return FailWithMessage[bool]("ID darf nicht leer sein")
	}

	if err := s.sessionRepo.Delete(id); err != nil {
		return Fail[bool](err)
	}

	return OK(true)
}

// UpdateNote aktualisiert die Notiz einer Session.
func (s *SessionService) UpdateNote(sessionID string, note string) Result[SessionDTO] {
	if sessionID == "" {
		return FailWithMessage[SessionDTO]("Session-ID darf nicht leer sein")
	}

	session, err := s.sessionRepo.Load(sessionID)
	if err != nil {
		return FailWithMessage[SessionDTO]("Session nicht gefunden")
	}

	session.SetNote(note)

	if err := s.sessionRepo.Save(session); err != nil {
		return Fail[SessionDTO](err)
	}

	return OK(SessionToDTO(session))
}
