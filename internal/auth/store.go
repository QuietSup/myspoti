package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Tokens is the persisted OAuth token set.
type Tokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "myspoti"), nil
}

func tokenPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tokens.json"), nil
}

func saveTokens(tokens Tokens) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path, err := tokenPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func loadTokens() (Tokens, error) {
	path, err := tokenPath()
	if err != nil {
		return Tokens{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Tokens{}, err
	}

	var tokens Tokens
	if err := json.Unmarshal(data, &tokens); err != nil {
		return Tokens{}, err
	}
	return tokens, nil
}

// AccessToken returns a valid access token, refreshing if needed.
func AccessToken() (string, error) {
	tokens, err := loadTokens()
	if err != nil {
		return "", fmt.Errorf("loading tokens: %w (run `myspoti auth` first)", err)
	}

	if time.Until(tokens.ExpiresAt) > 60*time.Second {
		return tokens.AccessToken, nil
	}

	refreshed, err := refreshAccessToken(tokens.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("refreshing token: %w", err)
	}
	if err := saveTokens(refreshed); err != nil {
		return "", fmt.Errorf("saving refreshed tokens: %w", err)
	}
	return refreshed.AccessToken, nil
}
