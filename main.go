package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	pokecache "github.com/KrupaH/golang-pokedex/internal/pokecache"
	pokedexApiHandler "github.com/KrupaH/golang-pokedex/internal/pokedexApiHandler"
)

type cliCommand struct {
	Command  string
	Help     string
	Callback func(*config, *pokecache.Cache, ...string) error
}

type config struct {
	Next     *string
	Previous *string
	Pokedex  map[string]pokedexApiHandler.Pokemon
}

func cleanInput(text string) []string {
	// Clean input text
	// The purpose of this function will be to split the users input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace. For example:

	// hello world -> ["hello", "world"]
	// Charmander Bulbasaur PIKACHU -> ["charmander", "bulbasaur", "pikachu"]
	splitText := strings.Split(text, " ")
	var output []string
	for _, item := range splitText {
		if item != "" {
			output = append(output, strings.ToLower(item))
		}
	}
	return output
}

func help(*config, *pokecache.Cache, ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func exit(*config, *pokecache.Cache, ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func getLocations(params *config, cache *pokecache.Cache, args ...string) error {
	if params.Next == nil {
		fmt.Println("No more locations")
		return nil
	}

	// Check cache
	cacheEntry, ok := cache.Get(*params.Next)
	var locations pokedexApiHandler.Locations
	if ok {
		err := json.Unmarshal(cacheEntry, &locations)
		if err != nil {
			panic("Cache value invalid")
		}
	} else {
		locations = pokedexApiHandler.GetLocationParseApi(*params.Next)
		jsonStrLocations, err := json.Marshal(locations)
		if err != nil {
			panic("response format not jsonable")
		}
		cache.Add(*params.Next, jsonStrLocations)
	}
	params.Next = locations.Next
	params.Previous = locations.Previous

	for _, row := range locations.Results {
		fmt.Println(row.Name)
	}

	return nil
}

func getLocationsBack(params *config, cache *pokecache.Cache, args ...string) error {
	if params.Previous == nil {
		fmt.Println("No previous locations")
		return nil
	}

	// Check cache
	cacheEntry, ok := cache.Get(*params.Previous)
	var locations pokedexApiHandler.Locations
	if ok {
		err := json.Unmarshal(cacheEntry, &locations)
		if err != nil {
			panic("Unable to unmarshal cache entry")
		}
	} else {
		locations = pokedexApiHandler.GetLocationParseApi(*params.Previous)
		jsonLocations, err := json.Marshal(locations)

		if err != nil {
			panic("unable to marshal response to json")
		}
		cache.Add(*params.Previous, jsonLocations)

	}

	params.Next = locations.Next
	params.Previous = locations.Previous

	for _, row := range locations.Results {
		fmt.Println(row.Name)
	}

	return nil
}

func exploreArea(params *config, cache *pokecache.Cache, args ...string) error {
	// Try cache first
	cacheEntry, ok := cache.Get(args[0])

	var pokemon []string
	if !ok {
		pokemon = pokedexApiHandler.GetPokemonForLocationApi(args[0])
	} else {
		pokemonBytes, _ := json.Marshal(cacheEntry)
		for _, pokemonByte := range pokemonBytes {
			pokemon = append(pokemon, string(pokemonByte))
		}
	}

	for _, pokemon_name := range pokemon {
		fmt.Println(pokemon_name)
	}
	return nil
}

func catchPokemon(config *config, cache *pokecache.Cache, args ...string) error {
	pokemonName := args[0]
	pokemon, isCaught := pokedexApiHandler.CatchPokemon(pokemonName)
	if isCaught {
		config.Pokedex[pokemonName] = pokemon
		fmt.Println("Caught, inspect with pokedex")
	} else {
		fmt.Printf("Did not catch %v\n", pokemonName)
	}
	return nil
}

func pokedex(config *config, cache *pokecache.Cache, args ...string) error {
	for key, _ := range config.Pokedex {
		fmt.Printf(" - %s\n", key)
	}
	return nil
}

func inspectPokemon(config *config, cache *pokecache.Cache, args ...string) error {
	pokemonName := args[0]
	pokemon := pokedexApiHandler.GetPokemonApi(pokemonName)
	var stats, types string
	for _, stat := range pokemon.Stats {
		stats = fmt.Sprintf("%s  - %s: %v\n", stats, stat.Stat.Name, stat.BaseStat)
	}

	for _, pokemonType := range pokemon.Types {
		types = fmt.Sprintf("%s - %s\n", types, pokemonType.Type.Name)
	}
	str := fmt.Sprintf("- Name: %s\n- Height: %d\n- Weight: %d\n- Stats: \n%s\n- Types: \n%s\n",
		pokemon.Name, pokemon.Height, pokemon.Weight, stats, types)
	fmt.Print(str)
	return nil
}

func main() {
	commands := map[string]cliCommand{
		"help":    cliCommand{Command: "help", Help: "Shows help text", Callback: help},
		"exit":    cliCommand{Command: "exit", Help: "Graceful shutdown of Pokedex", Callback: exit},
		"map":     cliCommand{Command: "map", Help: "Display 20 locations", Callback: getLocations},
		"mapb":    cliCommand{Command: "mapb", Help: "Display previous 20 locations", Callback: getLocationsBack},
		"explore": cliCommand{Command: "explore", Help: "Get list of pokemon in area", Callback: exploreArea},
		"catch":   cliCommand{Command: "catch", Help: "try to catch a pokemon", Callback: catchPokemon},
		"pokedex": cliCommand{Command: "pokedex", Help: "show all pokemon", Callback: pokedex},
		"inspect": cliCommand{Command: "inspect", Help: "show stats of a pokemon", Callback: inspectPokemon},
	}

	initialUrl := "https://pokeapi.co/api/v2/location/"
	mapConfig := config{Next: &initialUrl, Pokedex: make(map[string]pokedexApiHandler.Pokemon)}
	scanner := bufio.NewScanner(os.Stdin)

	cache := pokecache.NewCache(3 * time.Second)

	for true {
		fmt.Print("pokedex >")
		scanner.Scan()
		input := scanner.Text()
		cmd, ok := commands[strings.Fields(input)[0]]
		if !ok {
			fmt.Printf("Command %s not recognized\n", input)
			continue
		}

		switch cmd.Command {
		case "map":
			_ = cmd.Callback(&mapConfig, cache)
		case "mapb":
			_ = cmd.Callback(&mapConfig, cache)
		case "explore":
			_ = cmd.Callback(&mapConfig, cache, strings.Fields(input)[1])
		case "catch":
			_ = cmd.Callback(&mapConfig, cache, strings.Fields(input)[1])
		case "inspect":
			_ = cmd.Callback(&mapConfig, cache, strings.Fields(input)[1])
		case "pokedex":
			_ = cmd.Callback(&mapConfig, cache)
		default:
			fmt.Printf("Command:%v\n", strings.Fields(input))
			_ = cmd.Callback(nil, nil)
		}

	}
}
