# Website Checker 🚀

A comprehensive website analysis tool that evaluates websites across multiple dimensions including SEO, security, robots.txt compliance, and sitemap validation. Built with Go and powered by AI-driven recommendations.

## 🌟 Features

### Core Analysis Capabilities
- **🤖 Robots.txt Validation** - Check robots.txt file existence, accessibility, and syntax
- **🗺️ Sitemap Analysis** - Validate XML sitemaps and their discoverability
- **🏷️ SEO Metadata Assessment** - Analyze title tags, meta descriptions, heading structure
- **🛡️ Security Headers Audit** - Verify essential security headers implementation

### Multiple Interfaces
- **📱 Interactive TUI** - Beautiful terminal user interface with real-time progress *(to be completed)*
- **⚡ Command Line** - Fast CLI for automation and scripting
- **🌐 REST API** - HTTP API for integration with other systems
- **🤖 AI Recommendations** - Powered by Google Gemini for actionable insights

### Output Formats
- **Human-readable text reports** with emoji status indicators
- **Structured JSON output** for programmatic processing
- **File export capabilities** for report storage
- **Real-time progress visualization** in TUI mode *(to be completed)*

## 🚀 Quick Start

### Installation

#### From Source
```bash
git clone https://github.com/checkly-go/checkly.git
cd checkly
go mod download
go build -o checkly ./cmd/
```

#### Build All Components
```bash
# Build CLI tool
go build -o checkly ./cmd/

# Build TUI interface (to be completed)
go build -o checkly-tui ./cmd/tui/

# Build API server
go build -o server ./cmd/server/
```

### Basic Usage

#### CLI Mode
```bash
# Quick check of a website
./checkly -url https://example.com

# Specific checks only
./checkly -url https://example.com -checkers robots,seo

# JSON output to file
./checkly -url https://example.com -output json -o report.json

# Custom checker selection
./checkly -url https://example.com -checkers security,sitemap -output text
```

#### Interactive TUI Mode *(to be completed)*
```bash
# Launch beautiful terminal interface
./checkly -tui

# Or run TUI directly
./checkly-tui
```

#### API Server Mode
```bash
# Start the REST API server
./server

# Server runs on http://localhost:8080
```

## 📋 Available Checks

| Check Type | Description | Status Indicators |
|------------|-------------|------------------|
| **Robots.txt** | Validates robots.txt file existence, accessibility, and syntax | ✅ Found & Valid / 🟡 Issues / ❌ Missing |
| **Sitemap** | Checks XML sitemap presence and discoverability via robots.txt | ✅ Found / 🟡 Partial / ❌ Missing |
| **SEO Metadata** | Analyzes title tags, meta descriptions, heading structure | ✅ Optimized / 🟡 Needs Work / ❌ Missing |
| **Security Headers** | Audits security headers (HSTS, CSP, X-Frame-Options, etc.) | ✅ Secure / 🟡 Partial / ❌ Vulnerable |

## 🎯 Usage Examples

### Command Line Interface

```bash
# Complete website audit
./checkly -url https://mywebsite.com

# Security-focused check
./checkly -url https://mywebsite.com -checkers security

# SEO analysis with JSON export
./checkly -url https://mywebsite.com -checkers seo -output json -o seo-report.json

# Multiple checks with text output
./checkly -url https://mywebsite.com -checkers robots,sitemap,seo,security -output text
```

### TUI Interface Navigation *(to be completed)*

```bash
./checkly -tui
```

**Controls:**
- **Type** to enter website URL
- **Enter** to proceed to next step  
- **↑/↓** or **j/k** to navigate options
- **Space** to toggle checker selection
- **Ctrl+C** or **q** to quit

### API Usage

#### Start the Server
```bash
# Configure environment (optional)
export MONGO_URI="mongodb+srv://user:pass@cluster.mongodb.net/"
export GEMINI_API_KEY="your_gemini_api_key"

# Start server
./server
```

#### API Endpoints

```bash
# Submit a website check
curl -X POST http://localhost:8080/api/v1/check \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Get check results
curl http://localhost:8080/api/v1/check/{check-id}

# Get detailed report
curl http://localhost:8080/api/v1/check/{check-id}/report

# Get AI-powered recommendations
curl -X POST http://localhost:8080/api/v1/recommend \
  -H "Content-Type: application/json" \
  -d '{"check_id": "check-id", "focus": ["seo", "security"]}'

# Health check
curl http://localhost:8080/api/v1/health
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```bash
# MongoDB connection (for API server)
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority

# Google Gemini AI API Key (for recommendations)
GEMINI_API_KEY=your_gemini_api_key_here

# Server port (optional, defaults to 8080)
PORT=8080
```

### Command Line Options

```bash
Usage: checkly [options]

