package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/weirdwyrd/pokego/internal"
)

var commands = map[string]internal.CliCommand{
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
}

func main() {
	data, err := loadData()
	if err != nil {
		fmt.Println("Error fetching data, scanner will start regardless:", err)
	}
	startScanner(data)
}

func loadData() (internal.DataLoad, error) {
	pokeService := internal.NewPokeAPIService()
	locationAreas, err := pokeService.GetLocationAreas()
	if err != nil {
		return internal.DataLoad{}, fmt.Errorf("failed to load data: %w", err)
	}

	return internal.DataLoad{LocationAreaData: locationAreas}, nil
}

func startScanner(data internal.DataLoad) {
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

		err := command.Callback()
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
