package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
)

type ExploreLocationResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c *Client) ExploreLocation(area string) (ExploreLocationResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area

	// CACHE CHECK
	if data, ok := c.cache.Get(url); ok {
		fmt.Println("cache hit:", url)
		var resp ExploreLocationResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return ExploreLocationResponse{}, err
		}
		return resp, nil
	}

	// HTTP REQUEST
	res, err := c.httpClient.Get(url)
	if err != nil {
		return ExploreLocationResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ExploreLocationResponse{}, err
	}

	// STORE IN CACHE
	c.cache.Add(url, body)

	var resp ExploreLocationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ExploreLocationResponse{}, err
	}

	return resp, nil
}
