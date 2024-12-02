package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/pavletto/event_processor/internal/db"
	"github.com/pavletto/event_processor/internal/handlers"
)

func main() {
	db.Init()

	videoUploadPath := os.Getenv("VIDEO_UPLOAD_PATH")
	imageUploadPath := os.Getenv("IMAGE_UPLOAD_PATH")

	if _, err := os.Stat(videoUploadPath); os.IsNotExist(err) {
		os.MkdirAll(videoUploadPath, os.ModePerm)
	}
	if _, err := os.Stat(imageUploadPath); os.IsNotExist(err) {
		os.MkdirAll(imageUploadPath, os.ModePerm)
	}

	router := gin.Default()

	router.POST("/upload", handlers.UploadHandler)

	router.GET("/sources", handlers.GetSourcesHandler)
	router.GET("/source/:id", handlers.GetSourceHandler)

	router.GET("/video/:id", handlers.GetVideoHandler)

	router.GET("/image/:id", handlers.GetImageHandler)

	router.Run(":8080")
}
