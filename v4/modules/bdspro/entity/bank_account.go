package entity

import (
	"time"

	"gorm.io/gorm"
)

// BankAccount maps to bank_accounts (linked via deals.bank_account_id).
type BankAccount struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt     *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy     *uint64        `gorm:"column:created_by"`
	UpdatedBy     *uint64        `gorm:"column:updated_by"`
	OwnerID       uint64         `gorm:"column:owner_id;not null;default:0"`
	OwnerType     uint32         `gorm:"column:owner_type;not null;default:10"`
	BankName      string         `gorm:"column:bank_name;not null"`
	BankID        uint64         `gorm:"column:bank_id;default:0"`
	AccountNumber string         `gorm:"column:account_number;not null"`
	AccountName   string         `gorm:"column:account_name;not null"`
	AccountType   uint32         `gorm:"column:account_type;default:1"`
	AccountStatus bool           `gorm:"column:account_status;default:true"`
}

func (BankAccount) TableName() string { return "bank_accounts" }
