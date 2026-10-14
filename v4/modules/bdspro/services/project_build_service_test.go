package services

import (
	"context"
	"testing"

	_dto "common/domain/dto"
	_enums "common/domain/enum"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
)

type stubProjectBuildRepo struct {
	mine []entity.ProjectBuild
}

func (s stubProjectBuildRepo) Search(_ context.Context, _ _dto.Pagable, _ string) ([]entity.ProjectBuild, int64, error) {
	return nil, 0, nil
}

func (s stubProjectBuildRepo) SearchByCreatedBy(_ context.Context, _ uint64, _ _dto.Pagable, _ string) ([]entity.ProjectBuild, int64, error) {
	return s.mine, int64(len(s.mine)), nil
}

func (s stubProjectBuildRepo) GetByID(context.Context, uint64) (*entity.ProjectBuild, error) {
	return nil, nil
}

func (s stubProjectBuildRepo) Create(context.Context, *entity.ProjectBuild) (*entity.ProjectBuild, error) {
	return nil, nil
}

func (s stubProjectBuildRepo) Update(context.Context, *entity.ProjectBuild) (*entity.ProjectBuild, error) {
	return nil, nil
}

func (s stubProjectBuildRepo) Delete(context.Context, uint64) error {
	return nil
}

func TestLocalProjectBuildService_ListMine(t *testing.T) {
	build := entity.ProjectBuild{Name: "Tòa A"}
	build.ID = 9
	svc := NewLocalProjectBuildService(stubProjectBuildRepo{mine: []entity.ProjectBuild{build}})

	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))
	resp, err := svc.ListMine(ctx, dto.CatalogSearchRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.GetData()) != 1 || resp.GetData()[0].GetName() != "Tòa A" {
		t.Fatalf("unexpected: %+v", resp.GetData())
	}
}
