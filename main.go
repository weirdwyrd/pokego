package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/weirdwyrd/pokego/internal"
	"github.com/weirdwyrd/pokego/internal/pokecache"
)

func main() {
	cliState := initCli()
	startScanner(cliState)
}

func initCli() *internal.CliState {
	cache, err := pokecache.NewCache(5 * time.Second)
	if err != nil {
		fmt.Println("Error creating cache:", err)
		os.Exit(1)
		// fmt.Println("Continuing without cache...")
		// cache = nil
		// todo prompt user to continue without cache
	}

	return &internal.CliState{
		CurrentCommand: internal.CliCommand{},
		CurrentPage:    0,
		Cache:          cache,
		PageLength:     20,
		CommandHistory: []internal.CliEvent{},
		Pokedex:        make(map[string]internal.Pokemon),
		AvailableCommands: map[string]internal.CliCommand{
			"help": {
				Name:        "help",
				Description: "Shows the help message",
				Callback:    commandHelp,
			},
			"exit": {
				Name:        "exit",
				Description: "Exits the program",
				Callback:    commandExit,
			},
			"map": {
				Name:        "map",
				Description: "Shows the map of location-areas one page at a time. map again to go forward. mapb to go backward.",
				Callback:    commandMap,
			},
			"mapb": {
				Name:        "mapb",
				Description: "Goes back one page of the map",
				Callback:    commandMapBack,
			},
			"explore": {
				Name:        "explore",
				Description: "Explore the map",
				Callback:    commandExplore,
			},
			"catch": {
				Name:        "catch",
				Description: "Attempts to catch a Pokemon",
				Callback:    commandCatch,
			},
			"inspect": {
				Name:        "inspect",
				Description: "Inspect a Pokemon in your Pokedex",
				Callback:    commandInspect,
			},
			"pokedex": {
				Name:        "pokedex",
				Description: "Gotta catch em all!",
				Callback:    commandPokedex,
			},
			// "undo": {
			// 	Name:        "undo",
			// 	Description: "Undoes the last command",
			// 	Callback:    commandUndo,
			// },
		},
	}
}

func startScanner(cliState *internal.CliState) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)

		if len(cleaned) == 0 {
			continue
		}

		commandInput := cleaned[0]
		// fetch command
		command, exists := cliState.AvailableCommands[commandInput]
		if !exists {
			fmt.Println("Unknown command")
			cliState.CurrentCommand = internal.CliCommand{}
			continue
		}

		// update cli State
		cliState.CurrentCommand = command

		// check for command arguments
		var commandArgs []string
		if len(cleaned) > 1 {
			commandArgs = cleaned[1:]
		}

		output, err := command.Callback(cliState, commandArgs)

		if err != nil {
			fmt.Println("Error:", err)
		}
		if output != "" {
			fmt.Println(output)
		}

		// if command.Name != "undo" {
		// Record the event in history
		cliState.CommandHistory = append(cliState.CommandHistory, internal.CliEvent{
			Command:     command,
			CommandArgs: commandArgs,
			Page:        cliState.CurrentPage,
			Output:      output,
		})
		// }
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func commandHelp(cliState *internal.CliState, commandArgs []string) (string, error) {
	output := "Welcome to the Pokedex!\n"
	output += "Usage:\n\n"
	for _, command := range cliState.AvailableCommands {
		output += fmt.Sprintf("%s - %s\n", command.Name, command.Description)
	}
	return output, nil
}

func commandExit(cliState *internal.CliState, commandArgs []string) (string, error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return "", nil
}

func commandMap(cliState *internal.CliState, commandArgs []string) (string, error) {
	//increment page
	output, err := printLocationAreasPage(cliState)
	if err != nil {
		return "", err
	}
	cliState.CurrentPage++
	return output, nil
}

// if we want this to work with an undo system, we might need to do a jump back
func commandMapBack(cliState *internal.CliState, commandArgs []string) (string, error) {
	if cliState.CurrentPage > 1 {
		cliState.CurrentPage = cliState.CurrentPage - 2
	} else {
		return "no page to go back to", nil
	}
	output, err := printLocationAreasPage(cliState)
	if err != nil {
		return "", err
	}
	return output, nil
}

