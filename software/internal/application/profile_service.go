package application

import (
	"fmt"
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"metric-neo/internal/infrastructure/persistence"
	"path/filepath"
)

// ProfileService orchestriert Profile-Use-Cases.
//
// GO-KONZEPT: Service Layer / Application Layer
// Services sind STATELESS - sie laden immer aus dem Repository.
// Keine In-Memory-Caches (SSD ist schnell genug).
//
// ARCHITECTURE LAYERS:
// Wails Binding -> ProfileService -> Domain (Entities/ValueObjects) -> Repository -> JSON
type ProfileService struct {
	repo *persistence.ProfileRepository
}

// NewProfileService erstellt einen neuen ProfileService.
// dataDir ist das Root-Verzeichnis, Repository nutzt inventory/profiles/
func NewProfileService(dataDir string) *ProfileService {
	return &ProfileService{
		repo: persistence.NewProfileRepository(filepath.Join(dataDir, "inventory", "profiles")),
	}
}

// CreateProfile erstellt ein neues Profile und speichert es.
//
// VALIDATION: Erfolgt in Domain-Entities (ValueObjects + Entity-Constructor)
// INPUT: Primitive Typen von Wails
// OUTPUT: Result[ProfileDTO] für Wails
func (s *ProfileService) CreateProfile(
	name string,
	category string,
	barrelLengthMM float64,
	triggerWeightG float64,
	sightHeightMM float64,
) Result[ProfileDTO] {
	// Validierung: Name darf nicht leer sein
	if name == "" {
		return FailWithMessage[ProfileDTO]("Name darf nicht leer sein")
	}

	// Value Objects erstellen (mit Validierung!)
	barrelLength, err := valueobjects.NewLength(barrelLengthMM)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	triggerWeight, err := valueobjects.NewMass(triggerWeightG)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	sightHeight, err := valueobjects.NewLength(sightHeightMM)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	// Domain-Entity erstellen
	profile, err := entities.NewProfile(
		name,
		entities.ProfileCategory(category),
		barrelLength,
		triggerWeight,
		sightHeight,
	)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	// Persistieren
	if err := s.repo.Save(profile); err != nil {
		return Fail[ProfileDTO](err)
	}

	// Konvertiere zu DTO und return
	return OK(ProfileToDTO(profile))
}

// LoadProfile lädt ein Profile anhand der ID.
func (s *ProfileService) LoadProfile(id string) Result[ProfileDTO] {
	// Validierung: ID darf nicht leer sein
	if id == "" {
		return FailWithMessage[ProfileDTO]("ID darf nicht leer sein")
	}

	// Lade aus Repository
	profile, err := s.repo.Load(id)
	if err != nil {
		return FailWithMessage[ProfileDTO](fmt.Sprintf("Profile nicht gefunden: %s", id))
	}

	return OK(ProfileToDTO(profile))
}

// ListProfiles gibt alle Profile als DTOs zurück.
func (s *ProfileService) ListProfiles() Result[[]ProfileDTO] {
	// Liste IDs
	ids, err := s.repo.List()
	if err != nil {
		return Fail[[]ProfileDTO](err)
	}

	// Lade alle Profile
	dtos := make([]ProfileDTO, 0, len(ids))
	for _, id := range ids {
		profile, err := s.repo.Load(id)
		if err != nil {
			// Skip fehlerhafte Profile (nicht crashen!)
			continue
		}
		dtos = append(dtos, ProfileToDTO(profile))
	}

	return OK(dtos)
}

// DeleteProfile löscht ein Profile.
func (s *ProfileService) DeleteProfile(id string) Result[bool] {
	if id == "" {
		return FailWithMessage[bool]("ID darf nicht leer sein")
	}

	if err := s.repo.Delete(id); err != nil {
		return Fail[bool](err)
	}

	return OK(true)
}

// SetOptic fügt einem Profile eine Optik hinzu.
//
// GO-KONZEPT: Stateless Service Pattern
// 1. Lade Profile
// 2. Modifiziere in-memory
// 3. Speichere zurück
// 4. Return aktualisiertes DTO
func (s *ProfileService) SetOptic(
	profileID string,
	opticType string,
	modelName string,
	weightG float64,
	minMagnification float64,
	maxMagnification float64,
) Result[ProfileDTO] {
	// Lade Profile
	profile, err := s.repo.Load(profileID)
	if err != nil {
		return FailWithMessage[ProfileDTO]("Profile nicht gefunden")
	}

	// Erstelle SightingSystem
	weight, err := valueobjects.NewMass(weightG)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	minMag, err := valueobjects.NewMagnification(minMagnification)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	maxMag, err := valueobjects.NewMagnification(maxMagnification)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	optic, err := entities.NewSightingSystem(
		entities.SightingSystemType(opticType),
		modelName,
		weight,
		minMag,
		maxMag,
	)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	// Setze Optik
	profile.SetOptic(optic)

	// Speichere
	if err := s.repo.Save(profile); err != nil {
		return Fail[ProfileDTO](err)
	}

	return OK(ProfileToDTO(profile))
}

// RemoveOptic entfernt die Optik von einem Profile.
func (s *ProfileService) RemoveOptic(profileID string) Result[ProfileDTO] {
	profile, err := s.repo.Load(profileID)
	if err != nil {
		return FailWithMessage[ProfileDTO]("Profile nicht gefunden")
	}

	profile.SetOptic(nil)

	if err := s.repo.Save(profile); err != nil {
		return Fail[ProfileDTO](err)
	}

	return OK(ProfileToDTO(profile))
}

// SetTwistRate setzt den Drall (Twist Rate) für ein Profile.
func (s *ProfileService) SetTwistRate(profileID string, twistRateMM float64) Result[ProfileDTO] {
	profile, err := s.repo.Load(profileID)
	if err != nil {
		return FailWithMessage[ProfileDTO]("Profile nicht gefunden")
	}

	twistRate, err := valueobjects.NewLength(twistRateMM)
	if err != nil {
		return Fail[ProfileDTO](err)
	}

	profile.SetTwistRate(twistRate)

	if err := s.repo.Save(profile); err != nil {
		return Fail[ProfileDTO](err)
	}

	return OK(ProfileToDTO(profile))
}
