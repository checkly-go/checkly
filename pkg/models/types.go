package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckResult struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Status string

const (
	StatusPass    Status = "pass"
	StatusWarning Status = "warning"
	StatusFail    Status = "fail"
)

type WebsiteReport struct {
	URL          string        `json:"url"`
	Timestamp    time.Time     `json:"timestamp"`
	Duration     time.Duration `json:"duration"`
	Results      []CheckResult `json:"results"`
	OverallScore int           `json:"overall_score"`
}

// WebsiteCheck represents a stored check in the database.
type WebsiteCheck struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL       string             `bson:"url" json:"url"`
	Status    string             `bson:"status" json:"status"`
	Report    *WebsiteReport     `bson:"report,omitempty" json:"report,omitempty"`
	Error     string             `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Recommendation models for AI-powered recommendations

// RecommendationRequest represents the request payload for recommendations
type RecommendationRequest struct {
	CheckID string         `json:"check_id,omitempty"` // Reference existing check
	URL     string         `json:"url,omitempty"`      // Or provide data directly
	Report  *WebsiteReport `json:"report,omitempty"`   // Direct report data
	Focus   []string       `json:"focus,omitempty"`    // Focus areas: robots, seo, security, sitemap
}

// RecommendationResponse represents the AI-generated recommendations
type RecommendationResponse struct {
	URL             string                   `json:"url"`
	GeneratedAt     time.Time                `json:"generated_at"`
	Summary         string                   `json:"summary"`
	Recommendations []CategoryRecommendation `json:"recommendations"`
}

// CategoryRecommendation represents recommendations for a specific category
type CategoryRecommendation struct {
	Category     string        `json:"category"`            // robots, seo, security, sitemap
	Priority     string        `json:"priority"`            // high, medium, low
	Issues       []IssueDetail `json:"issues"`              // Detailed issue descriptions
	Improvements []string      `json:"improvements"`        // Actionable improvement steps
	Resources    []string      `json:"resources,omitempty"` // Optional resources/links
}

// IssueDetail represents a specific issue found in the analysis
type IssueDetail struct {
	Issue         string `json:"issue"`
	Impact        string `json:"impact"`
	CurrentStatus string `json:"current_status"`
}
