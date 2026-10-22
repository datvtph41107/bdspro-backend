package controller

// func propertyTypeListResult(out *services.PropertyTypeListOutcome) (*httpapi.Result, error) {
// 	if out == nil {
// 		return httpapi.OK(map[string]any{"data": []any{}, "totalElements": int64(0)}), nil
// 	}
// 	var opts []httpapi.OKOption
// 	if out.Timestamp > 0 {
// 		opts = append(opts, httpapi.WithTimestamp(out.Timestamp))
// 	}
// 	// if out.Unchanged {
// 	// 	return httpapi.OK(map[string]any{
// 	// 		"data":          []any{},
// 	// 		"totalElements": int64(0),
// 	// 		"unchanged":     true,
// 	// 	}, opts...), nil
// 	// }
// 	if out.List == nil {
// 		return httpapi.OK(map[string]any{"data": []any{}, "totalElements": int64(0)}, opts...), nil
// 	}
// 	// data, err := protoMessages(out.List.GetData())
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	return httpapi.OK(data, map[string]any{
// 		// "data":          data,
// 		"totalElements": out.List.GetTotalElements(),
// 		"timestamp":     out.Timestamp,
// 	}, opts...), nil
// }
