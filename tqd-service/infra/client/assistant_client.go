package client

import (
	"common"
	qhprorpc "common/rpc"
	"common/transporttrust"
	"context"
	"fmt"
	"time"

	assistantpb "pb/types/assistant"
	"tqd/env"
	"tqd/internal/interface/provider"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// assistantGrpcPort is the configured gRPC port exposed by assistant-service.
const assistantGrpcPort = "8218"

// assistantGrpcMaxMsgBytes allows document classification to send inline files
// up to roughly 15 MiB plus protobuf overhead.
const assistantGrpcMaxMsgBytes = 32 << 20

// AssistantClient owns the gRPC channel used by the classify worker.
// Construction is synchronous and does not spawn connection/monitor goroutines;
// the process closes the channel only after the classify actor has stopped.
type AssistantClient struct {
	conn    *grpc.ClientConn
	client  assistantpb.AssistantServiceClient
	initErr error
}

// NewAssistantClient returns provider.AssistantProvider directly so Wire does
// not need an additional bind. The constructor preserves the existing Wire
// signature while using the canonical Foundation RPC transport.
func NewAssistantClient() provider.AssistantProvider {
	target := fmt.Sprintf("%s:%s", common.GetPrefixProtobufUrl(env.ENV_RUNTIME, "assistant"), assistantGrpcPort)
	conn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       qhprorpc.TransportConfig{Trust: transporttrust.ConfigFromEnv()},
		DefaultCallOptions: []grpc.CallOption{
			grpc.MaxCallRecvMsgSize(assistantGrpcMaxMsgBytes),
			grpc.MaxCallSendMsgSize(assistantGrpcMaxMsgBytes),
		},
	})
	if err != nil {
		return &AssistantClient{initErr: fmt.Errorf("create assistant gRPC client: %w", err)}
	}
	return &AssistantClient{
		conn:   conn,
		client: assistantpb.NewAssistantServiceClient(conn),
	}
}

// Close releases the channel owned by this client. It is safe to call when
// construction failed or after a previous successful Close.
func (c *AssistantClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	c.client = nil
	return err
}

func (c *AssistantClient) ready() error {
	if c == nil {
		return fmt.Errorf("assistant client is nil")
	}
	if c.initErr != nil {
		return c.initErr
	}
	if c.client == nil {
		return fmt.Errorf("assistant client is closed or unavailable")
	}
	return nil
}

// GenerateContent calls GenerateContentWithGemini with a free-form prompt.
func (c *AssistantClient) GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, string, error) {
	if err := c.ready(); err != nil {
		return "", "", err
	}

	resp, err := c.client.GenerateContentWithGemini(ctx, &assistantpb.GenerateRequest{
		Prompt:      prompt,
		MaxTokens:   int32(maxTokens),
		Temperature: temperature,
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

// ClassifyDocument calls the assistant multimodal classification RPC.
func (c *AssistantClient) ClassifyDocument(ctx context.Context, prompt string, fileContent []byte, mimeType string) (string, string, error) {
	if err := c.ready(); err != nil {
		return "", "", err
	}

	resp, err := c.client.ClassifyDocumentWithGemini(ctx, &assistantpb.ClassifyDocumentRequest{
		Prompt:      prompt,
		FileContent: fileContent,
		MimeType:    mimeType,
		MaxTokens:   512,
		Temperature: 0.2,
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
