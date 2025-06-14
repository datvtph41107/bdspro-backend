package transaction

import (
	"time"
)

// PaymentMethod represents a payment method
type PaymentMethod struct {
	Id         uint32        `json:"id" gorm:"type:uint32;primary_key"`
	Name       string        `json:"name" gorm:"type:varchar(255)"`
	MethodCode string        `json:"method_code" gorm:"type:varchar(50)"`
	Status     PaymentStatus `json:"status" gorm:"type:varchar(20)"`
	CreatedAt  time.Time     `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  time.Time     `json:"updated_at" gorm:"type:timestamp"`
} 