package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"log"
	"strings"
)

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	Abilities              []Ability  `json:"abilities"`
	BaseExperience         int        `json:"base_experience"`
	Height                 int        `json:"height"`
	HeldItems              []HeldItem `json:"held_items"`
	LocationAreaEncounters string     `json:"location_area_encounters"`
	Moves                  []Move     `json:"moves"`
	Name                   string     `json:"name"`
	Stats                  []Stat     `json:"stats"`
	Types                  []Type     `json:"types"`
	Weight                 int        `json:"weight"`
}

type Ability struct {
	Ability  NamedAPIResource `json:"ability"`
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
}

type HeldItem struct {
	Item           NamedAPIResource `json:"item"`
	VersionDetails []VersionDetail  `json:"version_details"`
}

type VersionDetail struct {
	Rarity  int              `json:"rarity"`
	Version NamedAPIResource `json:"version"`
}

type Move struct {
	Move                NamedAPIResource     `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}

type VersionGroupDetail struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
	VersionGroup    NamedAPIResource `json:"version_group"`
}

type Stat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type Type struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

func fetchPokemonDetails(p string, errors chan<- error, results chan <- Pokemon) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", p)

	response, err := http.Get(url)
	if err != nil {
		errors <- fmt.Errorf("Error while fetching details of %s: %v", p, err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("\nError while reading response data for %s: %v", p, err)
	}

	var pokemonData Pokemon
	err = json.Unmarshal(responseData, &pokemonData)
	if err != nil {
		fmt.Printf("\nError while unmarshaling JSON for %s: %v", p, err)
	}

	results <- pokemonData
}

func printPokedexEntry(pokemon Pokemon) {
    fmt.Printf("Name: %s\n", strings.ToUpper(pokemon.Name))
    
    // Types
    types := make([]string, len(pokemon.Types))
    for i, t := range pokemon.Types {
        types[i] = strings.ToUpper(t.Type.Name)
    }
    fmt.Printf("Type: %s\n", strings.Join(types, "/"))
    
    // Abilities
    abilities := make([]string, len(pokemon.Abilities))
    for i, ability := range pokemon.Abilities {
        abilities[i] = ability.Ability.Name
    }
    fmt.Printf("Abilities: %s\n", strings.Join(abilities, ", "))
    
    // Stats
    fmt.Println("Stats:")
    for _, stat := range pokemon.Stats {
        fmt.Printf("  %s: %d\n", strings.ToUpper(stat.Stat.Name), stat.BaseStat)
    }
    
    // Other details
    fmt.Printf("Height: %.1f m\n", float64(pokemon.Height) / 10)
    fmt.Printf("Weight: %.1f kg\n", float64(pokemon.Weight) / 10)
    fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)
    
    fmt.Println("------------------------")
}


func main() {

	pokemons := []string{
		"mewtwo",
		"lugia",
		"rayquaza",
		"giratina",
		"arceus",
		"zekrom",
		"xerneas",
		"solgaleo",
		"eternatus",
	}

	var wg sync.WaitGroup

	errors := make(chan error, len(pokemons))
	results := make(chan Pokemon, len(pokemons))

	for _, pokemon := range pokemons{
		wg.Add(1)
		go func(pokemon string){
			defer wg.Done()
			fetchPokemonDetails(pokemon, errors, results)
		}(pokemon)
	}

	go func(){
		wg.Wait()
		close(results)
		close(errors)
	}()

	var responses []Pokemon

	for range pokemons{
		select{
		case result:= <-results:
			responses = append(responses, result)
		case err := <-errors:
			log.Println(err)
		}
	}

    for _, res := range responses {
        printPokedexEntry(res)
    }

}
