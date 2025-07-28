package checker

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/checkly-go/checkly/pkg/models"
)

// SecurityHeaders holds extracted security-related headers
type SecurityHeaders struct {
	StrictTransportSecurity string
	ContentSecurityPolicy   string
	XFrameOptions           string
	XContentTypeOptions     string
	ReferrerPolicy          string
	XSSProtection           string
	PermissionsPolicy       string
	Server                  string
	XPoweredBy              string
}

// CheckSecurityHeaders checks URL for security headers and returns multiple results
func CheckSecurityHeaders(url string) []models.CheckResult {
	start := time.Now()
	var results []models.CheckResult

	// Make HEAD request to get headers
	resp, err := http.Head(url)
	if err != nil {
		return []models.CheckResult{{
			Name:      "Security Headers",
			Status:    models.StatusFail,
			Message:   "Failed to fetch headers",
			Details:   err.Error(),
			Timestamp: start,
		}}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return []models.CheckResult{{
			Name:      "Security Headers",
			Status:    models.StatusFail,
			Message:   fmt.Sprintf("HTTP %d response", resp.StatusCode),
			Details:   fmt.Sprintf("Unable to check headers for %s", url),
			Timestamp: start,
		}}
	}

	// Extract security headers
	headers := extractSecurityHeaders(resp.Header)

	// Check individual security headers
	results = append(results, checkStrictTransportSecurity(headers.StrictTransportSecurity, start))
	results = append(results, checkContentSecurityPolicy(headers.ContentSecurityPolicy, start))
	results = append(results, checkXFrameOptions(headers.XFrameOptions, start))
	results = append(results, checkXContentTypeOptions(headers.XContentTypeOptions, start))
	results = append(results, checkReferrerPolicy(headers.ReferrerPolicy, start))
	results = append(results, checkXSSProtection(headers.XSSProtection, start))
	results = append(results, checkInformationDisclosure(headers, start))

	return results
}

// extractSecurityHeaders extracts security-related headers from HTTP response
func extractSecurityHeaders(httpHeaders http.Header) SecurityHeaders {
	return SecurityHeaders{
		StrictTransportSecurity: httpHeaders.Get("Strict-Transport-Security"),
		ContentSecurityPolicy:   httpHeaders.Get("Content-Security-Policy"),
		XFrameOptions:           httpHeaders.Get("X-Frame-Options"),
		XContentTypeOptions:     httpHeaders.Get("X-Content-Type-Options"),
		ReferrerPolicy:          httpHeaders.Get("Referrer-Policy"),
		XSSProtection:           httpHeaders.Get("X-XSS-Protection"),
		PermissionsPolicy:       httpHeaders.Get("Permissions-Policy"),
		Server:                  httpHeaders.Get("Server"),
		XPoweredBy:              httpHeaders.Get("X-Powered-By"),
	}
}

// checkStrictTransportSecurity validates HSTS header
func checkStrictTransportSecurity(hsts string, timestamp time.Time) models.CheckResult {
	if hsts == "" {
		return models.CheckResult{
			Name:      "HSTS (Strict-Transport-Security)",
			Status:    models.StatusFail,
			Message:   "Missing HSTS header",
			Details:   "Add Strict-Transport-Security header to enforce HTTPS connections",
			Timestamp: timestamp,
		}
	}

	// Parse max-age value
	var maxAge int
	var includeSubDomains bool
	var preload bool

	parts := strings.Split(hsts, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "max-age=") {
			if age, err := strconv.Atoi(part[8:]); err == nil {
				maxAge = age
			}
		} else if part == "includeSubDomains" {
			includeSubDomains = true
		} else if part == "preload" {
			preload = true
		}
	}

	var warnings []string
	if maxAge < 31536000 { // 1 year
		warnings = append(warnings, "max-age should be at least 1 year (31536000)")
	}
	if !includeSubDomains {
		warnings = append(warnings, "consider adding includeSubDomains")
	}
	if !preload {
		warnings = append(warnings, "consider adding preload")
	}

	status := models.StatusPass
	message := "HSTS header present"
	details := fmt.Sprintf("max-age=%d", maxAge)

	if includeSubDomains {
		details += ", includeSubDomains"
	}
	if preload {
		details += ", preload"
	}

	if len(warnings) > 0 {
		status = models.StatusWarning
		message = "HSTS header present but could be improved"
		details += ". Recommendations: " + strings.Join(warnings, "; ")
	}

	return models.CheckResult{
		Name:      "HSTS (Strict-Transport-Security)",
		Status:    status,
		Message:   message,
		Details:   details,
		Timestamp: timestamp,
	}
}

