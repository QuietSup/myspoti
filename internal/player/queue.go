package player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	recentLimit   = 5
	upcomingLimit = 5
)

// TrackInfo is a compact track summary for lists.
type TrackInfo struct {
	Name    string
	Artists string
	Album   string
	URI     string
}

func (t TrackInfo) Label() string {
	if t.Artists == "" {
		return t.Name
	}
	return t.Name + " — " + t.Artists
}

// QueueView is recent history, current track, and upcoming queue.
type QueueView struct {
	Recent   []TrackInfo
	Current  *TrackInfo
	Upcoming []TrackInfo
}

type spotifyTrack struct {
	Name    string `json:"name"`
	URI     string `json:"uri"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
}

func trackInfoFrom(t spotifyTrack) TrackInfo {
	artists := make([]string, 0, len(t.Artists))
	for _, a := range t.Artists {
		artists = append(artists, a.Name)
	}
	return TrackInfo{
		Name:    t.Name,
		Artists: strings.Join(artists, ", "),
		Album:   t.Album.Name,
		URI:     t.URI,
	}
}

type queueResponse struct {
	CurrentlyPlaying *spotifyTrack  `json:"currently_playing"`
	Queue            []spotifyTrack `json:"queue"`
}

type recentlyPlayedResponse struct {
	Items []struct {
		Track spotifyTrack `json:"track"`
	} `json:"items"`
}

// Queue returns recent plays, the current track, and upcoming queue items.
func Queue() (QueueView, error) {
	current, upcoming, err := fetchPlayerQueue()
	if err != nil {
		return QueueView{}, err
	}

	recent, err := fetchRecentlyPlayed(recentLimit)
	if err != nil {
		return QueueView{}, err
	}

	currentURI := ""
	if current != nil {
		currentURI = current.URI
	}

	// API returns newest-first; drop the current track and reverse so older → newer.
	filtered := make([]TrackInfo, 0, len(recent))
	for _, t := range recent {
		if currentURI != "" && t.URI == currentURI {
			continue
		}
		filtered = append(filtered, t)
	}
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	if len(upcoming) > upcomingLimit {
		upcoming = upcoming[:upcomingLimit]
	}

	return QueueView{
		Recent:   filtered,
		Current:  current,
		Upcoming: upcoming,
	}, nil
}

func fetchPlayerQueue() (*TrackInfo, []TrackInfo, error) {
	status, body, err := do(http.MethodGet, "/me/player/queue", nil)
	if err != nil {
		return nil, nil, err
	}

	switch status {
	case http.StatusNoContent:
		return nil, nil, nil
	case http.StatusOK:
		var raw queueResponse
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, nil, err
		}
		var current *TrackInfo
		if raw.CurrentlyPlaying != nil && raw.CurrentlyPlaying.URI != "" {
			t := trackInfoFrom(*raw.CurrentlyPlaying)
			current = &t
		}
		upcoming := make([]TrackInfo, 0, len(raw.Queue))
		for _, item := range raw.Queue {
			if item.URI == "" {
				continue
			}
			upcoming = append(upcoming, trackInfoFrom(item))
		}
		return current, upcoming, nil
	default:
		return nil, nil, apiError(status, body)
	}
}

func fetchRecentlyPlayed(limit int) ([]TrackInfo, error) {
	path := fmt.Sprintf("/me/player/recently-played?limit=%d", limit)
	status, body, err := do(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusNoContent, http.StatusOK:
		if len(body) == 0 {
			return nil, nil
		}
		var raw recentlyPlayedResponse
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, err
		}
		out := make([]TrackInfo, 0, len(raw.Items))
		for _, item := range raw.Items {
			if item.Track.URI == "" {
				continue
			}
			out = append(out, trackInfoFrom(item.Track))
		}
		return out, nil
	default:
		return nil, apiError(status, body)
	}
}
