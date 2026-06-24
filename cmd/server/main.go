package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hadid/backend-assessment-tj-hadid/internal/api"
	"github.com/hadid/backend-assessment-tj-hadid/internal/config"
	"github.com/hadid/backend-assessment-tj-hadid/internal/database"
	"github.com/hadid/backend-assessment-tj-hadid/internal/mqtt"
	"github.com/hadid/backend-assessment-tj-hadid/internal/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	log.Printf("Loading environment variables (.env)...")
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: error loading .env file: %v", err)
	}

	log.Printf("Loading application configuration...")
	cfg := config.Load()

	log.Printf("Connecting to database...")
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Printf("Database connection established successfully")

	log.Printf("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Printf("Database migrations completed successfully")

	repo := database.NewLocationRepository(db)

	producer, err := rabbitmq.NewProducer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	subscriber, err := mqtt.NewSubscriber(cfg, repo, producer)
	if err != nil {
		log.Fatalf("Failed to initialize MQTT subscriber: %v", err)
	}

	if err := subscriber.Subscribe(); err != nil {
		log.Fatalf("Failed to subscribe to MQTT topic: %v", err)
	}

	log.Println("MQTT subscriber started, listening for vehicle locations...")

	handler := api.NewHandler(repo)
	router := api.NewRouter(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.APIPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	log.Printf("API server started on port %s", cfg.APIPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
