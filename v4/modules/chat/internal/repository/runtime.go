package repository

import (
	configs "mdn4/modules/chat/config"
	"mdn4/modules/chat/internal/repository/orm"
	"mdn4/modules/chat/internal/usecases"
	"mdn4/modules/chat/models"

	"github.com/hyperledger/fabric/common/flogging"
	"gorm.io/gorm"
)

var repoLogger = flogging.MustGetLogger("repository")

func NewRepository(db *gorm.DB) usecases.Repository {
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
			repoLogger.Error("Failed to migrate user model")
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
			repoLogger.Error("Failed to migrate user model")
		}
		return orm.NewGormRepository(db)
	}
}
