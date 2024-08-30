package pokedexApiHandler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Locations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationParseApi(url string) Locations {
	res, err := http.Get(url)

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	locations := Locations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatalf("Failed to unmarshal json %s\n", err)
	}
	return locations
}
