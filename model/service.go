package model

import (
	"context"
	"os"

	"google.golang.org/genai"
)

// define struct
type GeminiClient struct {
	client *genai.Client
	model  string
	config *genai.GenerateContentConfig
}

// define constructor
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"decision": map[string]any{
					"type": "string",
					"enum": []string{"delete", "keep"},
				},
			},
			"required": []string{"decision"},
		},
	}
	return &GeminiClient{
		client: client,
		model:  "gemini-3-flash-preview",
		config: config,
	}, nil
}

func (c *GeminiClient) Analyze(prompt string) (string, error) {
	ctx := context.Background()
	result, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), c.config)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
