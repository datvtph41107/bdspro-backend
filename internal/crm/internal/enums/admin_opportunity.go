package enums

// EOpportunityStatus — trạng thái cơ hội IV.10.9.1
type EOpportunityStatus int32

const (
	OpportunityStatusNew             EOpportunityStatus = 10
	OpportunityStatusUnassigned      EOpportunityStatus = 20
	OpportunityStatusAssigned        EOpportunityStatus = 30
	OpportunityStatusContacted       EOpportunityStatus = 40
	OpportunityStatusConsulting      EOpportunityStatus = 50
	OpportunityStatusWaitingQuote    EOpportunityStatus = 60
	OpportunityStatusWaitingApproval EOpportunityStatus = 70
	OpportunityStatusQuoteSent       EOpportunityStatus = 80
	OpportunityStatusNegotiating     EOpportunityStatus = 90
	OpportunityStatusWaitingPayment  EOpportunityStatus = 100
	OpportunityStatusWon             EOpportunityStatus = 110
	OpportunityStatusLost            EOpportunityStatus = 120
	OpportunityStatusCancelled       EOpportunityStatus = 130
)

var OpportunityStatusMap = map[EOpportunityStatus]string{
	OpportunityStatusNew:             "Mới",
	OpportunityStatusUnassigned:      "Chưa gán",
	OpportunityStatusAssigned:        "Đã gán",
	OpportunityStatusContacted:       "Đã liên hệ",
	OpportunityStatusConsulting:      "Đang tư vấn",
	OpportunityStatusWaitingQuote:    "Chờ báo giá",
	OpportunityStatusWaitingApproval: "Chờ duyệt báo giá",
	OpportunityStatusQuoteSent:       "Đã gửi báo giá",
	OpportunityStatusNegotiating:     "Đang đàm phán",
	OpportunityStatusWaitingPayment:  "Chờ thanh toán",
	OpportunityStatusWon:             "Chốt thành công",
	OpportunityStatusLost:            "Thất bại",
	OpportunityStatusCancelled:       "Hủy",
}

func (e EOpportunityStatus) IsValid() bool {
	_, ok := OpportunityStatusMap[e]
	return ok
}

func (e EOpportunityStatus) IsClosed() bool {
	return e == OpportunityStatusWon || e == OpportunityStatusLost || e == OpportunityStatusCancelled
}

// ECustomerType — loại khách IV.10.9
type ECustomerType int32

const (
	CustomerTypeIndividual ECustomerType = 10
	CustomerTypeBusiness   ECustomerType = 20
	CustomerTypeAPIPartner ECustomerType = 30
)

var CustomerTypeMap = map[ECustomerType]string{
	CustomerTypeIndividual: "Cá nhân",
	CustomerTypeBusiness:   "Doanh nghiệp",
	CustomerTypeAPIPartner: "Đối tác API",
}

func (e ECustomerType) IsValid() bool {
	_, ok := CustomerTypeMap[e]
	return ok
}

// EChurnRisk — nguy cơ rời bỏ
type EChurnRisk int32

const (
	ChurnRiskNone   EChurnRisk = 0
	ChurnRiskLow    EChurnRisk = 10
	ChurnRiskMedium EChurnRisk = 20
	ChurnRiskHigh   EChurnRisk = 30
)

var ChurnRiskMap = map[EChurnRisk]string{
	ChurnRiskNone:   "Chưa xác định",
	ChurnRiskLow:    "Thấp",
	ChurnRiskMedium: "Trung bình",
	ChurnRiskHigh:   "Cao",
}

// EWinProbability band
type EWinProbability int32

const (
	WinProbabilityUnknown EWinProbability = 0
	WinProbabilityLow     EWinProbability = 10
	WinProbabilityMedium  EWinProbability = 20
	WinProbabilityHigh    EWinProbability = 30
)

var WinProbabilityMap = map[EWinProbability]string{
	WinProbabilityUnknown: "Chưa xác định",
	WinProbabilityLow:     "Thấp",
	WinProbabilityMedium:  "Trung bình",
	WinProbabilityHigh:    "Cao",
}

// EAdminSalesSource — nguồn phát sinh admin sales (SRS)
type EAdminSalesSource int32

const (
	AdminSourceApp      EAdminSalesSource = 10
	AdminSourceWeb      EAdminSalesSource = 20
	AdminSourceCSKH     EAdminSalesSource = 30
	AdminSourceAdmin    EAdminSalesSource = 40
	AdminSourceBDSPro   EAdminSalesSource = 50
	AdminSourceAPI      EAdminSalesSource = 60
	AdminSourceManual   EAdminSalesSource = 70
	AdminSourcePartner  EAdminSalesSource = 80
)

var AdminSalesSourceMap = map[EAdminSalesSource]string{
	AdminSourceApp:     "App",
	AdminSourceWeb:     "Web",
	AdminSourceCSKH:    "CSKH",
	AdminSourceAdmin:   "Admin",
	AdminSourceBDSPro:  "BDS Pro",
	AdminSourceAPI:     "API/Đối tác",
	AdminSourceManual:  "Nhập thủ công",
	AdminSourcePartner: "Đối tác",
}

func (e EAdminSalesSource) IsValid() bool {
	_, ok := AdminSalesSourceMap[e]
	return ok
}
