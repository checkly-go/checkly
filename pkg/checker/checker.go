package checker

import (
	"io"
	"net/http"
	"time"

	"github.com/hawkaii/website-checker.git/pkg/models"
)

type Checker struct {
	Config Config
}

type Config struct {
	Timeout    time.Duration
	UserAgent  string
	Concurrent bool
}

func NewChecker() *Checker {
	return &Checker{
		Config: Config{
			Timeout:    30 * time.Second,
			UserAgent:  "Website-Checker/1.0",
			Concurrent: true,
		},
	}
}

func (c *Checker) CheckWebsite(url string) (*models.WebsiteReport, error) {
	startTime := time.Now()

	// Create a basic report structure
	report := &models.WebsiteReport{
		URL:       url,
		Timestamp: startTime,
		Results:   []models.CheckResult{},
	}

	// Fetch HTML content for SEO checks
	htmlContent, err := c.fetchHTMLContent(url)
	if err != nil {
		// If we can't fetch HTML, we'll still run other checks
		htmlContent = ""
	}

	// Run individual checks
	robotsResult := CheckRobotsTxt(url)
	report.Results = append(report.Results, robotsResult)

	sitemapResult := CheckSitemap(url, "")
	report.Results = append(report.Results, sitemapResult)

	securityResults := CheckSecurityHeaders(url)
	report.Results = append(report.Results, securityResults...)

	// Only run SEO checks if we have HTML content
	if htmlContent != "" {
		seoResults := CheckSEOMetadata(htmlContent)
		report.Results = append(report.Results, seoResults...)
	}

	// Calculate duration
	report.Duration = time.Since(startTime)

	// Calculate overall score (simple implementation)
	passCount := 0
	for _, result := range report.Results {
		if result.Status == models.StatusPass {
			passCount++
		}
	}

	if len(report.Results) > 0 {
		report.OverallScore = (passCount * 100) / len(report.Results)
	}

	return report, nil
}

func (c *Checker) fetchHTMLContent(url string) (string, error) {
	client := &http.Client{
		Timeout: c.Config.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", c.Config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
