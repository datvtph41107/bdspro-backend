package controller

import (
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/services"
)

func costTypeListResult(out *services.CostTypeListOutcome) (*httpapi.Result, error) {
	if out == nil {
		return httpapi.OK(map[string]any{"data": []any{}, "total": uint32(0)}), nil
	}
	var opts []httpapi.OKOption
	if out.Timestamp > 0 {
		opts = append(opts, httpapi.WithTimestamp(out.Timestamp))
	}
	if out.Unchanged {
		return httpapi.OK(map[string]any{
			"data":      []any{},
			"total":     uint32(0),
			"unchanged": true,
		}, opts...), nil
	}
	if out.List == nil {
		return httpapi.OK(map[string]any{"data": []any{}, "total": uint32(0)}, opts...), nil
	}
	data, err := protoMessages(out.List.GetData())
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]any{
		"data":  data,
		"total": out.List.GetTotal(),
	}, opts...), nil
}
