package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"userv4/internal/httpapi"
)

type protoEnvelope struct {
	Data       json.RawMessage `json:"data"`
	Pagination *struct {
		Page  int32 `json:"page"`
		Size  int32 `json:"size"`
		Total int64 `json:"total"`
	} `json:"pagination"`
	Message string `json:"message"`
}

func serveProto(t *testing.T, mux *http.ServeMux, r *http.Request) (*httptest.ResponseRecorder, protoEnvelope) {
	t.Helper()
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	var env protoEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode %q: %v", w.Body.String(), err)
	}
	return w, env
}

func TestHandleProto_GetBindsPathAndQuery_ListEnvelope(t *testing.T) {
	var got *crmpb.GetRateListRequest
	mux := http.NewServeMux()
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/rate/{registry}/list/{ownerId}",
		func(_ context.Context, req *crmpb.GetRateListRequest) (*crmpb.GetRateListResponse, error) {
			got = req
			return &crmpb.GetRateListResponse{Data: []*crmpb.RateResponse{{Id: 7}}, Total: 21}, nil
		})

	w, env := serveProto(t, mux, httptest.NewRequest(http.MethodGet, "/v2/feedback/rate/product/list/42?page=2&size=10&parentId=5&registry=ignored", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if got.Registry != "product" || got.OwnerId != 42 || got.Page != 2 || got.Size != 10 || got.GetParentId() != 5 {
		t.Fatalf("bound request = %+v", got)
	}
	if env.Pagination == nil || env.Pagination.Total != 21 {
		t.Fatalf("pagination=%+v body=%s", env.Pagination, w.Body.String())
	}
	if !strings.Contains(string(env.Data), `"id":7`) {
		t.Fatalf("data=%s", env.Data)
	}
}

func TestHandleProto_PostBindsBodyThenPath(t *testing.T) {
	var got *crmpb.CreateRateRequest
	mux := http.NewServeMux()
	httpapi.HandleProto(mux, http.MethodPost, "/v2/feedback/rate/{registry}",
		func(_ context.Context, req *crmpb.CreateRateRequest) (*sharepb.SubmitResponse, error) {
			got = req
			return &sharepb.SubmitResponse{Id: 9, Message: "ok"}, nil
		})

	body := `{"ownerId":"42","score":5,"comment":"tốt","attachs":["a.png"],"registry":"ignored","unknown":1}`
	w, env := serveProto(t, mux, httptest.NewRequest(http.MethodPost, "/v2/feedback/rate/product", strings.NewReader(body)))

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if got.Registry != "product" || got.OwnerId != 42 || got.Score != 5 || got.Comment != "tốt" || len(got.Attachs) != 1 {
		t.Fatalf("bound request = %+v", got)
	}
	if !strings.Contains(string(env.Data), `"message":"ok"`) || env.Pagination != nil {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestHandleProto_InvalidPathParamIs400(t *testing.T) {
	mux := http.NewServeMux()
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/rate/detail/{id}",
		func(_ context.Context, _ *crmpb.RegistryIdRequest) (*crmpb.RateResponse, error) {
			t.Fatal("handler must not run")
			return nil, nil
		})

	w, _ := serveProto(t, mux, httptest.NewRequest(http.MethodGet, "/v2/feedback/rate/detail/abc", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestHandleProto_GRPCStatusMapsToHTTP(t *testing.T) {
	mux := http.NewServeMux()
	httpapi.HandleProto(mux, http.MethodGet, "/v2/enums/crm/{name}",
		func(_ context.Context, _ *crmpb.GetEnumRequest) (*sharepb.ItemV3PageProto, error) {
			return nil, status.Error(codes.InvalidArgument, "bad enum")
		})

	w, env := serveProto(t, mux, httptest.NewRequest(http.MethodGet, "/v2/enums/crm/x", nil))
	if w.Code != http.StatusBadRequest || env.Message != "bad enum" {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
