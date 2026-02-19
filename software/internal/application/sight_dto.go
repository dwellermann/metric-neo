package application

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
)

// SightDTO ist die Wails-kompatible Repräsentation einer Optik.
type SightDTO struct {
	ID               string  `json:"id"`
	Type             string  `json:"type"`
	ModelName        string  `json:"modelName"`
	WeightG          float64 `json:"weightG"`
	MinMagnification float64 `json:"minMagnification"`
	MaxMagnification float64 `json:"maxMagnification"`
}

// SightToDTO konvertiert Domain-Entity zu DTO.
func SightToDTO(s *entities.SightingSystem) SightDTO {
	return SightDTO{
		ID:               s.ID,
		Type:             string(s.Type),
		ModelName:        s.ModelName,
		WeightG:          s.Weight.Grams(),
		MinMagnification: s.MinMagnification.Factor(),
		MaxMagnification: s.MaxMagnification.Factor(),
	}
}

// DTOToSight konvertiert DTO zu Domain-Entity.
func DTOToSight(dto SightDTO) (*entities.SightingSystem, error) {
	weight, err := valueobjects.NewMass(dto.WeightG)
	if err != nil {
		return nil, err
	}

	minMag, err := valueobjects.NewMagnification(dto.MinMagnification)
	if err != nil {
		return nil, err
	}

	maxMag, err := valueobjects.NewMagnification(dto.MaxMagnification)
	if err != nil {
		return nil, err
	}

	sight, err := entities.NewSightingSystem(
		entities.SightingSystemType(dto.Type),
		dto.ModelName,
		weight,
		minMag,
		maxMag,
	)
	if err != nil {
		return nil, err
	}

	sight.ID = dto.ID
	return sight, nil
}
