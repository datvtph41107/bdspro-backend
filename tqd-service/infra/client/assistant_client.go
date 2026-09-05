package client

import (
	"context"
	"fmt"

	assistantpb "pb/types/assistant"
	"tqd/internal/interface/provider"
)

// AssistantClient adapts the process-owned assistant gRPC client to the narrow
// capability required by TQD classification. It does not own connection
// lifecycle; cmd/grpc closes the underlying connection after actors stop.
type AssistantClient struct {
	client assistantpb.AssistantServiceClient
}

func NewAssistantClient(client assistantpb.AssistantServiceClient) provider.AssistantProvider {
	return &AssistantClient{client: client}
}

// Close is retained only for the legacy provider contract. Resource ownership
// is at process composition, therefore the adapter has nothing to close.
func (c *AssistantClient) Close() error { return nil }

func (c *AssistantClient) ready() error {
	if c == nil || c.client == nil {
		return fmt.Errorf("assistant client is unavailable")
	}
	return nil
}

func (c *AssistantClient) GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, string, error) {
	if err := c.ready(); err != nil {
		return "", "", err
	}
	resp, err := c.client.GenerateContentWithGemini(ctx, &assistantpb.GenerateRequest{
		Prompt: prompt, MaxTokens: int32(maxTokens), Temperature: temperature,
	})
	if err != nil {
		return "", "", err
	}
	model := ""
	if resp.Metadata != nil {
		model = resp.Metadata.Model
	}
	return resp.Content, model, nil
}

func (c *AssistantClient) ClassifyDocument(ctx context.Context, prompt string, fileContent []byte, mimeType string) (string, string, error) {
	if err := c.ready(); err != nil {
		return "", "", err
	}
	resp, err := c.client.ClassifyDocumentWithGemini(ctx, &assistantpb.ClassifyDocumentRequest{
		Prompt: prompt, FileContent: fileContent, MimeType: mimeType,
		MaxTokens: 512, Temperature: 0.2,
	})
	if err != nil {
		return "", "", err
	}
	model := ""
	if resp.Metadata != nil {
		model = resp.Metadata.Model
	}
	return resp.Content, model, nil
}
