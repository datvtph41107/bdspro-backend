package services

import (
	"context"
	"testing"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/repository"
)

type stubCostTypeRepo struct{}

func (stubCostTypeRepo) Create(context.Context, *entity.TxCostType) error { return nil }
func (stubCostTypeRepo) Update(context.Context, *entity.TxCostType) error { return nil }
func (stubCostTypeRepo) Delete(context.Context, uint64) error             { return nil }
func (stubCostTypeRepo) GetByID(context.Context, uint64) (*entity.TxCostType, error) {
	return nil, nil
}
func (stubCostTypeRepo) List(context.Context, dto.CostTypeSearchRequest) ([]*entity.TxCostType, error) {
	return []*entity.TxCostType{{ID: 1, TypeName: "Demo", CostType: 10}}, nil
}

func TestLocalCostTypeService_List_NilSync(t *testing.T) {
	svc := NewLocalCostTypeService(stubCostTypeRepo{}, nil)
	out, err := svc.List(context.Background(), dto.CostTypeSearchRequest{CostType: 10})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if out == nil || out.List == nil || len(out.List.Data) != 1 {
		t.Fatalf("expected one row, got %#v", out)
	}
}

var _ repository.CostTypeRepository = stubCostTypeRepo{}
