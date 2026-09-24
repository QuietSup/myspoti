package main

import (
	"fmt"
	"os"

	"myspoti/internal/player"
)

func nextCmd() {
	if err := player.Next(); err != nil {
		fmt.Fprintf(os.Stderr, "next failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Skipped to next track")
}

func prevCmd() {
	if err := player.Previous(); err != nil {
		fmt.Fprintf(os.Stderr, "prev failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Skipped to previous track")
}

func playCmd() {
	if err := player.Play(); err != nil {
		fmt.Fprintf(os.Stderr, "play failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Playing")
}

func pauseCmd() {
	if err := player.Pause(); err != nil {
		fmt.Fprintf(os.Stderr, "pause failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Paused")
}

func nowCmd() {
	playback, err := player.Now()
	if err != nil {
		fmt.Fprintf(os.Stderr, "now failed: %v\n", err)
		os.Exit(1)
	}
	if playback == nil {
		fmt.Println("Nothing is playing")
		return
	}
	fmt.Println(playback)
}

func likeCmd() {
	playback, err := player.Like()
	if err != nil {
		fmt.Fprintf(os.Stderr, "like failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Liked: %s — %s\n", playback.Track, playback.Artists)
}

func unlikeCmd() {
	playback, err := player.Unlike()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unlike failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Unliked: %s — %s\n", playback.Track, playback.Artists)
}

func helpCmd() {
	fmt.Println(`myspoti — control Spotify from the terminal

Commands:
  auth    Log in with Spotify (PKCE)
  now     Show the currently playing track
  play    Resume playback
  pause   Pause playback
  next    Skip to next track
  prev    Skip to previous track
  like    Save current track to Your Library
  unlike  Remove current track from Your Library
  help    Show this help`)
}
