package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Radarr  ServiceConfig `yaml:"radarr"`
	Sonarr  ServiceConfig `yaml:"sonarr"`
	Lidarr  ServiceConfig `yaml:"lidarr"`
	Readarr ServiceConfig `yaml:"readarr"`
}

type ServiceConfig struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
	APIKey  string `yaml:"api_key"`
}

func Load() (*Config, error) {
	candidates := []string{
		"config.yaml",
	}

	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".config", "media-tui", "config.yaml"))
	}

	for _, path := range candidates {
		cfg, err := loadFile(path)
		if err == nil {
			return cfg, nil
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("loading %s: %w", path, err)
		}
	}

	return nil, fmt.Errorf("config file not found (tried: %v)", candidates)
}

func loadFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
