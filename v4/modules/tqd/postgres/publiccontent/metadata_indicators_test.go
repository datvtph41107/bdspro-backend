package publiccontent

import "testing"

func TestMetadataIndicatorsAcceptsDirectArray(t *testing.T) {
	got := metadataIndicators(`[
		{
			"code":"density",
			"name":"Mật độ xây dựng",
			"value":"40",
			"unit":"%"
		}
	]`)

	if len(got) != 1 {
		t.Fatalf("expected one indicator, got %d", len(got))
	}

	if got[0].Name != "Mật độ xây dựng" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].Value != "40" {
		t.Fatalf("unexpected value: %q", got[0].Value)
	}
	if got[0].Unit != "%" {
		t.Fatalf("unexpected unit: %q", got[0].Unit)
	}
}

func TestMetadataIndicatorsAcceptsMetadataObject(t *testing.T) {
	got := metadataIndicators(`{
		"planning_indicators":[
			{
				"label":"Tầng cao",
				"target":"12",
				"uom":"tầng"
			}
		]
	}`)

	if len(got) != 1 {
		t.Fatalf("expected one indicator, got %d", len(got))
	}

	if got[0].Name != "Tầng cao" {
		t.Fatalf("unexpected name: %q", got[0].Name)
	}
	if got[0].Value != "12" {
		t.Fatalf("unexpected value: %q", got[0].Value)
	}
}
