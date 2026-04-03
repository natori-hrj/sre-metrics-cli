package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ServiceConfig struct {
	Name             string `json:"name"`
	URL              string `json:"url"`
	VercelToken      string `json:"vercel_token,omitempty"`
	VercelProjectID  string `json:"vercel_project_id,omitempty"`
	VercelTeamID     string `json:"vercel_team_id,omitempty"`
	UptimeRobotKey   string `json:"uptimerobot_api_key,omitempty"`
	UptimeRobotMonID string `json:"uptimerobot_monitor_id,omitempty"`
}

func configPath(service string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	dir := filepath.Join(home, defaultDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create data directory: %w", err)
	}
	return filepath.Join(dir, service+".json"), nil
}

func SaveConfig(service string, cfg ServiceConfig) error {
	path, err := configPath(service)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

func LoadConfig(service string) (ServiceConfig, error) {
	path, err := configPath(service)
	if err != nil {
		return ServiceConfig{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ServiceConfig{}, fmt.Errorf("service %q not initialized. Run 'slo init --service %s' first", service, service)
		}
		return ServiceConfig{}, fmt.Errorf("read config: %w", err)
	}
	var cfg ServiceConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ServiceConfig{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}
