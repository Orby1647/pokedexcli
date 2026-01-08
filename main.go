package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/orby1647/pokedexcli/internal/pokeapi"
)

func main() {
	cfg := &config{}
	scanner := bufio.NewScanner(os.Stdin)

	client := pokeapi.NewClient(5*time.Second, 10*time.Second)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			return
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		cmd, ok := commandRegistry[cmdName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(cfg, &client); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
