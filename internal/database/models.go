package database

import (
	"errors"
	"time"
)

type VehicleLocation struct {
	ID        int64     `json:"-"`
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp int64     `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

type LocationPayload struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func (p *LocationPayload) Validate() error {
	if p.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}
	if p.Latitude < -90 || p.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if p.Longitude < -180 || p.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	if p.Timestamp <= 0 {
		return errors.New("timestamp must be greater than 0")
	}
	return nil
}
