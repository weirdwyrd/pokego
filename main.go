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
		LoadedData:     internal.DataLoad{},
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
	println(cliState.CurrentPage)
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
	println(cliState.CurrentPage)
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
	pageStartIndex := cliState.CurrentPage * cliState.PageLength
	pageEndIndex := pageStartIndex + cliState.PageLength

	// if we have not loaded the data before, we need to fetch it
	if cliState.LoadedData.LocationAreaData == nil || len(cliState.LoadedData.LocationAreaData) <= pageStartIndex {
		pokeService := internal.NewPokeAPIService()
		locationAreaPage, err := pokeService.GetLocationAreas(cliState.CurrentPage)
		if err != nil {
			return "", fmt.Errorf("failed to load data: %w", err)
		}
		if cliState.LoadedData.LocationAreaData == nil {
			cliState.LoadedData.LocationAreaData = []internal.LocationArea{}
		}
		// add the new data page to the loaded data
		cliState.LoadedData.LocationAreaData = append(cliState.LoadedData.LocationAreaData, locationAreaPage...)
	}

	// TODO edge case what if we reach the end of the data from API? do not want to append the last page forever, need to test.
	// also need to have way to inform user that there is no more to load / display

	// we know we at least have the start of the page, we need to see if we have a full page
	if len(cliState.LoadedData.LocationAreaData) < pageEndIndex {
		pageEndIndex = len(cliState.LoadedData.LocationAreaData)
	}

	// now print the page
	output := ""
	for _, locationArea := range cliState.LoadedData.LocationAreaData[pageStartIndex:pageEndIndex] {
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
