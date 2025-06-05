package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
