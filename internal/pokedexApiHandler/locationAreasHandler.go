package pokedexApiHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationAreas struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetPokemonForLocationApi(area string) []string {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area)
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to get location")
	}
	if res.StatusCode > 299 {
		log.Fatalf("Request failed with status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf("Failed to parse response body")
	}

	LocationAreas := LocationAreas{}
	err = json.Unmarshal(body, &LocationAreas)
	if err != nil {
		log.Fatalf("Unable to marshal response %s into LocationAreas struct", string(body))
	}

	listOfPokemon := []string{}
	for _, obj := range LocationAreas.PokemonEncounters {
		listOfPokemon = append(listOfPokemon, obj.Pokemon.Name)
	}
	return listOfPokemon
}
