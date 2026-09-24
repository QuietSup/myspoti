package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"time"
)

const (
	redirectURI = "http://127.0.0.1:8080/callback"
	authScope   = "user-modify-playback-state user-read-playback-state user-read-currently-playing user-library-modify"
)

type params struct {
	ClientID      string
	RedirectURI   string
	Scope         string
	CodeChallenge string
}

func newParams(codeChallenge string) (params, error) {
	id, err := clientID()
	if err != nil {
		return params{}, err
	}
	return params{
		ClientID:      id,
		RedirectURI:   redirectURI,
		Scope:         authScope,
		CodeChallenge: codeChallenge,
	}, nil
}

func (p params) buildURL() (string, error) {
	authURL, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", p.ClientID)
	query.Set("scope", p.Scope)
	query.Set("code_challenge_method", "S256")
	query.Set("code_challenge", p.CodeChallenge)
	query.Set("redirect_uri", p.RedirectURI)

	authURL.RawQuery = query.Encode()
	return authURL.String(), nil
}

func generateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func openBrowser(targetURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd, args = "open", []string{targetURL}
	case "windows":
		cmd, args = "rundll32", []string{"url.dll,FileProtocolHandler", targetURL}
	default:
		cmd, args = "xdg-open", []string{targetURL}
	}

	return exec.CommandContext(context.Background(), cmd, args...).Start()
}

// Run starts the PKCE login flow, exchanges the code, and saves tokens.
func Run() error {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return fmt.Errorf("generating code verifier: %w", err)
	}

	codeChallenge := generateCodeChallenge(codeVerifier)
	p, err := newParams(codeChallenge)
	if err != nil {
		return err
	}

	authURL, err := p.buildURL()
	if err != nil {
		return fmt.Errorf("building auth URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	type result struct {
		code string
		err  error
	}
	resCh := make(chan result, 1)
	go func() {
		code, err := waitForAuthCode(ctx)
		resCh <- result{code, err}
	}()

	time.Sleep(100 * time.Millisecond)

	fmt.Println("Opening browser for Spotify login...")
	if err := openBrowser(authURL); err != nil {
		fmt.Println("Couldn't open browser. Visit manually:", authURL)
	} else {
		fmt.Println("If the browser didn't open, visit:", authURL)
	}

	res := <-resCh
	if res.err != nil {
		return res.err
	}

	tokens, err := exchangeCode(res.code, codeVerifier)
	if err != nil {
		return fmt.Errorf("exchanging code: %w", err)
	}

	if err := saveTokens(tokens); err != nil {
		return fmt.Errorf("saving tokens: %w", err)
	}
	return nil
}
