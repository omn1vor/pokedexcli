package main

import (
	"bufio"
	"fmt"
	"os"
)

type command struct {
	name        string
	description string
	handler     func() error
}

var commands map[string]command

func main() {
	commands = initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	showPrompt()

	for scanner.Scan() {
		token := scanner.Text()
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
