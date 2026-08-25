package filegrpc

import (
	"context"
	"errors"

	"file/internal/filecontent"
	filepb "pb/types/file"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) PutOwnedFile(
	ctx context.Context,
	req *filepb.PutOwnedFileRequest,
) (*filepb.UploadFileResponse, error) {
	if req == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	saved, err :=
		h.files.PutOwnedFile(
			ctx,
			filecontent.PutOwnedInput{
				Content:        req.GetContent(),
				Filename:       req.GetFilename(),
				ContentType:    req.GetContentType(),
				Description:    req.GetDescription(),
				OwnerNamespace: req.GetOwnerNamespace(),
				OwnerKey:       req.GetOwnerKey(),
			},
		)
	if err != nil {
		switch {
		case errors.Is(
			err,
			filecontent.ErrInvalidOwnedFile,
		):
			return nil, status.Error(
				codes.InvalidArgument,
				"owned file request is invalid",
			)

		case errors.Is(
			err,
			filecontent.ErrOwnedFileConflict,
		):
			return nil, status.Error(
				codes.AlreadyExists,
				filecontent.ErrOwnedFileConflict.Error(),
			)

		default:
			return nil, status.Error(
				codes.Internal,
				"cannot store owned file",
			)
		}
	}

	return &filepb.UploadFileResponse{
		FileId:      saved.ID,
		Path:        saved.Path,
		Filename:    saved.Filename,
		ContentType: saved.ContentType,
		Size:        saved.Size,
	}, nil
}
