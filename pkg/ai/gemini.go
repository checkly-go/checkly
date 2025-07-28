package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/checkly-go/checkly/pkg/models"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(ctx context.Context) (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.7)
	model.SetTopP(0.8)
	model.SetTopK(40)
	model.SetMaxOutputTokens(2048)

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiClient) Close() {
	if g.client != nil {
		g.client.Close()
	}
}

func (g *GeminiClient) GenerateRecommendations(ctx context.Context, url string, report *models.WebsiteReport, focus []string) (*models.RecommendationResponse, error) {
	prompt := g.buildPrompt(url, report, focus)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	content := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			content += string(txt)
		}
	}

	// Extract JSON from the response
	jsonContent := g.extractJSON(content)

	// Parse the JSON response
	var result models.RecommendationResponse
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w. Raw content: %s", err, jsonContent)
	}

	return &result, nil
}

// extractJSON attempts to extract valid JSON from the AI response
func (g *GeminiClient) extractJSON(content string) string {
	content = strings.TrimSpace(content)

	// Remove markdown code blocks
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	// Find the first { and last } to extract JSON object
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")

	if start != -1 && end != -1 && end > start {
		return content[start : end+1]
	}

	return content
}

func (g *GeminiClient) buildPrompt(url string, report *models.WebsiteReport, focus []string) string {
	reportJSON, _ := json.MarshalIndent(report, "", "  ")

	focusStr := "all areas"
	if len(focus) > 0 {
		focusStr = strings.Join(focus, ", ")
	}

	return fmt.Sprintf(`You are a website optimization expert. Analyze the following website check report and provide actionable recommendations.

Website URL: %s
Focus Areas: %s

Website Check Report:
%s

IMPORTANT: You must respond with ONLY a valid JSON object. Do not include any markdown formatting, code blocks, or explanatory text. Start your response directly with { and end with }.

Required JSON format:
{
  "url": "%s",
  "generated_at": "%s",
  "summary": "Brief overall summary of main issues and priorities",
  "recommendations": [
    {
      "category": "robots|seo|security|sitemap",
      "priority": "high|medium|low",
      "issues": [
        {
          "issue": "Description of the specific issue",
          "impact": "What this affects (SEO ranking, security, etc.)",
          "current_status": "Current state description"
        }
      ],
      "improvements": [
        "Specific actionable step 1",
        "Specific actionable step 2"
      ],
      "resources": [
        "Valid URL to documentation or leave empty"
      ]
    }
  ]
}

Guidelines:
1. Focus on failed and warning status items first
2. Provide specific, actionable recommendations
3. Explain the business impact of each issue
4. Include code examples or specific configurations when relevant
5. Prioritize based on security risks and SEO impact
6. Keep recommendations concise but comprehensive
7. Only include valid URLs in resources or omit the field
8. Categories must be one of: robots, seo, security, sitemap
9. Priority must be one of: high, medium, low

Return ONLY the JSON object, no other text.`, url, focusStr, string(reportJSON), url, report.Timestamp.Format("2006-01-02T15:04:05Z"))
}
