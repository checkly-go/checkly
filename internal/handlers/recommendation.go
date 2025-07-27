package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/hawkaii/website-checker.git/pkg/ai"
	"github.com/hawkaii/website-checker.git/pkg/models"
)

// GetRecommendations handles POST /api/v1/recommend
// It accepts either a check_id to get recommendations for an existing check
// or direct report data for immediate analysis
func (s *Service) GetRecommendations(c *gin.Context) {
	var payload models.RecommendationRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var report *models.WebsiteReport
	var url string

	// If check_id is provided, fetch the existing check
	if payload.CheckID != "" {
		id, err := primitive.ObjectIDFromHex(payload.CheckID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID format"})
			return
		}

		ctx := context.Background()
		check, err := s.CheckRepo.GetCheck(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Check not found: " + err.Error()})
			return
		}

		if check.Report == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No report available for this check"})
			return
		}

		report = check.Report
		url = check.URL
	} else if payload.Report != nil && payload.URL != "" {
		// Use directly provided data
		report = payload.Report
		url = payload.URL
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either check_id or both url and report must be provided"})
		return
	}

	// Initialize Gemini client
	ctx := context.Background()
	geminiClient, err := ai.NewGeminiClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AI service: " + err.Error()})
		return
	}
	defer geminiClient.Close()

	// Generate recommendations
	recommendations, err := geminiClient.GenerateRecommendations(ctx, url, report, payload.Focus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recommendations: " + err.Error()})
		return
	}

	// Set generated timestamp
	recommendations.GeneratedAt = time.Now()

	c.JSON(http.StatusOK, recommendations)
}
