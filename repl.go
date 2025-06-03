package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		output := commands(command[0], &config)
		if output != nil {
			err := output()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
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
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *config) error {
	res, err := http.Get(config.Next)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %v", res.StatusCode)
	}
	if err != nil {
		return err
	}
	type Location struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	var locations Location
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s", err)
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	return nil
}

func commandMapb(config *config) error {
	if config.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	res, err := http.Get(config.Previous)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %v", res.StatusCode)
	}
	if err != nil {
		return err
	}
	type Location struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	var locations Location
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s", err)
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	return nil
}
