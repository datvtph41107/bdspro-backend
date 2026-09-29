package mapper

import (
	"bdspro/internal/domain"
	bdspropb "pb/types/bdspro"
)

type BankAccountMapper struct{}

func NewBankAccountMapper() *BankAccountMapper {
	return &BankAccountMapper{}
}

func (m *BankAccountMapper) EntityToPb(bankAccount *domain.BankAccount) *bdspropb.BankAccountItem {
	return &bdspropb.BankAccountItem{
		BankId:        bankAccount.BankId,
		Name:          bankAccount.BankName,
		BankName:      bankAccount.BankName,
		AccountNumber: bankAccount.AccountNumber,
		AccountName:   bankAccount.AccountName,
	}
}
