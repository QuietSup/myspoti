package player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Playback is a snapshot of what's currently playing.
type Playback struct {
	IsPlaying bool
	Track     string
	Artists   string
	Album     string
	URI       string
	Progress  time.Duration
	Duration  time.Duration
}

func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	if total < 0 {
		total = 0
	}
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}

type currentlyPlayingResponse struct {
	IsPlaying  bool `json:"is_playing"`
	ProgressMS int  `json:"progress_ms"`
	Item       *struct {
		Name       string `json:"name"`
		URI        string `json:"uri"`
		DurationMS int    `json:"duration_ms"`
		Artists    []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
	} `json:"item"`
}

// Now returns the currently playing track, or nil if nothing is playing.
func Now() (*Playback, error) {
	status, body, err := do(http.MethodGet, "/me/player/currently-playing", nil)
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusNoContent:
		return nil, nil
	case http.StatusOK:
		var raw currentlyPlayingResponse
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, err
		}
		if raw.Item == nil {
			return nil, nil
		}
		artists := make([]string, 0, len(raw.Item.Artists))
		for _, a := range raw.Item.Artists {
			artists = append(artists, a.Name)
		}
		return &Playback{
			IsPlaying: raw.IsPlaying,
			Track:     raw.Item.Name,
			Artists:   strings.Join(artists, ", "),
			Album:     raw.Item.Album.Name,
			URI:       raw.Item.URI,
			Progress:  time.Duration(raw.ProgressMS) * time.Millisecond,
			Duration:  time.Duration(raw.Item.DurationMS) * time.Millisecond,
		}, nil
	default:
		return nil, apiError(status, body)
	}
}
