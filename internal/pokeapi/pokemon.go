package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
)

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	// CACHE CHECK
	if data, ok := c.cache.Get(url); ok {
		fmt.Println("cache hit:", url)
		var p Pokemon
		if err := json.Unmarshal(data, &p); err != nil {
			return Pokemon{}, err
		}
		return p, nil
	}

	// HTTP REQUEST
	res, err := c.httpClient.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	// STORE IN CACHE
	c.cache.Add(url, body)

	var p Pokemon
	if err := json.Unmarshal(body, &p); err != nil {
		return Pokemon{}, err
	}

	return p, nil
}
