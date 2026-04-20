package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Manager struct {
	dataDir      string
	metadataPath string
	metadata     *Metadata
	mu           sync.Mutex
}

func NewManager(dataDir string) (*Manager, error) {
	versionsDir := filepath.Join(dataDir, "versions")
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create versions directory: %w", err)
	}

	metadataPath := filepath.Join(dataDir, "metadata.json")
	metadata, err := loadMetadata(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	return &Manager{
		dataDir:      dataDir,
		metadataPath: metadataPath,
		metadata:     metadata,
	}, nil
}

func (m *Manager) NextVersion() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	maxVer := 0
	for _, v := range m.metadata.Versions {
		if v.Version > maxVer {
			maxVer = v.Version
		}
	}
	return maxVer + 1
}

func (m *Manager) GetVersionPath(version int) string {
	return filepath.Join(m.dataDir, "versions", fmt.Sprintf("%d", version))
}

func (m *Manager) GetLatestVersion() *VersionInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.metadata.Versions) == 0 {
		return nil
	}

	latest := &m.metadata.Versions[0]
	for i := range m.metadata.Versions {
		if m.metadata.Versions[i].Version > latest.Version {
			latest = &m.metadata.Versions[i]
		}
	}
	return latest
}

func (m *Manager) AddVersion(info VersionInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metadata.Versions = append(m.metadata.Versions, info)
	return m.metadata.save(m.metadataPath)
}

func (m *Manager) DeleteVersion(version int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, v := range m.metadata.Versions {
		if v.Version == version {
			m.metadata.Versions = append(m.metadata.Versions[:i], m.metadata.Versions[i+1:]...)
			if err := m.metadata.save(m.metadataPath); err != nil {
				return err
			}
			return os.RemoveAll(m.GetVersionPath(version))
		}
	}
	return fmt.Errorf("version %d not found", version)
}

func (m *Manager) GetAllVersions() []VersionInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]VersionInfo{}, m.metadata.Versions...)
}

func (m *Manager) GetVersion(version int) *VersionInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.metadata.Versions {
		if m.metadata.Versions[i].Version == version {
			return &m.metadata.Versions[i]
		}
	}
	return nil
}
