package services

import (
	"context"
	"testing"

	"mdn4/modules/payment/entity"
)

type stubBankRepo struct {
	rows []*entity.Bank
	err  error
}

func (s stubBankRepo) ListActive(context.Context) ([]*entity.Bank, error) {
	return s.rows, s.err
}

func TestLocalBankService_ListAll(t *testing.T) {
	svc := NewLocalBankService(stubBankRepo{rows: []*entity.Bank{
		{ID: 1, Name: "VCB", Code: "VCB", Logo: "logo", Active: true, Description: "d"},
	}}, nil)
	out, err := svc.ListAll(context.Background(), 0)
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if out == nil || out.List == nil {
		t.Fatalf("got %+v", out)
	}
	resp := out.List
	if resp.GetTotal() != 1 || len(resp.GetData()) != 1 {
		t.Fatalf("got %+v", resp)
	}
	if resp.GetData()[0].GetCode() != "VCB" {
		t.Fatalf("bank code: %s", resp.GetData()[0].GetCode())
	}
}