// checkContentSecurityPolicy validates CSP header
func checkContentSecurityPolicy(csp string, timestamp time.Time) models.CheckResult {
	if csp == "" {
		return models.CheckResult{
			Name:      "Content Security Policy",
			Status:    models.StatusFail,
			Message:   "Missing CSP header",
			Details:   "Add Content-Security-Policy header to prevent XSS attacks",
			Timestamp: timestamp,
		}
	}

	// Basic CSP validation
	var issues []string
	var good []string

	cspLower := strings.ToLower(csp)

	// Check for dangerous directives
	if strings.Contains(cspLower, "'unsafe-inline'") {
		issues = append(issues, "contains 'unsafe-inline'")
	}
	if strings.Contains(cspLower, "'unsafe-eval'") {
		issues = append(issues, "contains 'unsafe-eval'")
	}
	if strings.Contains(cspLower, "*") && !strings.Contains(cspLower, "data:") {
		issues = append(issues, "contains wildcard (*)")
	}

	// Check for good directives
	if strings.Contains(cspLower, "default-src") {
		good = append(good, "default-src")
	}
	if strings.Contains(cspLower, "script-src") {
		good = append(good, "script-src")
	}
	if strings.Contains(cspLower, "object-src") {
		good = append(good, "object-src")
	}

	status := models.StatusPass
	message := "CSP header present"
	details := fmt.Sprintf("Directives: %s", strings.Join(good, ", "))

	if len(issues) > 0 {
		status = models.StatusWarning
		message = "CSP header present but has potential issues"
		details += ". Issues: " + strings.Join(issues, ", ")
	}

	return models.CheckResult{
		Name:      "Content Security Policy",
		Status:    status,
		Message:   message,
		Details:   details,
		Timestamp: timestamp,
	}
}

// checkXFrameOptions validates X-Frame-Options header
func checkXFrameOptions(xfo string, timestamp time.Time) models.CheckResult {
	if xfo == "" {
		return models.CheckResult{
			Name:      "X-Frame-Options",
			Status:    models.StatusFail,
			Message:   "Missing X-Frame-Options header",
			Details:   "Add X-Frame-Options header to prevent clickjacking attacks",
			Timestamp: timestamp,
		}
	}

	xfoLower := strings.ToLower(strings.TrimSpace(xfo))

	switch xfoLower {
	case "deny":
		return models.CheckResult{
			Name:      "X-Frame-Options",
			Status:    models.StatusPass,
			Message:   "X-Frame-Options properly configured",
			Details:   "Set to DENY - provides maximum protection",
			Timestamp: timestamp,
		}
	case "sameorigin":
		return models.CheckResult{
			Name:      "X-Frame-Options",
			Status:    models.StatusPass,
			Message:   "X-Frame-Options properly configured",
			Details:   "Set to SAMEORIGIN - allows framing by same origin",
			Timestamp: timestamp,
		}
	default:
		if strings.HasPrefix(xfoLower, "allow-from") {
			return models.CheckResult{
				Name:      "X-Frame-Options",
				Status:    models.StatusWarning,
				Message:   "X-Frame-Options uses deprecated ALLOW-FROM",
				Details:   "ALLOW-FROM is deprecated, consider using CSP frame-ancestors instead",
				Timestamp: timestamp,
			}
		}
		return models.CheckResult{
			Name:      "X-Frame-Options",
			Status:    models.StatusWarning,
			Message:   "X-Frame-Options has invalid value",
			Details:   fmt.Sprintf("Invalid value: %s. Use DENY or SAMEORIGIN", xfo),
			Timestamp: timestamp,
		}
	}
}

// checkXContentTypeOptions validates X-Content-Type-Options header
func checkXContentTypeOptions(xcto string, timestamp time.Time) models.CheckResult {
	if xcto == "" {
		return models.CheckResult{
			Name:      "X-Content-Type-Options",
			Status:    models.StatusFail,
			Message:   "Missing X-Content-Type-Options header",
			Details:   "Add X-Content-Type-Options: nosniff to prevent MIME type sniffing",
			Timestamp: timestamp,
		}
	}

	xctoLower := strings.ToLower(strings.TrimSpace(xcto))

	if xctoLower == "nosniff" {
		return models.CheckResult{
			Name:      "X-Content-Type-Options",
			Status:    models.StatusPass,
			Message:   "X-Content-Type-Options properly configured",
			Details:   "Set to nosniff - prevents MIME type sniffing attacks",
			Timestamp: timestamp,
		}
	}

	return models.CheckResult{
		Name:      "X-Content-Type-Options",
		Status:    models.StatusWarning,
		Message:   "X-Content-Type-Options has invalid value",
		Details:   fmt.Sprintf("Invalid value: %s. Should be 'nosniff'", xcto),
		Timestamp: timestamp,
	}
}

