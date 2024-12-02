package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/pavletto/event_processor/internal/db"
	"github.com/pavletto/event_processor/internal/models"
	"github.com/pavletto/event_processor/internal/utils"
)

func UploadHandler(c *gin.Context) {
	videoUploadPath := os.Getenv("VIDEO_UPLOAD_PATH")
	imageUploadPath := os.Getenv("IMAGE_UPLOAD_PATH")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "details": err.Error()})
		return
	}

	videoFiles := form.File["videos"]
	imageFiles := form.File["images"]
	jsonFiles := form.File["json"]

	if len(videoFiles) == 0 || len(jsonFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video files or JSON data not provided"})
		return
	}

	if len(videoFiles) != len(jsonFiles) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch between number of video files and JSON data"})
		return
	}

	if len(imageFiles) > 0 && len(imageFiles) != len(videoFiles) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch between number of video files and image files"})
		return
	}

	for i, file := range videoFiles {
		videoFilename := filepath.Base(file.Filename)
		videoSavePath := filepath.Join(videoUploadPath, videoFilename)

		if err := c.SaveUploadedFile(file, videoSavePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video file", "details": err.Error()})
			return
		}

		var imageSavePath string
		if len(imageFiles) > 0 {
			imageFile := imageFiles[i]
			imageFilename := filepath.Base(imageFile.Filename)
			imageSavePath = filepath.Join(imageUploadPath, imageFilename)

			if err := c.SaveUploadedFile(imageFile, imageSavePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file", "details": err.Error()})
				return
			}
		}

		jsonFile, err := jsonFiles[i].Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open JSON file", "details": err.Error()})
			return
		}
		defer jsonFile.Close()

		jsonBytes, err := io.ReadAll(jsonFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read JSON file", "details": err.Error()})
			return
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &raw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
			return
		}

		transformed := utils.TransformDates(raw)

		modifiedJSON, err := json.Marshal(transformed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reformat JSON", "details": err.Error()})
			return
		}

		var source models.Source
		decoder := json.NewDecoder(bytes.NewReader(modifiedJSON))
		if err := decoder.Decode(&source); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
			return
		}

		source.VideoKey = videoSavePath
		if imageSavePath != "" {
			source.ThumbnailKey = imageSavePath
		}

		if err := db.DB.Create(&source).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database", "details": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files and data successfully uploaded"})
}
