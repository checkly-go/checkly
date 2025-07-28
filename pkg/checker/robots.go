package checker

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/checkly-go/checkly/pkg/models"
)

// CheckRobotsTxt checks if robots.txt exists and is accessible
func CheckRobotsTxt(baseURL string) models.CheckResult {
	start := time.Now()

	u, err := url.Parse(baseURL)
	if err != nil {
		return models.CheckResult{
			Name:      "Robots.txt",
			Status:    models.StatusFail,
			Message:   "Invalid URL provided",
			Details:   err.Error(),
			Timestamp: start,
		}
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)
	fmt.Println(robotsURL)

	// Make HTTP GET request
	resp, err := http.Get(robotsURL)
	if err != nil {
		return models.CheckResult{
			Name:      "Robots.txt",
			Status:    models.StatusFail,
			Message:   "Failed to fetch robots.txt",
			Details:   err.Error(),
			Timestamp: start,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return models.CheckResult{
			Name:      "Robots.txt",
			Status:    models.StatusPass,
			Message:   "Found and accessible",
			Details:   fmt.Sprintf("HTTP %d", resp.StatusCode),
			Timestamp: start,
		}
	} else if resp.StatusCode == 404 {
		return models.CheckResult{
			Name:      "Robots.txt",
			Status:    models.StatusFail,
			Message:   "Missing robots.txt",
			Details:   "Create a robots.txt file in your website's root to guide search engines",
			Timestamp: start,
		}
	}

	return models.CheckResult{
		Name:      "Robots.txt",
		Status:    models.StatusWarning,
		Message:   fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
		Details:   "robots.txt returned an unexpected response",
		Timestamp: start,
	}
}
