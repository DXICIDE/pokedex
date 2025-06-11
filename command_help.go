package main

import "fmt"

// function for printing help into stdout
func commandHelp(cfg *Config, input []string, pokedex *pokedex) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, name := range commands {
		fmt.Printf("%v: %v\n", name.name, name.description)
	}
	return nil
}
