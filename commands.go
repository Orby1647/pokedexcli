package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
}

type locationAreaResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

var commandRegistry = map[string]cliCommand{}

func init() {
	commandRegistry["exit"] = cliCommand{
		name:        "exit",
		description: "Exits the Pokedex CLI",
		callback:    commandExit,
	}
	commandRegistry["help"] = cliCommand{
		name:        "help",
		description: "Displays this help message",
		callback:    commandHelp,
	}
	commandRegistry["map"] = cliCommand{
		name:        "map",
		description: "Display the next 20 location areas",
		callback:    commandMap,
	}
	commandRegistry["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display the previous 20 location areas",
		callback:    commandMapBack,
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // unreachable, but required by signature
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(cfg *config) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data locationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	resp, err := http.Get(cfg.Previous)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data locationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	return nil
}
