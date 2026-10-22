package mapper

import (
	"testing"
	"time"

	bdspropb "pb/types/bdspro"
)

func TestCreatePropertyRequestToBundle_Asset(t *testing.T) {
	title := "Nhà phố Q1"
	typeID := uint64(3)
	provinceID := uint64(79)
	area := 85.5
	issued := time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	purchase := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)

	req := &bdspropb.CreatePropertyRequest{
		Info:     &bdspropb.PropertyInfo{Title: &title, PropertyTypeId: &typeID},
		Location: &bdspropb.PropertyLocationCreate{ProvinceId: &provinceID, AddressDetail: "12 Lê Lợi"},
		LandInfo: &bdspropb.PropertyLandInfoCreate{AreaTotal: &area},
		Asset: &bdspropb.AssetCreateWithProduct{
			PurchasePrice: 5e9,
			PurchaseDate:  int32(purchase.Unix()),
			RentStatus:    20,
			Description:   "mua lại",
			LegalItems:    []*bdspropb.LegalItem{{DocumentName: "Sổ hồng", DocumentUrl: "https://x/1.jpg"}},
		},
		LegalInfo: []*bdspropb.AssetLegalItem{{DocumentName: "HĐ mua bán", IssuedDate: &issued}},
	}

	bundle := CreatePropertyRequestToBundle(req, 42)
	a := bundle.Asset
	if a == nil {
		t.Fatal("expected asset")
	}
	if !bundle.IncludeAsset {
		t.Fatal("IncludeAsset should be true")
	}
	if a.Name != title || a.PropertyTypeID == nil || *a.PropertyTypeID != typeID {
		t.Fatalf("name/type: %q %v", a.Name, a.PropertyTypeID)
	}
	if a.ProvinceID == nil || *a.ProvinceID != provinceID || a.Address != "12 Lê Lợi" || a.Area != area {
		t.Fatalf("location/area: %+v", a)
	}
	if a.OwnerID == nil || *a.OwnerID != 42 {
		t.Fatalf("owner: %v", a.OwnerID)
	}
	if a.PurchasePrice != 5e9 || a.RentStatus != 20 || a.LegalStatus != 10 {
		t.Fatalf("price/status: %+v", a)
	}
	if a.PurchaseDate == nil || !a.PurchaseDate.Equal(purchase) {
		t.Fatalf("purchaseDate: %v", a.PurchaseDate)
	}
	if len(a.LegalItems) != 2 {
		t.Fatalf("legal items: %d", len(a.LegalItems))
	}
	if a.LegalItems[1].IssuedDate == nil || a.LegalItems[1].IssuedDate.UnixMilli() != issued {
		t.Fatalf("issuedDate: %v", a.LegalItems[1].IssuedDate)
	}
}

func TestCreatePropertyRequestToBundle_NoAsset(t *testing.T) {
	bundle := CreatePropertyRequestToBundle(&bdspropb.CreatePropertyRequest{}, 1)
	if bundle.Asset != nil || bundle.IncludeAsset {
		t.Fatalf("unexpected asset: %+v", bundle.Asset)
	}
}
