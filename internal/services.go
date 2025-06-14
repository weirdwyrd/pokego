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

func (s *PokeAPIService) GetLocationArea(locationArea string) (LocationArea, error) {
	url := fmt.Sprintf("%s/location-area/%s", s.baseURL, locationArea)
	res, err := s.client.Get(url)
	if err != nil {
		fmt.Println("error", err)
		return LocationArea{}, fmt.Errorf("failed to get location area: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return LocationArea{}, fmt.Errorf("location area not found: %s", locationArea)
	}

	if res.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var decodedResponse LocationArea
	if err := json.NewDecoder(res.Body).Decode(&decodedResponse); err != nil {
		return LocationArea{}, fmt.Errorf("failed to decode location area: %w", err)
	}

	return decodedResponse, nil
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

	var decodedResponse LocationAreasAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&decodedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode location areas: %w", err)
	}

	locationAreas := decodedResponse.Results
	fmt.Println("Loaded", len(locationAreas), "location areas")
	return locationAreas, nil
}
