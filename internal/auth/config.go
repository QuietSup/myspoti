package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds persistent CLI settings.
type Config struct {
	ClientID string `json:"client_id"`
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func loadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func clientID() (string, error) {
	if id := os.Getenv("SPOTIFY_CLIENT_ID"); id != "" {
		return id, nil
	}

	cfg, err := loadConfig()
	if err == nil && cfg.ClientID != "" {
		return cfg.ClientID, nil
	}

	return "", fmt.Errorf("client ID not set: export SPOTIFY_CLIENT_ID or add client_id to ~/.config/myspoti/config.json")
}
