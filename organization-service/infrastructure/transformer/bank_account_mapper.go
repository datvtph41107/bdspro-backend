package transformer

import (
	"organization/internal/domain/entity"
	organizationpb "pb/types/organization"
)

type BankAccountMapper struct{}

func NewBankAccountMapper() *BankAccountMapper {
	return &BankAccountMapper{}
}

func (m *BankAccountMapper) EntityToPb(bankAccount *entity.BankAccount) *organizationpb.BankAccountItem {
	return &organizationpb.BankAccountItem{
		Name:          bankAccount.BankName,
		BankName:      bankAccount.BankName,
		AccountNumber: bankAccount.AccountNumber,
		AccountName:   bankAccount.AccountName,
	}
}
