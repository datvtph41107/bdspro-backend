package domain

import "time"

type Account struct {
	UserId        uint64     `gorm:"primaryKey;not null" json:"userId" binding:"required"`
	Fullname      string     `json:"fullName"`
	Phone         string     `json:"phone"`
	AccountNumber string     `json:"accountNumber"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

func (Account) TableName() string {
	return "accounts"
}
