package main

import (
	"context"
	"log"
	"os"

	"TGBot/client/telegramclient"
	"TGBot/config"
	"TGBot/consumer/eventconsumer"
	"TGBot/events/telegram"
	"TGBot/storage/sqlite"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.MustLoad()

	host, exists := os.LookupEnv("TGHOST")
	if !exists {
		log.Fatalf("failed to retrieve TG host")
	}

	storagePath, exists := os.LookupEnv("SQLSTORAGE")
	if !exists {
		log.Fatalf("failed to retrieve storage path")
	}

	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		log.Fatal("can't connect to storage %w", err)
	}

	storage.Init(context.Background())
	if err != nil {
		log.Fatal("context init troubles %w", err)
	}

	eventsDispatcher := telegram.NewDispatcher(
		telegramclient.NewClient(host, cfg.TgBotToken),
		storage)

	consumer := eventconsumer.NewConsumer(eventsDispatcher, eventsDispatcher, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}
