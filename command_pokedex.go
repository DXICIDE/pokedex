package main

import "fmt"

// function for printing caught pokemon
func commandPokedex(cfg *Config, input []string, pokedex *pokedex) error {
	fmt.Println("Your Pokedex:")
	for id := range pokedex.pokemon {
		fmt.Printf(" - %v\n", pokedex.pokemon[id].Name)
	}
	return nil
}
