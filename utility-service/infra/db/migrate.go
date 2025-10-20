package db

import (
	_db "common/db"
	"context"
	"log"
	"utility/internal/domain"
)

func AutoMigrate() {
	db := _db.GetDB(context.Background(), nil)
	
	err := db.AutoMigrate(
		&domain.EventQueueEntity{},
	)
	
	if err != nil {
		log.Fatalf("Auto migrate error: %v", err)
	}
	
	log.Println("Auto migrate success")
}

