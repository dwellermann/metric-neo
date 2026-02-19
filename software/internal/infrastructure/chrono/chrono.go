package chrono

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"
)

// ChronoService definiert die Schnittstelle für RS232-Chronograph-Kommunikation.
//
// GO-KONZEPT: Interface für Dependency Injection
// Konkrete Implementierung (SerialChrono) wird später hinzugefügt.
// Stub-Implementierung (MockChrono) für Testing.
type ChronoService interface {
	// Connect stellt Verbindung zum Chronograph her
	Connect(port string, baudRate int) error

	// Disconnect trennt Verbindung
	Disconnect() error

	// IsConnected prüft Verbindungsstatus
	IsConnected() bool

	// ReadVelocity liest einen einzelnen Messwert (blockierend)
	// Gibt float32 in m/s zurück oder error
	ReadVelocity() (float32, error)

	// StartAutoRead startet Hintergrund-Leseschleife
	// Messwerte werden in channel geschrieben
	StartAutoRead(velocityChannel chan<- float32, errorChannel chan<- error)

	// StopAutoRead stoppt Hintergrund-Leseschleife
	StopAutoRead()
}

// MockChrono ist eine Test-Implementierung (Stub).
// Nützlich für UI-Testing ohne echte Hardware.
type MockChrono struct {
	mu          sync.Mutex
	connected   bool
	autoRunning bool
	stopChan    chan struct{}
}

// NewMockChrono erstellt einen Mock-Chronograph für Testing.
func NewMockChrono() *MockChrono {
	return &MockChrono{
		connected:   false,
		autoRunning: false,
	}
}

// Connect simuliert Verbindung
func (m *MockChrono) Connect(port string, baudRate int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	m.connected = true
	return nil
}

// Disconnect simuliert Trennung
func (m *MockChrono) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connected = false
	if m.autoRunning {
		close(m.stopChan)
		m.autoRunning = false
	}
	return nil
}

// IsConnected gibt Mock-Verbindungsstatus zurück
func (m *MockChrono) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.connected
}

// ReadVelocity gibt simulierten Messwert zurück (175 m/s + Rauschen)
func (m *MockChrono) ReadVelocity() (float32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return 0, fmt.Errorf("chrono not connected")
	}

	// Simuliere Leseverzögerung
	time.Sleep(100 * time.Millisecond)

	// Gebe konstanten Wert zurück (+ kleines Rauschen für Realismus)
	return 175.0 + float32(time.Now().UnixNano()%10-5)/10.0, nil
}

// StartAutoRead startet Leseschleife mit 1 Messwert pro Sekunde
func (m *MockChrono) StartAutoRead(velocityChannel chan<- float32, errorChannel chan<- error) {
	m.mu.Lock()
	if m.autoRunning || !m.connected {
		m.mu.Unlock()
		if errorChannel != nil {
			errorChannel <- fmt.Errorf("already running or not connected")
		}
		return
	}

	m.autoRunning = true
	m.stopChan = make(chan struct{})
	m.mu.Unlock()

	// Starte Goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				return
			case <-ticker.C:
				velocity, err := m.ReadVelocity()
				if err != nil && errorChannel != nil {
					errorChannel <- err
				} else if velocityChannel != nil {
					velocityChannel <- velocity
				}
			}
		}
	}()
}

// StopAutoRead stoppt Leseschleife
func (m *MockChrono) StopAutoRead() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.autoRunning && m.stopChan != nil {
		close(m.stopChan)
		m.autoRunning = false
	}
}

// SerialChrono ist die echte RS232-Implementierung.
// Nutzt go.bug.st/serial für Kommunikation mit Chronograph.
type SerialChrono struct {
	mu          sync.Mutex
	connected   bool
	autoRunning bool
	port        serial.Port
	reader      *bufio.Reader
	stopChan    chan struct{}
}

// NewSerialChrono erstellt einen echten Chronograph.
func NewSerialChrono() *SerialChrono {
	return &SerialChrono{
		connected:   false,
		autoRunning: false,
	}
}

// Connect stellt RS232-Verbindung her.
func (s *SerialChrono) Connect(port string, baudRate int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if baudRate <= 0 {
		return fmt.Errorf("baud rate must be > 0")
	}

	if s.port != nil {
		return fmt.Errorf("already connected")
	}

	mode := &serial.Mode{
		BaudRate: baudRate,
	}

	serialPort, err := serial.Open(port, mode)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %w", port, err)
	}

	s.port = serialPort
	s.reader = bufio.NewReader(serialPort)
	s.connected = true
	return nil
}

// Disconnect trennt RS232-Verbindung.
func (s *SerialChrono) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.autoRunning && s.stopChan != nil {
		close(s.stopChan)
		s.autoRunning = false
	}

	if s.port != nil {
		err := s.port.Close()
		s.port = nil
		s.reader = nil
		s.connected = false
		return err
	}

	s.connected = false
	return nil
}

// IsConnected gibt Verbindungsstatus zurück.
func (s *SerialChrono) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.connected
}

// ReadVelocity liest von RS232.
// Erwartet Linie mit Geschwindigkeit in m/s oder fps (wird konvertiert).
// Kommas werden als Dezimaltrennzeichen erkannt (z.B. "175,23" oder "175.23").
func (s *SerialChrono) ReadVelocity() (float32, error) {
	s.mu.Lock()
	if !s.connected || s.reader == nil {
		s.mu.Unlock()
		return 0, fmt.Errorf("chrono not connected")
	}
	reader := s.reader
	s.mu.Unlock()

	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, fmt.Errorf("chrono disconnected (EOF)")
		}
		return 0, fmt.Errorf("failed to read from chrono: %w", err)
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return 0, fmt.Errorf("empty chrono line")
	}

	// Ersetze Komma durch Punkt für Dezimalzahl-Parsing
	line = strings.ReplaceAll(line, ",", ".")

	value, err := strconv.ParseFloat(line, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid chrono value format: %s (error: %w)", line, err)
	}

	// Validierung: Bounds-Check (0.1 - 5000 m/s für Debugging bis Extreme Rifles)
	if value < 0.1 || value > 5000 {
		return 0, fmt.Errorf("velocity out of bounds (0.1-5000 m/s), got: %.2f", value)
	}

	return float32(value), nil
}

// StartAutoRead startet Leseschleife im Hintergrund.
// Messwerte werden kontinuierlich in velocityChannel geschrieben.
func (s *SerialChrono) StartAutoRead(velocityChannel chan<- float32, errorChannel chan<- error) {
	s.mu.Lock()
	if s.autoRunning || !s.connected {
		s.mu.Unlock()
		if errorChannel != nil {
			errorChannel <- fmt.Errorf("already running or not connected")
		}
		return
	}

	s.autoRunning = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				velocity, err := s.ReadVelocity()
				if err != nil {
					if errorChannel != nil {
						errorChannel <- err
					}
					// Bei Fehler: kurze Pause, dann retry
					time.Sleep(100 * time.Millisecond)
					continue
				}
				if velocityChannel != nil {
					velocityChannel <- velocity
				}
			}
		}
	}()
}

// StopAutoRead stoppt Leseschleife.
func (s *SerialChrono) StopAutoRead() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.autoRunning && s.stopChan != nil {
		close(s.stopChan)
		s.autoRunning = false
	}
}
