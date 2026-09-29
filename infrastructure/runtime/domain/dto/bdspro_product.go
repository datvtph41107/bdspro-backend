package _dto

type ProductV3DTO struct {
	ID              uint64          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	PropertyType    *ItemDTO        `json:"propertyType"`
	DocType         *ItemDTO        `json:"docType"`
	Amenities       []ItemDTO       `json:"amenities"`
	Address         *AddressV3DTO   `json:"addressV3"`
	Area            float32         `json:"area"`
	TransactionType *ItemDTO        `json:"transactionType"`
	PriceData       *PriceV3DTO     `json:"priceData"`
	HouseInfo       *HouseInfoV3DTO `json:"houseInfo"`
	LegalDoc        *ItemDTO        `json:"legalDoc"`
	SourceType      *ItemDTO        `json:"sourceType"`
}

type PriceV3DTO struct {
	Currency           string  `json:"currency"`
	SalePrice          float32 `json:"salePrice"`
	SaleCommission     float32 `json:"saleCommission"`
	SaleCommissionType uint32  `json:"saleCommissionType"`
	RentPrice          float32 `json:"rentPrice"`
	RentCommission     float32 `json:"rentCommission"`
	RentCommissionType uint32  `json:"rentCommissionType"`
	RentPaymentCycle   uint32  `json:"rentPaymentCycle"`
}

type HouseInfoV3DTO struct {
	NumBedroom  int32   `json:"numBedroom"`
	NumBathroom int32   `json:"numBathroom"`
	NumFloor    int32   `json:"numFloor"`
	NumFront    float32 `json:"numFront"`
	NumCarPark  int32   `json:"numCarPark"`
	Orientation string  `json:"orientation"`
	Furniture   string  `json:"furniture"`
}
