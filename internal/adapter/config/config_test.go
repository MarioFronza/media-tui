package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MarioFronza/media-tui/internal/adapter/config"
)

func writeConfig(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

const validYAML = `
radarr:
  enabled: true
  url: http://homelab:7878
  api_key: abc123
sonarr:
  enabled: true
  url: http://homelab:8989
  api_key: def456
lidarr:
  enabled: false
  url: http://homelab:8686
  api_key: ""
readarr:
  enabled: false
  url: http://homelab:8787
  api_key: ""
`

func TestLoad_FromWorkingDir(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, validYAML)

	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if !cfg.Radarr.Enabled {
		t.Error("expected Radarr.Enabled = true")
	}
	if cfg.Radarr.URL != "http://homelab:7878" {
		t.Errorf("unexpected Radarr.URL: %s", cfg.Radarr.URL)
	}
	if cfg.Radarr.APIKey != "abc123" {
		t.Errorf("unexpected Radarr.APIKey: %s", cfg.Radarr.APIKey)
	}
	if cfg.Lidarr.Enabled {
		t.Error("expected Lidarr.Enabled = false")
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error when no config file exists")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "{unclosed: [bracket")

	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error on invalid YAML")
	}
}
