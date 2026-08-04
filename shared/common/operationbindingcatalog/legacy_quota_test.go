package operationbindingcatalog

import (
	_operationbinding "common/operationbinding"

	"errors"
	"testing"
)

func TestLegacyQuotaReviewCatalogIsReviewedButNotResolvable(t *testing.T) {
	t.Parallel()

	registry, err := _operationbinding.New(LegacyQuotaReviewV1(), nil)
	if err != nil {
		t.Fatalf("operationbinding.New() error = %v", err)
	}

	reviewed := []_operationbinding.Transport{
		_operationbinding.GRPCUnary("/tqdpb.ParcelService/CreateBatchParcels"),
		_operationbinding.GRPCUnary("/tqdpb.ParcelService/FindParcelByLocation"),
		_operationbinding.GRPCUnary("/tqdpb.ParcelService/FindParcelsByPolygon"),
		_operationbinding.GRPCUnary("/tqdpb.ParcelService/CheckingPolygon"),
		_operationbinding.GRPCUnary("/tqdpb.ParcelService/GetParcelPlanning"),
		_operationbinding.GRPCUnary("/userpb.ProfileService/GetProfileInfo"),
		_operationbinding.GRPCUnary("/userpb.ProfileService/GetProfileInfoV3"),
		_operationbinding.GRPCUnary("/userpb.ProfileService/GetProfileInfoMe"),
	}
	if err := registry.RequireReviewed(reviewed); err != nil {
		t.Fatalf("RequireReviewed() error = %v", err)
	}
	if err := registry.RequireResolvable(reviewed); !errors.Is(err, _operationbinding.ErrNotResolvable) {
		t.Fatalf("RequireResolvable() error = %v, want ErrNotResolvable", err)
	}

	missing, ok := registry.LookupCode("legacy.user_upload_avatar")
	if !ok || missing.Disposition != _operationbinding.DispositionSourceMissing {
		t.Fatalf("missing binding = %#v, %v", missing, ok)
	}
}
