package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/checkly-go/checkly/pkg/checker"
	"github.com/checkly-go/checkly/pkg/models"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	appStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type state int

const (
	inputView state = iota
	checkerView
	progressView
	resultsView
)

type model struct {
	state        state
	url          string
	checkers     map[string]bool
	cursor       int
	results      map[string][]models.CheckResult
	progress     float64
	currentCheck string
	err          error
}

func initialModel() model {
	return model{
		state: inputView,
		checkers: map[string]bool{
			"robots":   true,
			"sitemap":  true,
			"seo":      true,
			"security": true,
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case inputView:
				if m.url != "" {
					m.state = checkerView
				}
			case checkerView:
				m.state = progressView
				return m, m.runChecks()
			case resultsView:
				return m, tea.Quit
			}
		case "backspace":
			if m.state == inputView && len(m.url) > 0 {
				m.url = m.url[:len(m.url)-1]
			}
		case " ":
			if m.state == checkerView {
				checkerNames := []string{"robots", "sitemap", "seo", "security"}
				if m.cursor < len(checkerNames) {
					checker := checkerNames[m.cursor]
					m.checkers[checker] = !m.checkers[checker]
				}
			}
		case "up", "k":
			if m.state == checkerView && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == checkerView && m.cursor < 3 {
				m.cursor++
			}
		default:
			if m.state == inputView {
				m.url += msg.String()
			}
		}
	case progressMsg:
		m.progress = msg.progress
		m.currentCheck = msg.current
		if msg.progress >= 1.0 {
			m.state = resultsView
		}
	case resultsMsg:
		m.results = msg.results
	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m model) View() string {
	var s string

	switch m.state {
	case inputView:
		s = m.inputView()
	case checkerView:
		s = m.checkerView()
	case progressView:
		s = m.progressView()
	case resultsView:
		s = m.resultsView()
	}

	return appStyle.Render(s)
}

func (m model) inputView() string {
	title := titleStyle.Render("🌐 Website Checker")

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(50).
		Render(m.url + "│")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Enter the website URL to check • Press Enter to continue • Ctrl+C to quit")

	return fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s",
		title,
		"Enter website URL:",
		input,
		help,
	)
}

func (m model) checkerView() string {
	title := titleStyle.Render("🔧 Select Checkers")

	var choices []string
	checkerNames := []string{"robots", "sitemap", "seo", "security"}
	checkerLabels := map[string]string{
		"robots":   "📋 Robots.txt Check",
		"sitemap":  "🗺️  Sitemap Check",
		"seo":      "🏷️  SEO Metadata",
		"security": "🛡️  Security Headers",
	}

	for i, checker := range checkerNames {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.checkers[checker] {
			checked = "✓"
		}

		choices = append(choices, fmt.Sprintf("%s [%s] %s", cursor, checked, checkerLabels[checker]))
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓ to move • Space to toggle • Enter to start • Ctrl+C to quit")

	return fmt.Sprintf(
		"%s\n\nURL: %s\n\n%s\n\n%s",
		title,
		lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(m.url),
		lipgloss.JoinVertical(lipgloss.Left, choices...),
		help,
	)
}

func (m model) progressView() string {
	title := titleStyle.Render("⚡ Running Checks")

	progressBar := m.renderProgressBar()

	current := ""
	if m.currentCheck != "" {
		current = fmt.Sprintf("Currently checking: %s", m.currentCheck)
	}

	return fmt.Sprintf(
		"%s\n\nURL: %s\n\n%s\n\n%s",
		title,
		lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(m.url),
		progressBar,
		current,
	)
}

func (m model) resultsView() string {
	title := titleStyle.Render("✨ Results")

	if m.err != nil {
		return fmt.Sprintf("%s\n\nError: %v\n\nPress Enter to exit", title, m.err)
	}

	var results []string
	for checker, checkResults := range m.results {
		checkerTitle := fmt.Sprintf("📊 %s Results:", checker)
		results = append(results, checkerTitle)

		for _, result := range checkResults {
			statusEmoji := getStatusEmoji(result.Status)
			resultLine := fmt.Sprintf("  %s %s: %s", statusEmoji, result.Name, result.Message)
			results = append(results, resultLine)
		}
		results = append(results, "")
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Press Enter to exit • Ctrl+C to quit")

	return fmt.Sprintf(
		"%s\n\nURL: %s\n\n%s\n\n%s",
		title,
		lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(m.url),
		lipgloss.JoinVertical(lipgloss.Left, results...),
		help,
	)
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

func (m model) renderProgressBar() string {
	w := 50
	filled := int(m.progress * float64(w))

	bar := ""
	for i := 0; i < w; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	percentage := int(m.progress * 100)

	return fmt.Sprintf("[%s] %d%%", bar, percentage)
}

type progressMsg struct {
	progress float64
	current  string
}

type resultsMsg struct {
	results map[string][]models.CheckResult
}

type errMsg struct {
	err error
}

func (m model) runChecks() tea.Cmd {
	return func() tea.Msg {
		results := make(map[string][]models.CheckResult)

		enabledCheckers := []string{}
		for checker, enabled := range m.checkers {
			if enabled {
				enabledCheckers = append(enabledCheckers, checker)
			}
		}

		totalCheckers := len(enabledCheckers)
		for i, checkerName := range enabledCheckers {
			time.Sleep(500 * time.Millisecond)

			switch checkerName {
			case "robots":
				result := checker.CheckRobotsTxt(m.url)
				results["robots"] = []models.CheckResult{result}
			case "sitemap":
				result := checker.CheckSitemapWithRobotsURL(m.url)
				results["sitemap"] = []models.CheckResult{result}
			case "seo":
				seoResults := checker.CheckSEOMetadataFromURL(m.url)
				results["seo"] = seoResults
			case "security":
				securityResults := checker.CheckSecurityHeaders(m.url)
				results["security"] = securityResults
			}

			if i == totalCheckers-1 {
				break
			}
		}

		return resultsMsg{results: results}
	}
}
