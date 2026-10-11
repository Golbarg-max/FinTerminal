package main

import (
	"FinTerminal/cmd"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file; using environment variables %v\n", err)

	}
	cmd.Execute()
}
