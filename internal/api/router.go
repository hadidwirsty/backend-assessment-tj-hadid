package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/vehicles/:vehicle_id/location", handler.GetLastLocation)
	r.GET("/vehicles/:vehicle_id/history", handler.GetLocationHistory)

	return r
}
