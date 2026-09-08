package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type BankAccount struct {
	_models.BaseEntity
	OwnerID       uint64         `gorm:"not null"`
	OwnerType     enums.EOwnerOf `gorm:"not null;default:10"`
	BankName      string         `gorm:"not null"`
	BankId        uint64         ``
	AccountNumber string         `gorm:"not null"`
	AccountName   string         `gorm:"not null"`
	AccountType   uint           // 1: tài khoản, 2: thẻ
	AccountStatus bool           `gorm:"not null;default:true"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}
