package dashboard

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type DashboardStatsJob interface {
	Run(ctx context.Context) error
	RunDailyAtMidnight(ctx context.Context)
	GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*DashboardStats, *DashboardStats, error)
}

type dashboardStatsJob struct {
	Usecase   DashboardStatsUsecase
	TableName string
	GetValue  func(context.Context) (int64, error)
}

func NewDashboardStatsJob(db *gorm.DB, tableName string, getValue func(context.Context) int64) DashboardStatsJob {
	return NewDashboardStatsJobWithError(db, tableName, func(ctx context.Context) (int64, error) {
		return getValue(ctx), nil
	})
}

func NewDashboardStatsJobWithError(db *gorm.DB, tableName string, getValue func(context.Context) (int64, error)) DashboardStatsJob {
	// Auto migrate table trước khi tạo job
	// err := db.Table(tableName).AutoMigrate(&DashboardStats{})
	// if err != nil {
	// 	log.Printf("Error auto migrating table %s: %v", tableName, err)
	// } else {
	// 	log.Printf("Successfully auto migrated table %s", tableName)
	// }

	repo := NewDashboardStatsPostgresRepository(db)
	usecase := NewDashboardStatsUsecase(repo)
	return &dashboardStatsJob{
		Usecase:   usecase,
		TableName: tableName,
		GetValue:  getValue,
	}
}

func (j *dashboardStatsJob) Run(ctx context.Context) error {
	// Lấy giá trị từ function. A technical read failure is not a valid metric.
	count, err := j.GetValue(ctx)
	if err != nil {
		return err
	}

	// Tạo thời gian hiện tại (ngày hiện tại)
	now := time.Now()
	calculateTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Tạo stats object
	stats := &DashboardStats{
		Count:         count,
		CalculateTime: calculateTime,
	}

	// Lưu hoặc cập nhật vào database
	_, err = j.Usecase.UpdateOrCreate(ctx, j.TableName, stats)
	if err != nil {
		return err
	}

	log.Printf("Dashboard stats updated for table %s: count=%d, date=%s",
		j.TableName, count, calculateTime.Format("2006-01-02"))
	return nil
}

// RunDailyAtMidnight blocks as one process-owned actor. It schedules the next
// actual midnight instead of polling on an arbitrary process-start offset.
func (j *dashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			log.Printf("Dashboard stats job stopped for table %s", j.TableName)
			return
		case <-timer.C:
			if err := j.Run(ctx); err != nil {
				log.Printf("Error running dashboard stats job for table %s: %v", j.TableName, err)
			}
		}
	}
}

func (j *dashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*DashboardStats, *DashboardStats, error) {
	return j.Usecase.GetByTwoTime(ctx, j.TableName, firstDate, secondDate)
}
