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