func printLocationAreasPage(cliState *internal.CliState) (string, error) {
	// Check if we have the data in cache
	// if cliState.Cache != nil {
	cacheKey := fmt.Sprintf("location_areas_%d", cliState.CurrentPage)
	cachedData, exists := cliState.Cache.Get(cacheKey)

	if !exists {
		pokeService := internal.NewPokeAPIService()
		locationAreaPage, err := pokeService.GetLocationAreas(cliState.CurrentPage)
		if err != nil {
			return "", fmt.Errorf("failed to load data: %w", err)
		}

		// Convert the data to JSON for caching
		jsonData, err := json.Marshal(locationAreaPage)
		if err != nil {
			return "", fmt.Errorf("failed to marshal data for cache: %w", err)
		}

		// Store in cache
		cliState.Cache.Add(cacheKey, jsonData)
		cachedData = jsonData
	}

	// Unmarshal the cached data
	var locationAreas []internal.LocationArea
	if err := json.Unmarshal(cachedData, &locationAreas); err != nil {
		return "", fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	// Print the page
	output := ""
	// for _, locationArea := range locationAreas[pageStartIndex:pageEndIndex] { not needed with cache logic
	for _, locationArea := range locationAreas {
		output += locationArea.Name + "\n"
	}
	return output, nil
}

func commandExplore(cliState *internal.CliState, commandArgs []string) (string, error) {
	if len(commandArgs) == 0 {
		return "Please provide a location area to explore", nil
	}

	locationAreaName := commandArgs[0]
	fmt.Printf("exploring %s ...\n", locationAreaName)
	locationArea, err := getLocationArea(cliState, locationAreaName)
	if err != nil {
		return "", fmt.Errorf("explore failed, %w", err)
	}

	// Use a map operation to collect Pokemon names
	pokemonNames := make([]string, len(locationArea.PokemonEncounters))
	for i, encounter := range locationArea.PokemonEncounters {
		pokemonNames[i] = encounter.Pokemon.Name
	}

	return fmt.Sprintf("Found Pokemon:\n%s", strings.Join(pokemonNames, "\n")), nil
}

func getLocationArea(cliState *internal.CliState, locationAreaName string) (internal.LocationArea, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("location_area_%s", locationAreaName)
	cachedData, exists := cliState.Cache.Get(cacheKey)

	if !exists {
		pokeService := internal.NewPokeAPIService()
		locationAreaData, err := pokeService.GetLocationArea(locationAreaName)
		if err != nil {
			return internal.LocationArea{}, fmt.Errorf("explore failed, %w", err)
		}

		// Convert the data to JSON for caching
		jsonData, err := json.Marshal(locationAreaData)
		if err != nil {
			return internal.LocationArea{}, fmt.Errorf("failed to marshal data for cache: %w", err)
		}
		cliState.Cache.Add(cacheKey, jsonData)
		cachedData = jsonData
	}

	var locationArea internal.LocationArea
	if err := json.Unmarshal(cachedData, &locationArea); err != nil {
		return internal.LocationArea{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return locationArea, nil
}

func commandCatch(cliState *internal.CliState, commandArgs []string) (string, error) {
	pokemonName := commandArgs[0]
	if pokemonName == "" {
		return "What pokemon are you trying to catch?", nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemon, err := getPokemon(cliState, pokemonName)
	if err != nil {
		return "", err // err already formatted in getPokemon
	}

	chanceToCatch := rand.Intn(100) - min(95, (pokemon.BaseExperience/10))
	if chanceToCatch > 0 {
		cliState.Pokedex[pokemonName] = pokemon
		return fmt.Sprintf("You caught %s!\n", pokemonName), nil
	}
	return fmt.Sprintf("You missed %s!\n", pokemonName), nil
}

func getPokemon(cliState *internal.CliState, pokemonName string) (internal.Pokemon, error) {
	cacheKey := fmt.Sprintf("pokemon_%s", pokemonName)
	cachedData, exists := cliState.Cache.Get(cacheKey)

	if !exists {
		pokeService := internal.NewPokeAPIService()
		pokemonData, err := pokeService.GetPokemon(pokemonName)
		if err != nil {
			return internal.Pokemon{}, fmt.Errorf("failed to get pokemon: %w", err)
		}
		// fmt.Printf("pokemonData: %+v\n", pokemonData)
		jsonData, err := json.Marshal(pokemonData)
		if err != nil {
			return internal.Pokemon{}, fmt.Errorf("failed to marshal data for cache: %w", err)
		}
		cliState.Cache.Add(cacheKey, jsonData)
		cachedData = jsonData
	}
	var pokemon internal.Pokemon
	if err := json.Unmarshal(cachedData, &pokemon); err != nil {
		return internal.Pokemon{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return pokemon, nil
}

func commandInspect(cliState *internal.CliState, commandArgs []string) (string, error) {
	if len(commandArgs) == 0 {
		return "Please provide the name of a Pokemon to inspect", nil
	}
	pokemonName := commandArgs[0]
	pokemon, exists := cliState.Pokedex[pokemonName]
	if !exists {
		return fmt.Sprintf("You don't have %s in your Pokedex", pokemonName), nil
	}

	pokemonJson, _ := json.MarshalIndent(pokemon, "", "  ")
	return fmt.Sprintf("Pokemon %s:\n%s\n", pokemonName, pokemonJson), nil
}

func commandPokedex(cliState *internal.CliState, commandArgs []string) (string, error) {
	output := "Gotta catch em all!\n"
	for pokemonName, _ := range cliState.Pokedex {
		output += fmt.Sprintf(" - %s\n", pokemonName)
	}
	return output, nil
}

// func commandUndo(cliState *internal.CliState) (string, error) {
// 	history := cliState.CommandHistory
// 	if len(history) > 1 {
// 		// Remove the last event (current command)
// 		history = history[:len(history)-1]
// 		last := history[len(history)-1]
// 		cliState.CurrentCommand = last.Command
// 		cliState.CurrentPage = last.Page
// 		cliState.CommandHistory = history
// 		return last.Output, nil
// 	} else {
// 		return "", nil
// 	}
// }
