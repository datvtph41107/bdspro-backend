package enums

// ProAIJobType — loại job trên bảng pro_ai_jobs (mở rộng sau).
type ProAIJobType uint32

const (
	ProAIJobTypePlanningProject uint32 = 10 // Đồ án → gắn qh_planning_projects / documents
)

var ProAIJobTypeMap = map[uint32]string{
	ProAIJobTypePlanningProject: "Đồ án",
}

func GetProAIJobTypeLabel(jobType uint32) string {
	return ProAIJobTypeMap[jobType]
}

func IsSupportedProAIJobType(jobType uint32) bool {
	_, ok := ProAIJobTypeMap[jobType]
	return ok
}
