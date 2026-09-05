package admin

import (
	"context"
	"errors"
	"testing"
)

type repositoryProbe struct {
	query Query
	page  Page
	err   error
}

func (p *repositoryProbe) ListAdminGeneratedReports(_ context.Context, query Query) (Page, error) {
	p.query = query
	return p.page, p.err
}

func TestServiceNormalizesPaginationAndJobStatus(t *testing.T) {
	repository := &repositoryProbe{page: Page{Total: 1}}
	page, err := NewService(repository).List(context.Background(), Query{JobStatus: " RUNNING "})
	if err != nil {
		t.Fatal(err)
	}
	if repository.query.Page != 1 || repository.query.PageSize != 20 || repository.query.JobStatus != "running" {
		t.Fatalf("normalized query = %#v", repository.query)
	}
	if page.Total != 1 {
		t.Fatalf("total = %d", page.Total)
	}
}

func TestServiceRejectsInvalidQueryAndUnavailableRepository(t *testing.T) {
	if _, err := NewService(&repositoryProbe{}).List(context.Background(), Query{PageSize: 101}); err == nil {
		t.Fatal("oversized page accepted")
	}
	if _, err := NewService(&repositoryProbe{}).List(context.Background(), Query{JobStatus: "unknown"}); err == nil {
		t.Fatal("invalid job status accepted")
	}
	if _, err := NewService(nil).List(context.Background(), Query{}); err == nil {
		t.Fatal("nil repository accepted")
	}
	expected := errors.New("database unavailable")
	if _, err := NewService(&repositoryProbe{err: expected}).List(context.Background(), Query{}); !errors.Is(err, expected) {
		t.Fatalf("repository error = %v", err)
	}
}
