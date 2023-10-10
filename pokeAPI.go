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

type pokemonData struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name      string        `json:"name"`
	Order     int           `json:"order"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

const baseUrl string = "https://pokeapi.co/api/v2/location?offset=0&limit=20"
const locationUrl string = "https://pokeapi.co/api/v2/location/"
const pokemonUrl string = "https://pokeapi.co/api/v2/pokemon/"

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

func getPokemonData(name string) *pokemonData {
	url := pokemonUrl + name
	data := getData(url)

	pokemonData := unmarshal[pokemonData](data)
	return &pokemonData
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

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
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
