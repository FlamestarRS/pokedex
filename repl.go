package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FlamestarRS/pokedex/internal/pokecache"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{
		cache:    *pokecache.NewCache(5 * time.Second),
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Name:     "",
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
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	cache    pokecache.Cache
	Next     string
	Previous string
	Name     string
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

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(config *config) error {
	if val, ok := config.cache.Get(config.Next); ok {
		var locations Location
		err := json.Unmarshal(val, &locations)
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

	var locations Location
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s", err)
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	config.cache.Add(config.Next, body)
	config.Next = locations.Next
	config.Previous = locations.Previous

	return nil
}

func commandMapb(config *config) error {
	if config.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	if val, ok := config.cache.Get(config.Previous); ok {
		var locations Location
		err := json.Unmarshal(val, &locations)
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

	var locations Location
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s", err)
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	config.cache.Add(config.Previous, body)
	config.Next = locations.Next
	config.Previous = locations.Previous

	return nil
}

type Explore struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExplore(config *config) error {
	if config.Name == "" {
		fmt.Println("Please enter the name of an area")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + config.Name

	if val, ok := config.cache.Get(config.Name); ok {
		var explore Explore
		err := json.Unmarshal(val, &explore)
		if err != nil {
			return fmt.Errorf("error unmarshaling JSON: %s", err)
		}
		for _, loc := range explore.PokemonEncounters {
			fmt.Println(loc.Pokemon.Name)
		}
		return nil
	}

	res, err := http.Get(url)
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

	var explore Explore
	err = json.Unmarshal(body, &explore)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s", err)
	}
	for _, loc := range explore.PokemonEncounters {
		fmt.Println(loc.Pokemon.Name)
	}

	return nil
}
