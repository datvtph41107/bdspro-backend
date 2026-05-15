package types

// Config - Cấu hình cho resolver engine
type Config struct {
	LegalWeight          float64
	PriorityWeight       float64
	AreaWeight           float64
	LegalScores          map[uint32]float64
	ConflictMatrix       map[string]map[string]bool
	DefaultGroupPriority map[string]int
}

// DefaultConfig - Config mặc định (fallback)
func DefaultConfig() *Config {
	return &Config{
		LegalWeight:    0.60,
		PriorityWeight: 0.30,
		AreaWeight:     0.10,
		LegalScores: map[uint32]float64{
			1000: 100,
			900:  70,
			800:  30,
			700:  15,
			400:  10,
			300:  10,
		},
		ConflictMatrix:       make(map[string]map[string]bool),
		DefaultGroupPriority: make(map[string]int),
	}
}
