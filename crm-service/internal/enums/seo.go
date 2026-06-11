package enums

type ESeoClassify uint32

const (
	// Chứa thông tin chi tiết, trang detail phục vụ long-tail SEO.
	ESeoClassifyA ESeoClassify = 10
	// Chứa điều hướng, guide/category/hub page.
	ESeoClassifyB ESeoClassify = 20
	// Kết quả dữ liệu tìm kiếm/faceted search.
	ESeoClassifyC ESeoClassify = 30
	// Link chia sẻ/viral/preview không phải page SEO chính.
	ESeoClassifyD ESeoClassify = 40
	// Không phục vụ SEO: login, register, admin, 404, 500...
	ESeoClassifyN ESeoClassify = 50
)

func IsValidSeoClassify(v uint32) bool {
	switch ESeoClassify(v) {
	case ESeoClassifyA, ESeoClassifyB, ESeoClassifyC, ESeoClassifyD, ESeoClassifyN:
		return true
	default:
		return false
	}
}

func SeoClassifyFromUint32(v uint32) ESeoClassify {
	if IsValidSeoClassify(v) {
		return ESeoClassify(v)
	}
	return ESeoClassifyN
}

func (e ESeoClassify) Uint32() uint32 {
	return uint32(e)
}
