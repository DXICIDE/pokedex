package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DXICIDE/pokedex/internal/pokecache"
)

var commands = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *Config) error
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreasApi struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Areas    []LocationArea `json:"results"`
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
				err := commands[command].callback(config)
				if err != nil {
					fmt.Printf("%v", err)
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

// function for exiting program
func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// function for printing help into stdout
func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, name := range commands {
		fmt.Printf("%v: %v\n", name.name, name.description)
	}
	return nil
}

// function for command Map
func commandMap(cfg *Config) error {
	var res *http.Response
	var err error
	var body []byte
	var ok bool

	//if its the first one, use without offset, otherwise use the next in line
	if cfg.Next != nil {
		//try to find it in cache, if not use the url
		body, ok = cfg.cache.Get(*cfg.Next)
		if !ok {
			res, err = http.Get(*cfg.Next)
		}
	} else {
		body, ok = cfg.cache.Get("https://pokeapi.co/api/v2/location-area/?limit=20&offset=0")
		if !ok {
			res, err = http.Get("https://pokeapi.co/api/v2/location-area/?limit=20&offset=0")
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	//if it wasnt cached the body is read
	if len(body) == 0 {
		body, err = io.ReadAll(res.Body)
	}

	//completes if it wasnt cached, added to chace and error checks
	if res != nil {
		cfg.cache.Add(res.Request.URL.String(), body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	//unmashaling json and inserting into structures
	locationAreasApi := LocationAreasApi{}
	err = json.Unmarshal(body, &locationAreasApi)
	if err != nil {
		log.Fatal(err)
	}

	//sets the config
	cfg.Next = locationAreasApi.Next
	cfg.Previous = locationAreasApi.Previous

	//printing the areas
	for id := range locationAreasApi.Areas {
		fmt.Println(locationAreasApi.Areas[id].Name)
	}
	return nil
}

// almost exact same function as commandMap expept it goes backwards. if its on the at the start, it just print that the user is on the first page
func commandMapb(cfg *Config) error {
	var res *http.Response
	var err error
	var body []byte
	var ok bool

	if cfg.Previous != nil {
		//try to find it in cache, if not use the url
		body, ok = cfg.cache.Get(*cfg.Previous)
		if !ok {
			res, err = http.Get(*cfg.Previous)
		}
	} else {
		fmt.Println("youre on the first page, type map to print it")
		return nil
	}

	if err != nil {
		log.Fatal(err)
	}

	//if it wasnt cached the body is read
	if len(body) == 0 {
		body, err = io.ReadAll(res.Body)
	}

	//completes if it wasnt cached, added to chace and error checks
	if res != nil {
		cfg.cache.Add(res.Request.URL.String(), body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	//unmashaling json and inserting into structures
	locationAreasApi := LocationAreasApi{}
	err = json.Unmarshal(body, &locationAreasApi)
	if err != nil {
		log.Fatal(err)
	}

	//set the config
	cfg.Next = locationAreasApi.Next
	cfg.Previous = locationAreasApi.Previous

	//prints the areas
	for id := range locationAreasApi.Areas {
		fmt.Println(locationAreasApi.Areas[id].Name)
	}
	return nil
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
	}
	return commands
}
