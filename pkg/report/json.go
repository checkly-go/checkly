package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hawkaii/website-checker.git/pkg/models"
)

// JSONReporter handles JSON output formatting for website check reports
type JSONReporter struct {
	PrettyPrint bool
	Writer      io.Writer
}

// NewJSONReporter creates a new JSON reporter with the specified options
func NewJSONReporter(writer io.Writer, prettyPrint bool) *JSONReporter {
	return &JSONReporter{
		PrettyPrint: prettyPrint,
		Writer:      writer,
	}
}

// GenerateReport creates a comprehensive JSON report from individual check results
func (r *JSONReporter) GenerateReport(url string, results map[string][]models.CheckResult) error {
	start := time.Now()

	// Flatten all results into a single slice
	var allResults []models.CheckResult
	for _, resultSet := range results {
		allResults = append(allResults, resultSet...)
	}

	// Calculate overall score
	score := calculateOverallScore(allResults)

	// Create website report
	report := models.WebsiteReport{
		URL:          url,
		Timestamp:    start,
		Duration:     time.Since(start),
		Results:      allResults,
		OverallScore: score,
	}

	return r.WriteReport(report)
}

// WriteReport writes a WebsiteReport to the output in JSON format
func (r *JSONReporter) WriteReport(report models.WebsiteReport) error {
	var output []byte
	var err error

	if r.PrettyPrint {
		output, err = json.MarshalIndent(report, "", "  ")
	} else {
		output, err = json.Marshal(report)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	_, err = r.Writer.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

// WriteSummaryReport writes a condensed summary report in JSON format
func (r *JSONReporter) WriteSummaryReport(url string, results map[string][]models.CheckResult) error {
	start := time.Now()

	// Create summary structure
	summary := struct {
		URL       string                     `json:"url"`
		Timestamp time.Time                  `json:"timestamp"`
		Duration  time.Duration              `json:"duration"`
		Summary   map[string]CategorySummary `json:"summary"`
		Score     int                        `json:"overall_score"`
	}{
		URL:       url,
		Timestamp: start,
		Duration:  time.Since(start),
		Summary:   make(map[string]CategorySummary),
	}

	// Generate category summaries
	totalResults := 0
	for category, resultSet := range results {
		categorySummary := generateCategorySummary(resultSet)
		summary.Summary[category] = categorySummary
		totalResults += len(resultSet)
	}

	// Calculate overall score
	var allResults []models.CheckResult
	for _, resultSet := range results {
		allResults = append(allResults, resultSet...)
	}
	summary.Score = calculateOverallScore(allResults)

	// Marshal and write
	var output []byte
	var err error

	if r.PrettyPrint {
		output, err = json.MarshalIndent(summary, "", "  ")
	} else {
		output, err = json.Marshal(summary)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal summary to JSON: %w", err)
	}

	_, err = r.Writer.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write JSON summary: %w", err)
	}

	return nil
}

// CategorySummary represents a summary of checks within a category
type CategorySummary struct {
	TotalChecks int           `json:"total_checks"`
	Passed      int           `json:"passed"`
	Warnings    int           `json:"warnings"`
	Failed      int           `json:"failed"`
	Score       int           `json:"score"`
	Status      models.Status `json:"overall_status"`
	Issues      []string      `json:"issues,omitempty"`
}

// generateCategorySummary creates a summary for a category of checks
func generateCategorySummary(results []models.CheckResult) CategorySummary {
	summary := CategorySummary{
		TotalChecks: len(results),
		Issues:      []string{},
	}

	for _, result := range results {
		switch result.Status {
		case models.StatusPass:
			summary.Passed++
		case models.StatusWarning:
			summary.Warnings++
			summary.Issues = append(summary.Issues, fmt.Sprintf("%s: %s", result.Name, result.Message))
		case models.StatusFail:
			summary.Failed++
			summary.Issues = append(summary.Issues, fmt.Sprintf("%s: %s", result.Name, result.Message))
		}
	}

	// Calculate category score (0-100)
	if summary.TotalChecks == 0 {
		summary.Score = 0
	} else {
		passScore := summary.Passed * 100
		warningScore := summary.Warnings * 50
		totalPossible := summary.TotalChecks * 100
		summary.Score = (passScore + warningScore) / totalPossible
	}

	// Determine overall status
	if summary.Failed > 0 {
		summary.Status = models.StatusFail
	} else if summary.Warnings > 0 {
		summary.Status = models.StatusWarning
	} else {
		summary.Status = models.StatusPass
	}

	return summary
}

// WriteDetailedReport writes a detailed JSON report with additional metadata
func (r *JSONReporter) WriteDetailedReport(url string, results map[string][]models.CheckResult) error {
	start := time.Now()

	// Create detailed structure
	detailed := struct {
		URL        string                      `json:"url"`
		Timestamp  time.Time                   `json:"timestamp"`
		Duration   time.Duration               `json:"duration"`
		Categories map[string]DetailedCategory `json:"categories"`
		Summary    map[string]CategorySummary  `json:"summary"`
		Score      int                         `json:"overall_score"`
		Metadata   ReportMetadata              `json:"metadata"`
	}{
		URL:        url,
		Timestamp:  start,
		Duration:   time.Since(start),
		Categories: make(map[string]DetailedCategory),
		Summary:    make(map[string]CategorySummary),
		Metadata: ReportMetadata{
			Version:     "1.0.0",
			Generator:   "website-checker",
			TotalChecks: 0,
		},
	}

	// Process each category
	for category, resultSet := range results {
		detailed.Categories[category] = DetailedCategory{
			Name:        category,
			Description: getCategoryDescription(category),
			Results:     resultSet,
		}
		detailed.Summary[category] = generateCategorySummary(resultSet)
		detailed.Metadata.TotalChecks += len(resultSet)
	}

	// Calculate overall score
	var allResults []models.CheckResult
	for _, resultSet := range results {
		allResults = append(allResults, resultSet...)
	}
	detailed.Score = calculateOverallScore(allResults)

	// Marshal and write
	var output []byte
	var err error

	if r.PrettyPrint {
		output, err = json.MarshalIndent(detailed, "", "  ")
	} else {
		output, err = json.Marshal(detailed)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal detailed report to JSON: %w", err)
	}

	_, err = r.Writer.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write JSON detailed report: %w", err)
	}

	return nil
}

