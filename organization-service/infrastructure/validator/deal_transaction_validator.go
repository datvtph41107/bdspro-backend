package validator

// type DealTransactionValidator interface {
// 	ValidateCreateRequest(req *pb.CreateDealTransactionRequest) error
// 	ValidateUpdateRequest(req *pb.UpdateDealTransactionRequest) error
// }

// type dealTransactionValidator struct{}

// func NewDealTransactionValidator() DealTransactionValidator {
// 	return &dealTransactionValidator{}
// }

// func (v *dealTransactionValidator) ValidateCreateRequest(req *pb.CreateDealTransactionRequest) error {
// 	if req == nil {
// 		return fmt.Errorf("request cannot be nil")
// 	}

// 	if req.DealId == 0 {
// 		return fmt.Errorf("deal ID is required")
// 	}

// 	if req.Type == pb.TransactionType_TRANSACTION_TYPE_UNKNOWN {
// 		return fmt.Errorf("transaction type is required")
// 	}

// 	if req.Amount <= 0 {
// 		return fmt.Errorf("amount must be greater than 0")
// 	}

// 	if req.Description == "" {
// 		return fmt.Errorf("description is required")
// 	}

// 	if req.TransactionDate == "" {
// 		return fmt.Errorf("transaction date is required")
// 	}

// 	if req.Category == "" {
// 		return fmt.Errorf("category is required")
// 	}

// 	return nil
// }

// func (v *dealTransactionValidator) ValidateUpdateRequest(req *pb.UpdateDealTransactionRequest) error {
// 	if req == nil {
// 		return fmt.Errorf("request cannot be nil")
// 	}

// 	if req.Id == 0 {
// 		return fmt.Errorf("transaction ID is required")
// 	}

// 	if req.Amount <= 0 {
// 		return fmt.Errorf("amount must be greater than 0")
// 	}

// 	if req.Description == "" {
// 		return fmt.Errorf("description is required")
// 	}

// 	if req.TransactionDate == "" {
// 		return fmt.Errorf("transaction date is required")
// 	}

// 	if req.Category == "" {
// 		return fmt.Errorf("category is required")
// 	}

// 	return nil
// }
