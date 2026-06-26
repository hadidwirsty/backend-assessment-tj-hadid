# Fleet Management System — Transjakarta Backend Assessment

## Architecture Overview
Sistem ini terdiri dari 6 komponen yang berjalan dalam Docker:
- **Mosquitto** — MQTT broker, menerima data lokasi kendaraan
- **PostgreSQL** — Database penyimpanan data lokasi
- **RabbitMQ** — Message broker untuk geofence events
- **Server** — REST API (Gin) + MQTT Subscriber
- **Worker** — Geofence alert consumer dari RabbitMQ
- **Publisher** — Mock script pengirim data lokasi setiap 2 detik

Flow data:
Publisher → Mosquitto → Server (subscriber) → PostgreSQL
                                              → RabbitMQ → Worker

## Demo Video
[▶ Watch on YouTube](https://youtu.be/TW_JuyE_258)

## Notes
- The PDF specification contains a typo in the table name (`vehicle_loctions`). 
  This implementation uses the correct spelling `vehicle_locations` for maintainability.

## Prerequisites
- Docker
- Docker Compose v2+

## How to Run

### 1. Clone repository
git clone git clone https://github.com/hadidwirsty/backend-assessment-tj-hadid.git
cd backend-assessment-tj-hadid

### 2. Start all services
docker compose up --build

### 3. Verify services are running
docker compose ps

## API Endpoints

### Health Check
GET http://localhost:8080/health

### Get Last Known Location
GET http://localhost:8080/vehicles/{vehicle_id}/location

Example:
curl http://localhost:8080/vehicles/B1234XYZ/location

Expected response:
{
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": 1715003456
}

### Get Location History
GET http://localhost:8080/vehicles/{vehicle_id}/history?start={unix_ts}&end={unix_ts}

Example:
curl "http://localhost:8080/vehicles/B1234XYZ/history?start=1700000000&end=9999999999"

## RabbitMQ Management UI
URL: http://localhost:15672
Username: guest
Password: guest
Queue to monitor: geofence_alerts

## Mock Vehicle IDs (from publisher)
- B1234XYZ
- B5678ABC
- B9999DEF

## Geofence Configuration
Center: -6.2088, 106.8456 (Jakarta)
Radius: 50 meters
The publisher randomly triggers exact geofence coordinates (1 in 5 chance per vehicle per tick)

## Stopping the System
docker compose down

To remove all data volumes:
docker compose down -v
