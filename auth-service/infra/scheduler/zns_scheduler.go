package scheduler

import (
	"auth/env"
	"auth/internal/interface/provider"
	"context"
	"log"
	"time"
)

// ZnsScheduler quản lý các scheduled jobs cho ZNS
type ZnsScheduler struct {
	znsProvider provider.IZnsProvider
	stopChan    chan struct{}
}

// NewZnsScheduler tạo instance mới của ZnsScheduler
func NewZnsScheduler(znsProvider provider.IZnsProvider) *ZnsScheduler {
	return &ZnsScheduler{
		znsProvider: znsProvider,
		stopChan:    make(chan struct{}),
	}
}

// Start khởi động scheduler
func (s *ZnsScheduler) Start() {
	log.Println("🔧 [ZNS Scheduler] Starting ZNS Token Refresh Scheduler...")

	// Chạy ngay lần đầu tiên khi khởi động
	if env.ENV_RUNTIME != "" {
		s.refreshTokenJob()
	}

	// Khởi động goroutine để schedule job chạy vào 00:00 mỗi ngày
	go s.scheduleDaily()
}

// Stop dừng scheduler
func (s *ZnsScheduler) Stop() {
	log.Println("🛑 [ZNS Scheduler] Stopping ZNS Token Refresh Scheduler...")
	close(s.stopChan)
}

// scheduleDaily schedule job chạy vào 00:00 mỗi ngày
func (s *ZnsScheduler) scheduleDaily() {
	for {
		// Tính toán thời gian còn lại đến 00:00 ngày tiếp theo
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := next.Sub(now)

		log.Printf("⏰ [ZNS Scheduler] Next refresh token job will run at: %s (in %s)",
			next.Format("2006-01-02 15:04:05"),
			duration.Round(time.Second))

		// Tạo timer để chờ đến 00:00
		timer := time.NewTimer(duration)

		select {
		case <-timer.C:
			// Đến 00:00, chạy job
			s.refreshTokenJob()

		case <-s.stopChan:
			// Nhận tín hiệu stop, dừng scheduler
			timer.Stop()
			log.Println("✅ [ZNS Scheduler] Scheduler stopped gracefully")
			return
		}
	}
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

	log.Printf("✅ [ZNS Scheduler] Successfully refreshed ZNS token at %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Printf("📝 [ZNS Scheduler] New Access Token: %s...", newAccessToken[:50])
	log.Printf("📝 [ZNS Scheduler] New Refresh Token: %s...", newRefreshToken[:50])

	// TODO: Lưu token mới vào config file hoặc database nếu cần persist
	// Hiện tại token đã được cập nhật trong memory của znsProvider
}
