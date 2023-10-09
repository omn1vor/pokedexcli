package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/omn1vor/pokedexcli/internal/pokecache"
)

type command struct {
	name          string
	description   string
	handler       func(string)
	needsArgument bool
}

var commands map[string]command
var ctx *navigator
var cache = pokecache.NewCache(2 * time.Minute)
var collection = map[string]*pokemonData{}

func startRepl() {
	commands = initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showPrompt()

		if !scanner.Scan() {
			break
		}

		tokens := cleanInput(scanner.Text())
		if len(tokens) == 0 {
			continue
		}

		cmd, ok := commands[tokens[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		arg := ""
		if cmd.needsArgument {
			if len(tokens) < 2 {
				fmt.Println("This command needs an argument!")
				continue
			}
			arg = tokens[1]
		}
		cmd.handler(arg)
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
		"explore": {
			name: "explore",
			description: "explores the location, listing the Pokemon's names from the areas in this location. " +
				"Requires a location name as a parameter",
			handler:       explore,
			needsArgument: true,
		},
		"catch": {
			name: "catch",
			description: "tries to catch a pokemon with given name. Chance to catch a pokemon depends on its base experience " +
				"(more experience -> harder to catch)",
			handler:       catch,
			needsArgument: true,
		},
		"inspect": {
			name:          "inspect",
			description:   "shows characteristics of a given pokemon. It works only on pokemons that you have caught",
			handler:       inspect,
			needsArgument: true,
		},
	}

}

func showPrompt() {
	fmt.Print("\npokedex > ")
}

func cleanInput(input string) []string {
	trimmed := strings.TrimSpace(input)
	lowered := strings.ToLower(trimmed)
	return strings.Fields(lowered)
}

func displayHelpMessage(_ string) {
	fmt.Println()
	fmt.Println("Pokedex CLI usage:")
	for _, v := range commands {
		fmt.Println(v.name, "-", v.description)
	}
	fmt.Println()
}

func exit(_ string) {
	os.Exit(0)
}

func nextMaps(_ string) {
	url := baseUrl
	if ctx != nil {
		if ctx.next == nil {

			return
		}
		url = *ctx.next
	}
	ctx = getMaps(url)

	for _, m := range ctx.maps {
		fmt.Println(m)
	}
}

func prevMaps(_ string) {
	if ctx == nil || ctx.prev == nil {
		fmt.Println("You have reached the end of locations list")
		return
	}
	ctx = getMaps(*ctx.prev)

	for _, m := range ctx.maps {
		fmt.Println(m)
	}
}

func explore(location string) {
	pokemons := getPokemonsFromLocation(location)
	for _, p := range pokemons {
		fmt.Println(p)
	}
}

func catch(name string) {
	const skill = 35
	pokemonData := getPokemonData(name)
	check := rand.Intn(pokemonData.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	fmt.Printf("Your skill is %d, pokemon check is %d (out of %d)\n", skill, check, pokemonData.BaseExperience)
	if skill >= check {
		collection[name] = pokemonData
		fmt.Println(name, "was caught!")
	} else {
		fmt.Println(name, "escaped!")
	}
}

func inspect(name string) {
	pokemonData, ok := collection[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return
	}

	fmt.Println("Name: ", pokemonData.Name)
	fmt.Println("Height: ", pokemonData.Height)
	fmt.Println("Weight: ", pokemonData.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemonData.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemonData.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}
}
