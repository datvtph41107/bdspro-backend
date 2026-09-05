package wallet

import "time"

type DashboardMetric struct {
	Id              uint32
	Time            time.Time
	ProcessingCount uint64
	CompletedCount  uint64
}
