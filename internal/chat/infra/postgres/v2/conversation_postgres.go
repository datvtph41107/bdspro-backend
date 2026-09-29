package postgres_v2

import (
	"gorm.io/gorm"
)

type ConversationPostgres struct {
	db *gorm.DB
}

func NewConversationPostgres(db *gorm.DB) *ConversationPostgres {
	return &ConversationPostgres{db}
}
