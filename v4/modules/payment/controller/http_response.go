package controller

import (
	"encoding/json"

	"common/httpapi"
	"mdn4/modules/payment/services"

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

func protoDataResult(msg proto.Message) (*httpapi.Result, error) {
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
}

func bankListResult(out *services.BankListOutcome) (*httpapi.Result, error) {
	if out == nil {
		return httpapi.OK(map[string]any{"data": []any{}, "total": int64(0)}), nil
	}
	if out.Unchanged {
		var opts []httpapi.OKOption
		if out.Timestamp > 0 {
			opts = append(opts, httpapi.WithTimestamp(out.Timestamp))
		}
		return httpapi.OK(map[string]any{
			"data":      []any{},
			"total":     int64(0),
			"unchanged": true,
		}, opts...), nil
	}
	if out.List == nil {
		var opts []httpapi.OKOption
		if out.Timestamp > 0 {
			opts = append(opts, httpapi.WithTimestamp(out.Timestamp))
		}
		return httpapi.OK(map[string]any{"data": []any{}, "total": int64(0)}, opts...), nil
	}
	data, err := protoToAny(out.List)
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if out.Timestamp > 0 {
		opts = append(opts, httpapi.WithTimestamp(out.Timestamp))
	}
	return httpapi.OK(data, opts...), nil
}
