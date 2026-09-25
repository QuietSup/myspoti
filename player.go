package main

import (
	"fmt"
	"os"
	"time"

	"myspoti/internal/player"
	"myspoti/internal/view"
)

func nextCmd() {
	if err := player.Next(); err != nil {
		fmt.Fprintf(os.Stderr, "next failed: %v\n", err)
		os.Exit(1)
	}
	printNowAfterSkip()
}

func prevCmd() {
	if err := player.Previous(); err != nil {
		fmt.Fprintf(os.Stderr, "prev failed: %v\n", err)
		os.Exit(1)
	}
	printNowAfterSkip()
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
	printNow()
}

func queueCmd() {
	q, err := player.Queue()
	if err != nil {
		fmt.Fprintf(os.Stderr, "queue failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(view.Queue(q))
}

func likeCmd() {
	playback, err := player.Like()
	if err != nil {
		fmt.Fprintf(os.Stderr, "like failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("‹𝟹 Liked")
	fmt.Println(view.Playback(*playback))
}

func unlikeCmd() {
	playback, err := player.Unlike()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unlike failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("‹/𝟹 Unliked")
	fmt.Println(view.Playback(*playback))
}

func printNowAfterSkip() {
	// Spotify's currently-playing endpoint can lag briefly after a skip.
	time.Sleep(300 * time.Millisecond)
	printNow()
}

func printNow() {
	playback, err := player.Now()
	if err != nil {
		fmt.Fprintf(os.Stderr, "now failed: %v\n", err)
		os.Exit(1)
	}
	if playback == nil {
		fmt.Println("Nothing is playing")
		return
	}
	fmt.Println(view.Playback(*playback))
}

func helpCmd() {
	fmt.Println(`myspoti — control Spotify from the terminal

Commands:
  auth    Log in with Spotify (PKCE)
  now     Show the currently playing track
  queue   Show recent, current, and up next
  play    Resume playback
  pause   Pause playback
  next    Skip to next track
  prev    Skip to previous track
  like    Save current track to Your Library
  unlike  Remove current track from Your Library
  help    Show this help`)
}
