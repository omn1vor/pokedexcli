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

type result struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const baseUrl string = "https://pokeapi.co/api/v2/location?limit=20"

func getMaps(url string) *navigator {
	mapData, ok := cache.Get(url)
	if !ok {
		mapData = getMapDataFromAPI(url)
	}

	result := result{}
	if err := json.Unmarshal(mapData, &result); err != nil {
		log.Fatalln("Can't deserialize results from PokeAPI:", err)
	}

	maps := make([]string, len(result.Results))
	for i, m := range result.Results {
		maps[i] = m.Name
	}
	return &navigator{prev: result.Previous, next: result.Next, maps: maps}
}

func getMapDataFromAPI(url string) []byte {
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
