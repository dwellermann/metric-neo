package application

import (
	"metric-neo/internal/domain/entities"
	"metric-neo/internal/domain/valueobjects"
	"time"
)

// SessionDTO ist die Wails-kompatible Repräsentation einer Session.
//
// WICHTIG: Enthält SNAPSHOTS als eingebettete DTOs!
// ProfileSnapshot und ProjectileSnapshot sind "frozen" zum Zeitpunkt
// der Session-Erstellung - historische Daten bleiben unverändert.
type SessionDTO struct {
	ID                 string        `json:"id"`
	ProfileSnapshot    ProfileDTO    `json:"profileSnapshot"`    // Embedded snapshot
	ProjectileSnapshot ProjectileDTO `json:"projectileSnapshot"` // Embedded snapshot
	Shots              []ShotDTO     `json:"shots"`
	TemperatureCelsius *float64      `json:"temperatureCelsius,omitempty"`
	Note               string        `json:"note"`
	CreatedAt          string        `json:"createdAt"` // ISO 8601 timestamp
}

// ShotDTO ist die Wails-kompatible Repräsentation eines Shot.
type ShotDTO struct {
	VelocityMPS  float64 `json:"velocityMPS"` // Meters per second
	EnergyJoules float64 `json:"energyJoules"`
	Timestamp    string  `json:"timestamp"` // ISO 8601
	Valid        bool    `json:"valid"`
}

// SessionMetaDTO ist eine leichtgewichtige Version für Listen-Ansichten.
// Enthält KEINE Shots, nur Metadata für Performance.
type SessionMetaDTO struct {
	ID              string  `json:"id"`
	ProfileName     string  `json:"profileName"`    // Aus Snapshot
	ProjectileName  string  `json:"projectileName"` // Aus Snapshot
	ShotCount       int     `json:"shotCount"`
	ValidShotCount  int     `json:"validShotCount"`
	CreatedAt       string  `json:"createdAt"`
	Note            string  `json:"note"`
	AvgVelocityMPS  float64 `json:"avgVelocityMPS,omitempty"`
	AvgEnergyJoules float64 `json:"avgEnergyJoules,omitempty"`
}

// StatisticsDTO enthält berechnete Statistiken einer Session.
type StatisticsDTO struct {
	AvgVelocityMPS    float64 `json:"avgVelocityMPS"`
	StandardDeviation float64 `json:"standardDeviation"`
	MinVelocityMPS    float64 `json:"minVelocityMPS"`
	MaxVelocityMPS    float64 `json:"maxVelocityMPS"`
	ExtremeSpread     float64 `json:"extremeSpread"`
	AvgEnergyJoules   float64 `json:"avgEnergyJoules"`
	ValidShotCount    int     `json:"validShotCount"`
	TotalShotCount    int     `json:"totalShotCount"`
}

// SessionToDTO konvertiert Domain-Session zu vollständigem DTO (mit Shots).
func SessionToDTO(s *entities.Session) SessionDTO {
	dto := SessionDTO{
		ID:                 s.ID,
		ProfileSnapshot:    profileSnapshotToDTO(s.ProfileSnapshot),
		ProjectileSnapshot: projectileSnapshotToDTO(s.ProjectileSnapshot),
		Shots:              make([]ShotDTO, 0, len(s.Shots)),
		Note:               s.Note,
		CreatedAt:          s.CreatedAt.Format(time.RFC3339),
	}

	// Konvertiere Temperature (optional)
	if s.Temperature != nil {
		temp := s.Temperature.Celsius()
		dto.TemperatureCelsius = &temp
	}

	// Konvertiere alle Shots
	for _, shot := range s.Shots {
		dto.Shots = append(dto.Shots, shotToDTO(shot, s.ProjectileSnapshot.Weight))
	}

	return dto
}

