package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
)

var baseLocationAreaURL = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

type LocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResponse, error) {
	url := baseLocationAreaURL
	if pageURL != nil && *pageURL != "" {
		url = *pageURL
	}

	// CACHE CHECK
	if data, ok := c.cache.Get(url); ok {
		fmt.Println("cache hit:", url)

		var resp LocationAreasResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return LocationAreasResponse{}, err
		}

		return resp, nil
	}

	// MAKE HTTP REQUEST
	res, err := c.httpClient.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	// STORE IN CACHE
	c.cache.Add(url, body)

	// UNMARSHAL
	var resp LocationAreasResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return LocationAreasResponse{}, err
	}

	return resp, nil
}
