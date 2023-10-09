package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type navigator struct {
	prev *string
	next *string
	maps []string
}

type locationData struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
}

type areaData struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type mapsChunk struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const baseUrl string = "https://pokeapi.co/api/v2/location?offset=0&limit=20"
const locationUrl string = "https://pokeapi.co/api/v2/location/"

func getMaps(url string) *navigator {
	data := getData(url)
	mapsChunk := unmarshal[mapsChunk](data)

	maps := make([]string, len(mapsChunk.Results))
	for i, m := range mapsChunk.Results {
		maps[i] = m.Name
	}
	return &navigator{prev: mapsChunk.Previous, next: mapsChunk.Next, maps: maps}
}

func getPokemonsFromLocation(location string) []string {
	url := locationUrl + location
	data := getData(url)

	locationData := unmarshal[locationData](data)
	pokemonNames := []string{}
	for _, area := range locationData.Areas {
		data = getData(area.URL)
		areaData := unmarshal[areaData](data)
		for _, pokemon := range areaData.PokemonEncounters {
			pokemonNames = append(pokemonNames, pokemon.Pokemon.Name)
		}
	}

	return pokemonNames
}

func getData(url string) []byte {
	data, ok := cache.Get(url)
	if !ok {
		data = getDataFromAPI(url)
		cache.Add(url, data)
	}
	return data
}

func getDataFromAPI(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("Can't get results from PokeAPI:", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("PokeAPI response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatalln("Can't get results from PokeAPI:", err)
	}

	return body
}

func unmarshal[T any](data []byte) T {
	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		log.Fatalln("Can't deserialize results from PokeAPI:", err)
	}
	return obj
}
