package report

import "github.com/checkly-go/checkly/pkg/models"

// calculateOverallScore computes an overall score from 0-100 based on all check results
func calculateOverallScore(results []models.CheckResult) int {
	if len(results) == 0 {
		return 0
	}

	totalScore := 0
	for _, result := range results {
		switch result.Status {
		case models.StatusPass:
			totalScore += 100
		case models.StatusWarning:
			totalScore += 50
		case models.StatusFail:
			totalScore += 0
		}
	}

	return totalScore / len(results)
}
