package filegrpc

import (
	"context"
	"errors"
	"io"

	filepb "pb/types/file"

	"google.golang.org/grpc"
)

// BuildReportClient deliberately returns the concrete File client because
// Report requires PutOwnedFile in addition to legacy generic FileProvider.
func BuildReportClient(
	client filepb.InternalFileServiceClient,
) *Client {
	return &Client{client: client}
}

func (c *Client) PutOwnedFile(
	ctx context.Context,
	ownerNamespace string,
	ownerKey string,
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
		return "", errors.New(
			"file content exceeds 64 MiB",
		)
	}

	resp, err :=
		c.client.PutOwnedFile(
			ctx,
			&filepb.PutOwnedFileRequest{
				Content:        content,
				Filename:       filename,
				ContentType:    contentType,
				OwnerNamespace: ownerNamespace,
				OwnerKey:       ownerKey,
			},
			grpc.MaxCallSendMsgSize(
				maxFileMessageBytes,
			),
		)
	if err != nil {
		return "", err
	}

	return resp.GetPath(), nil
}
