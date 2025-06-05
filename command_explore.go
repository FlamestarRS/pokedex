package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
