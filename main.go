package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		// LoadedData:     internal.DataLoad{},
		Cache:          cache,
		PageLength:     20,
		CommandHistory: []internal.CliEvent{},
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

		command, exists := cliState.AvailableCommands[commandInput]
		if !exists {
			fmt.Println("Unknown command")
			cliState.CurrentCommand = internal.CliCommand{}
			continue
		}

		// update cli State
		cliState.CurrentCommand = command

		output, err := command.Callback(cliState)

		if err != nil {
			fmt.Println("Error:", err)
		}
		if output != "" {
			fmt.Println(output)
		}

		// if command.Name != "undo" {
		// 	// Record the event in history
		// 	cliState.CommandHistory = append(cliState.CommandHistory, internal.CliEvent{
		// 		Command: command,
		// 		Page:    cliState.CurrentPage,
		// 		Output:  output,
		// 	})
		// }
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func commandHelp(cliState *internal.CliState) (string, error) {
	output := "Welcome to the Pokedex!\n"
	output += "Usage:\n\n"
	for _, command := range cliState.AvailableCommands {
		output += fmt.Sprintf("%s - %s\n", command.Name, command.Description)
	}
	return output, nil
}

func commandExit(cliState *internal.CliState) (string, error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return "", nil
}

func commandMap(cliState *internal.CliState) (string, error) {
	//increment page
	output, err := printPage(cliState)
	if err != nil {
		return "", err
	}
	cliState.CurrentPage++
	return output, nil
}

// if we want this to work with an undo system, we might need to do a jump back
func commandMapBack(cliState *internal.CliState) (string, error) {
	if cliState.CurrentPage > 1 {
		cliState.CurrentPage = cliState.CurrentPage - 2
	} else {
		return "no page to go back to", nil
	}
	output, err := printPage(cliState)
	if err != nil {
		return "", err
	}
	return output, nil
}

func printPage(cliState *internal.CliState) (string, error) {
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
