package filegrpc

import (
	"common/identitycontext"
	"context"
	"errors"
	"file/internal/filecontent"
	filepb "pb/types/file"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	filepb.UnimplementedInternalFileServiceServer
	files *filecontent.Service
}

func NewHandler(files *filecontent.Service) *Handler {
	return &Handler{files: files}
}

/**
 * UploadFile nhận bytes qua gRPC nội bộ và chuyển sang file service thuần.
 */
func (h *Handler) UploadFile(
	ctx context.Context,
	req *filepb.UploadFileRequest,
) (*filepb.UploadFileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	var profileID uint64
	if actor, ok := identitycontext.FromContext(ctx); ok {
		profileID = actor.ProfileID
	}

	saved, err := h.files.UploadFile(ctx, filecontent.UploadInput{
		Content:     req.GetContent(),
		Filename:    req.GetFilename(),
		ContentType: req.GetContentType(),
		Description: req.GetDescription(),
		Secure:      req.GetSecure(),
		ProfileID:   profileID,
	})
	if err != nil {
		if errors.Is(err, filecontent.ErrNoActor) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &filepb.UploadFileResponse{
		FileId:      saved.ID,
		Path:        saved.Path,
		Filename:    saved.Filename,
		ContentType: saved.ContentType,
		Size:        saved.Size,
	}, nil
}

func (h *Handler) ReadFile(
	ctx context.Context,
	req *filepb.ReadFileRequest,
) (*filepb.ReadFileResponse, error) {
	if req == nil || req.GetPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "path is required")
	}
	var profileID uint64
	if actor, ok := identitycontext.FromContext(ctx); ok {
		profileID = actor.ProfileID
	}

	file, err := h.files.ReadFile(ctx, req.GetPath(), profileID)
	if err != nil {
		switch {
		case errors.Is(err, filecontent.ErrInvalidPath):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, filecontent.ErrNoActor):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, filecontent.ErrAccessDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.NotFound, "file not found")
	}
	return &filepb.ReadFileResponse{
		Content:     file.Content,
		ContentType: file.ContentType,
	}, nil
}
