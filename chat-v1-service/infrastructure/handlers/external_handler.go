package handlers

import (
	"context"
	"time"

	"chat/models"
	chatpb "pb/types/chat"
)

type externalHandler struct {
	chatpb.UnimplementedBackgroundImageServiceServer
	backgroundImageUsecases backgroundImageUsecases
}

func NewExternalHandler(backgroundImageUsecases backgroundImageUsecases) *externalHandler {
	return &externalHandler{
		backgroundImageUsecases: backgroundImageUsecases,
	}
}

func (h *externalHandler) CreateBackgroundImage(ctx context.Context, req *chatpb.CreateBackgroundImageRequest) (*chatpb.CreateBackgroundImageResponse, error) {
	model, err := h.backgroundImageUsecases.CreateBackgroundImage(ctx, &models.BackgroundImageModel{
		ImageURL:    req.ImageUrl,
		Title:       req.Title,
		Background:  req.Background,
		Text:        req.Text,
		Tint:        req.Tint,
		ImageType:   models.ImageType(req.ImageType),
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
			Background:  model.Background,
			Text:        model.Text,
			Tint:        model.Tint,
			ImageType:   chatpb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		},
	}, nil
}

func (h *externalHandler) GetBackgroundImages(ctx context.Context, req *chatpb.GetBackgroundImageRequest) (*chatpb.GetBackgroundImagesResponse, error) {
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
			Background:  model.Background,
			Text:        model.Text,
			Tint:        model.Tint,
			ImageType:   chatpb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		})
	}
	return response, nil
}

func (h *externalHandler) UpdateBackgroundImage(ctx context.Context, req *chatpb.UpdateBackgroundImageRequest) (*chatpb.UpdateBackgroundImageResponse, error) {
	model, err := h.backgroundImageUsecases.UpdateBackgroundImage(ctx, &models.BackgroundImageModel{
		ID:          req.Id,
		ImageURL:    req.ImageUrl,
		Title:       req.Title,
		Background:  req.Background,
		Text:        req.Text,
		Tint:        req.Tint,
		ImageType:   models.ImageType(req.ImageType),
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
			Background:  model.Background,
			Text:        model.Text,
			Tint:        model.Tint,
			ImageType:   chatpb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		},
	}, nil
}

func (h *externalHandler) DeleteBackgroundImage(ctx context.Context, req *chatpb.DeleteBackgroundImageRequest) (*chatpb.DeleteBackgroundImageResponse, error) {
	err := h.backgroundImageUsecases.DeleteBackgroundImage(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &chatpb.DeleteBackgroundImageResponse{
		Success: true,
	}, nil
}
