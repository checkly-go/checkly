# Website Checker Go Project Plan

## Project Overview

A comprehensive Go-based website checker that provides a library, CLI tool, and backend service for analyzing websites across multiple dimensions: SEO, performance, security, and accessibility.

## Project Structure

```
website-checker/
├── cmd/
│   ├── cli/                    # CLI application
│   │   └── main.go
│   └── server/                 # Backend service
│       └── main.go
├── pkg/
│   ├── checker/                # Core checking library
│   │   ├── robots.go
│   │   ├── sitemap.go
│   │   ├── favicon.go
│   │   ├── seo.go
│   │   ├── performance.go
│   │   ├── accessibility.go
│   │   ├── security.go
│   │   └── checker.go
│   ├── api/                    # External API integrations
│   │   ├── pagespeed.go
│   │   └── accessibility.go
│   ├── report/                 # Report generation
│   │   ├── console.go
│   │   ├── json.go
│   │   └── html.go
│   └── models/                 # Data structures
│       └── types.go
├── internal/
│   ├── config/                 # Configuration management
│   ├── server/                 # HTTP server logic
│   └── storage/                # Data persistence (future)
├── web/                        # Frontend assets (future)
├── docs/                       # Documentation
├── examples/                   # Usage examples
├── scripts/                    # Build and deployment scripts
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── Plans.md
```

## Phase 1: Core Library Development

### Step 1: Project Initialization
1. **Initialize Go module**
   ```bash
   go mod init github.com/your-username/website-checker
   ```

2. **Create basic project structure**
   ```bash
   mkdir -p cmd/{cli,server} pkg/{checker,api,report,models} internal/{config,server,storage} web docs examples scripts
   ```

3. **Setup dependencies**
   ```go
   // Key dependencies to add:
   - golang.org/x/net/html        // HTML parsing
   - github.com/spf13/cobra       // CLI framework
   - github.com/gin-gonic/gin     // HTTP server
   - github.com/spf13/viper       // Configuration
   - github.com/fatih/color       // Terminal colors
   - github.com/olekukonko/tablewriter // Table formatting
   ```

### Step 2: Core Data Models
Create `pkg/models/types.go`:
```go
type CheckResult struct {
    Name        string    `json:"name"`
    Status      Status    `json:"status"`
    Message     string    `json:"message"`
    Details     string    `json:"details,omitempty"`
    Timestamp   time.Time `json:"timestamp"`
}

type Status string
const (
    StatusPass    Status = "pass"    // ✅
    StatusWarning Status = "warning" // 🟡
    StatusFail    Status = "fail"    // ❌
)

type WebsiteReport struct {
    URL          string        `json:"url"`
    Timestamp    time.Time     `json:"timestamp"`
    Duration     time.Duration `json:"duration"`
    Results      []CheckResult `json:"results"`
    OverallScore int           `json:"overall_score"`
}
```

### Step 3: Individual Checkers Implementation

#### 3.1 Robots.txt Checker (`pkg/checker/robots.go`)
```go
func CheckRobotsTxt(baseURL string) CheckResult {
    // 1. Construct robots.txt URL
    // 2. Make HTTP GET request
    // 3. Check status code
    // 4. Return appropriate result
}
```

#### 3.2 Sitemap Checker (`pkg/checker/sitemap.go`)
```go
func CheckSitemap(baseURL string, robotsContent string) CheckResult {
    // 1. Parse robots.txt for sitemap location
    // 2. Try default sitemap.xml location
    // 3. Validate sitemap accessibility
    // 4. Return result
}
```

#### 3.3 Favicon Checker (`pkg/checker/favicon.go`)
```go
func CheckFavicon(htmlContent string) CheckResult {
    // 1. Parse HTML content
    // 2. Look for link[rel="icon"] or link[rel="shortcut icon"]
    // 3. Check for modern variants (apple-touch-icon)
    // 4. Return result
}
```

#### 3.4 SEO Metadata Checker (`pkg/checker/seo.go`)
```go
func CheckSEOMetadata(htmlContent string) []CheckResult {
    // 1. Parse HTML
    // 2. Check for <title> tag
    // 3. Check for meta description
    // 4. Check for Open Graph tags
    // 5. Validate content length and quality
    // 6. Return multiple results
}
```

#### 3.5 Performance Checker (`pkg/checker/performance.go`)
```go
func CheckPerformance(url string, apiKey string) CheckResult {
    // 1. Call Google PageSpeed Insights API
    // 2. Parse Core Web Vitals
    // 3. Extract performance score
    // 4. Return result with recommendations
}
```

#### 3.6 Security Headers Checker (`pkg/checker/security.go`)
```go
func CheckSecurityHeaders(url string) []CheckResult {
    // 1. Make HEAD request
    // 2. Check for security headers:
    //    - Strict-Transport-Security
    //    - Content-Security-Policy
    //    - X-Frame-Options
    //    - X-Content-Type-Options
    // 3. Return results for each header
}
```

