package main

import (
	"log"
	"os"

	"TGBot/client/telegramclient"
	"TGBot/config"
	"TGBot/consumer/eventconsumer"
	"TGBot/events/telegram"
	"TGBot/storage/files"

	"github.com/joho/godotenv"
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
	storagePath, exists := os.LookupEnv("STORAGE")
	if !exists {
		log.Fatalf("failed to retrieve storage path")
	}
	storage := files.NewPath(storagePath)
	tgClient := telegramclient.NewClient(host, cfg.TgBotToken)
	eventsProccesor := telegram.NewDispatcher(tgClient, storage)

	log.Print("service started")
	consumer := eventconsumer.NewConsumer(eventsProccesor, eventsProccesor, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}
