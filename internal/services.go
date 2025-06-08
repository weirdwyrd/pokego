package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PokeAPIService struct {
	baseURL string
	client  *http.Client
}

func NewPokeAPIService() *PokeAPIService {
	return &PokeAPIService{
		baseURL: "https://pokeapi.co/api/v2",
		client:  &http.Client{},
	}
}

func (s *PokeAPIService) GetLocationAreas(pageIndex int) ([]LocationArea, error) {
	offset := pageIndex * 20
	url := fmt.Sprintf("%s/location-area?offset=%d", s.baseURL, offset)
	res, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get location areas: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var decodedResponse LocationAreaAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&decodedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode location areas: %w", err)
	}

	locationAreas := decodedResponse.Results
	fmt.Println("Loaded", len(locationAreas), "location areas")
	return locationAreas, nil
}
