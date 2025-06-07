package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/weirdwyrd/pokego/internal"
)

func main() {
	cliState := initCli()
	// data, err := loadData()
	// if err != nil {
	// 	fmt.Println("Error fetching data, scanner will start regardless:", err)
	// }
	startScanner(cliState)
}

func initCli() *internal.CliState {
	return &internal.CliState{
		CurrentCommand: internal.CliCommand{},
		CurrentPage:    0,
		PrevCommand:    internal.CliCommand{},
		PrevPage:       0,
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
			"undo": {
				Name:        "undo",
				Description: "Undoes the last command",
				Callback:    commandUndo,
			},
		},
	}
}

// func loadData() (internal.DataLoad, error) {
// 	pokeService := internal.NewPokeAPIService()
// 	locationAreas, err := pokeService.GetLocationAreas()
// 	if err != nil {
// 		return internal.DataLoad{}, fmt.Errorf("failed to load data: %w", err)
// 	}

// 	return internal.DataLoad{LocationAreaData: locationAreas}, nil
// }

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
		cliState.PrevCommand = cliState.CurrentCommand
		cliState.PrevPage = cliState.CurrentPage
		cliState.CurrentCommand = command
		cliState.CurrentPage = 0

		err := command.Callback(cliState)
		if err != nil {
			fmt.Println("Error:", err)
		}

	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func commandHelp(cliState *internal.CliState) error {

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range cliState.AvailableCommands {
		fmt.Printf("%s - %s\n", command.Name, command.Description)
	}
	return nil
}

func commandExit(cliState *internal.CliState) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cliState *internal.CliState) error {
	fmt.Println("Map")
	return nil
}

func commandMapBack(cliState *internal.CliState) error {
	fmt.Println("Map Back")
	return nil
}

func commandUndo(cliState *internal.CliState) error {
	if cliState.PrevCommand.Name != "" {
		// reset state
		cliState.CurrentCommand = cliState.PrevCommand
		cliState.CurrentPage = cliState.PrevPage
		// call the past command again
		cliState.CurrentCommand.Callback(cliState)
	} else {
		fmt.Println("No previous command to undo")
	}
	return nil
}
