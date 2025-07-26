package models

import "time"

// Status represents the status of a check
type Status string

const (
	StatusPass    Status = "pass"    // ✅
	StatusWarning Status = "warning" // 🟡
	StatusFail    Status = "fail"    // ❌
)

// CheckResult represents the result of a single check
type CheckResult struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// WebsiteReport represents the complete report for a website
type WebsiteReport struct {
	URL          string        `json:"url"`
	Timestamp    time.Time     `json:"timestamp"`
	Duration     time.Duration `json:"duration"`
	Results      []CheckResult `json:"results"`
	OverallScore int           `json:"overall_score"`
}
