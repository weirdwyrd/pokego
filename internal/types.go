package internal

type DataLoad struct {
	LocationAreaData []LocationArea `json:"location_area_data"`
}

type Config struct {
	NextURL string `json:"next"`
	PrevURL string `json:"prev"`
}

type LocationArea struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GameIndex int    `json:"game_index"`
	// EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	// Location             NamedAPIResource      `json:"location"`
	// Names                []Name               `json:"names"`
	// PokemonEncounters    []PokemonEncounter   `json:"pokemon_encounters"`
}
