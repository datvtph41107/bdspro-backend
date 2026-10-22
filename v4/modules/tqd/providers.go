package tqd

import (
	repo "mdn4/modules/tqd/repository"
	"mdn4/modules/tqd/services"
	genreportapp "mdn4/modules/tqd/services/generatedreport/application"
	"mdn4/modules/tqd/services/quota"
)

// ProvideReportUsecase builds the report usecase without user/auth RPC enrichment on v4.
// @wire:provider
func ProvideReportUsecase(reportRepo repo.ReportRepository) services.ReportUsecase {
	return services.NewReportUsecase(reportRepo, nil, nil, nil)
}

// ProvideGeneratedReportService is optional until full report runtime is wired.
// @wire:provider
func ProvideGeneratedReportService() *genreportapp.Service {
	return nil
}

// ProvideProfileUsageService is optional until access/durable/runtime stores are wired.
// @wire:provider
func ProvideProfileUsageService() *quota.ProfileUsageService {
	return nil
}
