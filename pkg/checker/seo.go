package checker

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hawkaii/website-checker.git/pkg/models"
	"golang.org/x/net/html"
)

// SEOMetadata holds extracted SEO-related metadata
type SEOMetadata struct {
	Title           string
	MetaDescription string
	OpenGraphTitle  string
	OpenGraphDesc   string
	OpenGraphImage  string
	OpenGraphURL    string
	TwitterCard     string
	TwitterTitle    string
	TwitterDesc     string
	TwitterImage    string
	MetaKeywords    string
	Canonical       string
	MetaRobots      string
}

// CheckSEOMetadata checks HTML content for SEO metadata and returns multiple results
func CheckSEOMetadata(htmlContent string) []models.CheckResult {
	start := time.Now()
	var results []models.CheckResult

	// Parse HTML content
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return []models.CheckResult{{
			Name:      "SEO Metadata",
			Status:    models.StatusFail,
			Message:   "Failed to parse HTML",
			Details:   err.Error(),
			Timestamp: start,
		}}
	}

	metadata := extractSEOMetadata(doc)

	results = append(results, checkTitle(metadata.Title, start))

	results = append(results, checkMetaDescription(metadata.MetaDescription, start))

	results = append(results, checkOpenGraphTags(metadata, start))

	results = append(results, checkTwitterCardTags(metadata, start))

	results = append(results, checkOtherMetaTags(metadata, start))

	return results
}

// CheckSEOMetadataFromURL fetches HTML content from URL and checks SEO metadata
func CheckSEOMetadataFromURL(url string) []models.CheckResult {
	start := time.Now()

	// Fetch HTML content
	resp, err := http.Get(url)
	if err != nil {
		return []models.CheckResult{{
			Name:      "SEO Metadata",
			Status:    models.StatusFail,
			Message:   "Failed to fetch page",
			Details:   err.Error(),
			Timestamp: start,
		}}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []models.CheckResult{{
			Name:      "SEO Metadata",
			Status:    models.StatusFail,
			Message:   fmt.Sprintf("HTTP %d response", resp.StatusCode),
			Details:   fmt.Sprintf("Unable to fetch page content from %s", url),
			Timestamp: start,
		}}
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.CheckResult{{
			Name:      "SEO Metadata",
			Status:    models.StatusFail,
			Message:   "Failed to read page content",
			Details:   err.Error(),
			Timestamp: start,
		}}
	}

	return CheckSEOMetadata(string(htmlBytes))
}

// extractSEOMetadata parses HTML and extracts SEO-related metadata
func extractSEOMetadata(doc *html.Node) SEOMetadata {
	metadata := SEOMetadata{}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					metadata.Title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				extractMetaTag(n, &metadata)
			case "link":
				extractLinkTag(n, &metadata)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return metadata
}

// extractMetaTag extracts information from meta tags
func extractMetaTag(n *html.Node, metadata *SEOMetadata) {
	var name, property, content string

	for _, attr := range n.Attr {
		switch attr.Key {
		case "name":
			name = strings.ToLower(attr.Val)
		case "property":
			property = strings.ToLower(attr.Val)
		case "content":
			content = attr.Val
		}
	}

	// Standard meta tags
	switch name {
	case "description":
		metadata.MetaDescription = content
	case "keywords":
		metadata.MetaKeywords = content
	case "robots":
		metadata.MetaRobots = content
	}

	// Open Graph tags
	switch property {
	case "og:title":
		metadata.OpenGraphTitle = content
	case "og:description":
		metadata.OpenGraphDesc = content
	case "og:image":
		metadata.OpenGraphImage = content
	case "og:url":
		metadata.OpenGraphURL = content
	}

	// Twitter Card tags
	switch name {
	case "twitter:card":
		metadata.TwitterCard = content
	case "twitter:title":
		metadata.TwitterTitle = content
	case "twitter:description":
		metadata.TwitterDesc = content
	case "twitter:image":
		metadata.TwitterImage = content
	}
}

// extractLinkTag extracts information from link tags
func extractLinkTag(n *html.Node, metadata *SEOMetadata) {
	var rel, href string

	for _, attr := range n.Attr {
		switch attr.Key {
		case "rel":
			rel = strings.ToLower(attr.Val)
		case "href":
			href = attr.Val
		}
	}

	if rel == "canonical" {
		metadata.Canonical = href
	}
}

// checkTitle validates the title tag
func checkTitle(title string, timestamp time.Time) models.CheckResult {
	if title == "" {
		return models.CheckResult{
			Name:      "Title Tag",
			Status:    models.StatusFail,
			Message:   "Missing title tag",
			Details:   "Add a descriptive <title> tag to improve SEO",
			Timestamp: timestamp,
		}
	}

	length := len(title)
	if length < 30 {
		return models.CheckResult{
			Name:      "Title Tag",
			Status:    models.StatusWarning,
			Message:   "Title too short",
			Details:   fmt.Sprintf("Title is %d characters. Recommended: 30-60 characters", length),
			Timestamp: timestamp,
		}
	}

	if length > 60 {
		return models.CheckResult{
			Name:      "Title Tag",
			Status:    models.StatusWarning,
			Message:   "Title too long",
			Details:   fmt.Sprintf("Title is %d characters. May be truncated in search results", length),
			Timestamp: timestamp,
		}
	}

	return models.CheckResult{
		Name:      "Title Tag",
		Status:    models.StatusPass,
		Message:   "Title tag present and good length",
		Details:   fmt.Sprintf("Title: \"%s\" (%d characters)", title, length),
		Timestamp: timestamp,
	}
}

