package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hawkaii/website-checker.git/pkg/checker"
	"github.com/hawkaii/website-checker.git/pkg/models"
	"github.com/hawkaii/website-checker.git/pkg/report"
)

type Config struct {
	URL        string
	Checkers   []string
	Output     string
	OutputFile string
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
		var writer io.Writer = os.Stdout

		// If output file is specified, create the file
		if config.OutputFile != "" {
			file, err := os.Create(config.OutputFile)
			if err != nil {
				log.Fatalf("Error creating output file: %v", err)
			}
			defer file.Close()
			writer = file
			fmt.Printf("Writing JSON report to: %s\n", config.OutputFile)
		}

		jsonReporter := report.NewJSONReporter(writer, true)
		err := jsonReporter.GenerateReport(config.URL, allResults)
		if err != nil {
			log.Printf("Error generating JSON report: %v", err)
		}

		// If writing to file, also show success message
		if config.OutputFile != "" {
			fmt.Printf("JSON report generated successfully!\n")
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
	var tuiMode bool

	flag.StringVar(&config.URL, "url", "", "URL to check (required)")
	flag.StringVar(&config.URL, "link", "", "URL to check (alias for -url)")
	flag.BoolVar(&tuiMode, "tui", false, "Run in TUI mode (interactive terminal UI)")

	var checkersFlag string
	flag.StringVar(&checkersFlag, "checkers", "robots,sitemap,seo,security", "Comma-separated list of checkers to run (robots,sitemap,seo,security)")

	flag.StringVar(&config.Output, "output", "text", "Output format (text or json)")
	flag.StringVar(&config.OutputFile, "o", "", "Output file path (for JSON reports)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -url https://example.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -tui\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -link https://example.com -checkers robots,seo -output json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -url https://example.com -output json -o report.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -url https://example.com -checkers security -output text\n", os.Args[0])
	}

	flag.Parse()

	// If TUI mode is requested, launch TUI
	if tuiMode {
		runTUI()
		os.Exit(0)
	}

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

	// Validate output file usage
	if config.OutputFile != "" && config.Output != "json" {
		fmt.Println("Warning: Output file (-o) only supported with JSON output format. Setting output to 'json'.")
		config.Output = "json"
	}

	return config
}

func runTUI() {
	// Build and run the TUI binary
	fmt.Println("🚀 Starting Website Checker TUI...")

	// Try to run the TUI binary directly if it exists
	cmd := exec.Command("go", "run", "./cmd/tui/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		fmt.Println("Please ensure the TUI binary is built or run: go run ./cmd/tui/main.go")
	}
}
