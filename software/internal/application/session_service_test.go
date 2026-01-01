package application

import (
	"testing"
)

func TestSessionService_CreateSession(t *testing.T) {
	dir := t.TempDir()

	// Setup: Erstelle Profile und Projectile mit neuer inventory/ Struktur
	profileService := NewProfileService(dir)
	projectileService := NewProjectileService(dir)
	sessionService := NewSessionService(dir)

	profileResult := profileService.CreateProfile("Test Profile", "air_rifle", 420.0, 500.0, 50.0)
	projectileResult := projectileService.CreateProjectile("JSB Exact", 0.547, 0.024)

	if !profileResult.Success || !projectileResult.Success {
		t.Fatal("Setup failed")
	}

	// Erstelle Session
	temp := 21.5
	result := sessionService.CreateSession(
		profileResult.Data.ID,
		projectileResult.Data.ID,
		&temp,
		"Test session",
	)

	if !result.Success {
		t.Fatalf("CreateSession failed: %s", result.Error)
	}

	dto := result.Data
	if dto.ID == "" {
		t.Error("UUID not generated")
	}

	if dto.ProfileSnapshot.Name != "Test Profile" {
		t.Error("Profile snapshot incorrect")
	}

	if dto.ProjectileSnapshot.Name != "JSB Exact" {
		t.Error("Projectile snapshot incorrect")
	}

	if dto.TemperatureCelsius == nil || *dto.TemperatureCelsius != 21.5 {
		t.Error("Temperature not set")
	}

	t.Logf("✓ Session created with snapshots")
}

func TestSessionService_SnapshotIsolation(t *testing.T) {
	dir := t.TempDir()

	profileService := NewProfileService(dir)
	projectileService := NewProjectileService(dir)
	sessionService := NewSessionService(dir)

	// Erstelle Profile und Projectile
	profileResult := profileService.CreateProfile("Original Profile", "air_rifle", 420.0, 500.0, 50.0)
	projectileResult := projectileService.CreateProjectile("Original Pellet", 0.547, 0.024)

	profileID := profileResult.Data.ID
	projectileID := projectileResult.Data.ID

	// Erstelle Session (macht Snapshots!)
	sessionResult := sessionService.CreateSession(profileID, projectileID, nil, "Snapshot test")
	sessionID := sessionResult.Data.ID

	// ÄNDERE Originals NACH Session-Erstellung
	// Simuliert: User ändert Profile-Name und Projectile-BC
	profileService.CreateProfile("CHANGED PROFILE", "air_rifle", 420.0, 500.0, 50.0)
	projectileService.UpdateBC(projectileID, 0.999)

	// Lade Session zurück
	loadResult := sessionService.LoadSession(sessionID)
	if !loadResult.Success {
		t.Fatalf("LoadSession failed: %s", loadResult.Error)
	}

	// Prüfe: Snapshots sind UNVERÄNDERT!
	if loadResult.Data.ProfileSnapshot.Name != "Original Profile" {
		t.Error("❌ ProfileSnapshot changed! Snapshot isolation broken!")
	}

	if loadResult.Data.ProjectileSnapshot.BC != 0.024 {
		t.Error("❌ ProjectileSnapshot changed! Snapshot isolation broken!")
	}

	t.Log("✓ Snapshot isolation works - session preserves historical data")
}

func TestSessionService_RecordShots(t *testing.T) {
	dir := t.TempDir()

	profileService := NewProfileService(dir)
	projectileService := NewProjectileService(dir)
	sessionService := NewSessionService(dir)

	// Setup
	profileResult := profileService.CreateProfile("Test", "air_rifle", 420.0, 500.0, 50.0)
	projectileResult := projectileService.CreateProjectile("Test", 0.547, 0.024)

	sessionResult := sessionService.CreateSession(
		profileResult.Data.ID,
		projectileResult.Data.ID,
		nil,
		"",
	)
	sessionID := sessionResult.Data.ID

	// Füge Schüsse hinzu
	velocities := []float64{175.2, 176.1, 174.8, 175.5, 176.3}

	for i, v := range velocities {
		result := sessionService.RecordShot(sessionID, v)
		if !result.Success {
			t.Fatalf("RecordShot %d failed: %s", i, result.Error)
		}
	}

	// Lade Session und prüfe Shots
	loadResult := sessionService.LoadSession(sessionID)
	if len(loadResult.Data.Shots) != 5 {
		t.Errorf("Shot count = %d, want 5", len(loadResult.Data.Shots))
	}

	// Prüfe ersten Shot
	firstShot := loadResult.Data.Shots[0]
	if firstShot.VelocityMPS != 175.2 {
		t.Errorf("First shot velocity = %.1f, want 175.2", firstShot.VelocityMPS)
	}

	if firstShot.EnergyJoules < 8.0 || firstShot.EnergyJoules > 9.0 {
		t.Errorf("Energy = %.2f J, expected ~8.4 J", firstShot.EnergyJoules)
	}

	t.Logf("✓ Recorded %d shots", len(loadResult.Data.Shots))
}

