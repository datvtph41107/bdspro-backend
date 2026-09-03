package handler_grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	tqdpb "pb/types/tqd"
	reportadmin "tqd/internal/usecase/generatedreport/admin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type reportAdminRepositoryStub struct {
	query reportadmin.Query
	page  reportadmin.Page
}

func (s *reportAdminRepositoryStub) ListAdminGeneratedReports(_ context.Context, query reportadmin.Query) (reportadmin.Page, error) {
	s.query = query
	return s.page, nil
}

type usageAuthorizerStub struct{ err error }

func (s usageAuthorizerStub) HasPermissions(_ context.Context, _ []string) error { return s.err }

func TestListGeneratedReportsAuthorizesFiltersAndMapsEvidence(t *testing.T) {
	now := time.Date(2026, 8, 31, 10, 0, 0, 0, time.UTC)
	repository := &reportAdminRepositoryStub{page: reportadmin.Page{
		Total: 1, Page: 2, PageSize: 10,
		Reports: []reportadmin.Report{{
			ReportID: 41, ProfileID: 7, ReportType: 1, Status: 20,
			Title: "Parcel report", ParcelID: 99, PDFURL: "https://files/report.pdf",
			JobID: "job-41", JobStatus: "running", JobAttempts: 2,
			JobClaimVersion: 3, CreatedAt: now, UpdatedAt: now,
		}},
	}}
	handler := NewAdminUsageGrpcHandler(nil, reportadmin.NewService(repository), usageAuthorizerStub{}, nil)

	response, err := handler.ListGeneratedReports(context.Background(), &tqdpb.ListAdminGeneratedReportsRequest{
		Page: 2, PageSize: 10, ProfileId: 7, Status: 20, JobStatus: "RUNNING",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repository.query.JobStatus != "running" || repository.query.ProfileID != 7 || repository.query.Status != 20 {
		t.Fatalf("query = %#v", repository.query)
	}
	if response.GetTotal() != 1 || len(response.GetReports()) != 1 {
		t.Fatalf("response = %#v", response)
	}
	item := response.GetReports()[0]
	if item.GetReportId() != 41 || item.GetProfileId() != 7 || item.GetJobClaimVersion() != 3 || item.GetPdfUrl() == "" {
		t.Fatalf("projection = %#v", item)
	}
}

func TestListGeneratedReportsFailsClosedWithoutPermission(t *testing.T) {
	handler := NewAdminUsageGrpcHandler(nil, reportadmin.NewService(&reportAdminRepositoryStub{}), usageAuthorizerStub{err: errors.New("denied")}, nil)
	_, err := handler.ListGeneratedReports(context.Background(), &tqdpb.ListAdminGeneratedReportsRequest{})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("code = %v, want PermissionDenied", status.Code(err))
	}

	handler = NewAdminUsageGrpcHandler(nil, reportadmin.NewService(&reportAdminRepositoryStub{}), nil, nil)
	_, err = handler.ListGeneratedReports(context.Background(), &tqdpb.ListAdminGeneratedReportsRequest{})
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("code = %v, want Unavailable", status.Code(err))
	}
}
