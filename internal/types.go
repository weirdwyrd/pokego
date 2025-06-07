package internal

// main types

type CliState struct {
	CurrentCommand CliCommand
	CurrentPage    int

	PrevCommand CliCommand
	PrevPage    int

	AvailableCommands map[string]CliCommand
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*CliState) error
}

type DataLoad struct {
	LocationAreaData []LocationArea `json:"location_area_data"`
}

// data types

type LocationArea struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GameIndex int    `json:"game_index"`
	// EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	// Location             NamedAPIResource      `json:"location"`
	// Names                []Name               `json:"names"`
	// PokemonEncounters    []PokemonEncounter   `json:"pokemon_encounters"`
}

// api service types

type LocationAreaAPIResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}
