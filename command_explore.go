package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// func for exploring an area to see which pokemon are there, practicaly the same as command map
func commandExplore(cfg *Config, input []string, pokedex *pokedex) error {
	var res *http.Response
	var err error
	var body []byte
	var ok bool

	if len(input) < 2 {
		err := errors.New("no location area specified")
		return err
	}

	body, ok = cfg.cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", input[1]))

	if !ok {
		res, err = http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", input[1]))
		if err != nil {
			return err
		}
	}

	fmt.Printf("Exploring %v...\n", input[1])

	if len(body) == 0 {
		body, err = io.ReadAll(res.Body)
	}

	if res != nil {
		cfg.cache.Add(res.Request.URL.String(), body)
		res.Body.Close()
		if res.StatusCode > 299 {
			err = fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
			return err
		}
		if err != nil {
			return err
		}
	}

	locationAreasNamed := LocationAreaNamed{}
	err = json.Unmarshal(body, &locationAreasNamed)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for id := range locationAreasNamed.PokemonEncounters {
		fmt.Printf(" - %v\n", locationAreasNamed.PokemonEncounters[id].Pokemon.Name)
	}
	return nil
}
