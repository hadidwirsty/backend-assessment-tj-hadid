package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hadid/backend-assessment-tj-hadid/internal/config"
	"github.com/hadid/backend-assessment-tj-hadid/internal/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables")
	}

	cfg := config.Load()

	log.Println("Starting Geofence Worker...")

	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	log.Println("Geofence Worker connected to RabbitMQ, waiting for alerts...")

	go func() {
		if err := consumer.Consume(); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Geofence Worker shutting down...")
}
