package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadid/backend-assessment-tj-hadid/internal/database"
)

type Handler struct {
	repo *database.LocationRepository
}

func NewHandler(repo *database.LocationRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	loc, err := h.repo.GetLastLocation(c.Request.Context(), vehicleID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicle_id": loc.VehicleID,
		"latitude":   loc.Latitude,
		"longitude":  loc.Longitude,
		"timestamp":  loc.Timestamp,
	})
}

func (h *Handler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")

	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end query params are required"})
		return
	}

	startTs, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end must be valid unix timestamps"})
		return
	}

	endTs, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end must be valid unix timestamps"})
		return
	}

	if startTs >= endTs {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start must be less than end"})
		return
	}

	locations, err := h.repo.GetLocationHistory(c.Request.Context(), vehicleID, startTs, endTs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if locations == nil {
		locations = []*database.VehicleLocation{}
	}

	c.JSON(http.StatusOK, locations)
}
