package player

import (
	"fmt"
	"net/url"
)

// Like saves the currently playing track to the user's library.
func Like() (*Playback, error) {
	return mutateLibrary(put)
}

// Unlike removes the currently playing track from the user's library.
func Unlike() (*Playback, error) {
	return mutateLibrary(del)
}

func mutateLibrary(fn func(string) error) (*Playback, error) {
	playback, err := Now()
	if err != nil {
		return nil, err
	}
	if playback == nil {
		return nil, fmt.Errorf("nothing is playing")
	}
	if playback.URI == "" {
		return nil, fmt.Errorf("current item has no track URI")
	}

	path := "/me/library?uris=" + url.QueryEscape(playback.URI)
	if err := fn(path); err != nil {
		return nil, err
	}
	return playback, nil
}
