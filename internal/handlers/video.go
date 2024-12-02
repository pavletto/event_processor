package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/pavletto/event_processor/internal/db"
	"github.com/pavletto/event_processor/internal/models"
)

func GetVideoHandler(c *gin.Context) {
	id := c.Param("id")

	var source models.Source

	if err := db.DB.First(&source, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	videoPath := source.VideoKey

	videoPath = filepath.Clean(videoPath)

	absVideoPath, err := filepath.Abs(videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid video path"})
		return
	}

	uploadsDir, err := filepath.Abs("uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid uploads directory"})
		return
	}

	if !filepath.HasPrefix(absVideoPath, uploadsDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.File(absVideoPath)
}
