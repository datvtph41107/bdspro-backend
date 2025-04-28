package repository

import (
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/tuandq3/go-chat-service/configs"
	"github.com/tuandq3/go-chat-service/internal/repository/orm"
	"github.com/tuandq3/go-chat-service/models"
	"gorm.io/gorm"
)

var repoLogger = flogging.MustGetLogger("repository")

func NewRepository(db *gorm.DB) repository {
	switch configs.DB_TYPE {
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
