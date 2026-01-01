package application

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
)

// ProjectileDTO ist die Wails-kompatible Repräsentation eines Projectile.
//
// GO-KONZEPT: DTOs (Data Transfer Objects)
// DTOs verwenden primitive Typen statt Value Objects, damit Wails
// sie als JSON über die JavaScript-Bridge senden kann.
//
// Beispiel: weightGrams (float64) statt Weight (valueobjects.Mass)
type ProjectileDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	WeightGrams float64 `json:"weightGrams"`
	BC          float64 `json:"bc"` // Ballistic Coefficient
}

// ProjectileToDTO konvertiert Domain-Entity zu DTO.
func ProjectileToDTO(p *entities.Projectile) ProjectileDTO {
	return ProjectileDTO{
		ID:          p.ID,
		Name:        p.Name,
		WeightGrams: p.Weight.Grams(),
		BC:          p.BC,
	}
}

// DTOToProjectile konvertiert DTO zu Domain-Entity.
// WICHTIG: Validierung erfolgt durch Value Objects!
func DTOToProjectile(dto ProjectileDTO) (*entities.Projectile, error) {
	// Value Object erstellen (mit Validierung)
	weight, err := valueobjects.NewMass(dto.WeightGrams)
	if err != nil {
		return nil, err
	}

	// NewProjectile benötigt name, weight UND bc
	projectile, err := entities.NewProjectile(dto.Name, weight, dto.BC)
	if err != nil {
		return nil, err
	}

	// Setze ID (für Updates)
	projectile.ID = dto.ID

	return projectile, nil
}
