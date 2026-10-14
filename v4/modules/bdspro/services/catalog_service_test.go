package services

import (
	"context"
	"testing"

	_dto "common/domain/dto"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
)

type stubCatalogRepo struct {
	projects []entity.ProjectItem
	builds   []entity.ProjectBuild
}

func (s stubCatalogRepo) SearchProjects(_ context.Context, _ _dto.Pagable, _ string) ([]entity.ProjectItem, int64, error) {
	return s.projects, int64(len(s.projects)), nil
}

func (s stubCatalogRepo) SearchProjectBuilds(_ context.Context, _ _dto.Pagable, _ string) ([]entity.ProjectBuild, int64, error) {
	return s.builds, int64(len(s.builds)), nil
}

func TestLocalCatalogService_ListProjects(t *testing.T) {
	svc := NewLocalCatalogService(stubCatalogRepo{
		projects: []entity.ProjectItem{{ID: 1, Name: "Dự án A"}},
	})
	resp, err := svc.ListProjects(context.Background(), dto.CatalogSearchRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.GetData()) != 1 || resp.GetData()[0].GetName() != "Dự án A" {
		t.Fatalf("unexpected: %+v", resp.GetData())
	}
	if resp.GetTotal() != 1 {
		t.Fatalf("total: %d", resp.GetTotal())
	}
}
