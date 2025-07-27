package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hawkaii/website-checker.git/pkg/checker"
)

func main() {
	fmt.Println("Website Checker - Robots.txt Example")
	fmt.Println("=====================================")

	// Example URLs to test
	testURLs := []string{
		"https://google.com",
		"https://github.com",
		"https://hawkaii.netlify.app",
		"https://nonexistent-site-12345.com",
	}

	for _, url := range testURLs {
		fmt.Printf("\nChecking robots.txt for: %s\n", url)
		fmt.Println("----------------------------------------")

		result := checker.CheckRobotsTxt(url)

		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Printf("Error marshaling result: %v", err)
			continue
		}

		fmt.Println(string(jsonResult))

		fmt.Printf("Result: %s - %s\n", result.Status, result.Message)
		if result.Details != "" {
			fmt.Printf("Details: %s\n", result.Details)
		}
	}
}
