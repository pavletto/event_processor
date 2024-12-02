package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/pavletto/event_processor/internal/db"
	"github.com/pavletto/event_processor/internal/models"
)

func GetImageHandler(c *gin.Context) {
	id := c.Param("id")

	var source models.Source

	if err := db.DB.First(&source, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	imagePath := source.ThumbnailKey
	if imagePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	imagePath = filepath.Clean(imagePath)

	absImagePath, err := filepath.Abs(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid image path"})
		return
	}

	imageDir, err := filepath.Abs(os.Getenv("IMAGE_UPLOAD_PATH"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid images directory"})
		return
	}

	if !filepath.HasPrefix(absImagePath, imageDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.File(absImagePath)
}
