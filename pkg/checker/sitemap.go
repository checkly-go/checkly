package checker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hawkaii/website-checker.git/pkg/models"
)

// CheckSitemap checks if sitemap exists and is accessible
// It first looks for sitemap reference in robots.txt content, then tries default locations
func CheckSitemap(baseURL string, robotsContent string) models.CheckResult {
	start := time.Now()

	// Parse the base URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusFail,
			Message:   "Invalid URL provided",
			Details:   err.Error(),
			Timestamp: start,
		}
	}

	// First, try to find sitemap URL from robots.txt content
	sitemapURLs := parseSitemapFromRobots(robotsContent, u)

	//  default locations
	if len(sitemapURLs) == 0 {
		defaultSitemapURL := fmt.Sprintf("%s://%s/sitemap.xml", u.Scheme, u.Host)
		sitemapURLs = append(sitemapURLs, defaultSitemapURL)
	}

	// Test each sitemap URL
	for _, sitemapURL := range sitemapURLs {
		result := testSitemapURL(sitemapURL, start)
		if result.Status == models.StatusPass {
			return result
		}
	}

	// If we get here, no sitemap was found
	return models.CheckResult{
		Name:      "Sitemap",
		Status:    models.StatusFail,
		Message:   "No accessible sitemap found",
		Details:   "Create a sitemap.xml file and reference it in robots.txt to help search engines index your site",
		Timestamp: start,
	}
}

// parseSitemapFromRobots extracts sitemap URLs from robots.txt content
func parseSitemapFromRobots(robotsContent string, baseURL *url.URL) []string {
	var sitemapURLs []string

	lines := strings.Split(robotsContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			sitemapURL := strings.TrimSpace(line[8:]) // Remove "sitemap:" prefix

			// If it's a relative URL, make it absolute
			if strings.HasPrefix(sitemapURL, "/") {
				sitemapURL = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, sitemapURL)
			}

			sitemapURLs = append(sitemapURLs, sitemapURL)
		}
	}

	return sitemapURLs
}

// testSitemapURL tests if a specific sitemap URL is accessible and valid
func testSitemapURL(sitemapURL string, startTime time.Time) models.CheckResult {
	resp, err := http.Get(sitemapURL)
	if err != nil {
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusFail,
			Message:   "Failed to fetch sitemap",
			Details:   fmt.Sprintf("Error accessing %s: %v", sitemapURL, err),
			Timestamp: startTime,
		}
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return models.CheckResult{
				Name:      "Sitemap",
				Status:    models.StatusFail,
				Message:   "Sitemap not found",
				Details:   fmt.Sprintf("HTTP %d for %s", resp.StatusCode, sitemapURL),
				Timestamp: startTime,
			}
		}
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusWarning,
			Message:   fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
			Details:   fmt.Sprintf("Sitemap at %s returned HTTP %d", sitemapURL, resp.StatusCode),
			Timestamp: startTime,
		}
	}

	// Read content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusWarning,
			Message:   "Could not read sitemap content",
			Details:   fmt.Sprintf("Error reading content from %s: %v", sitemapURL, err),
			Timestamp: startTime,
		}
	}

	// Basic validation that it's XML and contains sitemap elements
	contentStr := string(content)
	if !strings.Contains(contentStr, "<?xml") {
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusWarning,
			Message:   "Sitemap doesn't appear to be valid XML",
			Details:   fmt.Sprintf("Content from %s doesn't start with XML declaration", sitemapURL),
			Timestamp: startTime,
		}
	}

	// Check for sitemap-specific elements
	hasSitemapElements := strings.Contains(contentStr, "<urlset") ||
		strings.Contains(contentStr, "<sitemapindex") ||
		strings.Contains(contentStr, "<url>") ||
		strings.Contains(contentStr, "<sitemap>")

	if !hasSitemapElements {
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusWarning,
			Message:   "Sitemap XML doesn't contain expected elements",
			Details:   fmt.Sprintf("File at %s appears to be XML but missing sitemap elements", sitemapURL),
			Timestamp: startTime,
		}
	}

	// checking url count
	urlCount := strings.Count(contentStr, "<url>")
	sitemapCount := strings.Count(contentStr, "<sitemap>")

	var details string
	if urlCount > 0 {
		details = fmt.Sprintf("Found valid sitemap with %d URLs at %s", urlCount, sitemapURL)
	} else if sitemapCount > 0 {
		details = fmt.Sprintf("Found sitemap index with %d sitemaps at %s", sitemapCount, sitemapURL)
	} else {
		details = fmt.Sprintf("Found valid sitemap at %s", sitemapURL)
	}

	return models.CheckResult{
		Name:      "Sitemap",
		Status:    models.StatusPass,
		Message:   "Sitemap found and accessible",
		Details:   details,
		Timestamp: startTime,
	}
}

func CheckSitemapWithRobotsURL(baseURL string) models.CheckResult {
	start := time.Now()

	u, err := url.Parse(baseURL)
	if err != nil {
		return models.CheckResult{
			Name:      "Sitemap",
			Status:    models.StatusFail,
			Message:   "Invalid URL provided",
			Details:   err.Error(),
			Timestamp: start,
		}
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)

	// checking the content
	var robotsContent string
	resp, err := http.Get(robotsURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		if content, err := io.ReadAll(resp.Body); err == nil {
			robotsContent = string(content)
		}
	}

	return CheckSitemap(baseURL, robotsContent)
}
