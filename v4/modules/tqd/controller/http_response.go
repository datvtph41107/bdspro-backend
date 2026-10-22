package controller

import (
	"encoding/json"
	"net/http"

	"common/httpapi"

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

func protoResult(msg proto.Message, opts ...httpapi.OKOption) (*httpapi.Result, error) {
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, opts...), nil
}

func protoListResult[T proto.Message](items []T, total int64, page, size uint32) (*httpapi.Result, error) {
	data, err := protoMessages(items)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{
		Page:  int32(page),
		Size:  int32(size),
		Total: total,
	})), nil
}

func jsonListResult(items any, total int64, page, size uint32) (*httpapi.Result, error) {
	if items == nil {
		items = []any{}
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{
		Page:  int32(page),
		Size:  int32(size),
		Total: total,
	})), nil
}

func jsonOK(data any, opts ...httpapi.OKOption) (*httpapi.Result, error) {
	return httpapi.OK(data, opts...), nil
}
