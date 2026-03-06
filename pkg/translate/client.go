package translate

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

type Client struct {
	genaiClient *genai.Client
	model       string
}

func NewClient(ctx context.Context, apiKey, model string) (*Client, error) {

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		genaiClient: client,
		model:       model,
	}, nil
}

func (c *Client) TranslateStream(ctx context.Context, w io.Writer, text, from, to, systemPrompt string) error {
	baseSystemInstruction := "You are a professional translator. Output ONLY the translation result. No explanations, no preamble, no markdown code blocks."

	if systemPrompt != "" {
		baseSystemInstruction = fmt.Sprintf("%s\n\n%s", baseSystemInstruction, systemPrompt)
	}
	userPrompt := fmt.Sprintf("From: %s\nTo: %s\nText:\n%s", from, to, text)

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: baseSystemInstruction}},
		},
		Temperature: pointer(0.2),
	}

	iter := c.genaiClient.Models.GenerateContentStream(ctx, c.model, genai.Text(userPrompt), config)

	for res, err := range iter {
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}

		if len(res.Candidates) > 0 && res.Candidates[0].Content != nil {
			for _, part := range res.Candidates[0].Content.Parts {
				if part.Text != "" {
					fmt.Fprint(w, part.Text)
				}
			}
		}
	}

	return nil
}

func pointer(f float32) *float32 {
	return &f
}
