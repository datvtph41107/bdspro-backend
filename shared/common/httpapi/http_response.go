package httpapi

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var protoJSON = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: false,
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

type ResultResponse struct {
	Data      any
	Total     int64
	Page      int32
	Size      int32
	Timestamp int64
}

func OKResponse(resp *ResultResponse) (*Result, error) {
	if resp == nil {
		return OK([]any{}), nil
	}
	opts := []OKOption{
		WithPagination(&Pagination{
			Page:  resp.Page,
			Size:  resp.Size,
			Total: resp.Total,
			// TotalPages: resp.TotalPages,
		}),
	}
	if resp.Timestamp > 0 {
		opts = append(opts, WithTimestamp(resp.Timestamp))
	}
	return OK(resp.Data, opts...), nil
}
