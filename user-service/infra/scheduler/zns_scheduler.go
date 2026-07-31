package scheduler

import (
	"context"
	"log"
	"time"
	"user/internal/interface/providers"
)

// ZnsScheduler quản lý các scheduled jobs cho ZNS
type ZnsScheduler struct {
	znsProvider providers.IZnsProvider
	stopChan    chan struct{}
}

const (
	// Chạy mỗi ngày một lần vào 00:00
	dailyInterval = 24 * time.Hour
	timeLayout    = "2006-01-02 15:04:05"
)

// NewZnsScheduler tạo instance mới của ZnsScheduler
func NewZnsScheduler(znsProvider providers.IZnsProvider) *ZnsScheduler {
	return &ZnsScheduler{
		znsProvider: znsProvider,
		stopChan:    make(chan struct{}),
	}
}

// Start khởi động scheduler
func (s *ZnsScheduler) Start() {
	log.Println("🔧 [ZNS Scheduler] Starting ZNS Token Refresh Scheduler...")

	// Nếu muốn chạy ngay khi khởi động, mở dòng sau:
	// s.refreshTokenJob()

	// Khởi động goroutine để schedule job chạy vào 00:00 mỗi ngày
	go s.scheduleDailyAtMidnight()
}

// Stop dừng scheduler
func (s *ZnsScheduler) Stop() {
	log.Println("🛑 [ZNS Scheduler] Stopping ZNS Token Refresh Scheduler...")
	close(s.stopChan)
}

// scheduleDailyAtMidnight schedule job chạy vào 00:00 mỗi ngày
func (s *ZnsScheduler) scheduleDailyAtMidnight() {
	for {
		next := s.getNextMidnightTime()
		duration := time.Until(next)

		log.Printf("⏰ [ZNS Scheduler] Next refresh token job will run at: %s (in %s)",
			next.Format(timeLayout), duration.Round(time.Second),
		)

		// Tạo timer để chờ đến 00:00 tiếp theo
		timer := time.NewTimer(duration)

		select {
		case <-timer.C:
			// Đến 00:00, thực hiện job
			s.refreshTokenJob()
			// Sau khi chạy xong, lặp lại tiếp vào 00:00 ngày hôm sau
		case <-s.stopChan:
			// Nhận tín hiệu stop, dừng scheduler
			timer.Stop()
			log.Println("✅ [ZNS Scheduler] Scheduler stopped gracefully")
			return
		}
	}
}

// getNextMidnightTime tính thời điểm 00:00 tiếp theo
func (s *ZnsScheduler) getNextMidnightTime() time.Time {
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return nextMidnight
}

// refreshTokenJob job refresh ZNS access token
func (s *ZnsScheduler) refreshTokenJob() {
	log.Println("🔄 [ZNS Scheduler] Starting refresh token job...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gọi refresh token
	newAccessToken, newRefreshToken, err := s.znsProvider.RefreshAccessToken(ctx)
	if err != nil {
		log.Printf("❌ [ZNS Scheduler] Failed to refresh token: %v", err)
		return
	}

	log.Printf("✅ [ZNS Scheduler] Successfully refreshed ZNS token at %s", time.Now().Format(timeLayout))
	log.Printf("📝 [ZNS Scheduler] New Access Token: %s...", prefix(newAccessToken, 50))
	log.Printf("📝 [ZNS Scheduler] New Refresh Token: %s...", prefix(newRefreshToken, 50))

	// TODO: Lưu token mới vào config file hoặc database nếu cần persist
	// Hiện tại token đã được cập nhật trong memory của znsProvider
}

// prefix giới hạn chuỗi đầu ra cho log
func prefix(str string, n int) string {
	if len(str) <= n {
		return str
	}
	return str[:n]
}
