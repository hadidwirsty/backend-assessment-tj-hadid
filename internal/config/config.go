package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	MQTTBroker           string
	MQTTClientID         string
	RabbitMQURL          string
	GeofenceLat          float64
	GeofenceLon          float64
	GeofenceRadiusMeters float64
	APIPort              string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error reading .env file, relying on environment variables")
	}

	geofenceLat, _ := strconv.ParseFloat(os.Getenv("GEOFENCE_LAT"), 64)
	geofenceLon, _ := strconv.ParseFloat(os.Getenv("GEOFENCE_LON"), 64)
	geofenceRadius, _ := strconv.ParseFloat(os.Getenv("GEOFENCE_RADIUS_METERS"), 64)

	return &Config{
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		MQTTBroker:           os.Getenv("MQTT_BROKER"),
		MQTTClientID:         os.Getenv("MQTT_CLIENT_ID"),
		RabbitMQURL:          os.Getenv("RABBITMQ_URL"),
		GeofenceLat:          geofenceLat,
		GeofenceLon:          geofenceLon,
		GeofenceRadiusMeters: geofenceRadius,
		APIPort:              os.Getenv("API_PORT"),
	}
}
