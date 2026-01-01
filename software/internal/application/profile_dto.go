package application

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
)

type ProfileDTO struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"` // "air_rifle", "air_pistol", etc.
	BarrelLengthMM float64   `json:"barrelLengthMM"`
	TriggerWeightG float64   `json:"triggerWeightG"`
	SightHeightMM  float64   `json:"sightHeightMM"`
	Optic          *OpticDTO `json:"optic,omitempty"`
	TwistRateMM    *float64  `json:"twistRateMM,omitempty"`
	DefaultAmmoID  *string   `json:"defaultAmmoID,omitempty"`
	TotalWeightG   float64   `json:"totalWeightG"` // Calculated
}

type OpticDTO struct {
	Type             string  `json:"type"` // "scope", "red_dot", etc.
	ModelName        string  `json:"modelName"`
	WeightG          float64 `json:"weightG"`
	MinMagnification float64 `json:"minMagnification"`
	MaxMagnification float64 `json:"maxMagnification"`
}

// ProfileToDTO konvertiert Domain-Entity zu DTO
func ProfileToDTO(p *entities.Profile) ProfileDTO {
	dto := ProfileDTO{
		ID:             p.ID,
		Name:           p.Name,
		Category:       string(p.Category),
		BarrelLengthMM: p.BarrelLength.Millimeters(),
		TriggerWeightG: p.TriggerWeight.Grams(),
		SightHeightMM:  p.SightHeight.Millimeters(),
		TotalWeightG:   p.TotalWeight().Grams(),
	}

	if p.Optic != nil {
		dto.Optic = &OpticDTO{
			Type:             string(p.Optic.Type),
			ModelName:        p.Optic.ModelName,
			WeightG:          p.Optic.Weight.Grams(),
			MinMagnification: p.Optic.MinMagnification.Factor(),
			MaxMagnification: p.Optic.MaxMagnification.Factor(),
		}
	}

	if p.TwistRate != nil {
		mm := p.TwistRate.Millimeters()
		dto.TwistRateMM = &mm
	}

	if p.DefaultAmmoID != nil {
		dto.DefaultAmmoID = p.DefaultAmmoID
	}

	return dto
}

// DTOToProfile konvertiert DTO zu Domain-Entity
func DTOToProfile(dto ProfileDTO) (*entities.Profile, error) {
	barrelLength, err := valueobjects.NewLength(dto.BarrelLengthMM)
	if err != nil {
		return nil, err
	}

	triggerWeight, err := valueobjects.NewMass(dto.TriggerWeightG)
	if err != nil {
		return nil, err
	}

	sightHeight, err := valueobjects.NewLength(dto.SightHeightMM)
	if err != nil {
		return nil, err
	}

	profile, err := entities.NewProfile(
		dto.Name,
		entities.ProfileCategory(dto.Category),
		barrelLength,
		triggerWeight,
		sightHeight,
	)
	if err != nil {
		return nil, err
	}

	// Setze optionale Felder
	if dto.Optic != nil {
		optic, err := opticDTOToEntity(dto.Optic)
		if err != nil {
			return nil, err
		}
		profile.SetOptic(optic)
	}

	if dto.TwistRateMM != nil {
		twist, err := valueobjects.NewLength(*dto.TwistRateMM)
		if err != nil {
			return nil, err
		}
		profile.SetTwistRate(twist)
	}

	if dto.DefaultAmmoID != nil {
		profile.SetDefaultAmmo(*dto.DefaultAmmoID)
	}

	// Setze ID (für Updates)
	profile.ID = dto.ID

	return profile, nil
}

func opticDTOToEntity(dto *OpticDTO) (*entities.SightingSystem, error) {
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

	return entities.NewSightingSystem(
		entities.SightingSystemType(dto.Type),
		dto.ModelName,
		weight,
		minMag,
		maxMag,
	)
}
