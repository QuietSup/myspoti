package player

// Next skips to the next track.
func Next() error {
	return post("/me/player/next")
}

// Previous skips to the previous track.
func Previous() error {
	return post("/me/player/previous")
}
