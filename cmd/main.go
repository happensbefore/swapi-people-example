package main

import (
	"fmt"
	"os"

	"swapi/internal/app"
)

func main() {

	if len(os.Args) < 2 {
		usage()
		return
	}

	command := os.Args[1]

	switch command {
	case "start":
		fmt.Println("\n---=== LOADING DATA.... ===---\n")

		app.New().Start()
	default:
		usage()
		return
	}

	fmt.Println("\n---=== DONE ===---")
}

func usage() {
	fmt.Println("Please provide start command")
}
