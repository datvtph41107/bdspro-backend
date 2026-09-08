package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_routes "common/routes"
	"context"
)

func (s *ProductUsecase) MapProductStatusList(c context.Context, product *dto.ProductListResponse) error {
	return nil
}
func (s *ProductUsecase) MapProductStatus(c context.Context, product *domain.Product) error {
	// if product.SaleTransaction != nil {
	// 	switch product.SaleTransaction.TransactionStatus {
	// 	case enums.TransactionStatusDraft:
	// 		product.SaleStatus = enums.EProductNotSold
	// 	case enums.TransactionStatusDeposite:
	// 		product.SaleStatus = enums.EProductSelling
	// 	case enums.TransactionStatusSigned:
	// 		product.SaleStatus = enums.EProductSold
	// 	case enums.TransactionStatusCompleted:
	// 		product.SaleStatus = enums.EProductSold
	// 	case enums.TransactionStatusCanceled:
	// 		product.SaleStatus = enums.EProductNotSold
	// 	}
	// } else {
	// 	product.SaleStatus = enums.EProductNotSold
	// }

	// if product.RentTransaction != nil {
	// 	switch product.RentTransaction.TransactionStatus {
	// 	case enums.TransactionStatusDraft:
	// 		product.RentStatus = enums.EProductNotRent
	// 	case enums.TransactionStatusDeposite:
	// 		product.RentStatus = enums.EProductRenting
	// 	case enums.TransactionStatusSigned:
	// 		product.RentStatus = enums.EProductRenting
	// 	case enums.TransactionStatusCompleted:
	// 		product.RentStatus = enums.EProductRented
	// 	case enums.TransactionStatusCanceled:
	// 		product.RentStatus = enums.EProductNotRent
	// 	}
	// } else {
	// 	product.RentStatus = enums.EProductNotRent
	// }

	return nil
}

func (s *ProductUsecase) SaleStatusUpdate(c context.Context, dto dto.StatusUpdate) error {
	// todo: không thể cập nhật trạng thái sản phẩm thành cọc ở api này. dung api deposite riêng
	newSaleStatus := dto.SaleStatus
	productId := dto.ProductID
	if newSaleStatus == enums.EProductSelling {
		return &_routes.Except{
			Code:    400,
			Message: "Không thể cập nhật trạng thái sản phẩm thành cọc",
		}
	}
	err := s.CheckPermission(c, productId, enums.PermissionProductUpdateStatus)
	if err != nil {
		return err
	}
	product, err := s.ProductRepo.GetByIDContext(c, *productId)
	if err != nil {
		return err
	}
	// todo: transaction
	// if product.SaleTransactionID != nil {
	// 	product.SaleTransaction, _ = s.TransactionRepo.GetByID(c, product.SaleTransactionID)
	// }

	s.MapProductStatus(c, product)
	if product.SaleStatus == newSaleStatus {
		return &_routes.Except{
			Code:    400,
			Message: "Trạng thái trùng với trạng thái cũ",
		}
	}

	contextTx := s.Transaction.StartTransaction(c)

	if newSaleStatus == enums.EProductNotSold {
		err = s.SaleStatusToNotSold(contextTx, product)
	} else if newSaleStatus == enums.EProductSold {
		err = s.SaleStatusToSold(contextTx, product, dto)
	}
	if err != nil {
		s.Transaction.RollbackTransaction(contextTx)
		return err
	}

	if err := s.Transaction.CommitTransaction(contextTx); err != nil {
		return err
	}

	return nil
}

