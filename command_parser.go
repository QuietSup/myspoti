package main

import (
	"fmt"
	"os"
)

func parseCommand() func() {
	commands := map[string]func(){
		"help": helpCmd,
		"next": nextCmd,
		"prev": prevCmd,
		"now":  nowCmd,
		"auth": authCmd,
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: command_parser <command>")
		os.Exit(1)
	}

	command := os.Args[1]
	fmt.Println("Command:", command)

	action := commands[command]
	if action == nil {
		fmt.Println("Invalid command:", command)
		os.Exit(1)
	}

	return action
}