Options:
  -url string
        URL to check (required)
  -link string
        URL to check (alias for -url)
  -tui
        Run in TUI mode (interactive terminal UI) [to be completed]
  -checkers string
        Comma-separated list of checkers to run (default "robots,sitemap,seo,security")
        Options: robots, sitemap, seo, security
  -output string
        Output format (text or json) (default "text")
  -o string
        Output file path (for JSON reports)

Examples:
  checkly -url https://example.com
  checkly -tui                    # (to be completed)
  checkly -link https://example.com -checkers robots,seo -output json
  checkly -url https://example.com -output json -o report.json
  checkly -url https://example.com -checkers security -output text
```

## 🏗️ Architecture

### Project Structure

```
checkly/
├── cmd/                          # Application entry points
│   ├── main.go                   # Main CLI application
│   ├── server/main.go            # REST API server
│   └── tui/main.go              # Terminal UI application (to be completed)
├── internal/                     # Private application code
│   ├── handlers/                 # HTTP request handlers
│   │   ├── service.go           # Main service handlers
│   │   └── recommendation.go    # AI recommendation handlers
│   └── storage/                  # Database layer
│       └── mongo.go             # MongoDB implementation
├── pkg/                         # Public library code
│   ├── ai/                      # AI integration
│   │   └── gemini.go           # Google Gemini client
│   ├── checker/                 # Core checking logic
│   │   ├── checker.go          # Main checker orchestrator
│   │   ├── robots.go           # Robots.txt validation
│   │   ├── sitemap.go          # Sitemap analysis
│   │   ├── seo.go              # SEO metadata checks
│   │   └── security.go         # Security headers audit
│   ├── models/                  # Data models
│   │   └── types.go            # Shared types and structures
│   └── report/                  # Report generation
│       ├── json.go             # JSON report formatter
│       └── score.go            # Scoring algorithms
├── .env.example                 # Environment configuration template
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Dockerfile                   # Container configuration
└── README.md                    # This file
```

### Core Components

#### 1. Checker Engine (`pkg/checker/`)
The heart of the application that orchestrates all website analysis:
- **checker.go**: Main coordinator that runs all checks
- **robots.go**: Validates robots.txt files
- **sitemap.go**: Analyzes XML sitemaps  
- **seo.go**: Evaluates SEO metadata
- **security.go**: Audits security headers

#### 2. AI Integration (`pkg/ai/`)
- **gemini.go**: Google Gemini integration for intelligent recommendations
- Analyzes check results and provides actionable insights
- Generates prioritized improvement suggestions

#### 3. Storage Layer (`internal/storage/`)
- **mongo.go**: MongoDB implementation for persistent storage
- Stores check results, reports, and user data
- Supports future authentication and user management

#### 4. API Layer (`internal/handlers/`)
- **service.go**: Core HTTP handlers for check operations
- **recommendation.go**: AI-powered recommendation endpoints
- RESTful API design with JSON responses

#### 5. Report Generation (`pkg/report/`)
- **json.go**: Structured JSON report generation
- **score.go**: Overall scoring algorithms
- Extensible format support

### Data Models

#### CheckResult
```go
type CheckResult struct {
    Name      string    `json:"name"`
    Status    Status    `json:"status"`     // pass, warning, fail
    Message   string    `json:"message"`
    Details   string    `json:"details,omitempty"`
    Timestamp time.Time `json:"timestamp"`
}
```

#### WebsiteReport
```go
type WebsiteReport struct {
    URL          string        `json:"url"`
    Timestamp    time.Time     `json:"timestamp"`
    Duration     time.Duration `json:"duration"`
    Results      []CheckResult `json:"results"`
    OverallScore int           `json:"overall_score"`
}
```

## 🤖 AI-Powered Recommendations

The tool integrates with Google Gemini to provide intelligent, actionable recommendations based on check results.

### Features
- **Contextual Analysis**: Understands the impact of each issue
- **Prioritized Suggestions**: Ranks recommendations by importance
- **Specific Guidance**: Provides concrete steps to resolve issues
- **Category Grouping**: Organizes suggestions by area (SEO, Security, etc.)

### Example Request
```bash
curl -X POST http://localhost:8080/api/v1/recommend \
  -H "Content-Type: application/json" \
  -d '{
    "check_id": "507f1f77bcf86cd799439011",
    "focus": ["seo", "security"]
  }'
