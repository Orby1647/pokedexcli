package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/orby1647/pokedexcli/internal/pokecache"
)

func TestListLocationAreas_UsesBaseURL(t *testing.T) {
	// Fake API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := LocationAreasResponse{
			Next:     "next-url",
			Previous: "",
			Results: []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				{Name: "test-area", URL: "url"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Override the base URL for testing
	// NOTE: This works because tests are in the same package (pokeapi)
	originalBase := baseLocationAreaURL
	baseLocationAreaURL = server.URL
	defer func() { baseLocationAreaURL = originalBase }()

	client := Client{
		httpClient: http.Client{Timeout: 2 * time.Second},
		cache:      pokecache.NewCache(1 * time.Minute),
	}

	resp, err := client.ListLocationAreas(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Results) != 1 || resp.Results[0].Name != "test-area" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestListLocationAreas_UsesCache(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := LocationAreasResponse{
			Next:     "",
			Previous: "",
			Results: []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				{Name: "test-area", URL: "url"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	originalBase := baseLocationAreaURL
	baseLocationAreaURL = server.URL
	defer func() { baseLocationAreaURL = originalBase }()

	client := Client{
		httpClient: http.Client{Timeout: 2 * time.Second},
		cache:      pokecache.NewCache(1 * time.Minute),
	}

	// First call → hits server
	_, err := client.ListLocationAreas(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call → should hit cache, NOT server
	_, err = client.ListLocationAreas(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected server to be called once, got %d", callCount)
	}
}
