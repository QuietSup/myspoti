package main

import (
	"fmt"
	"os"

	"myspoti/internal/auth"
)

func authCmd() {
	if err := auth.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "auth failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Authenticated. Tokens saved to ~/.config/myspoti/tokens.json")
}
