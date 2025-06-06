package entity

import "time"

type GroupSetting struct {
	Id uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32
	
	GroupId uint32
	ConfigKey string
	ConfigValue string
}