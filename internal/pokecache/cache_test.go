package pokecache

import (
	"fmt"
	"testing"
	"time"
)

// TestAddGet tests the basic functionality of adding and retrieving items from the cache
func TestAddGet(t *testing.T) {
	// Define the cache expiration interval
	const interval = 5 * time.Second

	// Define test cases with different keys and values
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "location_areas_0", // First page of location areas
			val: []byte(`{"count": 20, "next": "https://pokeapi.co/api/v2/location-area?offset=20", "previous": null, "results": [{"name": "canalave-city-area", "url": "https://pokeapi.co/api/v2/location-area/1/"}]}`),
		},
		{
			key: "location_areas_1", // Second page of location areas
			val: []byte(`{"count": 20, "next": "https://pokeapi.co/api/v2/location-area?offset=40", "previous": "https://pokeapi.co/api/v2/location-area?offset=0", "results": [{"name": "eterna-city-area", "url": "https://pokeapi.co/api/v2/location-area/2/"}]}`),
		},
	}

	// Run each test case in a subtest
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			// Create a new cache instance for each test
			cache, err := NewCache(interval)
			if err != nil {
				t.Fatalf("Failed to create cache: %v", err)
			}

			// Add the test data to the cache
			cache.Add(c.key, c.val)

			// Try to retrieve the data
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("Cache.Get(%q) returned false, expected true", c.key)
				return
			}

			// Verify the retrieved value matches what we stored
			if string(val) != string(c.val) {
				t.Errorf("Cache.Get(%q) = %q, want %q", c.key, string(val), string(c.val))
				return
			}
		})
	}
}

// TestReapLoop tests that the cache properly removes expired entries
func TestReapLoop(t *testing.T) {
	// Define time intervals for testing
	const baseTime = 5 * time.Millisecond          // Cache entry lifetime
	const waitTime = baseTime + 5*time.Millisecond // Time to wait for expiration

	// Create a new cache with a short expiration time for testing
	cache, err := NewCache(baseTime)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Add test data to the cache
	testKey := "location_areas_0"
	testVal := []byte(`{"count": 20, "next": "https://pokeapi.co/api/v2/location-area?offset=20", "previous": null, "results": [{"name": "canalave-city-area", "url": "https://pokeapi.co/api/v2/location-area/1/"}]}`)
	cache.Add(testKey, testVal)

	// Verify the data is immediately available
	val, ok := cache.Get(testKey)
	if !ok {
		t.Errorf("Cache.Get(%q) returned false immediately after Add, expected true", testKey)
		return
	}
	if string(val) != string(testVal) {
		t.Errorf("Cache.Get(%q) = %q, want %q", testKey, string(val), string(testVal))
		return
	}

	// Wait for the cache entry to expire
	time.Sleep(waitTime)

	// Verify the data has been removed
	_, ok = cache.Get(testKey)
	if ok {
		t.Errorf("Cache.Get(%q) returned true after expiration, expected false", testKey)
		return
	}
}
