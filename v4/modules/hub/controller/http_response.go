package controller

import (
	"encoding/json"
	"net/http"

	hubpb "pb/types/hub"
	"userv4/internal/httpapi"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var protoJSON = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: false,
}

func writeOK(w http.ResponseWriter, data any, opts ...httpapi.OKOption) {
	httpapi.WriteOK(w, data, opts...)
}

func writeHTTPError(w http.ResponseWriter, httpStatus int, message string) {
	httpapi.WriteError(w, httpStatus, message)
}

func writeServiceError(w http.ResponseWriter, err error) {
	httpapi.WriteServiceError(w, err)
}

func protoToAny(msg proto.Message) (any, error) {
	if msg == nil {
		return nil, nil
	}
	body, err := protoJSON.Marshal(msg)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func writeProtoData(w http.ResponseWriter, msg proto.Message) {
	data, err := protoToAny(msg)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data)
}

func protoDataResult(msg proto.Message) (*httpapi.Result, error) {
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
}

func emptyProtoListResult() (*httpapi.Result, error) {
	return httpapi.OK([]any{}, httpapi.WithPagination(&httpapi.Pagination{Total: 0})), nil
}

func protoListWithTotalResult[T proto.Message](items []T, total int64) (*httpapi.Result, error) {
	data, err := protoMessages(items)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []any{}
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func wrappedDataResult(msg proto.Message) (*httpapi.Result, error) {
	if msg == nil {
		return httpapi.OK(nil), nil
	}
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	if m, ok := data.(map[string]any); ok {
		if inner, has := m["data"]; has && len(m) == 1 {
			return httpapi.OK(inner), nil
		}
	}
	return httpapi.OK(data), nil
}

func allUserGuidesByKeyResult(resp *hubpb.GetAllUserGuidesByKeyResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}

func writeListWithTotal(w http.ResponseWriter, items []any, total int64) {
	writeOK(w, items, httpapi.WithPagination(&httpapi.Pagination{Total: total}))
}

func writeProtoListWithTotal[T proto.Message](w http.ResponseWriter, items []T, total int64) {
	data, err := protoMessages(items)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeListWithTotal(w, data, total)
}

func writeAllUserGuidesByKey(w http.ResponseWriter, resp *hubpb.GetAllUserGuidesByKeyResponse) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}

func writeWrappedData(w http.ResponseWriter, msg proto.Message) {
	if msg == nil {
		writeOK(w, nil)
		return
	}
	data, err := protoToAny(msg)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Unwrap {data: ...} if the proto response only wraps a single data field.
	if m, ok := data.(map[string]any); ok {
		if inner, has := m["data"]; has && len(m) == 1 {
			writeOK(w, inner)
			return
		}
	}
	writeOK(w, data)
}

func protoMessages[T proto.Message](items []T) ([]any, error) {
	out := make([]any, 0, len(items))
	for _, item := range items {
		v, err := protoToAny(item)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}
