package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hawkaii/website-checker.git/pkg/checker"
	"github.com/hawkaii/website-checker.git/pkg/models"
	"github.com/hawkaii/website-checker.git/pkg/report"
)

type Config struct {
	URL      string
	Checkers []string
	Output   string
}

func main() {
	config := parseFlags()

	if config.URL == "" {
		fmt.Println("Error: URL is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Website Checker - Analyzing: %s\n", config.URL)
	fmt.Println("=========================================")

	// Collect all results by category
	allResults := make(map[string][]models.CheckResult)

	// Run checkers based on flags
	for _, checkerName := range config.Checkers {
		switch checkerName {
		case "robots":
			result := checker.CheckRobotsTxt(config.URL)
			allResults["robots"] = []models.CheckResult{result}
			if config.Output == "text" {
				fmt.Println("\n📋 Robots.txt Check:")
				fmt.Println("--------------------")
				printTextResult(result)
			}

		case "sitemap":
			result := checker.CheckSitemapWithRobotsURL(config.URL)
			allResults["sitemap"] = []models.CheckResult{result}
			if config.Output == "text" {
				fmt.Println("\n🗺️  Sitemap Check:")
				fmt.Println("------------------")
				printTextResult(result)
			}

		case "seo":
			results := checker.CheckSEOMetadataFromURL(config.URL)
			allResults["seo"] = results
			if config.Output == "text" {
				fmt.Println("\n🏷️  SEO Metadata Checks:")
				fmt.Println("------------------------")
				for _, result := range results {
					printTextResult(result)
					fmt.Println()
				}
			}

		case "security":
			results := checker.CheckSecurityHeaders(config.URL)
			allResults["security"] = results
			if config.Output == "text" {
				fmt.Println("\n🛡️  Security Headers Checks:")
				fmt.Println("----------------------------")
				for _, result := range results {
					printTextResult(result)
					fmt.Println()
				}
			}
		}
	}

	// Output results
	if config.Output == "json" {
		jsonReporter := report.NewJSONReporter(os.Stdout, true)
		err := jsonReporter.GenerateReport(config.URL, allResults)
		if err != nil {
			log.Printf("Error generating JSON report: %v", err)
		}
	}
}

func printTextResult(result models.CheckResult) {
	statusEmoji := getStatusEmoji(result.Status)
	fmt.Printf("%s %s: %s\n", statusEmoji, result.Name, result.Message)
	if result.Details != "" {
		fmt.Printf("   Details: %s\n", result.Details)
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

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.URL, "url", "", "URL to check (required)")
	flag.StringVar(&config.URL, "link", "", "URL to check (alias for -url)")

	var checkersFlag string
	flag.StringVar(&checkersFlag, "checkers", "robots,sitemap,seo,security", "Comma-separated list of checkers to run (robots,sitemap,seo,security)")

	flag.StringVar(&config.Output, "output", "text", "Output format (text or json)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -url https://example.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -link https://example.com -checkers robots,seo -output json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -url https://example.com -checkers security -output text\n", os.Args[0])
	}

	flag.Parse()

	// Parse checkers
	if checkersFlag != "" {
		config.Checkers = strings.Split(checkersFlag, ",")
		for i, checker := range config.Checkers {
			config.Checkers[i] = strings.TrimSpace(checker)
		}
	}

	// Validate checkers
	validCheckers := map[string]bool{
		"robots":   true,
		"sitemap":  true,
		"seo":      true,
		"security": true,
	}

	var filteredCheckers []string
	for _, checker := range config.Checkers {
		if validCheckers[checker] {
			filteredCheckers = append(filteredCheckers, checker)
		} else {
			fmt.Printf("Warning: Unknown checker '%s' ignored\n", checker)
		}
	}
	config.Checkers = filteredCheckers

	// Validate output format
	if config.Output != "text" && config.Output != "json" {
		fmt.Printf("Warning: Unknown output format '%s', defaulting to 'text'\n", config.Output)
		config.Output = "text"
	}

	return config
}
