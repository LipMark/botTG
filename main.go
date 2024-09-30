package main

import (
	"fmt"
	"log"
	"os"

	"TGBot/client/telegram"

	"github.com/joho/godotenv"
)

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	token, exists := os.LookupEnv("TGTOKEN")
	if !exists {
		log.Fatalf("failed to retrieve TG token")
	}
	host, exists := os.LookupEnv("TGHOST")
	if !exists {
		log.Fatalf("failed to retrieve TG host")
	}
	tgClient := telegram.NewClient(host, token)
	fmt.Println(tgClient)
}
