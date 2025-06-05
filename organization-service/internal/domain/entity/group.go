package entity

import "time"

type Group struct {
	Id uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   uint32
	UpdatedBy   uint32
	
	Name        string
	Description string
	AvatarUrl   string
}
