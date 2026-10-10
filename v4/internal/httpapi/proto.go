package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtoMethod is a unary gRPC-style method: proto request in, proto response out.
type ProtoMethod[Req, Resp proto.Message] func(ctx context.Context, req Req) (Resp, error)

var (
	patternParamRe   = regexp.MustCompile(`\{([^}/]+)\}`)
	protoBodyDecoder = protojson.UnmarshalOptions{DiscardUnknown: true}
	protoDataEncoder = protojson.MarshalOptions{UseProtoNames: false, EmitUnpopulated: false}
)

// HandleProto mounts a gRPC-style method with the v2 gateway (grpc-gateway, body "*")
// binding: path params fill fields by name; GET/DELETE read the query string,
// other methods decode the JSON body. The response is written like the gateway
// marshaler did (Data + pagination, otherwise the whole message as data).
func HandleProto[Req, Resp proto.Message](mux *http.ServeMux, method, pattern string, h ProtoMethod[Req, Resp]) {
	params := patternParams(pattern)
	mux.HandleFunc(normalize(method, pattern), func(w http.ResponseWriter, r *http.Request) {
		var zero Req
		req := zero.ProtoReflect().New().Interface().(Req)
		if err := bindProto(r, req, params); err != nil {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		resp, err := h(r.Context(), req)
		if err != nil {
			writeServiceError(r.Context(), w, err, r)
			return
		}
		WriteProto(w, resp)
	})
}

func patternParams(pattern string) []string {
	matches := patternParamRe.FindAllStringSubmatch(pattern, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, strings.TrimSuffix(m[1], "..."))
	}
	return out
}

func bindProto(r *http.Request, msg proto.Message, params []string) error {
	switch r.Method {
	case http.MethodGet, http.MethodDelete:
		filter := make([][]string, 0, len(params))
		for _, p := range params {
			filter = append(filter, []string{p})
		}
		if err := runtime.PopulateQueryParameters(msg, r.URL.Query(), utilities.NewDoubleArray(filter)); err != nil {
			return fmt.Errorf("invalid query: %w", err)
		}
	default:
		if r.Body != nil && r.Body != http.NoBody {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			if len(strings.TrimSpace(string(body))) > 0 {
				if err := protoBodyDecoder.Unmarshal(body, msg); err != nil {
					return fmt.Errorf("invalid body: %w", err)
				}
			}
		}
	}
	for _, p := range params {
		if err := runtime.PopulateFieldFromPath(msg, p, r.PathValue(p)); err != nil {
			return fmt.Errorf("invalid %s: %w", p, err)
		}
	}
	return nil
}

// WriteProto writes a proto response in the gateway envelope shape.
func WriteProto(w http.ResponseWriter, msg proto.Message) {
	v := reflect.ValueOf(msg)
	if msg == nil || (v.Kind() == reflect.Ptr && v.IsNil()) {
		WriteOK(w, nil)
		return
	}
	elem := v.Elem()
	if data := elem.FieldByName("Data"); data.IsValid() && data.CanInterface() {
		WriteOK(w, data.Interface(),
			WithPagination(protoPagination(elem)),
			WithTimestamp(intField(elem, "Timestamp")))
		return
	}
	raw, err := protoDataEncoder.Marshal(msg)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(data) == 0 {
		WriteOK(w, nil)
		return
	}
	WriteOK(w, data)
}

func protoPagination(elem reflect.Value) *Pagination {
	total := intField(elem, "Total")
	if total == 0 {
		total = intField(elem, "TotalElements")
	}
	size := intField(elem, "Size")
	if size == 0 {
		size = intField(elem, "PageSize")
	}
	p := &Pagination{Page: int32(intField(elem, "Page")), Size: int32(size), Total: total}
	if p.Page == 0 && p.Size == 0 && p.Total == 0 {
		return nil
	}
	if p.Size > 0 && p.Total > 0 {
		p.TotalPages = int32((p.Total + int64(p.Size) - 1) / int64(p.Size))
	}
	return p
}

func intField(elem reflect.Value, name string) int64 {
	f := elem.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(f.Uint())
	}
	return 0
}
