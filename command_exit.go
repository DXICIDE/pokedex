package main

import (
	"fmt"
	"os"
)

// function for exiting program
func commandExit(cfg *Config, input []string, pokedex *pokedex) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
