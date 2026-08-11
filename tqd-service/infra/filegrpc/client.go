package filegrpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	filepb "pb/types/file"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"

	"google.golang.org/grpc"
)

const maxFileMessageBytes = (64 << 20) + (1 << 20)

/**
 * Client triển khai FileProvider qua gRPC.
 *
 * TQD không gọi HTTP trực tiếp sang file-service trong target flow.
 */
type Client struct {
	client filepb.InternalFileServiceClient
}

func NewClient(client filepb.InternalFileServiceClient) provider.FileProvider {
	return &Client{client: client}
}

func (c *Client) UploadFile(
	ctx context.Context,
	file io.Reader,
	filename string,
	contentType string,
) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	if len(content) > 64<<20 {
		return "", errors.New("file content exceeds 64 MiB")
	}
	resp, err := c.client.UploadFile(
		ctx,
		&filepb.UploadFileRequest{
			Content:     content,
			Filename:    filename,
			ContentType: contentType,
		},
		grpc.MaxCallSendMsgSize(maxFileMessageBytes),
	)
	if err != nil {
		return "", err
	}
	return resp.GetPath(), nil
}

func (c *Client) GetFile(ctx context.Context, path string) ([]byte, error) {
	resp, err := c.client.ReadFile(
		ctx,
		&filepb.ReadFileRequest{Path: path},
		grpc.MaxCallRecvMsgSize(maxFileMessageBytes),
	)
	if err != nil {
		return nil, err
	}
	return bytes.Clone(resp.GetContent()), nil
}

func (c *Client) DeleteFile(context.Context, string) error { return nil }

func (c *Client) GetFileURL(context.Context, string) string {
	// URL public phải được Gateway/file delivery layer tạo; TQD không tự ghép HTTP URL.
	return ""
}

func (c *Client) ValidateFile(context.Context, io.Reader, string, int64, []string) error {
	return nil
}

func (c *Client) GetFileInfo(context.Context, string) (*dto.FileInfoDTO, error) {
	return nil, nil
}

func (c *Client) GenerateThumbnail(context.Context, string) (string, error) {
	return "", nil
}
