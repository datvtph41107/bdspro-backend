package controller

import (
	"encoding/json"
	"net/http"

	notificationpb "pb/types/notification"
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

func writeCountResponse(w http.ResponseWriter, resp *notificationpb.CountResponse) {
	if resp == nil {
		writeOK(w, map[string]any{"count": 0, "unread": 0})
		return
	}
	data := map[string]any{
		"count":     resp.GetCount(),
		"unread":    resp.GetUnread(),
		"unchanged": resp.GetUnchanged(),
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}

func writeListResponse(w http.ResponseWriter, resp *notificationpb.ListResponse) {
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

func writeProtoData(w http.ResponseWriter, msg proto.Message) {
	data, err := protoToAny(msg)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data)
}

func countResult(resp *notificationpb.CountResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK(map[string]any{"count": 0, "unread": 0}), nil
	}
	data := map[string]any{
		"count":     resp.GetCount(),
		"unread":    resp.GetUnread(),
		"unchanged": resp.GetUnchanged(),
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}

func listResult(resp *notificationpb.ListResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(resp.GetTotal())}),
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}

func protoDataResult(msg proto.Message) (*httpapi.Result, error) {
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
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
