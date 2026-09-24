package player

// Play resumes playback on the active device.
func Play() error {
	return put("/me/player/play")
}

// Pause pauses playback on the active device.
func Pause() error {
	return put("/me/player/pause")
}
