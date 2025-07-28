package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/checkly-go/checkly/internal/storage"
	"github.com/checkly-go/checkly/pkg/checker"
	"github.com/checkly-go/checkly/pkg/models"
)

// Service holds the shared objects needed by the HTTP handlers.
type Service struct {
	Checker   *checker.Checker
	CheckRepo storage.CheckRepository
	UserRepo  storage.UserRepository // for future authentication
}

// SubmitCheck handles POST /api/v1/check.
// It expects a JSON payload with a URL field.
func (s *Service) SubmitCheck(c *gin.Context) {
	var payload struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Run website check (calls core checker logic)
	report, err := s.Checker.CheckWebsite(payload.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check website: " + err.Error()})
		return
	}

	// Construct a WebsiteCheck model instance.
	check := &models.WebsiteCheck{
		URL:       payload.URL,
		Status:    "completed", // Could be an enum
		Report:    report,
		CreatedAt: time.Now(),
	}

	// Save the check in the database.
	ctx := context.Background()
	if err := s.CheckRepo.CreateCheck(ctx, check); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store check: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, check)
}

// GetCheck handles GET /api/v1/check/:id to retrieve a check.
func (s *Service) GetCheck(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx := context.Background()
	check, err := s.CheckRepo.GetCheck(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, check)
}

// GetCheckReport handles GET /api/v1/check/:id/report to retrieve the check report.
func (s *Service) GetCheckReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx := context.Background()
	check, err := s.CheckRepo.GetCheck(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found: " + err.Error()})
		return
	}

	if check.Report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not available"})
		return
	}

	c.JSON(http.StatusOK, check.Report)
}

// GetLeaderboard handles GET /api/v1/leaderboard to retrieve top websites by score.
func (s *Service) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	ctx := context.Background()
	leaderboard, err := s.CheckRepo.GetLeaderboard(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
		"total":       len(leaderboard),
		"limit":       limit,
	})
}

// GetAllChecks handles GET /api/v1/debug/checks to see all checks (debug endpoint).
func (s *Service) GetAllChecks(c *gin.Context) {
	ctx := context.Background()
	checks, err := s.CheckRepo.GetAllChecks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch checks: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checks": checks,
		"total":  len(checks),
	})
}
