package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/weirdwyrd/pokego/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{
	"help": {
		name:        "help",
		description: "Shows the help message",
		callback:    commandHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exits the program",
		callback:    commandExit,
	},
}

func main() {
	locationAreas, err := loadData()
	if err != nil {
		fmt.Println("Error fetching data, scanner will start regardless:", err)
	}
	startScanner()
}

func loadData() (*internal.DataLoad, error) {
	pokeService := internal.NewPokeAPIService()
	locationAreas, err := pokeService.GetLocationAreas()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return &internal.DataLoad{LocationAreaData: locationAreas}, nil
}

func startScanner() {
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

		command, exists := commands[commandInput]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback()
		if err != nil {
			fmt.Println("Error:", err)
		}

	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help - Shows this help message")
	fmt.Println("exit - Exits the program")
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
