package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DXICIDE/pokedex/internal/pokecache"
)

var commands = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *Config, input []string, pokedex *pokedex) error
}

type Config struct {
	Next     *string
	Previous *string
	cache    pokecache.Cache
}

func main() {
	// init the commands
	commandList()
	config := new(Config)
	config.cache = *pokecache.NewCache(30 * time.Second)
	pokedex := new(pokedex)
	pokedex.pokemon = make(map[string]PokemonEndPoint)
	//create new scanner for expected input
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) < 1 {
			err := errors.New("no input")
			fmt.Println(err)
			continue
		}

		found := false
		for command := range commands {
			if input[0] == command {
				err := commands[command].callback(config, input, pokedex)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				found = true
			}
		}
		if !found {
			fmt.Println("Unknown command")
		}
	}
}

// function for cleaning input
func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	clean := strings.Fields(lowered)
	return clean
}

// collection of commands in the program
func commandList() map[string]cliCommand {
	commands = map[string]cliCommand{
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
			description: "Displays the names of 20 locations and its subsequent use prints the next 20",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore []",
			description: "Lists all of the Pokemon in a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch []",
			description: "Catches a chosen pokemon and puts him into the pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect []",
			description: "Inspects a caught pokemon and prints its stats",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex []",
			description: "Shows caught pokemon",
			callback:    commandPokedex,
		},
	}
	return commands
}
