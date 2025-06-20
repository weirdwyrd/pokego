package internal

// Cache defines the interface for caching operations
type Cache interface {
	Get(key string) ([]byte, bool)
	Add(key string, val []byte)
}

// main types

type CliState struct {
	CurrentCommand CliCommand
	CurrentPage    int
	CommandHistory []CliEvent
	// LoadedData        DataLoad
	Cache             Cache
	PageLength        int
	AvailableCommands map[string]CliCommand
}

type CliEvent struct {
	Command     CliCommand
	CommandArgs []string
	Page        int
	Output      string
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*CliState, []string) (string, error) // accepts current state and command arguments
}

// data types

type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             Location              `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod EncounterMethod `json:"encounter_method"`
	VersionDetails  []VersionDetail `json:"version_details"`
}

type VersionDetail struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Name struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type Language struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Official bool   `json:"official"`
	ISO639   string `json:"iso639"`
	ISO3166  string `json:"iso3166"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon                  `json:"pokemon"`
	VersionDetails []VersionEncounterDetail `json:"version_details"`
}

type Pokemon struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionEncounterDetail struct {
	Version          Version     `json:"version"`
	MaxChance        int         `json:"max_chance"`
	EncounterDetails []Encounter `json:"encounter_details"`
}

type Version struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Encounter struct {
	MinLevel        int                       `json:"min_level"`
	MaxLevel        int                       `json:"max_level"`
	ConditionValues []EncounterConditionValue `json:"condition_values"`
	Chance          int                       `json:"chance"`
	Method          EncounterMethod           `json:"method"`
}

type EncounterMethod struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Order int    `json:"order"`
}

type EncounterConditionValue struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	URL       string             `json:"url"`
	Condition EncounterCondition `json:"condition"`
}

type EncounterCondition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// api service types

type LocationAreasAPIResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}
