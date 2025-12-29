package instance

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofrs/flock"
	"github.com/varavelio/tribar/internal/config"
)

const (
	lockFileName = "tribar.lock"
	portFileName = "server.port"
)

// ErrAlreadyRunning is returned when another instance of Tribar is running.
var ErrAlreadyRunning = errors.New("tribar is already running")

// Manager handles single-instance enforcement and port discovery.
type Manager struct {
	lock     *flock.Flock
	lockPath string
	portPath string
	listener net.Listener
	port     int
}

// NewManager creates a new instance manager.
// Must be called after config.EnsureDirectories.
func NewManager() *Manager {
	return &Manager{
		lockPath: filepath.Join(config.DirectoryData, lockFileName),
		portPath: filepath.Join(config.DirectoryData, portFileName),
	}
}

// AcquireLock attempts to acquire an exclusive file lock.
// Returns ErrAlreadyRunning if another instance is running.
func (m *Manager) AcquireLock() error {
	m.lock = flock.New(m.lockPath)

	locked, err := m.lock.TryLock()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !locked {
		return ErrAlreadyRunning
	}

	return nil
}

// CreateListener creates a TCP listener on the appropriate port.
// If portExplicit is true, it will only try the specified port and fail if busy.
// If portExplicit is false, it will auto-discover a free port in the range.
func (m *Manager) CreateListener(host string, port int, portExplicit bool) (net.Listener, error) {
	if portExplicit {
		addr := fmt.Sprintf("%s:%d", host, port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("port %d is already in use (you explicitly requested this port)", port)
		}
		m.listener = listener
		m.port = port
		return listener, nil
	}

	for p := port; p < port+config.PortRangeSize; p++ {
		addr := fmt.Sprintf("%s:%d", host, p)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		m.listener = listener
		m.port = p
		return listener, nil
	}

	return nil, fmt.Errorf("could not find a free port in range %d-%d", port, port+config.PortRangeSize-1)
}

// WritePortFile writes the current port to the discovery file.
func (m *Manager) WritePortFile() error {
	content := strconv.Itoa(m.port)
	if err := os.WriteFile(m.portPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write port file: %w", err)
	}
	return nil
}

// Port returns the bound port number.
func (m *Manager) Port() int {
	return m.port
}

// Cleanup releases the lock and removes the port file.
// Safe to call multiple times or on a partially initialized manager.
func (m *Manager) Cleanup() {
	if m.listener != nil {
		_ = m.listener.Close()
	}

	_ = os.Remove(m.portPath)

	if m.lock != nil {
		_ = m.lock.Unlock()
	}
}

// ReadServerPort reads the port number from the discovery file.
// Must be called after config.EnsureDirectories.
func ReadServerPort() (int, error) {
	portPath := filepath.Join(config.DirectoryData, portFileName)
	content, err := os.ReadFile(portPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, errors.New("tribar is not running (no port file found)")
		}
		return 0, fmt.Errorf("failed to read port file: %w", err)
	}

	port, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, fmt.Errorf("invalid port file content: %w", err)
	}

	return port, nil
}
