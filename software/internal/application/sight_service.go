package application

import (
	"fmt"
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"metric-neo/internal/infrastructure/persistence"
	"path/filepath"
)

// SightService orchestriert SightingSystem-Use-Cases.
type SightService struct {
	repo *persistence.SightRepository
}

// NewSightService erstellt einen neuen SightService.
func NewSightService(dataDir string) *SightService {
	return &SightService{
		repo: persistence.NewSightRepository(filepath.Join(dataDir, "inventory", "sights")),
	}
}

// CreateSight erstellt eine neue Optik.
func (s *SightService) CreateSight(
	typeName string,
	modelName string,
	weightG float64,
	minMagnification float64,
	maxMagnification float64,
) Result[SightDTO] {
	if modelName == "" {
		return FailWithMessage[SightDTO]("Model name darf nicht leer sein")
	}

	weight, err := valueobjects.NewMass(weightG)
	if err != nil {
		return Fail[SightDTO](err)
	}

	minMag, err := valueobjects.NewMagnification(minMagnification)
	if err != nil {
		return Fail[SightDTO](err)
	}

	maxMag, err := valueobjects.NewMagnification(maxMagnification)
	if err != nil {
		return Fail[SightDTO](err)
	}

	sight, err := entities.NewSightingSystem(
		entities.SightingSystemType(typeName),
		modelName,
		weight,
		minMag,
		maxMag,
	)
	if err != nil {
		return Fail[SightDTO](err)
	}

	if err := s.repo.Save(sight); err != nil {
		return Fail[SightDTO](err)
	}

	return OK(SightToDTO(sight))
}

// LoadSight lädt eine Optik anhand der ID.
func (s *SightService) LoadSight(id string) Result[SightDTO] {
	if id == "" {
		return FailWithMessage[SightDTO]("ID darf nicht leer sein")
	}

	sight, err := s.repo.Load(id)
	if err != nil {
		return FailWithMessage[SightDTO](fmt.Sprintf("SightingSystem nicht gefunden: %s", id))
	}

	return OK(SightToDTO(sight))
}

// ListSights gibt alle Optiken als DTOs zurück.
func (s *SightService) ListSights() Result[[]SightDTO] {
	sights, err := s.repo.List()
	if err != nil {
		return Fail[[]SightDTO](err)
	}

	dtos := make([]SightDTO, 0, len(sights))
	for _, sight := range sights {
		dtos = append(dtos, SightToDTO(sight))
	}

	return OK(dtos)
}

// UpdateSight aktualisiert eine bestehende Optik.
func (s *SightService) UpdateSight(
	id string,
	typeName string,
	modelName string,
	weightG float64,
	minMagnification float64,
	maxMagnification float64,
) Result[SightDTO] {
	if id == "" {
		return FailWithMessage[SightDTO]("ID darf nicht leer sein")
	}

	sight, err := s.repo.Load(id)
	if err != nil {
		return FailWithMessage[SightDTO](fmt.Sprintf("SightingSystem nicht gefunden: %s", id))
	}

	if modelName == "" {
		return FailWithMessage[SightDTO]("Model name darf nicht leer sein")
	}

	weight, err := valueobjects.NewMass(weightG)
	if err != nil {
		return Fail[SightDTO](err)
	}

	minMag, err := valueobjects.NewMagnification(minMagnification)
	if err != nil {
		return Fail[SightDTO](err)
	}

	maxMag, err := valueobjects.NewMagnification(maxMagnification)
	if err != nil {
		return Fail[SightDTO](err)
	}

	// Update fields
	sight.Type = entities.SightingSystemType(typeName)
	sight.ModelName = modelName
	sight.Weight = weight
	sight.MinMagnification = minMag
	sight.MaxMagnification = maxMag

	if err := s.repo.Save(sight); err != nil {
		return Fail[SightDTO](err)
	}

	return OK(SightToDTO(sight))
}

// DeleteSight löscht eine Optik.
func (s *SightService) DeleteSight(id string) Result[bool] {
	if id == "" {
		return FailWithMessage[bool]("ID darf nicht leer sein")
	}

	if err := s.repo.Delete(id); err != nil {
		return Fail[bool](err)
	}

	return OK(true)
}
