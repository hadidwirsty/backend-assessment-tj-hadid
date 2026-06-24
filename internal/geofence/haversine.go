package geofence

import "math"

const earthRadiusMeters = 6371000.0

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

func IsInsideGeofence(vehicleLat, vehicleLon, centerLat, centerLon, radiusMeters float64) bool {
	dist := Distance(vehicleLat, vehicleLon, centerLat, centerLon)
	return dist <= radiusMeters
}
