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

func helpCmd() {
	fmt.Println(`myspoti — control Spotify from the terminal

Commands:
  auth   Log in with Spotify (PKCE)
  next   Skip to next track
  prev   Skip to previous track
  help   Show this help`)
}
