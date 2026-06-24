package main

import (
	"log"

	"github.com/hadid/backend-assessment-tj-hadid/internal/config"
	"github.com/hadid/backend-assessment-tj-hadid/internal/database"
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

	log.Println("Fleet Management Server starting...")
}
