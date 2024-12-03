package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/pavletto/event_processor/internal/db"
	"github.com/pavletto/event_processor/internal/models"
)

type CreateEventRequest struct {
	SourceID  string  `json:"source_id" binding:"required"`
	StartTime float32 `json:"start_time" binding:"required"`
	EndTime   float32 `json:"end_time" binding:"required"`
	Tag       string  `json:"tag" binding:"required,oneof=danger safe"`
	Comment   string  `json:"comment"`
}

func CreateEventHandler(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EndTime < req.StartTime {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EndTime must be after StartTime"})
		return
	}

	var source models.Source
	if err := db.DB.First(&source, "id = ?", req.SourceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	event := models.Event{
		SourceID:  req.SourceID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Tag:       req.Tag,
		Comment:   req.Comment,
	}

	if err := db.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	response := gin.H{
		"id":         event.ID,
		"source_id":  event.SourceID,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
		"tag":        event.Tag,
		"comment":    event.Comment,
		"created_at": event.CreatedAt.Unix(),
		"updated_at": event.UpdatedAt.Unix(),
	}

	c.JSON(http.StatusOK, response)
}

func GetEventHandler(c *gin.Context) {
	id := c.Param("id")

	var event models.Event

	if err := db.DB.First(&event, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	response := gin.H{
		"id":         event.ID,
		"source_id":  event.SourceID,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
		"tag":        event.Tag,
		"comment":    event.Comment,
		"created_at": event.CreatedAt.Unix(),
		"updated_at": event.UpdatedAt.Unix(),
	}

	c.JSON(http.StatusOK, response)
}

func GetEventsHandler(c *gin.Context) {
	var events []models.Event

	tag := c.Query("tag")
	minRisk := c.Query("minRisk")
	maxRisk := c.Query("maxRisk")
	minSpeed := c.Query("minSpeed")
	maxSpeed := c.Query("maxSpeed")

	query := db.DB.Model(&models.Event{})

	if tag != "" {
		query = query.Where("tag = ?", tag)
	}

	if minRisk != "" || maxRisk != "" {
		query = query.Joins("JOIN risk_data ON risk_data.source_id = events.source_id")
		if minRisk != "" {
			query = query.Where("risk_data.prob >= ?", minRisk)
		}
		if maxRisk != "" {
			query = query.Where("risk_data.prob <= ?", maxRisk)
		}
	}

	if minSpeed != "" || maxSpeed != "" {
		query = query.Joins("JOIN location_data ON location_data.source_id = events.source_id")
		if minSpeed != "" {
			query = query.Where("location_data.speed_kmh >= ?", minSpeed)
		}
		if maxSpeed != "" {
			query = query.Where("location_data.speed_kmh <= ?", maxSpeed)
		}
	}

	if err := query.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	var response []gin.H
	for _, event := range events {
		respEvent := gin.H{
			"id":         event.ID,
			"source_id":  event.SourceID,
			"start_time": event.StartTime,
			"end_time":   event.EndTime,
			"tag":        event.Tag,
			"comment":    event.Comment,
			"created_at": event.CreatedAt.Unix(),
			"updated_at": event.UpdatedAt.Unix(),
		}
		response = append(response, respEvent)
	}

	c.JSON(http.StatusOK, response)
}
