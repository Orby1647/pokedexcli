package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"

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
	Pokedex  map[string]pokeapi.Pokemon
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
	commandRegistry["catch"] = cliCommand{
		name:        "catch <pokemon_name>",
		description: "Catch a specific Pokemon",
		callback:    commandCatch,
	}
	commandRegistry["inspect"] = cliCommand{
		name:        "inspect <pokemon_name>",
		description: "Inspect a caught Pokemon",
		callback:    commandInspect,
	}
	commandRegistry["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List all caught Pokemon",
		callback:    commandPokedex,
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

func commandCatch(cfg *config, client *pokeapi.Client, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: catch <pokemon>")
		return nil
	}

	name := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	// Fetch Pokémon data
	p, err := client.GetPokemon(name)
	if err != nil {
		return fmt.Errorf("could not find Pokémon %s", name)
	}

	// Already caught?
	if _, exists := cfg.Pokedex[p.Name]; exists {
		fmt.Printf("%s is already in your Pokedex!\n", p.Name)
		return nil
	}

	// Catch chance based on base experience
	// Higher base experience = harder to catch
	chance := 100 - p.BaseExperience
	if chance < 10 {
		chance = 10 // always at least 10% chance
	}

	roll := rand.Intn(100)

	if roll < chance {
		fmt.Printf("%s was caught!\n", p.Name)
		cfg.Pokedex[p.Name] = p
	} else {
		fmt.Printf("%s escaped!\n", p.Name)
	}

	return nil
}

func commandInspect(cfg *config, client *pokeapi.Client, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: inspect <pokemon>")
		return nil
	}

	name := args[0]

	// Check if Pokémon is in the Pokedex
	p, ok := cfg.Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// Print details
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)

	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, client *pokeapi.Client, args []string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokémon!")
		return nil
	}

	fmt.Println("Your Pokedex:")

	// Collect names and sort them
	names := make([]string, 0, len(cfg.Pokedex))
	for name := range cfg.Pokedex {
		names = append(names, name)
	}
	sort.Strings(names)

	// Print each Pokémon
	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
