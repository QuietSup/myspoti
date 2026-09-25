package player

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"myspoti/internal/auth"
)

const apiBase = "https://api.spotify.com/v1"

func post(path string) error {
	return emptyBody(http.MethodPost, path)
}

func put(path string) error {
	return emptyBody(http.MethodPut, path)
}

func del(path string) error {
	return emptyBody(http.MethodDelete, path)
}

func emptyBody(method, path string) error {
	status, body, err := do(method, path, nil)
	if err != nil {
		return err
	}

	switch status {
	case http.StatusNoContent, http.StatusOK:
		return nil
	default:
		return apiError(status, body)
	}
}

func do(method, path string, body io.Reader) (int, []byte, error) {
	token, err := auth.AccessToken()
	if err != nil {
		return 0, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, apiBase+path, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

func apiError(status int, body []byte) error {
	switch status {
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized — run `myspoti auth` again")
	case http.StatusForbidden:
		msg := string(body)
		if strings.Contains(msg, "Insufficient client scope") {
			return fmt.Errorf("missing Spotify permission — run `myspoti auth` to grant new scopes")
		}
		return fmt.Errorf("forbidden — Premium required, or missing permission (try `myspoti auth`)")
	case http.StatusNotFound:
		return fmt.Errorf("no active device — open Spotify on a device first")
	default:
		return fmt.Errorf("spotify API %d: %s", status, trimBody(body))
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
