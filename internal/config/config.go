package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Storage StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Port            int    `yaml:"port"`
	PatchServerURL  string `yaml:"patch_server_url"`
}

type AuthConfig struct {
	APIKey string `yaml:"api_key"`
}

type StorageConfig struct {
	DataDir         string `yaml:"data_dir"`
	MaxUploadSizeMB int    `yaml:"max_upload_size_mb"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Server.Port)
	}
	if c.Server.PatchServerURL == "" {
		return fmt.Errorf("patch_server_url is required")
	}
	if c.Auth.APIKey == "" {
		return fmt.Errorf("api_key is required")
	}
	if c.Storage.DataDir == "" {
		return fmt.Errorf("data_dir is required")
	}
	if c.Storage.MaxUploadSizeMB <= 0 {
		return fmt.Errorf("max_upload_size_mb must be positive")
	}
	return nil
}
