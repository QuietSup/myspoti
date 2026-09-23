package player

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"myspoti/internal/auth"
)

const apiBase = "https://api.spotify.com/v1"

// Next skips to the next track.
func Next() error {
	return post("/me/player/next")
}

// Previous skips to the previous track.
func Previous() error {
	return post("/me/player/previous")
}

func post(path string) error {
	token, err := auth.AccessToken()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusNoContent, http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized — run `myspoti auth` again")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden: %s (Premium required, or missing scope)", trimBody(body))
	case http.StatusNotFound:
		return fmt.Errorf("no active device — open Spotify on a device first")
	default:
		return fmt.Errorf("spotify API %d: %s", resp.StatusCode, trimBody(body))
	}
}

func trimBody(body []byte) string {
	s := string(body)
	if len(s) > 200 {
		return s[:200] + "..."
	}
	if s == "" {
		return "(empty body)"
	}
	return s
}
