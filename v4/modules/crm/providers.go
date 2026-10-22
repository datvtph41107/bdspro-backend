package crm

import (
	"context"

	dashboard "base/dashboard"

	"gorm.io/gorm"

	"mdn4/modules/crm/services"
)

// ProvideAppointmentStatsJob reads appointment_stats snapshots; the daily writer is not scheduled in v4.
// @wire:provider
func ProvideAppointmentStatsJob(crm *gorm.DB, appointments services.AppointmentUsecase) dashboard.DashboardStatsJob {
	return dashboard.NewDashboardStatsJob(crm, "appointment_stats", func(ctx context.Context) int64 {
		count, err := appointments.CountCurrent(ctx)
		if err != nil {
			return 0
		}
		return count
	})
}
