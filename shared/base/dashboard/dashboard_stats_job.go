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
	GetValue  func() int64 // Function để lấy giá trị cần lưu
}

func NewDashboardStatsJob(db *gorm.DB, tableName string, getValue func() int64) DashboardStatsJob {
	// Auto migrate table trước khi tạo job
	err := db.Table(tableName).AutoMigrate(&DashboardStats{})
	if err != nil {
		log.Printf("Error auto migrating table %s: %v", tableName, err)
	} else {
		log.Printf("Successfully auto migrated table %s", tableName)
	}

	repo := NewDashboardStatsPostgresRepository(db)
	usecase := NewDashboardStatsUsecase(repo)
	return &dashboardStatsJob{
		Usecase:   usecase,
		TableName: tableName,
		GetValue:  getValue,
	}
}

func (j *dashboardStatsJob) Run(ctx context.Context) error {
	// Lấy giá trị từ function
	count := j.GetValue()

	// Tạo thời gian hiện tại (ngày hiện tại)
	now := time.Now()
	calculateTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Tạo stats object
	stats := &DashboardStats{
		Count:         count,
		CalculateTime: calculateTime,
	}

	// Lưu hoặc cập nhật vào database
	_, err := j.Usecase.UpdateOrCreate(ctx, j.TableName, stats)
	if err != nil {
		return err
	}

	log.Printf("Dashboard stats updated for table %s: count=%d, date=%s",
		j.TableName, count, calculateTime.Format("2006-01-02"))
	return nil
}

// RunDailyAtMidnight chạy job hàng ngày vào 0h00
func (j *dashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Check mỗi phút
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Dashboard stats job stopped for table %s", j.TableName)
			return
		case now := <-ticker.C:
			// Kiểm tra nếu đã qua 0h00 của ngày mới
			if now.Hour() == 0 && now.Minute() == 0 {
				log.Printf("Running dashboard stats job for table %s at %s",
					j.TableName, now.Format("2006-01-02 15:04:05"))

				if err := j.Run(ctx); err != nil {
					log.Printf("Error running dashboard stats job for table %s: %v",
						j.TableName, err)
				}

				// Sleep 1 phút để tránh chạy nhiều lần trong cùng 1 phút
				time.Sleep(1 * time.Minute)
			}
		}
	}
}

func (j *dashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*DashboardStats, *DashboardStats, error) {
	return j.Usecase.GetByTwoTime(ctx, j.TableName, firstDate, secondDate)
}
