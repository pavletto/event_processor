package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pavletto/event_processor/internal/db"
	"github.com/pavletto/event_processor/internal/models"
)

func GetSourcesHandler(c *gin.Context) {
	var sources []models.Source

	if err := db.DB.Preload("Accelerometer").Preload("Location").Preload("Risk").Find(&sources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sources"})
		return
	}

	c.JSON(http.StatusOK, sources)
}
func GetSourceHandler(c *gin.Context) {
	id := c.Param("id")

	var source models.Source

	if err := db.DB.Preload("Accelerometer").Preload("Location").Preload("Risk").First(&source, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}
