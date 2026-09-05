package rule

import (
	"context"
	"fmt"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// HistoricalRiskRepository — interface cho Historical Risk Engine
type HistoricalRiskRepository interface {
	GetExpiringLayers(ctx context.Context, parcelID uint64, monthsAhead int) (int, error)
	GetUpcomingLayers(ctx context.Context, parcelID uint64, monthsAhead int) (int, error)
}

// HistoricalRiskEngine — Engine phân tích xu hướng rủi ro lịch sử
type HistoricalRiskEngine struct {
	config      *types.HistoricalRiskConfig
	riskRepo    HistoricalRiskRepository
	timelineEng *TimelineEngine
}

// NewHistoricalRiskEngine tạo HistoricalRiskEngine
func NewHistoricalRiskEngine(
	config *types.HistoricalRiskConfig,
	riskRepo HistoricalRiskRepository,
	timelineEng *TimelineEngine,
) *HistoricalRiskEngine {
	if config == nil {
		config = types.DefaultHistoricalRiskConfig()
	}
	return &HistoricalRiskEngine{
		config:      config,
		riskRepo:    riskRepo,
		timelineEng: timelineEng,
	}
}

// Name trả về tên engine
func (e *HistoricalRiskEngine) Name() string { return "HistoricalRiskEngine" }

// AnalyzeTrend phân tích xu hướng rủi ro
func (e *HistoricalRiskEngine) AnalyzeTrend(
	ctx context.Context, parcelID uint64, currentRisk *types.RiskSnapshot,
) (*types.RiskTrendAnalysis, error) {
	now := time.Now()

	// B1: Tính risk tại thời điểm previous (VD: 3 tháng trước)
	previousDate := now.AddDate(0, -e.config.PreviousMonths, 0)
	previousRisk := e.computeRiskAt(ctx, parcelID, previousDate)

	// B2: Tính risk tại baseline (VD: 1 năm trước)
	baselineDate := now.AddDate(0, -e.config.BaselineMonths, 0)
	baselineRisk := e.computeRiskAt(ctx, parcelID, baselineDate)

	// B3: Xác định trend
	trend, detail := e.determineTrend(baselineRisk, previousRisk, currentRisk)

	// B4: Dự báo
	forecast := e.forecast(ctx, parcelID, currentRisk)

	return &types.RiskTrendAnalysis{
		ParcelID:     parcelID,
		CurrentRisk:  currentRisk,
		PreviousRisk: previousRisk,
		BaselineRisk: baselineRisk,
		Trend:        trend,
		TrendDetail:  detail,
		Forecast:     forecast,
	}, nil
}

// computeRiskAt tính risk tại 1 thời điểm (simplified)
func (e *HistoricalRiskEngine) computeRiskAt(ctx context.Context, parcelID uint64, at time.Time) *types.RiskSnapshot {
	return &types.RiskSnapshot{
		Date:      at,
		RiskScore: 0,
		RiskLevel: "UNKNOWN",
	}
}

// determineTrend xác định xu hướng dựa trên 3 điểm
func (e *HistoricalRiskEngine) determineTrend(
	baseline, previous, current *types.RiskSnapshot,
) (types.RiskTrend, string) {
	if baseline == nil || current == nil {
		return types.TrendStable, "Không đủ dữ liệu lịch sử để xác định xu hướng"
	}

	delta := int(current.RiskScore) - int(baseline.RiskScore)

	switch {
	case delta <= -20:
		return types.TrendImproving, fmt.Sprintf(
			"Rủi ro giảm %d điểm so với %s (từ %s xuống %s)",
			-delta, baseline.Date.Format("02/01/2006"), baseline.RiskLevel, current.RiskLevel,
		)
	case delta >= 20:
		return types.TrendWorsening, fmt.Sprintf(
			"Rủi ro tăng %d điểm so với %s (từ %s lên %s)",
			delta, baseline.Date.Format("02/01/2006"), baseline.RiskLevel, current.RiskLevel,
		)
	default:
		return types.TrendStable, fmt.Sprintf("Rủi ro ổn định ở mức %s", current.RiskLevel)
	}
}

// forecast dự báo risk trong tương lai
func (e *HistoricalRiskEngine) forecast(
	ctx context.Context, parcelID uint64, current *types.RiskSnapshot,
) *types.RiskForecast {
	expiringCount, _ := e.riskRepo.GetExpiringLayers(ctx, parcelID, 6)
	upcomingCount, _ := e.riskRepo.GetUpcomingLayers(ctx, parcelID, 6)

	predictedScore := current.RiskScore
	var assumptions []string

	if expiringCount > 0 {
		predictedScore -= uint32(expiringCount * 5)
		assumptions = append(assumptions,
			fmt.Sprintf("%d layer sắp hết hạn trong 6 tháng tới", expiringCount))
	}
	if upcomingCount > 0 {
		predictedScore += uint32(upcomingCount * 10)
		assumptions = append(assumptions,
			fmt.Sprintf("%d layer sắp có hiệu lực trong 6 tháng tới", upcomingCount))
	}

	return &types.RiskForecast{
		PredictedScore: predictedScore,
		PredictedLevel: e.classifyRiskLevel(predictedScore),
		Confidence:     0.7,
		Horizon:        "6 months",
		Assumptions:    assumptions,
	}
}

// classifyRiskLevel phân loại risk level từ score
func (e *HistoricalRiskEngine) classifyRiskLevel(score uint32) string {
	switch {
	case score >= 85:
		return "CRITICAL"
	case score >= 60:
		return "HIGH"
	case score >= 30:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// Validate kiểm tra config
func (e *HistoricalRiskEngine) Validate() error {
	return nil
}
