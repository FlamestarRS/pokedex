package main

import (
	"fmt"
)

func commandInspect(config *config) error {
	if _, ok := config.CaughtPokemon[config.Name]; !ok {
		fmt.Println("You haven't caught that Pokemon yet.")
		return nil
	}

	pkmn := config.CaughtPokemon[config.Name]
	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\n", pkmn.Name, pkmn.Height, pkmn.Weight)
	fmt.Println("Stats:")
	for _, stat := range pkmn.Stats {
		fmt.Printf("- %v\n", stat.Stat.Name)
	}

	fmt.Println("Types:")
	for _, pkmntype := range pkmn.Types {
		fmt.Printf("- %v\n", pkmntype.Type.Name)
	}

	return nil

}
