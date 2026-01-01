package application

import (
	"fmt"
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"metric-neo/internal/infrastructure/persistence"
	"path/filepath"
)

// ProjectileService orchestriert Projectile-Use-Cases.
//
// EINFACHER als ProfileService da Projectile keine komplexen
// optionalen Felder oder Relationen haben.
type ProjectileService struct {
	repo *persistence.ProjectileRepository
}

// NewProjectileService erstellt einen neuen ProjectileService.
// dataDir ist das Root-Verzeichnis, Repository nutzt inventory/projectiles/
func NewProjectileService(dataDir string) *ProjectileService {
	return &ProjectileService{
		repo: persistence.NewProjectileRepository(filepath.Join(dataDir, "inventory", "projectiles")),
	}
}

// CreateProjectile erstellt ein neues Projectile.
func (s *ProjectileService) CreateProjectile(
	name string,
	weightGrams float64,
	bc float64,
) Result[ProjectileDTO] {
	// Validierung
	if name == "" {
		return FailWithMessage[ProjectileDTO]("Name darf nicht leer sein")
	}

	if bc < 0 || bc > 1 {
		return FailWithMessage[ProjectileDTO]("BC muss zwischen 0 und 1 liegen")
	}

	// Value Object erstellen
	weight, err := valueobjects.NewMass(weightGrams)
	if err != nil {
		return Fail[ProjectileDTO](err)
	}

	// Entity erstellen (mit bc Parameter!)
	projectile, err := entities.NewProjectile(name, weight, bc)
	if err != nil {
		return Fail[ProjectileDTO](err)
	}

	// Speichern
	if err := s.repo.Save(projectile); err != nil {
		return Fail[ProjectileDTO](err)
	}

	return OK(ProjectileToDTO(projectile))
}

// LoadProjectile lädt ein Projectile anhand der ID.
func (s *ProjectileService) LoadProjectile(id string) Result[ProjectileDTO] {
	if id == "" {
		return FailWithMessage[ProjectileDTO]("ID darf nicht leer sein")
	}

	projectile, err := s.repo.Load(id)
	if err != nil {
		return FailWithMessage[ProjectileDTO](fmt.Sprintf("Projectile nicht gefunden: %s", id))
	}

	return OK(ProjectileToDTO(projectile))
}

// ListProjectiles gibt alle Projectiles als DTOs zurück.
func (s *ProjectileService) ListProjectiles() Result[[]ProjectileDTO] {
	// Liste alle Projectiles (gibt direkt Entities zurück!)
	projectiles, err := s.repo.List()
	if err != nil {
		return Fail[[]ProjectileDTO](err)
	}

	// Konvertiere zu DTOs
	dtos := make([]ProjectileDTO, 0, len(projectiles))
	for _, projectile := range projectiles {
		dtos = append(dtos, ProjectileToDTO(projectile))
	}

	return OK(dtos)
}

// UpdateBC aktualisiert den Ballistic Coefficient eines Projectiles.
func (s *ProjectileService) UpdateBC(id string, newBC float64) Result[ProjectileDTO] {
	if id == "" {
		return FailWithMessage[ProjectileDTO]("ID darf nicht leer sein")
	}

	if newBC < 0 || newBC > 1 {
		return FailWithMessage[ProjectileDTO]("BC muss zwischen 0 und 1 liegen")
	}

	// Lade Projectile
	projectile, err := s.repo.Load(id)
	if err != nil {
		return FailWithMessage[ProjectileDTO]("Projectile nicht gefunden")
	}

	// Update BC
	projectile.UpdateBC(newBC)

	// Speichern
	if err := s.repo.Save(projectile); err != nil {
		return Fail[ProjectileDTO](err)
	}

	return OK(ProjectileToDTO(projectile))
}

// DeleteProjectile löscht ein Projectile.
func (s *ProjectileService) DeleteProjectile(id string) Result[bool] {
	if id == "" {
		return FailWithMessage[bool]("ID darf nicht leer sein")
	}

	if err := s.repo.Delete(id); err != nil {
		return Fail[bool](err)
	}

	return OK(true)
}
