package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hawkaii/website-checker.git/pkg/checker"
	"github.com/hawkaii/website-checker.git/pkg/models"
)

func main() {
	fmt.Println("Website Checker - Robots.txt, Sitemap & SEO Example")
	fmt.Println("====================================================")

	// Example URLs to test
	testURLs := []string{
		"https://google.com",
		"https://github.com",
		"https://hawkaii.netlify.app",
	}

	for _, url := range testURLs {
		fmt.Printf("\n🔍 Checking: %s\n", url)
		fmt.Println("========================================")

		// Check robots.txt
		fmt.Println("\n📋 Robots.txt Check:")
		fmt.Println("--------------------")
		robotsResult := checker.CheckRobotsTxt(url)
		printResult(robotsResult)

		// Check sitemap
		fmt.Println("\n🗺️  Sitemap Check:")
		fmt.Println("------------------")
		sitemapResult := checker.CheckSitemapWithRobotsURL(url)
		printResult(sitemapResult)

		// Check SEO metadata
		fmt.Println("\n🏷️  SEO Metadata Checks:")
		fmt.Println("------------------------")
		seoResults := checker.CheckSEOMetadataFromURL(url)
		for _, result := range seoResults {
			printResult(result)
			fmt.Println() // Add spacing between SEO checks
		}
	}
}

func printResult(result models.CheckResult) {
	// Pretty print the result as JSON
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error marshaling result: %v", err)
		return
	}

	fmt.Println(string(jsonResult))

	// Also display in a user-friendly format
	statusEmoji := getStatusEmoji(result.Status)
	fmt.Printf("\n%s Result: %s - %s\n", statusEmoji, result.Status, result.Message)
	if result.Details != "" {
		fmt.Printf("Details: %s\n", result.Details)
	}
}

func getStatusEmoji(status models.Status) string {
	switch status {
	case models.StatusPass:
		return "✅"
	case models.StatusWarning:
		return "🟡"
	case models.StatusFail:
		return "❌"
	default:
		return "❓"
	}
}
