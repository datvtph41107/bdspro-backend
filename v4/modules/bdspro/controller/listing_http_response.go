package controller

import (
	"common/httpapi"
	bdspropb "pb/types/bdspro"
	"mdn4/modules/bdspro/services"
)

func listingListResult(out *services.ListingListOutcome, page, size uint32) (*httpapi.Result, error) {
	if out == nil {
		return httpapi.OK([]any{}), nil
	}
	var opts []httpapi.OKOption
	if out.Timestamp > 0 {
		opts = append(opts, httpapi.WithTimestamp(out.Timestamp))
	}
	pag := &httpapi.Pagination{
		Page: int32(page),
		Size: int32(size),
	}
	if out.Unchanged {
		pag.Total = 0
		opts = append(opts, httpapi.WithPagination(pag))
		return httpapi.OK([]any{}, opts...), nil
	}
	return postListArrayResult(out.List, page, size, opts...)
}

func postListArrayResult(resp *bdspropb.PostListResponse, page, size uint32, extraOpts ...httpapi.OKOption) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}, extraOpts...), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []any{}
	}
	opts := append(extraOpts, httpapi.WithPagination(&httpapi.Pagination{
		Page:  int32(page),
		Size:  int32(size),
		Total: int64(resp.GetTotal()),
	}))
	return httpapi.OK(data, opts...), nil
}
