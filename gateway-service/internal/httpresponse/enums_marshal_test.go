package httpresponse

import (
	"encoding/json"
	"testing"

	authpb "pb/types/auth"
	bdspropb "pb/types/bdspro"
	tqdpb "pb/types/tqd"
)

func TestMarshalGetEnumsResponse_TimestampAtRoot(t *testing.T) {
	m := &Marshaler{}
	msg := &bdspropb.GetEnumsResponse{
		Timestamp: 1790239512536,
		Lastest:   true,
		Enums: []*bdspropb.GetEnumResponse{
			{Key: "status", Data: []*bdspropb.EnumItem{{Value: 1, Name: "A"}}},
		},
	}

	raw, err := m.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}

	if got["code"].(float64) != 0 {
		t.Fatalf("code=%v", got["code"])
	}
	if got["timestamp"].(float64) != 1790239512536 {
		t.Fatalf("timestamp=%v want root-level 1790239512536; body=%s", got["timestamp"], raw)
	}
	data, ok := got["data"].([]any)
	if !ok || len(data) != 1 {
		t.Fatalf("data=%v want enums slice; body=%s", got["data"], raw)
	}
}

func TestMarshalGetEnumsResponse_UnchangedOnlyTimestamp(t *testing.T) {
	m := &Marshaler{}
	msg := &bdspropb.GetEnumsResponse{
		Timestamp: 1790239512536,
		Lastest:   false,
	}

	raw, err := m.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["timestamp"].(float64) != 1790239512536 {
		t.Fatalf("timestamp=%v; body=%s", got["timestamp"], raw)
	}
	if _, hasData := got["data"]; hasData {
		t.Fatalf("expected no data on unchanged; body=%s", raw)
	}
}

func TestMarshalPermissionKeysResponse_KeysAsDataAndTimestamp(t *testing.T) {
	m := &Marshaler{}
	msg := &authpb.PermissionKeysResponse{
		Keys:      []string{"product.read", "deal.write"},
		Timestamp: 1727160000000,
	}

	raw, err := m.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["timestamp"].(float64) != 1727160000000 {
		t.Fatalf("timestamp=%v; body=%s", got["timestamp"], raw)
	}
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("data=%v want object with keys; body=%s", got["data"], raw)
	}
	keys, ok := data["keys"].([]any)
	if !ok || len(keys) != 2 || keys[0] != "product.read" {
		t.Fatalf("data.keys=%v; body=%s", data["keys"], raw)
	}
}

func TestMarshalPermissionKeysResponse_EmptyKeysStillTimestamp(t *testing.T) {
	m := &Marshaler{}
	msg := &authpb.PermissionKeysResponse{
		Timestamp: 1727160000001,
	}

	raw, err := m.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["timestamp"].(float64) != 1727160000001 {
		t.Fatalf("timestamp=%v; body=%s", got["timestamp"], raw)
	}
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("data=%v want object; body=%s", got["data"], raw)
	}
	keys, ok := data["keys"].([]any)
	if !ok || len(keys) != 0 {
		t.Fatalf("data.keys=%v want []; body=%s", data["keys"], raw)
	}
}

// Unchanged catalog sync returns Data + nil Pagination + Timestamp — must not panic
// (was *reflect.ValueError on Pagination.Elem() for /v2/tqd/client/wards).
func TestMarshalClientListWardsResponse_NilPagination(t *testing.T) {
	m := &Marshaler{}
	msg := &tqdpb.ClientListWardsResponse{
		Timestamp: 1790241753289,
	}

	raw, err := m.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["timestamp"].(float64) != 1790241753289 {
		t.Fatalf("timestamp=%v; body=%s", got["timestamp"], raw)
	}
}

func TestMarshalClientListWardsResponse_WithPagination(t *testing.T) {
	m := &Marshaler{}
	msg := &tqdpb.ClientListWardsResponse{
		Data: []*tqdpb.Ward{{Id: 1, FullName: "A"}},
		Pagination: &tqdpb.LocationPagination{
			Page:  1,
			Size:  100,
			Total: 322,
		},
		Timestamp: 1790241753289,
	}

	raw, err := m.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	pag, ok := got["pagination"].(map[string]any)
	if !ok {
		t.Fatalf("pagination missing; body=%s", raw)
	}
	if pag["page"].(float64) != 1 || pag["size"].(float64) != 100 || pag["total"].(float64) != 322 {
		t.Fatalf("pagination=%v; body=%s", pag, raw)
	}
}
