package main

import (
	"fmt"
)

func commandPokedex(config *config) error {
	if len(config.CaughtPokemon) <= 0 {
		fmt.Println("You haven't caught any Pokemon yet.")
		return nil
	}

	fmt.Printf("You have caught %v Pokemon\n", len(config.CaughtPokemon))
	for _, pokemon := range config.CaughtPokemon {
		fmt.Printf("- #%v %s\n", pokemon.ID, pokemon.Name)
	}

	return nil

}
