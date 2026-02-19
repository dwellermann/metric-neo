package application

// ChronoPollResultDTO beschreibt das Ergebnis eines Chrono-Polls.
// Recorded=true bedeutet, dass ein neuer Schuss aufgezeichnet wurde.
type ChronoPollResultDTO struct {
	Recorded    bool        `json:"recorded"`
	VelocityMPS *float64    `json:"velocityMPS,omitempty"`
	Session     *SessionDTO `json:"session,omitempty"`
}
