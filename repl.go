package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FlamestarRS/pokedex/internal/pokecache"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{
		cache:         *pokecache.NewCache(5 * time.Second),
		Next:          "https://pokeapi.co/api/v2/location-area/",
		Previous:      "",
		Name:          "",
		CaughtPokemon: map[string]Pokemon{},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if len(command) > 1 {
			config.Name = command[1]
		}
		output := commands(command[0], &config)
		if output != nil {
			err := output()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
		config.Name = ""
	}
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	split := strings.Fields(lowered)
	return split
}

func commands(command string, config *config) func() error {
	value, exists := getCommands()[command]
	if !exists {
		fmt.Println("Unknown command")
		return nil
	}
	return func() error {
		return value.callback(config)
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays all Pokemon in an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	cache         pokecache.Cache
	Next          string
	Previous      string
	Name          string
	CaughtPokemon map[string]Pokemon
}
