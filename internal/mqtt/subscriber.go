package mqtt

import (
	"context"
	"encoding/json"
	"log"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hadid/backend-assessment-tj-hadid/internal/config"
	"github.com/hadid/backend-assessment-tj-hadid/internal/database"
	"github.com/hadid/backend-assessment-tj-hadid/internal/geofence"
	"github.com/hadid/backend-assessment-tj-hadid/internal/rabbitmq"
)

type Subscriber struct {
	client         pahomqtt.Client
	repo           *database.LocationRepository
	producer       *rabbitmq.Producer
	geofenceLat    float64
	geofenceLon    float64
	geofenceRadius float64
}

func NewSubscriber(cfg *config.Config, repo *database.LocationRepository, producer *rabbitmq.Producer) (*Subscriber, error) {
	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTTBroker)
	opts.SetClientID(cfg.MQTTClientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(func(c pahomqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v\n", err)
	})

	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Subscriber{
		client:         client,
		repo:           repo,
		producer:       producer,
		geofenceLat:    cfg.GeofenceLat,
		geofenceLon:    cfg.GeofenceLon,
		geofenceRadius: cfg.GeofenceRadiusMeters,
	}, nil
}

func (s *Subscriber) Subscribe() error {
	topic := "/fleet/vehicle/+/location"
	token := s.client.Subscribe(topic, 1, s.handleMessage)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Subscriber) handleMessage(client pahomqtt.Client, msg pahomqtt.Message) {
	var payload database.LocationPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("Failed to parse JSON payload: %v", err)
		return
	}

	if err := payload.Validate(); err != nil {
		log.Printf("Payload validation failed: %v", err)
		return
	}

	loc := &database.VehicleLocation{
		VehicleID: payload.VehicleID,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		Timestamp: payload.Timestamp,
	}

	log.Printf("Received location for vehicle %s", payload.VehicleID)

	err := s.repo.InsertLocation(context.Background(), loc)
	if err != nil {
		log.Printf("Failed to insert location: %v", err)
		// Tetap lanjut ke geofence check sesuai instruksi
	}

	isInside := geofence.IsInsideGeofence(payload.Latitude, payload.Longitude, s.geofenceLat, s.geofenceLon, s.geofenceRadius)
	if isInside {
		log.Printf("Geofence triggered for vehicle %s", payload.VehicleID)

		event := rabbitmq.GeofenceEvent{
			VehicleID: payload.VehicleID,
			Event:     "geofence_entry",
			Timestamp: payload.Timestamp,
		}
		event.Location.Latitude = payload.Latitude
		event.Location.Longitude = payload.Longitude

		err = s.producer.PublishGeofenceEvent(event)
		if err != nil {
			log.Printf("Failed to publish geofence event: %v", err)
		}
	}
}
