package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type command struct {
	name        string
	description string
	handler     func() error
}

var commands map[string]command

func startRepl() {
	commands = initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	showPrompt()

	for scanner.Scan() {
		token := cleanInput(scanner.Text())
		cmd, ok := commands[token]

		if !ok {
			fmt.Println("Unknown command")
			showPrompt()
			continue
		}

		cmd.handler()
		showPrompt()
	}
}

func initCommands() map[string]command {
	return map[string]command{
		"help": {
			name:        "help",
			description: "displays help message",
			handler:     displayHelpMessage,
		},
		"exit": {
			name:        "exit",
			description: "exits the Pokedex",
			handler:     exit,
		},
	}

}

func showPrompt() {
	fmt.Print("pokedex > ")
}

func cleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func displayHelpMessage() error {
	fmt.Println()
	fmt.Println("Pokedex CLI usage:")
	for _, v := range commands {
		fmt.Println(v.name, "-", v.description)
	}
	fmt.Println()
	return nil
}

func exit() error {
	os.Exit(0)
	return nil
}
