package entity

import "time"

type GroupDocument struct {
	Id          uint32
	GroupId     uint32
	Name        string
	Description string
	FileUrl     string
	FileType    string
	FileSize    uint64
	CreatedBy   uint32
	UpdatedBy   uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Uploader    string
} 