// checkReferrerPolicy validates Referrer-Policy header
func checkReferrerPolicy(rp string, timestamp time.Time) models.CheckResult {
	if rp == "" {
		return models.CheckResult{
			Name:      "Referrer Policy",
			Status:    models.StatusWarning,
			Message:   "Missing Referrer-Policy header",
			Details:   "Consider adding Referrer-Policy header to control referrer information",
			Timestamp: timestamp,
		}
	}

	rpLower := strings.ToLower(strings.TrimSpace(rp))

	secureValues := []string{
		"no-referrer",
		"no-referrer-when-downgrade",
		"origin",
		"origin-when-cross-origin",
		"same-origin",
		"strict-origin",
		"strict-origin-when-cross-origin",
	}

	for _, secure := range secureValues {
		if rpLower == secure {
			return models.CheckResult{
				Name:      "Referrer Policy",
				Status:    models.StatusPass,
				Message:   "Referrer-Policy properly configured",
				Details:   fmt.Sprintf("Set to %s", rp),
				Timestamp: timestamp,
			}
		}
	}

	return models.CheckResult{
		Name:      "Referrer Policy",
		Status:    models.StatusWarning,
		Message:   "Referrer-Policy has potentially unsafe value",
		Details:   fmt.Sprintf("Value: %s. Consider using stricter policy", rp),
		Timestamp: timestamp,
	}
}

// checkXSSProtection validates X-XSS-Protection header
func checkXSSProtection(xss string, timestamp time.Time) models.CheckResult {
	if xss == "" {
		return models.CheckResult{
			Name:      "X-XSS-Protection",
			Status:    models.StatusWarning,
			Message:   "Missing X-XSS-Protection header",
			Details:   "Consider adding X-XSS-Protection header (though CSP is preferred)",
			Timestamp: timestamp,
		}
	}

	xssLower := strings.ToLower(strings.TrimSpace(xss))

	if xssLower == "1; mode=block" || xssLower == "1;mode=block" {
		return models.CheckResult{
			Name:      "X-XSS-Protection",
			Status:    models.StatusPass,
			Message:   "X-XSS-Protection properly configured",
			Details:   "Set to block XSS attacks (CSP is still preferred)",
			Timestamp: timestamp,
		}
	}

	if xssLower == "0" {
		return models.CheckResult{
			Name:      "X-XSS-Protection",
			Status:    models.StatusWarning,
			Message:   "X-XSS-Protection disabled",
			Details:   "XSS protection is disabled. Ensure strong CSP is in place",
			Timestamp: timestamp,
		}
	}

	return models.CheckResult{
		Name:      "X-XSS-Protection",
		Status:    models.StatusWarning,
		Message:   "X-XSS-Protection misconfigured",
		Details:   fmt.Sprintf("Value: %s. Recommended: '1; mode=block'", xss),
		Timestamp: timestamp,
	}
}

// checkInformationDisclosure checks for information disclosure in headers
func checkInformationDisclosure(headers SecurityHeaders, timestamp time.Time) models.CheckResult {
	var issues []string
	var good []string

	// Check Server header
	if headers.Server != "" {
		if strings.Contains(strings.ToLower(headers.Server), "apache") ||
			strings.Contains(strings.ToLower(headers.Server), "nginx") ||
			strings.Contains(strings.ToLower(headers.Server), "iis") {
			issues = append(issues, fmt.Sprintf("Server header reveals software: %s", headers.Server))
		}
	} else {
		good = append(good, "Server header hidden")
	}

	// Check X-Powered-By header
	if headers.XPoweredBy != "" {
		issues = append(issues, fmt.Sprintf("X-Powered-By header reveals technology: %s", headers.XPoweredBy))
	} else {
		good = append(good, "X-Powered-By header hidden")
	}

	status := models.StatusPass
	message := "No information disclosure detected"
	details := strings.Join(good, ", ")

	if len(issues) > 0 {
		status = models.StatusWarning
		message = "Potential information disclosure"
		if details != "" {
			details += ". "
		}
		details += "Issues: " + strings.Join(issues, ", ")
	}

	return models.CheckResult{
		Name:      "Information Disclosure",
		Status:    status,
		Message:   message,
		Details:   details,
		Timestamp: timestamp,
	}
}