#### 3.7 Accessibility Checker (`pkg/checker/accessibility.go`)
```go
func CheckAccessibility(url string) CheckResult {
    // 1. Integrate with axe-core API or similar
    // 2. Run accessibility audit
    // 3. Count critical violations
    // 4. Return summary result
}
```

### Step 4: Main Checker Orchestrator (`pkg/checker/checker.go`)
```go
type Checker struct {
    Config Config
}

func (c *Checker) CheckWebsite(url string) (*WebsiteReport, error) {
    // 1. Validate URL
    // 2. Fetch HTML content
    // 3. Run all checks concurrently
    // 4. Aggregate results
    // 5. Calculate overall score
    // 6. Return comprehensive report
}
```

## Phase 2: CLI Application Development

### Step 5: CLI Implementation (`cmd/cli/main.go`)

#### 5.1 Cobra CLI Structure
```go
// Root command
var rootCmd = &cobra.Command{
    Use:   "webcheck",
    Short: "A comprehensive website checker",
    Long:  "Check websites for SEO, performance, security, and accessibility issues",
}

// Check command
var checkCmd = &cobra.Command{
    Use:   "check [URL]",
    Short: "Check a website",
    Args:  cobra.ExactArgs(1),
    Run:   runCheck,
}

// Subcommands for specific checks
var seoCmd = &cobra.Command{Use: "seo [URL]"}
var perfCmd = &cobra.Command{Use: "performance [URL]"}
var secCmd = &cobra.Command{Use: "security [URL]"}
var a11yCmd = &cobra.Command{Use: "accessibility [URL]"}
```

#### 5.2 CLI Flags and Options
```go
var (
    outputFormat string  // json, table, html
    configFile   string  // custom config file
    verbose      bool    // detailed output
    apiKey       string  // for external APIs
    timeout      int     // request timeout
    parallel     bool    // run checks in parallel
)
```

#### 5.3 Output Formatting (`pkg/report/`)
- **Console output**: Colored, formatted table
- **JSON output**: Machine-readable format
- **HTML output**: Detailed report with recommendations

## Phase 3: Backend Service Development

### Step 6: HTTP API Server (`cmd/server/main.go`)

#### 6.1 API Endpoints
```go
// RESTful API endpoints
POST /api/v1/check              // Submit URL for checking
GET  /api/v1/check/{id}         // Get check result
GET  /api/v1/check/{id}/report  // Get formatted report
GET  /api/v1/health             // Health check
GET  /api/v1/metrics            // Prometheus metrics
```

#### 6.2 Service Architecture
```go
type Service struct {
    checker   *checker.Checker
    storage   storage.Storage  // Redis/PostgreSQL for results
    queue     queue.Queue      // Background job processing
}
```

#### 6.3 Background Job Processing
- Use Redis or database queue for async processing
- Support webhook notifications when checks complete
- Rate limiting and API key management

## Phase 4: Advanced Features

### Step 7: Configuration Management
```yaml
# config.yaml
api:
  pagespeed_key: "your-api-key"
  accessibility_endpoint: "https://axe-api.example.com"

checks:
  enabled:
    - robots
    - sitemap
    - favicon
    - seo
    - performance
    - security
    - accessibility
  
timeout: 30s
parallel: true

output:
  format: table  # table, json, html
  verbose: false
```

### Step 8: Error Handling and Resilience
- Retry logic for external API calls
- Graceful degradation when services are unavailable
- Comprehensive error messages and logging

### Step 9: Testing Strategy
```
tests/
├── unit/           # Unit tests for each checker
├── integration/    # Integration tests with real websites
├── fixtures/       # Test HTML files and responses
└── benchmarks/     # Performance benchmarks
```

### Step 10: Documentation and Examples
- API documentation (OpenAPI/Swagger)
- CLI usage examples
- Library integration examples
- Best practices guide

## Implementation Steps Summary

1. **Week 1**: Core library development (Steps 1-4)
   - Project setup and dependencies
   - Implement individual checkers
   - Create main orchestrator

2. **Week 2**: CLI application (Step 5)
   - Cobra CLI implementation
   - Output formatting
   - Configuration management

3. **Week 3**: Backend service (Step 6)
   - HTTP API server
   - Async job processing
   - Database integration

4. **Week 4**: Polish and advanced features (Steps 7-10)
   - Error handling
   - Testing
   - Documentation
   - Performance optimization

## Build and Deployment

### Makefile
```makefile
.PHONY: build test clean install

build:
	go build -o bin/webcheck cmd/cli/main.go
	go build -o bin/webcheck-server cmd/server/main.go

test:
	go test ./...

install:
	go install cmd/cli/main.go

docker:
	docker build -t webcheck .

release:
	goreleaser release --rm-dist
```

### Docker Support
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o webcheck cmd/cli/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/webcheck /usr/local/bin/
CMD ["webcheck"]
```

This plan provides a comprehensive roadmap for building a professional-grade website checker that can serve as both a library and standalone applications.