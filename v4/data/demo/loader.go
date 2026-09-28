package demo

import (
	_ "embed"
	"encoding/json"
	"fmt"

	bdspropb "pb/types/bdspro"

	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed products.json
var productsJSON []byte

//go:embed contacts.json
var contactsJSON []byte

//go:embed medias.json
var mediasJSON []byte

// DemoContact là snapshot user_profile dùng để seed liên hệ CRM khi bật demo.
type DemoContact struct {
	ProfileID uint64 `json:"profileId"`
	FullName  string `json:"fullName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
}

// DemoMedia là media_url tái sử dụng từ product_media (id 803→788).
type DemoMedia struct {
	ID        uint64 `json:"id"`
	MediaURL  string `json:"mediaUrl"`
	MediaType string `json:"mediaType"`
}

// LoadDemoProducts đọc danh sách body tạo sản phẩm từ data/demo/products.json.
func LoadDemoProducts() ([]*bdspropb.ProductSaveRequest, error) {
	var rawItems []json.RawMessage
	if err := json.Unmarshal(productsJSON, &rawItems); err != nil {
		return nil, fmt.Errorf("parse demo products.json: %w", err)
	}

	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	products := make([]*bdspropb.ProductSaveRequest, 0, len(rawItems))
	for i, raw := range rawItems {
		var product bdspropb.ProductSaveRequest
		if err := opts.Unmarshal(raw, &product); err != nil {
			return nil, fmt.Errorf("parse demo product[%d]: %w", i, err)
		}
		products = append(products, &product)
	}
	return products, nil
}

// LoadDemoContacts đọc danh sách profile dùng để tạo liên hệ CRM khi seed demo.
func LoadDemoContacts() ([]DemoContact, error) {
	var contacts []DemoContact
	if err := json.Unmarshal(contactsJSON, &contacts); err != nil {
		return nil, fmt.Errorf("parse demo contacts.json: %w", err)
	}
	return contacts, nil
}

// LoadDemoMedias đọc media_url snapshot từ product_media id 803→788.
func LoadDemoMedias() ([]DemoMedia, error) {
	var medias []DemoMedia
	if err := json.Unmarshal(mediasJSON, &medias); err != nil {
		return nil, fmt.Errorf("parse demo medias.json: %w", err)
	}
	return medias, nil
}

// MediaItemsForProduct gán 1–2 ảnh demo theo index sản phẩm (round-robin trên pool).
func MediaItemsForProduct(medias []DemoMedia, productIndex int) []*bdspropb.MediaItem {
	if len(medias) == 0 {
		return nil
	}
	primary := medias[productIndex%len(medias)]
	items := []*bdspropb.MediaItem{
		{
			MediaUrl:  primary.MediaURL,
			MediaType: primary.MediaType,
			IsMain:    true,
			Order:     1,
		},
	}
	if len(medias) > 1 {
		secondary := medias[(productIndex+len(medias)/2)%len(medias)]
		if secondary.MediaURL != primary.MediaURL {
			items = append(items, &bdspropb.MediaItem{
				MediaUrl:  secondary.MediaURL,
				MediaType: secondary.MediaType,
				IsMain:    false,
				Order:     2,
			})
		}
	}
	return items
}
