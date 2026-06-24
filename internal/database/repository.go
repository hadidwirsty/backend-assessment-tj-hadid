package database

import (
	"context"
	"database/sql"
)

type LocationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) InsertLocation(ctx context.Context, loc *VehicleLocation) error {
	query := `INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp)
	return err
}

func (r *LocationRepository) GetLastLocation(ctx context.Context, vehicleID string) (*VehicleLocation, error) {
	query := `SELECT id, vehicle_id, latitude, longitude, timestamp, created_at FROM vehicle_locations WHERE vehicle_id = $1 ORDER BY timestamp DESC LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, vehicleID)

	var loc VehicleLocation
	err := row.Scan(&loc.ID, &loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp, &loc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *LocationRepository) GetLocationHistory(ctx context.Context, vehicleID string, start, end int64) ([]*VehicleLocation, error) {
	query := `SELECT id, vehicle_id, latitude, longitude, timestamp, created_at FROM vehicle_locations WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3 ORDER BY timestamp ASC`
	rows, err := r.db.QueryContext(ctx, query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]*VehicleLocation, 0)
	for rows.Next() {
		var loc VehicleLocation
		err := rows.Scan(&loc.ID, &loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp, &loc.CreatedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, &loc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
