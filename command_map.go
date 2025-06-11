package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// function for command Map
func commandMap(cfg *Config, input []string, pokedex *pokedex) error {
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
		return err
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
			err = fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
			return err
		}
		if err != nil {
			return err
		}
	}

	//unmashaling json and inserting into structures
	locationAreasApi := LocationAreasApi{}
	err = json.Unmarshal(body, &locationAreasApi)
	if err != nil {
		return err
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
func commandMapb(cfg *Config, input []string, pokedex *pokedex) error {
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
		return err
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
			err = fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
			return err
		}
		if err != nil {
			return err
		}
	}

	//unmashaling json and inserting into structures
	locationAreasApi := LocationAreasApi{}
	err = json.Unmarshal(body, &locationAreasApi)
	if err != nil {
		return err
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
