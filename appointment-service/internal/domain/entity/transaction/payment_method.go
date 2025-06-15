package transaction

import (
	"time"
)

// PaymentMethod represents a payment method
type PaymentMethod struct {
	Id         uint32        
	Name       string        
	MethodCode string        
	Status     PaymentStatus 
	CreatedAt  time.Time     
	UpdatedAt  time.Time     
} 