func TestSessionService_GetStatistics(t *testing.T) {
	dir := t.TempDir()

	profileService := NewProfileService(dir)
	projectileService := NewProjectileService(dir)
	sessionService := NewSessionService(dir)

	// Setup
	profileResult := profileService.CreateProfile("Test", "air_rifle", 420.0, 500.0, 50.0)
	projectileResult := projectileService.CreateProjectile("Test", 0.547, 0.024)

	sessionResult := sessionService.CreateSession(
		profileResult.Data.ID,
		projectileResult.Data.ID,
		nil,
		"",
	)
	sessionID := sessionResult.Data.ID

	// Füge bekannte Werte hinzu: 170, 175, 180 -> Mean=175, SD≈4.08
	sessionService.RecordShot(sessionID, 170.0)
	sessionService.RecordShot(sessionID, 175.0)
	sessionService.RecordShot(sessionID, 180.0)

	// Hole Statistiken
	statsResult := sessionService.GetStatistics(sessionID)
	if !statsResult.Success {
		t.Fatalf("GetStatistics failed: %s", statsResult.Error)
	}

	stats := statsResult.Data

	// Prüfe Durchschnitt
	if stats.AvgVelocityMPS != 175.0 {
		t.Errorf("Avg = %.1f, want 175.0", stats.AvgVelocityMPS)
	}

	// Prüfe SD (sollte ~4.08 sein)
	expectedSD := 4.08
	tolerance := 0.5
	if stats.StandardDeviation < expectedSD-tolerance || stats.StandardDeviation > expectedSD+tolerance {
		t.Errorf("SD = %.2f, expected ~%.2f", stats.StandardDeviation, expectedSD)
	}

	// Prüfe Min/Max
	if stats.MinVelocityMPS != 170.0 {
		t.Errorf("Min = %.1f, want 170.0", stats.MinVelocityMPS)
	}

	if stats.MaxVelocityMPS != 180.0 {
		t.Errorf("Max = %.1f, want 180.0", stats.MaxVelocityMPS)
	}

	// Prüfe Extreme Spread
	if stats.ExtremeSpread != 10.0 {
		t.Errorf("ES = %.1f, want 10.0", stats.ExtremeSpread)
	}

	t.Logf("✓ Statistics:")
	t.Logf("  Avg: %.2f m/s", stats.AvgVelocityMPS)
	t.Logf("  SD: %.2f m/s", stats.StandardDeviation)
	t.Logf("  ES: %.2f m/s", stats.ExtremeSpread)
	t.Logf("  Avg Energy: %.2f J", stats.AvgEnergyJoules)
}

func TestSessionService_ListSessions(t *testing.T) {
	dir := t.TempDir()

	profileService := NewProfileService(dir)
	projectileService := NewProjectileService(dir)
	sessionService := NewSessionService(dir)

	// Setup
	profileResult := profileService.CreateProfile("Test", "air_rifle", 420.0, 500.0, 50.0)
	projectileResult := projectileService.CreateProjectile("Test", 0.547, 0.024)

	// Erstelle mehrere Sessions
	sessionService.CreateSession(profileResult.Data.ID, projectileResult.Data.ID, nil, "Session 1")
	sessionService.CreateSession(profileResult.Data.ID, projectileResult.Data.ID, nil, "Session 2")

	// Liste (sollte Metadata zurückgeben, keine Shots)
	listResult := sessionService.ListSessions()
	if !listResult.Success {
		t.Fatalf("ListSessions failed: %s", listResult.Error)
	}

	if len(listResult.Data) != 2 {
		t.Errorf("Session count = %d, want 2", len(listResult.Data))
	}

	// Prüfe Metadata
	first := listResult.Data[0]
	if first.ProfileName != "Test" {
		t.Error("Profile name not in metadata")
	}

	if first.ProjectileName != "Test" {
		t.Error("Projectile name not in metadata")
	}

	t.Logf("✓ Listed %d sessions (metadata only)", len(listResult.Data))
}

func TestSessionService_MarkShotInvalid(t *testing.T) {
	dir := t.TempDir()

	profileService := NewProfileService(dir)
	projectileService := NewProjectileService(dir)
	sessionService := NewSessionService(dir)

	// Setup
	profileResult := profileService.CreateProfile("Test", "air_rifle", 420.0, 500.0, 50.0)
	projectileResult := projectileService.CreateProjectile("Test", 0.547, 0.024)

	sessionResult := sessionService.CreateSession(
		profileResult.Data.ID,
		projectileResult.Data.ID,
		nil,
		"",
	)
	sessionID := sessionResult.Data.ID

	// Füge Schüsse hinzu
	sessionService.RecordShot(sessionID, 175.0)
	sessionService.RecordShot(sessionID, 0.0) // Fehlmessung!
	sessionService.RecordShot(sessionID, 176.0)

	// Markiere zweiten Shot als invalid
	markResult := sessionService.MarkShotInvalid(sessionID, 1)
	if !markResult.Success {
		t.Fatalf("MarkShotInvalid failed: %s", markResult.Error)
	}

	// Hole Statistiken - sollten nur 2 valide Shots berücksichtigen
	statsResult := sessionService.GetStatistics(sessionID)
	stats := statsResult.Data

	if stats.ValidShotCount != 2 {
		t.Errorf("ValidShotCount = %d, want 2", stats.ValidShotCount)
	}

	if stats.TotalShotCount != 3 {
		t.Errorf("TotalShotCount = %d, want 3", stats.TotalShotCount)
	}

	// Durchschnitt sollte (175 + 176) / 2 = 175.5 sein
	if stats.AvgVelocityMPS != 175.5 {
		t.Errorf("Avg = %.1f, want 175.5 (invalid shot excluded)", stats.AvgVelocityMPS)
	}

	t.Log("✓ Invalid shot excluded from statistics")
}
