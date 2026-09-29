package mapper

import (
	"chat/models"
	chatpb "pb/types/chat"
	"time"
)

type BackgroundMapper struct {
}

func NewBackgroundMapper() *BackgroundMapper {
	return &BackgroundMapper{}
}

func (m *BackgroundMapper) BackgroundImageToPb(backgroundImage *models.BackgroundImageModel) *chatpb.BackgroundImage {
	if backgroundImage == nil {
		return nil
	}

	return &chatpb.BackgroundImage{
		Id:          backgroundImage.ID,
		ImageUrl:    backgroundImage.ImageURL,
		Title:       backgroundImage.Title,
		Color:       backgroundImage.Color,
		Background:  backgroundImage.Background,
		Text:        backgroundImage.Text,
		Tint:        backgroundImage.Tint,
		ImageType:   chatpb.ImageType(backgroundImage.ImageType),
		Description: backgroundImage.Description,
		CreatedAt:   backgroundImage.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   backgroundImage.UpdatedAt.Format(time.RFC3339),
	}
}
