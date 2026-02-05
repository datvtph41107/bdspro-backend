package handler

import (
	"chat/internal/domain"
	"chat/internal/usecases"
	"context"
	"time"

	chatpb "pb/types/chat"
)

type BackgroundHandler struct {
	chatpb.UnimplementedBackgroundImageServiceServer
	backgroundImageUsecases *usecases.BackgroundImageUsecases
}

func NewBackgroundHandler(backgroundImageUsecases *usecases.BackgroundImageUsecases) *BackgroundHandler {
	return &BackgroundHandler{
		backgroundImageUsecases: backgroundImageUsecases,
	}
}

func (h *BackgroundHandler) CreateBackgroundImage(ctx context.Context, req *chatpb.CreateBackgroundImageRequest) (*chatpb.CreateBackgroundImageResponse, error) {
	model, err := h.backgroundImageUsecases.CreateBackgroundImage(ctx, &domain.BackgroundImage{
		ImageURL:    req.ImageUrl,
		Title:       req.Title,
		ImageType:   domain.ImageType(req.ImageType),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &chatpb.CreateBackgroundImageResponse{
		BackgroundImage: &chatpb.BackgroundImage{
			Id:          model.ID,
			ImageUrl:    model.ImageURL,
			Title:       model.Title,
			ImageType:   chatpb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		},
	}, nil
}

func (h *BackgroundHandler) GetBackgroundImages(ctx context.Context, req *chatpb.GetBackgroundImageRequest) (*chatpb.GetBackgroundImagesResponse, error) {
	models, total, err := h.backgroundImageUsecases.GetBackgroundImages(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	response := &chatpb.GetBackgroundImagesResponse{
		Total: uint64(total),
		Data:  make([]*chatpb.BackgroundImage, 0),
	}

	for _, model := range models {
		response.Data = append(response.Data, &chatpb.BackgroundImage{
			Id:          model.ID,
			ImageUrl:    model.ImageURL,
			Title:       model.Title,
			ImageType:   chatpb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		})
	}
	return response, nil
}

func (h *BackgroundHandler) UpdateBackgroundImage(ctx context.Context, req *chatpb.UpdateBackgroundImageRequest) (*chatpb.UpdateBackgroundImageResponse, error) {
	model, err := h.backgroundImageUsecases.UpdateBackgroundImage(ctx, &domain.BackgroundImage{
		ImageURL:    req.ImageUrl,
		Title:       req.Title,
		ImageType:   domain.ImageType(req.ImageType),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &chatpb.UpdateBackgroundImageResponse{
		BackgroundImage: &chatpb.BackgroundImage{
			Id:          model.ID,
			ImageUrl:    model.ImageURL,
			Title:       model.Title,
			ImageType:   chatpb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		},
	}, nil
}

func (h *BackgroundHandler) DeleteBackgroundImage(ctx context.Context, req *chatpb.DeleteBackgroundImageRequest) (*chatpb.DeleteBackgroundImageResponse, error) {
	err := h.backgroundImageUsecases.DeleteBackgroundImage(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &chatpb.DeleteBackgroundImageResponse{
		Success: true,
	}, nil
}
