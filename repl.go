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
	handler     func()
}

var commands map[string]command
var ctx *navigator

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
		"map": {
			name:        "map",
			description: "gets next 20 Pokemon locations",
			handler:     nextMaps,
		},
		"mapb": {
			name:        "mapb",
			description: "gets previous 20 Pokemon locations",
			handler:     prevMaps,
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

func displayHelpMessage() {
	fmt.Println()
	fmt.Println("Pokedex CLI usage:")
	for _, v := range commands {
		fmt.Println(v.name, "-", v.description)
	}
	fmt.Println()
}

func exit() {
	os.Exit(0)
}

func nextMaps() {
	url := baseUrl
	if ctx != nil {
		if ctx.next == nil {
			fmt.Println("You have reached the end of locations list")
			return
		}
		url = *ctx.next
	}
	ctx = getMaps(url)

	for _, m := range ctx.maps {
		fmt.Println(m)
	}
}

func prevMaps() {
	if ctx == nil || ctx.prev == nil {
		fmt.Println("You have reached the end of locations list")
		return
	}
	ctx = getMaps(*ctx.prev)

	for _, m := range ctx.maps {
		fmt.Println(m)
	}
}