// SessionToMetaDTO konvertiert Session zu Metadata-Only DTO (OHNE Shots).
func SessionToMetaDTO(s *entities.Session) SessionMetaDTO {
	dto := SessionMetaDTO{
		ID:             s.ID,
		ProfileName:    s.ProfileSnapshot.Name,
		ProjectileName: s.ProjectileSnapshot.Name,
		ShotCount:      s.ShotCount(),
		ValidShotCount: s.ValidShotCount(),
		CreatedAt:      s.CreatedAt.Format(time.RFC3339),
		Note:           s.Note,
	}

	// Berechne Durchschnittswerte wenn Shots vorhanden
	if s.ValidShotCount() > 0 {
		if avg, err := s.CalculateAverageVelocity(); err == nil {
			dto.AvgVelocityMPS = avg.MetersPerSecond()
		}
		if avgEnergy, err := s.CalculateAverageEnergy(); err == nil {
			dto.AvgEnergyJoules = avgEnergy.Joules()
		}
	}

	return dto
}

// GetStatistics berechnet Statistiken aus Session.
func GetStatistics(s *entities.Session) (StatisticsDTO, error) {
	stats := StatisticsDTO{
		ValidShotCount: s.ValidShotCount(),
		TotalShotCount: s.ShotCount(),
	}

	// Keine gültigen Shots -> Return empty stats
	if s.ValidShotCount() == 0 {
		return stats, nil
	}

	// Berechne alle Statistiken
	if avg, err := s.CalculateAverageVelocity(); err == nil {
		stats.AvgVelocityMPS = avg.MetersPerSecond()
	}

	if sd, err := s.CalculateStandardDeviation(); err == nil {
		stats.StandardDeviation = sd
	}

	if min, err := s.MinVelocity(); err == nil {
		stats.MinVelocityMPS = min.MetersPerSecond()
	}

	if max, err := s.MaxVelocity(); err == nil {
		stats.MaxVelocityMPS = max.MetersPerSecond()
	}

	if es, err := s.ExtremeSpread(); err == nil {
		stats.ExtremeSpread = es
	}

	if avgE, err := s.CalculateAverageEnergy(); err == nil {
		stats.AvgEnergyJoules = avgE.Joules()
	}

	return stats, nil
}

// Helper: Konvertiert ProfileSnapshot (ist bereits ProfileCopy)
func profileSnapshotToDTO(snapshot *entities.Profile) ProfileDTO {
	if snapshot == nil {
		return ProfileDTO{}
	}

	dto := ProfileDTO{
		ID:             snapshot.ID,
		Name:           snapshot.Name,
		Category:       string(snapshot.Category),
		BarrelLengthMM: snapshot.BarrelLength.Millimeters(),
		TriggerWeightG: snapshot.TriggerWeight.Grams(),
		SightHeightMM:  snapshot.SightHeight.Millimeters(),
		TotalWeightG:   snapshot.TotalWeight().Grams(),
	}

	if snapshot.Optic != nil {
		dto.Optic = &OpticDTO{
			Type:             string(snapshot.Optic.Type),
			ModelName:        snapshot.Optic.ModelName,
			WeightG:          snapshot.Optic.Weight.Grams(),
			MinMagnification: snapshot.Optic.MinMagnification.Factor(),
			MaxMagnification: snapshot.Optic.MaxMagnification.Factor(),
		}
	}

	if snapshot.TwistRate != nil {
		mm := snapshot.TwistRate.Millimeters()
		dto.TwistRateMM = &mm
	}

	if snapshot.DefaultAmmoID != nil {
		dto.DefaultAmmoID = snapshot.DefaultAmmoID
	}

	return dto
}

// Helper: Konvertiert ProjectileSnapshot
func projectileSnapshotToDTO(snapshot *entities.Projectile) ProjectileDTO {
	if snapshot == nil {
		return ProjectileDTO{}
	}

	return ProjectileDTO{
		ID:          snapshot.ID,
		Name:        snapshot.Name,
		WeightGrams: snapshot.Weight.Grams(),
		BC:          snapshot.BC,
	}
}

// Helper: Konvertiert Shot zu DTO mit Energy-Berechnung
func shotToDTO(shot *entities.Shot, projectileWeight valueobjects.Mass) ShotDTO {
	return ShotDTO{
		VelocityMPS:  shot.Velocity.MetersPerSecond(),
		EnergyJoules: shot.CalculateEnergy(projectileWeight).Joules(),
		Timestamp:    shot.Timestamp.Format(time.RFC3339),
		Valid:        shot.Valid,
	}
}
