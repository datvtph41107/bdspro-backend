package controller

import (
	"encoding/json"
	"net/http"

	crmpb "pb/types/crm"
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

func writeContactListResponse(w http.ResponseWriter, resp *crmpb.ContactListDTO) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(resp.GetTotal())}),
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}
