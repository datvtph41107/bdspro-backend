package handlers

import (
	"context"
	"time"

	"github.com/tuandq3/go-chat-service/models"
	pb "github.com/tuandq3/go-chat-service/proto/types/external"
)

type externalHandler struct {
	pb.UnimplementedBackgroundImageServiceServer
	backgroundImageUsecases backgroundImageUsecases
}

func NewExternalHandler(backgroundImageUsecases backgroundImageUsecases) *externalHandler {
	return &externalHandler{
		backgroundImageUsecases: backgroundImageUsecases,
	}
}

func (h *externalHandler) CreateBackgroundImage(ctx context.Context, req *pb.CreateBackgroundImageRequest) (*pb.CreateBackgroundImageResponse, error) {
	model, err := h.backgroundImageUsecases.CreateBackgroundImage(ctx, &models.BackgroundImageModel{
		ImageURL:    req.ImageUrl,
		Title:       req.Title,
		ImageType:   models.ImageType(req.ImageType),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateBackgroundImageResponse{
		BackgroundImage: &pb.BackgroundImage{
			Id:          model.ID,
			ImageUrl:    model.ImageURL,
			Title:       model.Title,
			ImageType:   pb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		},
	}, nil
}

func (h *externalHandler) GetBackgroundImages(ctx context.Context, req *pb.GetBackgroundImageRequest) (*pb.GetBackgroundImagesResponse, error) {
	models, total, err := h.backgroundImageUsecases.GetBackgroundImages(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	response := &pb.GetBackgroundImagesResponse{
		Total: uint64(total),
		Data:  make([]*pb.BackgroundImage, 0),
	}

	for _, model := range models {
		response.Data = append(response.Data, &pb.BackgroundImage{
			Id:          model.ID,
			ImageUrl:    model.ImageURL,
			Title:       model.Title,
			ImageType:   pb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		})
	}
	return response, nil
}

func (h *externalHandler) UpdateBackgroundImage(ctx context.Context, req *pb.UpdateBackgroundImageRequest) (*pb.UpdateBackgroundImageResponse, error) {
	model, err := h.backgroundImageUsecases.UpdateBackgroundImage(ctx, &models.BackgroundImageModel{
		ID:          req.Id,
		ImageURL:    req.ImageUrl,
		Title:       req.Title,
		ImageType:   models.ImageType(req.ImageType),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateBackgroundImageResponse{
		BackgroundImage: &pb.BackgroundImage{
			Id:          model.ID,
			ImageUrl:    model.ImageURL,
			Title:       model.Title,
			ImageType:   pb.ImageType(model.ImageType),
			Description: model.Description,
			CreatedAt:   model.CreatedAt.Local().Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Local().Format(time.RFC3339),
		},
	}, nil
}

func (h *externalHandler) DeleteBackgroundImage(ctx context.Context, req *pb.DeleteBackgroundImageRequest) (*pb.DeleteBackgroundImageResponse, error) {
	err := h.backgroundImageUsecases.DeleteBackgroundImage(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteBackgroundImageResponse{
		Success: true,
	}, nil
}
