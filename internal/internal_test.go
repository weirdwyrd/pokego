package internal

import (
	"testing"
)

// TestGetLocationAreas tests the PokeAPI service's ability to fetch location areas
func TestGetLocationAreas(t *testing.T) {
	service := NewPokeAPIService()

	// Test first page
	areas, err := service.GetLocationAreas(0)
	if err != nil {
		t.Fatalf("Failed to get location areas: %v", err)
	}

	if len(areas) == 0 {
		t.Error("Expected non-empty location areas list")
	}

	// Test second page
	areas2, err := service.GetLocationAreas(1)
	if err != nil {
		t.Fatalf("Failed to get second page of location areas: %v", err)
	}

	if len(areas2) == 0 {
		t.Error("Expected non-empty second page of location areas")
	}

	// Verify pages are different
	if areas[0].Name == areas2[0].Name {
		t.Error("Expected different location areas in different pages")
	}
}
