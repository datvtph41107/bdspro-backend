package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"time"

	assistantpb "pb/types/assistant"
	"tqd/env"
	"tqd/internal/interface/provider"

	"google.golang.org/grpc"
)

// assistantGrpcPort — port gRPC cố định của assistant-service (config app.port.grpc).
const assistantGrpcPort = "8218"

// assistantGrpcMaxMsgBytes — đủ cho ClassifyDocument gửi file inline tới ~15MB (+ overhead).
const assistantGrpcMaxMsgBytes = 32 << 20

// AssistantClient gọi sang assistant-service để sinh nội dung AI (Gemini), dùng cho job
// phân loại tài liệu quy hoạch trong QHPlanningClassifyUsecase.
type AssistantClient struct {
	client assistantpb.AssistantServiceClient
}

// NewAssistantClient trả về provider.AssistantProvider trực tiếp nên không cần wire.Bind.
func NewAssistantClient() provider.AssistantProvider {
	assistantClient := &AssistantClient{}
	go func() {
		name := "assistant"
		url := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", url, assistantGrpcPort),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(assistantGrpcMaxMsgBytes),
				grpc.MaxCallSendMsgSize(assistantGrpcMaxMsgBytes),
			),
		)
		if err != nil {
			log.Printf("[AssistantClient] không dial được assistant-service: %v", err)
			return
		}

		go _middleware.MonitorConnection(context.Background(), name, conn)

		assistantClient.client = assistantpb.NewAssistantServiceClient(conn)
	}()
	return assistantClient
}

// GenerateContent gọi GenerateContentWithGemini bên assistant-service với 1 prompt tự do (không template).
func (c *AssistantClient) GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, string, error) {
	if c.client == nil {
		return "", "", fmt.Errorf("assistant client chưa kết nối tới assistant-service")
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

// ClassifyDocument gọi ClassifyDocumentWithGemini (multimodal) bên assistant-service.
func (c *AssistantClient) ClassifyDocument(ctx context.Context, prompt string, fileContent []byte, mimeType string) (string, string, error) {
	if c.client == nil {
		return "", "", fmt.Errorf("assistant client chưa kết nối tới assistant-service")
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
