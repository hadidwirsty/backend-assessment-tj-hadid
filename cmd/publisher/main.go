package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hadid/backend-assessment-tj-hadid/internal/config"
	"github.com/hadid/backend-assessment-tj-hadid/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables")
	}

	cfg := config.Load()

	vehicleIDs := []string{"B1234XYZ", "B5678ABC", "B9999DEF"}
	baseLat := -6.2088
	baseLon := 106.8456

	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTTBroker)
	opts.SetClientID("fleet-publisher-" + strconv.Itoa(rand.Intn(10000)))
	opts.SetCleanSession(true)

	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Println("MQTT Publisher started, sending location data every 2 seconds...")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, vehicleID := range vehicleIDs {
			lat := baseLat + (rand.Float64()-0.5)*0.001
			lon := baseLon + (rand.Float64()-0.5)*0.001

			if rand.Intn(5) == 0 {
				lat = baseLat
				lon = baseLon
			}

			payload := database.LocationPayload{
				VehicleID: vehicleID,
				Latitude:  lat,
				Longitude: lon,
				Timestamp: time.Now().Unix(),
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Error marshalling payload: %v", err)
				continue
			}

			topic := "/fleet/vehicle/" + vehicleID + "/location"
			token := client.Publish(topic, 1, false, payloadBytes)
			token.Wait()

			if token.Error() != nil {
				log.Printf("Error publishing to topic %s: %v", topic, token.Error())
			} else {
				log.Printf("Published location for %s: lat=%.4f, lon=%.4f", vehicleID, lat, lon)
			}
		}
	}
}