// DetailedCategory represents a category with additional metadata
type DetailedCategory struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Results     []models.CheckResult `json:"results"`
}

// ReportMetadata contains metadata about the report generation
type ReportMetadata struct {
	Version     string    `json:"version"`
	Generator   string    `json:"generator"`
	TotalChecks int       `json:"total_checks"`
	GeneratedAt time.Time `json:"generated_at"`
}

// getCategoryDescription returns a description for the given category
func getCategoryDescription(category string) string {
	descriptions := map[string]string{
		"robots":        "Robots.txt file accessibility and configuration",
		"sitemap":       "XML sitemap availability and structure validation",
		"seo":           "Search Engine Optimization metadata and best practices",
		"security":      "Security headers and information disclosure prevention",
		"performance":   "Page loading speed and Core Web Vitals",
		"accessibility": "Web accessibility compliance and best practices",
	}

	if desc, exists := descriptions[category]; exists {
		return desc
	}
	return fmt.Sprintf("Checks for %s category", category)
}

// WriteRawResults writes just the raw results without any additional processing
func (r *JSONReporter) WriteRawResults(results []models.CheckResult) error {
	var output []byte
	var err error

	if r.PrettyPrint {
		output, err = json.MarshalIndent(results, "", "  ")
	} else {
		output, err = json.Marshal(results)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	_, err = r.Writer.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write JSON results: %w", err)
	}

	return nil
}