```

### Example Response
```json
{
  "url": "https://example.com",
  "generated_at": "2024-01-15T10:30:00Z",
  "summary": "Critical security headers missing, SEO metadata needs optimization",
  "recommendations": [
    {
      "category": "security",
      "priority": "high", 
      "issues": [
        {
          "issue": "Missing Content Security Policy",
          "impact": "Vulnerable to XSS attacks",
          "current_status": "No CSP header detected"
        }
      ],
      "improvements": [
        "Add Content-Security-Policy header",
        "Configure CSP to restrict resource loading"
      ]
    }
  ]
}
```

## 🐳 Docker Deployment

### Build Container
```bash
docker build -t checkly .
```

### Run CLI in Container
```bash
docker run --rm checkly -url https://example.com
```

### Run API Server in Container
```bash
docker run -p 8080:8080 \
  -e MONGO_URI="your_mongo_connection" \
  -e GEMINI_API_KEY="your_api_key" \
  checkly-server
```

## 🧪 Testing

### Manual Testing Script
```bash
# Test recommendation endpoint
./test_recommend.sh
```

### Running Checks
```bash
# Test different websites
./checkly -url https://google.com
./checkly -url https://github.com  
./checkly -url https://stackoverflow.com

# Test different output formats
./checkly -url https://example.com -output json
./checkly -url https://example.com -output text

# Test specific checkers
./checkly -url https://example.com -checkers security
./checkly -url https://example.com -checkers robots,sitemap
```

## 🔧 Development

### Prerequisites
- Go 1.24.1 or later
- MongoDB (for API server)
- Google Gemini API key (for AI recommendations)

### Development Setup
```bash
# Clone repository
git clone https://github.com/checkly-go/checkly.git
cd checkly

# Install dependencies
go mod download

# Run tests
go test ./...

# Run locally
go run ./cmd/ -url https://example.com

# Run TUI locally (to be completed)
go run ./cmd/tui/

# Run API server locally
go run ./cmd/server/
```

### Adding New Checkers

1. **Create checker file** in `pkg/checker/`
2. **Implement check function** returning `[]models.CheckResult`
3. **Add to main checker** in `checker.go`
4. **Update CLI flags** in `cmd/main.go`
5. **Add tests** for the new checker

Example:
```go
// pkg/checker/performance.go
func CheckPerformance(url string) []models.CheckResult {
    // Implementation
    return []models.CheckResult{{
        Name: "Page Load Speed",
        Status: models.StatusPass,
        Message: "Page loads in 1.2s",
        Timestamp: time.Now(),
    }}
}
```

## 📈 Roadmap

### Planned Features
- [ ] **TUI Interface Enhancement** - Complete interactive terminal UI (in development)
- [ ] **Performance Analysis** - Page speed, Core Web Vitals
- [ ] **Accessibility Audit** - WCAG compliance checking  
- [ ] **Mobile Responsiveness** - Mobile-first design validation
- [ ] **Content Analysis** - Readability, keyword density
- [ ] **Link Validation** - Broken link detection
- [ ] **Schema Markup** - Structured data validation
- [ ] **Social Media Meta** - Open Graph, Twitter Cards
- [ ] **Analytics Integration** - Google Analytics, GTM validation
- [ ] **User Authentication** - API key management
- [ ] **Scheduled Checks** - Automated monitoring
- [ ] **Alerting System** - Email/Slack notifications
- [ ] **Historical Trending** - Track improvements over time
- [ ] **Bulk URL Processing** - Batch website analysis
- [ ] **Plugin System** - Custom checker extensions
- [ ] **Web Dashboard** - Browser-based interface

### Technical Improvements
- [ ] **Caching Layer** - Redis integration for performance
- [ ] **Rate Limiting** - API protection and fair usage
- [ ] **Metrics & Monitoring** - Prometheus/Grafana integration
- [ ] **Container Orchestration** - Kubernetes deployment
- [ ] **CI/CD Pipeline** - Automated testing and deployment
- [ ] **Load Testing** - Performance benchmarking
- [ ] **Security Scanning** - Vulnerability assessment
- [ ] **Code Coverage** - Comprehensive test coverage

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes** with tests
4. **Run the test suite** (`go test ./...`)
5. **Commit your changes** (`git commit -m 'Add amazing feature'`)
6. **Push to the branch** (`git push origin feature/amazing-feature`)
7. **Open a Pull Request**

### Code Style
- Follow Go conventions and `gofmt`
- Add comprehensive tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Bubble Tea** - Excellent TUI framework
- **Gin** - Fast HTTP web framework
- **MongoDB** - Reliable database solution
- **Google Gemini** - Powerful AI capabilities
- **Go Community** - Amazing ecosystem and tools

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/checkly-go/checkly/issues)
- **Discussions**: [GitHub Discussions](https://github.com/checkly-go/checkly/discussions)
- **Documentation**: This README and inline code comments

---

Built with ❤️ using Go. Made for developers who care about website quality.
