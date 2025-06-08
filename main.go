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
		LoadedData:     internal.DataLoad{},
		PageLength:     20,
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
		// if command.Name != "undo" {
		//	undo isn't working quite right, TODO fix
		// }

		cliState.PrevCommand = cliState.CurrentCommand
		cliState.CurrentCommand = command
		cliState.PrevPage = cliState.CurrentPage

		// if the last command was same as the current command,
		if cliState.CurrentCommand.Name == cliState.PrevCommand.Name {
			//increment the page number, to we can fetch the next page of data
			cliState.CurrentPage++
		} else {
			// otherwise reset the page number
			cliState.CurrentPage = 0
		}

		fmt.Println("prev page", cliState.PrevPage, "current page", cliState.CurrentPage)
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
	// need to handle page bounds
	pageStartIndex := cliState.CurrentPage * cliState.PageLength
	pageEndIndex := pageStartIndex + cliState.PageLength

	// if we have not loaded the data before, we need to fetch it
	if cliState.LoadedData.LocationAreaData == nil || len(cliState.LoadedData.LocationAreaData) <= pageStartIndex {
		pokeService := internal.NewPokeAPIService()
		locationAreaPage, err := pokeService.GetLocationAreas(cliState.CurrentPage)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}
		// need to check to make sure array exists
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
	for _, locationArea := range cliState.LoadedData.LocationAreaData[pageStartIndex:pageEndIndex] {
		fmt.Println(locationArea.Name)
	}

	return nil
}

// this doesn't really work either, because we are just looking back one command,
// instead we need to either store a full history and then use this to jump all the way back to the last map
// kinda like a undo-jump.
// alternatively we can keep a MapPage somewhere and then just use map commands to move those pages, and not
// interact with the full undo history
func commandMapBack(cliState *internal.CliState) error {
	fmt.Println(cliState.PrevCommand.Name, cliState.PrevPage)
	if cliState.PrevCommand.Name == "map" && cliState.PrevPage > 0 {
		cliState.CurrentPage = cliState.PrevPage - 1
		commandMap(cliState)
	} else {
		fmt.Println("No previous page to go back to")
	}
	return nil
}

// undo does not really work because we only go back one page, and the page handling is messed up
// if we want to fix this, we should create an array or a tree like structure to store history of
// command inputs and outputs, that way rather than re-run the command we can just replicate the output?
// if we don't store the outputs at the least we store inputs
func commandUndo(cliState *internal.CliState) error {
	fmt.Println("prev page", cliState.PrevPage, "current page", cliState.CurrentPage)
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
