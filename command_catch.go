package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type pokedex struct {
	pokemon map[string]PokemonEndPoint
}

// function for catching pokemon
func commandCatch(cfg *Config, input []string, pokedex *pokedex) error {

	if len(input) < 2 {
		err := errors.New("no Pokemon specified")
		return err
	}

	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", input[1]))
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", input[1])

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		err = fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
		return err
	}
	if err != nil {
		return err
	}

	pokemon := PokemonEndPoint{}
	err = json.Unmarshal(body, &pokemon)

	if err != nil {
		return err
	}

	//base expirience max is 255 so theres 45 numbers above to have atlest decent chance to chatch it
	catchrate := rand.Intn(300)

	if catchrate > pokemon.BaseExperience {
		fmt.Printf("%v was caught!\n", input[1])
		fmt.Println("You may now inspect it with the inspect command.")
		pokedex.pokemon[pokemon.Name] = pokemon
		return nil
	} else {
		fmt.Printf("%v escaped!\n", input[1])
		return nil
	}
}