// checkMetaDescription validates the meta description
func checkMetaDescription(description string, timestamp time.Time) models.CheckResult {
	if description == "" {
		return models.CheckResult{
			Name:      "Meta Description",
			Status:    models.StatusFail,
			Message:   "Missing meta description",
			Details:   "Add a meta description to improve click-through rates from search results",
			Timestamp: timestamp,
		}
	}

	length := len(description)
	if length < 120 {
		return models.CheckResult{
			Name:      "Meta Description",
			Status:    models.StatusWarning,
			Message:   "Meta description too short",
			Details:   fmt.Sprintf("Description is %d characters. Recommended: 120-160 characters", length),
			Timestamp: timestamp,
		}
	}

	if length > 160 {
		return models.CheckResult{
			Name:      "Meta Description",
			Status:    models.StatusWarning,
			Message:   "Meta description too long",
			Details:   fmt.Sprintf("Description is %d characters. May be truncated in search results", length),
			Timestamp: timestamp,
		}
	}

	return models.CheckResult{
		Name:      "Meta Description",
		Status:    models.StatusPass,
		Message:   "Meta description present and good length",
		Details:   fmt.Sprintf("Description: \"%s\" (%d characters)", description, length),
		Timestamp: timestamp,
	}
}

// checkOpenGraphTags validates Open Graph metadata
func checkOpenGraphTags(metadata SEOMetadata, timestamp time.Time) models.CheckResult {
	var missing []string
	var present []string

	checks := map[string]string{
		"og:title":       metadata.OpenGraphTitle,
		"og:description": metadata.OpenGraphDesc,
		"og:image":       metadata.OpenGraphImage,
		"og:url":         metadata.OpenGraphURL,
	}

	for tag, value := range checks {
		if value == "" {
			missing = append(missing, tag)
		} else {
			present = append(present, tag)
		}
	}

	if len(missing) == 4 {
		return models.CheckResult{
			Name:      "Open Graph Tags",
			Status:    models.StatusFail,
			Message:   "No Open Graph tags found",
			Details:   "Add og:title, og:description, og:image, and og:url for better social media sharing",
			Timestamp: timestamp,
		}
	}

	if len(missing) > 0 {
		return models.CheckResult{
			Name:      "Open Graph Tags",
			Status:    models.StatusWarning,
			Message:   fmt.Sprintf("Missing %d Open Graph tags", len(missing)),
			Details:   fmt.Sprintf("Missing: %s. Present: %s", strings.Join(missing, ", "), strings.Join(present, ", ")),
			Timestamp: timestamp,
		}
	}

	return models.CheckResult{
		Name:      "Open Graph Tags",
		Status:    models.StatusPass,
		Message:   "All essential Open Graph tags present",
		Details:   fmt.Sprintf("Found: %s", strings.Join(present, ", ")),
		Timestamp: timestamp,
	}
}

// checkTwitterCardTags validates Twitter Card metadata
func checkTwitterCardTags(metadata SEOMetadata, timestamp time.Time) models.CheckResult {
	if metadata.TwitterCard == "" {
		return models.CheckResult{
			Name:      "Twitter Card",
			Status:    models.StatusWarning,
			Message:   "No Twitter Card found",
			Details:   "Add twitter:card meta tag for better Twitter sharing",
			Timestamp: timestamp,
		}
	}

	var present []string
	present = append(present, "twitter:card")

	if metadata.TwitterTitle != "" {
		present = append(present, "twitter:title")
	}
	if metadata.TwitterDesc != "" {
		present = append(present, "twitter:description")
	}
	if metadata.TwitterImage != "" {
		present = append(present, "twitter:image")
	}

	return models.CheckResult{
		Name:      "Twitter Card",
		Status:    models.StatusPass,
		Message:   fmt.Sprintf("Twitter Card configured (%s)", metadata.TwitterCard),
		Details:   fmt.Sprintf("Found tags: %s", strings.Join(present, ", ")),
		Timestamp: timestamp,
	}
}

// checkOtherMetaTags validates other important meta tags
func checkOtherMetaTags(metadata SEOMetadata, timestamp time.Time) models.CheckResult {
	var issues []string
	var good []string

	// Check canonical URL
	if metadata.Canonical != "" {
		good = append(good, "canonical URL")
	} else {
		issues = append(issues, "missing canonical URL")
	}

	// Check meta robots
	if metadata.MetaRobots != "" {
		if strings.Contains(strings.ToLower(metadata.MetaRobots), "noindex") {
			issues = append(issues, "page set to noindex")
		} else {
			good = append(good, "robots directive")
		}
	}

	// Check meta keywords (generally not recommended anymore)
	if metadata.MetaKeywords != "" {
		issues = append(issues, "meta keywords present (outdated)")
	}

	if len(issues) > 2 {
		return models.CheckResult{
			Name:      "Additional Meta Tags",
			Status:    models.StatusWarning,
			Message:   fmt.Sprintf("%d potential issues found", len(issues)),
			Details:   fmt.Sprintf("Issues: %s", strings.Join(issues, ", ")),
			Timestamp: timestamp,
		}
	}

	status := models.StatusPass
	message := "Additional meta tags look good"
	details := ""

	if len(good) > 0 {
		details = fmt.Sprintf("Found: %s", strings.Join(good, ", "))
	}

	if len(issues) > 0 {
		status = models.StatusWarning
		message = "Minor meta tag issues"
		if details != "" {
			details += ". "
		}
		details += fmt.Sprintf("Issues: %s", strings.Join(issues, ", "))
	}

	return models.CheckResult{
		Name:      "Additional Meta Tags",
		Status:    status,
		Message:   message,
		Details:   details,
		Timestamp: timestamp,
	}
}
