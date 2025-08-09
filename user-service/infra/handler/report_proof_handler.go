package handler

import (
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type ReportProofHandler struct {
	userpb.UnimplementedReportProofServiceServer
	reportUsecase *usecases.ReportUsecase
	reportMapper  *mapper.ReportMapper
}

func NewReportProofHandler(
	reportUsecase *usecases.ReportUsecase,
	reportMapper *mapper.ReportMapper,
) *ReportProofHandler {
	return &ReportProofHandler{
		reportUsecase: reportUsecase,
		reportMapper:  reportMapper,
	}
}
