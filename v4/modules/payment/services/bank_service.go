package services

import (
	"context"
	"time"

	_utils "common/utils"
	paymentpb "pb/types/payment"

	"mdn4/modules/payment/repository"
)

// BankListOutcome is the bank catalog response + sync metadata.
type BankListOutcome struct {
	List      *paymentpb.BankList
	Timestamp int64
	Unchanged bool
}

// BankService serves client bank catalog (parity with payment BankService.GetAll).
type BankService interface {
	ListAll(ctx context.Context, clientTimestamp int64) (*BankListOutcome, error)
}

// @bind: mdn4/modules/payment/services.BankService
type LocalBankService struct {
	banks repository.BankRepository
	sync  *_utils.SyncUtil
}

func NewLocalBankService(banks repository.BankRepository, sync *_utils.SyncUtil) *LocalBankService {
	return &LocalBankService{banks: banks, sync: sync}
}

func (s *LocalBankService) ListAll(ctx context.Context, clientTimestamp int64) (*BankListOutcome, error) {
	if s.sync != nil {
		catalogTS := s.sync.EnsureTimestamp(ctx, _utils.SyncKeyPaymentBankAll)
		if clientTimestamp > 0 && !s.sync.HasUpdated(ctx, _utils.SyncKeyPaymentBankAll, clientTimestamp) {
			return &BankListOutcome{Timestamp: catalogTS, Unchanged: true}, nil
		}
	}

	rows, err := s.banks.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*paymentpb.Bank, 0, len(rows))
	for _, b := range rows {
		if b == nil {
			continue
		}
		out = append(out, &paymentpb.Bank{
			Id:          b.ID,
			Name:        b.Name,
			Code:        b.Code,
			Logo:        b.Logo,
			Active:      b.Active,
			Description: b.Description,
		})
	}

	ts := time.Now().UnixMilli()
	if s.sync != nil {
		ts = s.sync.EnsureTimestamp(ctx, _utils.SyncKeyPaymentBankAll)
	}

	return &BankListOutcome{
		List: &paymentpb.BankList{
			Data:  out,
			Total: int64(len(out)),
		},
		Timestamp: ts,
	}, nil
}
