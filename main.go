package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg := &config{}
	scanner := bufio.NewScanner(os.Stdin)

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

		if err := cmd.callback(cfg); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
