package client

// BdsproClient is a thin adapter placeholder used by CRM mappers.
// Wire real pb/clients when CRM product enrichment is enabled in v4.
type BdsproClient struct{}

func NewBdsproClient() *BdsproClient {
	return &BdsproClient{}
}
