package main

import (
	"fmt"
	"os"

	"github.com/orby1647/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokeapi.Client, []string) error
}

type config struct {
	Next     string
	Previous string
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
	commandRegistry["explore"] = cliCommand{
		name:        "explore <location_area>",
		description: "Explore a specific location area",
		callback:    commandExplore,
	}
}

func commandExit(cfg *config, client *pokeapi.Client, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, client *pokeapi.Client, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(cfg *config, client *pokeapi.Client, args []string) error {
	resp, err := client.ListLocationAreas(&cfg.Next)
	if err != nil {
		return err
	}

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = resp.Next
	cfg.Previous = resp.Previous

	return nil
}

func commandMapBack(cfg *config, client *pokeapi.Client, args []string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	resp, err := client.ListLocationAreas(&cfg.Previous)
	if err != nil {
		return err
	}

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = resp.Next
	cfg.Previous = resp.Previous

	return nil
}

func commandExplore(cfg *config, client *pokeapi.Client, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: explore <location-area>")
		return nil
	}

	area := args[0]

	fmt.Printf("Exploring %s...\n", area)

	resp, err := client.ExploreLocation(area)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, encounter := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
