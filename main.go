package main

import (
	"github.com/joho/godotenv"
	"log"
	"orderworker/database"
	"orderworker/handlers"
	"orderworker/messaging"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables if available")
	}

	database.Connect()
	defer database.ScyllaSession.Close()

	messaging.Connect()
	defer messaging.Close()

	messaging.Consume(handlers.ProcessEvent)
}
