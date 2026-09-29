package types

import "time"

// RiskTrend — xu hướng rủi ro
type RiskTrend string

const (
	TrendImproving RiskTrend = "improving"
	TrendStable    RiskTrend = "stable"
	TrendWorsening RiskTrend = "worsening"
	TrendVolatile  RiskTrend = "volatile"
)

// RiskTrendAnalysis — phân tích xu hướng rủi ro
type RiskTrendAnalysis struct {
	ParcelID     uint64
	CurrentRisk  *RiskSnapshot
	PreviousRisk *RiskSnapshot
	BaselineRisk *RiskSnapshot
	Trend        RiskTrend
	TrendDetail  string
	Anomalies    []*RiskAnomaly
	Forecast     *RiskForecast
}

// RiskAnomaly — điểm bất thường trong risk trend
type RiskAnomaly struct {
	Date          time.Time
	RiskScore     uint32
	ExpectedScore uint32
	Deviation     int
	Cause         string
	Severity      string // "minor", "major", "critical"
}

// RiskForecast — dự báo rủi ro
type RiskForecast struct {
	PredictedScore uint32
	PredictedLevel string
	Confidence     float64
	Horizon        string
	Assumptions    []string
}

// HistoricalRiskConfig — cấu hình cho Historical Risk Engine
type HistoricalRiskConfig struct {
	BaselineMonths   int
	PreviousMonths   int
	AnomalyThreshold int
}

// DefaultHistoricalRiskConfig — cấu hình mặc định
func DefaultHistoricalRiskConfig() *HistoricalRiskConfig {
	return &HistoricalRiskConfig{
		BaselineMonths:   12,
		PreviousMonths:   3,
		AnomalyThreshold: 30,
	}
}
