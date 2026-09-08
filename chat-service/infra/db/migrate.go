package db

import (
	"chat/internal/domain"
	_db "common/db"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		&domain.Conversation{},
		&domain.ChatTimeline{},
		&domain.ChatEvent{},
		&domain.Membership{},
		&domain.Approval{},
		&domain.MessageReaction{},
		&domain.BackgroundImage{},
	)
}
