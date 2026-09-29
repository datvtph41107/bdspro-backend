package repository

import (
	"context"

	configs "chat/config"
	"chat/internal/repository/orm"
	"chat/internal/usecases"
	"chat/models"
	"common/logging"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) usecases.Repository {
	logger := logging.WithComponent(context.Background(), "repository")
	switch configs.AppProperties.Database.Type {
	case "POSTGRES":
		err := db.AutoMigrate(&models.ConversationModel{},
			&models.MessageModel{},
			&models.ParticipantModel{},
			&models.ReadReceptModel{},
			&models.MessageReactionModel{},
			&models.BackgroundImageModel{},
		)
		if err != nil {
			logger.Error("Failed to migrate user model")
		}
		return orm.NewGormRepository(db)
	default:
		err := db.AutoMigrate(&models.ConversationModel{},
			&models.MessageModel{},
			&models.ParticipantModel{},
			&models.ReadReceptModel{},
			&models.MessageReactionModel{},
			&models.BackgroundImageModel{},
		)
		if err != nil {
			logger.Error("Failed to migrate user model")
		}
		return orm.NewGormRepository(db)
	}
}
