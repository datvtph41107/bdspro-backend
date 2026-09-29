package dto

// DeepseekRequest - DTO để gọi API DeepSeek
type DeepseekRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepseekMessage `json:"messages"`
	Stream      bool              `json:"stream"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float32           `json:"temperature,omitempty"`
}

// DeepseekMessage - Message trong request gửi đến DeepSeek
type DeepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepseekResponse - Response từ DeepSeek API
type DeepseekResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []DeepseekChoice `json:"choices"`
	Usage   DeepseekUsage    `json:"usage"`
}

// DeepseekChoice - Choice trong response
type DeepseekChoice struct {
	Index        int             `json:"index"`
	Message      DeepseekMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

// DeepseekUsage - Thông tin về usage
type DeepseekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// AddressInfo - Thông tin địa chỉ
type AddressInfo struct {
	ProvinceID *uint64 `json:"provinceId,omitempty"`
	Province   *string `json:"province,omitempty"`
	District   *string `json:"district,omitempty"`
	Ward       *string `json:"ward,omitempty"`
	Detail     *string `json:"detail,omitempty"`
	WardID     *uint64 `json:"wardId,omitempty"`
}

// PriceDataAnalysis - Thông tin giá cả
type PriceDataAnalysis struct {
	SalePrice        *float64 `json:"salePrice,omitempty"`
	SaleCommission   *float64 `json:"saleCommission,omitempty"`
	Deposite         *float64 `json:"deposite,omitempty"`
	RentPrice        *float64 `json:"rentPrice,omitempty"`
	RentCommission   *float64 `json:"rentCommission,omitempty"`
	RentPaymentCycle *uint32  `json:"rentPaymentCycle,omitempty"`
	Currency         *string  `json:"currency,omitempty"`
}

// HouseInfoAnalysis - Thông tin nhà
type HouseInfoAnalysis struct {
	NumBedroom  *int32  `json:"numBedroom,omitempty"`
	NumBathroom *int32  `json:"numBathroom,omitempty"`
	NumFloor    *int32  `json:"numFloor,omitempty"`
	NumFront    *int32  `json:"numFront,omitempty"`
	NumCarPark  *int32  `json:"numCarPark,omitempty"`
	Furniture   *string `json:"furniture,omitempty"`
	Orientation *string `json:"orientation,omitempty"`
}

// ProductPrivateAnalysis - Thông tin nội bộ
type ProductPrivateAnalysis struct {
	ImportPrice   *float64 `json:"importPrice,omitempty"`
	OperatingCost *float64 `json:"operatingCost,omitempty"`
	InternalNote  *string  `json:"internalNote,omitempty"`
	TargetProfit  *float64 `json:"targetProfit,omitempty"`
}

// ProductAnalysisResult - Kết quả phân tích sản phẩm BĐS từ AI
type ProductAnalysisResult struct {
	Name            *string                 `json:"name,omitempty"`
	Area            *float64                `json:"area,omitempty"`
	Description     *string                 `json:"description,omitempty"`
	Note            *string                 `json:"note,omitempty"`
	TransactionType *int32                  `json:"transactionType,omitempty"`
	GoogleMapLink   *string                 `json:"googleMapLink,omitempty"`
	Address         *AddressInfo            `json:"address,omitempty"`
	PropertyType    *string                 `json:"propertyType,omitempty"`
	LegalDoc        *string                 `json:"legalDoc,omitempty"`
	PriceData       *PriceDataAnalysis      `json:"priceData,omitempty"`
	HouseInfo       *HouseInfoAnalysis      `json:"houseInfo,omitempty"`
	ProductPrivate  *ProductPrivateAnalysis `json:"productPrivate,omitempty"`
	Amenities       []string                `json:"amenities,omitempty"`
}
