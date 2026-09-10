package client

import (
	"context"
	"fmt"
	"strings"

	assistantpb "pb/types/assistant"
	"tqd/infra/rpc"
	"tqd/internal/interface/provider"

	"github.com/spf13/viper"
)

// AssistantClient adapts assistant gRPC to the narrow classify capability.
type AssistantClient struct {
	client assistantpb.AssistantServiceClient
	close  func()
}

// NewAssistantClient dials assistant when classify is enabled; otherwise returns
// a no-op adapter so Wire can compose the graph without a process-owned dial.
func NewAssistantClient() (provider.AssistantProvider, func(), error) {
	if !viper.GetBool("qh_planning_classify.enabled") {
		return &AssistantClient{}, func() {}, nil
	}
	target := strings.TrimSpace(viper.GetString("rpc.assistant.address"))
	if target == "" {
		target = "127.0.0.1:8218"
	}
	rpcClient, cleanup, err := rpc.OpenAssistantRPCClient(target)
	if err != nil {
		return nil, nil, fmt.Errorf("open assistant gRPC client: %w", err)
	}
	return &AssistantClient{client: rpcClient, close: cleanup}, cleanup, nil
}

func (c *AssistantClient) Close() error {
	if c != nil && c.close != nil {
		c.close()
		c.close = nil
	}
	return nil
}

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
