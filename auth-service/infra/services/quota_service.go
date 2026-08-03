package services

import "context"

// cái này chứa toàn bộ thiết lập ở db trong bảng pack_item với bảng pack_map
type QuotaItem struct {
	// Quota  uint64
	// Method string
	// Path   string
}

type QuotaService struct {
	// nạp toàn bộ db liên quan cấu hình gói như api nào
	// được gọi bao nhiêu lần trong đây
	packItems map[int]QuotaItem

	// lưu lượng quota còn lại của user
	// profileId -> quota
	// ngâm cứu dùng kiểu mảng để tối ưu bộ nhớ,
	// key là profileId, value là ds apiName được phép gọi bao nhiêu lần nữa
	// value là mảng số lần gọi còn lại thì dựa vào quota_map.yml xem có bao nhiêu api config
	// thì tạo bấy nhiêu phần tử
	users map[uint64][]uint32
}

type QuotaRequest struct {
	ProfileId uint64
	ApiName   string
	Method    string
	Path      string
}

func NewQuotaService() *QuotaService {
	// viết đoạn load hết nối gói dịch vụ với keyPath ở đây
	return &QuotaService{}
}

// kiểm tra xem user có đủ quota để gọi api không
func (s *QuotaService) Verify(ctx context.Context, request QuotaRequest) (uint64, error) {
	// packId thêm vào token
	// lấy packId trong context -> truy cập vào packItems
	// có được gói hiện tại đang dùng thì lấy được ds API bị tính quota
	// xem ds quota có cho phép gọi api không, nếu không thì reject
	// nếu có thì lấy profileId từ context -> truy cập vào users
	// rồi trừ trực tiếp luôn, nhớ dùng syncKey theo id đảm bảo
	// 1 user chỉ gọi 1 lần

	//

	// nhớ tạo 1 goroutine lưu luôn quota vào trong file log ở đây

	// tạo schedule

	return 0, nil
}

func (s *QuotaService) GetQuota(ctx context.Context, apiName string) (uint64, error) {
	return 0, nil
}

// gọi về user service để cập nhật quota
func (s *QuotaService) StorageQuota(ctx context.Context, apiName string) error {
	return nil
}