func (s *ProductUsecase) EmptySaleTransaction(c context.Context, product *domain.Product) error {
	if product.SaleTransactionID != nil {
		err := s.ProductRepo.UpdateSaleTransactionID(c, product.ID, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ProductUsecase) SaleStatusToNotSold(ctx context.Context, product *domain.Product) error {
	// if product.SaleTransaction != nil &&
	// 	product.SaleTransaction.TransactionStatus != enums.TransactionStatusCompleted {
	// 	product.SaleTransaction.TransactionStatus = enums.TransactionStatusCanceled
	// 	// todo: transaction
	// 	// _, err := s.TransactionRepo.Update(ctx, product.SaleTransaction)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	err := s.EmptySaleTransaction(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

// todo: không có api này. phải gọi qua api cọc
func (s *ProductUsecase) SaleStatusToSelling(c context.Context, dto dto.StatusUpdate) error {
	return nil
}

func (s *ProductUsecase) SaleStatusToSold(ctx context.Context, product *domain.Product, dto dto.StatusUpdate) error {
	// trans := product.SaleTransaction

	// if (trans == nil || trans.ContactID == nil) && dto.ContactID == nil {
	// 	return &_routes.Except{
	// 		Code:    400,
	// 		Message: "Không thể chuyển trạng thái sản phẩm thành đã bán khi không có người liên hệ",
	// 	}
	// }

	// if product.SaleTransaction != nil {
	// 	// todo: transaction
	// 	// product.SaleTransaction.TransactionStatus = enums.TransactionStatusCompleted
	// 	// s.AttachContactToTransaction(product.SaleTransaction, dto)
	// 	// _, err := s.TransactionRepo.Update(ctx, product.SaleTransaction)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// } else {
	// 	// todo: transaction
	// 	// trans = s.EmptyTransaction(product, enums.TransactionTypeSale)
	// 	// trans.TransactionStatus = enums.TransactionStatusCompleted
	// 	// trans = s.AttachContactToTransaction(trans, dto)
	// 	// _, err := s.TransactionRepo.Create(ctx, trans)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	// if product.SaleTransactionID != &trans.ID {
	// 	product.SaleTransactionID = &trans.ID
	// 	err := s.ProductRepo.UpdateProduct(ctx, product.ID, product)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// func (s *ProductUsecase) EmptyTransaction(product *domain.Product, transactionType enums.TransactionType) *tx_domain.Tx {
// 	return &tx_domain.Tx{
// 		ProductID:         &product.ID,
// 		OwnerId:           product.OwnerID,
// 		OwnerType:         product.OwnerType,
// 		TransactionType:   transactionType,
// 		TransactionStatus: enums.TransactionStatusDraft,
// 	}
// }

// func (s *ProductUsecase) AttachContactToTransaction(trans *tx_domain.Tx, dto dto.StatusUpdate) *tx_domain.Tx {
// 	if trans.ContactID == nil {
// 		trans.ContactID = dto.ContactID
// 		trans.DepositeNote = &dto.Note
// 		trans.Amount = dto.Amount
// 		trans.ContactPhone = &dto.Phone
// 	}
// 	return trans
// }

func (s *ProductUsecase) RentStatusUpdate(c context.Context, dto dto.StatusUpdate) error {
	// todo: không thể cập nhật trạng thái sản phẩm thành cọc ở api này. dung api deposite riêng
	newRentStatus := dto.RentStatus
	productId := dto.ProductID
	if newRentStatus == enums.EProductRenting {
		return &_routes.Except{
			Code:    400,
			Message: "Không thể cập nhật trạng thái sản phẩm thành cọc",
		}
	}
	err := s.CheckPermission(c, productId, enums.PermissionProductUpdateStatus)
	if err != nil {
		return err
	}
	product, err := s.ProductRepo.GetByIDContext(c, *productId)
	if err != nil {
		return err
	}
	if product.RentTransactionID != nil {
		// todo: transaction
		// product.RentTransaction, _ = s.TransactionRepo.GetByID(c, product.RentTransactionID)
	}

	s.MapProductStatus(c, product)
	if product.RentStatus == newRentStatus {
		return &_routes.Except{
			Code:    400,
			Message: "Trạng thái trùng với trạng thái cũ",
		}
	}

	contextTx := s.Transaction.StartTransaction(c)

	if newRentStatus == enums.EProductNotRent {
		err = s.RentStatusToNotRent(contextTx, product)
	} else if newRentStatus == enums.EProductRented {
		err = s.RentStatusToRented(contextTx, product, dto)
	}
	if err != nil {
		s.Transaction.RollbackTransaction(contextTx)
		return err
	}

	if err := s.Transaction.CommitTransaction(contextTx); err != nil {
		return err
	}

	return nil
}

func (s *ProductUsecase) RentStatusToNotRent(ctx context.Context, product *domain.Product) error {
	// if product.RentTransaction != nil &&
	// 	product.RentTransaction.TransactionStatus != enums.TransactionStatusCompleted {
	// 	product.RentTransaction.TransactionStatus = enums.TransactionStatusCanceled
	// 	// todo: transaction
	// 	// _, err := s.TransactionRepo.Update(ctx, product.RentTransaction)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	err := s.EmptyRentTransaction(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductUsecase) RentStatusToRented(ctx context.Context, product *domain.Product, dto dto.StatusUpdate) error {
	// trans := product.RentTransaction

	// if (trans == nil || trans.ContactID == nil) && dto.ContactID == nil {
	// 	return &_routes.Except{
	// 		Code:    400,
	// 		Message: "Không thể chuyển trạng thái sản phẩm thành đã cho thuê khi không có người liên hệ",
	// 	}
	// }

	// if product.RentTransaction != nil {
	// 	// todo: transaction
	// 	// product.RentTransaction.TransactionStatus = enums.TransactionStatusCompleted
	// 	// s.AttachContactToTransaction(product.RentTransaction, dto)
	// 	// _, err := s.TransactionRepo.Update(ctx, product.RentTransaction)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// } else {
	// 	// todo: transaction
	// 	// trans = s.EmptyTransaction(product, enums.TransactionTypeRent)
	// 	// trans.TransactionStatus = enums.TransactionStatusCompleted
	// 	// trans = s.AttachContactToTransaction(trans, dto)
	// 	// _, err := s.TransactionRepo.Create(ctx, trans)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	// if product.RentTransactionID != &trans.ID {
	// 	product.RentTransactionID = &trans.ID
	// 	err := s.ProductRepo.UpdateProduct(ctx, product.ID, product)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (s *ProductUsecase) EmptyRentTransaction(c context.Context, product *domain.Product) error {
	if product.RentTransactionID != nil {
		err := s.ProductRepo.UpdateRentTransactionID(c, product.ID, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
