package planningclient

import (
	"strconv"
	"testing"

	_utils "common/utils"
)

func TestProvinceIDFilterValuesDecodesLocationIDs(t *testing.T) {
	t.Parallel()
	encoded, err := _utils.EncodeLocationID(38, _utils.LocationIDKindProvince)
	if err != nil {
		t.Fatalf("EncodeLocationID() error = %v", err)
	}
	got := provinceIDFilterValues("38," + strconv.FormatUint(encoded, 10))
	if len(got) != 1 || got[0] != "38" {
		t.Fatalf("provinceIDFilterValues() = %v, want [38]", got)
	}
}

func TestProvinceIDFilterValuesKeepsNonNumericCodes(t *testing.T) {
	t.Parallel()
	got := provinceIDFilterValues("HN,38")
	if len(got) != 2 || got[0] != "HN" || got[1] != "38" {
		t.Fatalf("provinceIDFilterValues() = %v, want [HN 38]", got)
	}
}
