package main

import "fmt"

// function for inspecting a pokemon, prints out the info
func commandInspect(cfg *Config, input []string, pokedex *pokedex) error {
	pokemon, ok := pokedex.pokemon[input[1]]
	if !ok {
		fmt.Println("Pokemon does not exist or has not been caught yet")
		return nil
	}
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for stat := range pokemon.Stats {
		fmt.Printf("  -%v: %v\n", pokemon.Stats[stat].Stat.Name, pokemon.Stats[stat].BaseStat)
	}
	fmt.Printf("Types:\n")
	for pokeType := range pokemon.Types {
		fmt.Printf("  - %v\n", pokemon.Types[pokeType].Type.Name)
	}
	return nil
}
