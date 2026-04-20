package storage

import (
	"encoding/json"
	"os"
	"time"
)

type Metadata struct {
	Versions []VersionInfo `json:"versions"`
}

type VersionInfo struct {
	Version    int        `json:"version"`
	UploadTime time.Time  `json:"upload_time"`
	TotalSize  int64      `json:"total_size"`
	FileCount  int        `json:"file_count"`
	Files      []FileInfo `json:"files"`
}

type FileInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func loadMetadata(path string) (*Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Metadata{Versions: []VersionInfo{}}, nil
		}
		return nil, err
	}

	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func (m *Metadata) save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
