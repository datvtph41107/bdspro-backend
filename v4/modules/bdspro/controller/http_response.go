package controller

import (
	"encoding/json"
	"net/http"

	"common/httpapi"
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	"mdn4/modules/bdspro/services"

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

func writeSearchResponse(w http.ResponseWriter, resp *bdspropb.SearchResponse) {
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
		httpapi.WithPagination(&httpapi.Pagination{
			Page:       resp.GetPage(),
			Size:       resp.GetSize(),
			Total:      resp.GetTotalElements(),
			TotalPages: resp.GetTotalPages(),
		}),
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}

func writeGetDealsResponse(w http.ResponseWriter, resp *bdspropb.GetDealsResponse) {
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

func writeOrganizationDealsResponse(w http.ResponseWriter, resp *bdspropb.GetOrganizationDealsResponse) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	writeProtoList(w, resp.GetData(), int64(resp.GetTotal()))
}

func writeProtoList[T proto.Message](w http.ResponseWriter, items []T, total int64) {
	data, err := protoMessages(items)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data, httpapi.WithPagination(&httpapi.Pagination{Total: total}))
}

func writePropertyListResponse(w http.ResponseWriter, resp *bdspropb.PropertyListResponse) {
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
	if p := resp.GetPagination(); p != nil {
		hasNext := p.GetHasNext()
		hasPrev := p.GetHasPrevious()
		opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{
			Page:        p.GetPage(),
			Size:        p.GetSize(),
			Total:       p.GetTotal(),
			TotalPages:  p.GetPages(),
			HasNext:     &hasNext,
			HasPrevious: &hasPrev,
		}))
	}
	writeOK(w, data, opts...)
}

func writeProtoData(w http.ResponseWriter, msg proto.Message, opts ...httpapi.OKOption) {
	data, err := protoToAny(msg)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data, opts...)
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

func searchResponseResult(resp *bdspropb.SearchResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{
			Page:       resp.GetPage(),
			Size:       resp.GetSize(),
			Total:      resp.GetTotalElements(),
			TotalPages: resp.GetTotalPages(),
		}),
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}

func organizationDealsResult(resp *bdspropb.GetOrganizationDealsResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	return protoListResult(resp.GetData(), int64(resp.GetTotal()))
}

func propertyListResult(resp *bdspropb.PropertyListResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if p := resp.GetPagination(); p != nil {
		hasNext := p.GetHasNext()
		hasPrev := p.GetHasPrevious()
		opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{
			Page:        p.GetPage(),
			Size:        p.GetSize(),
			Total:       p.GetTotal(),
			TotalPages:  p.GetPages(),
			HasNext:     &hasNext,
			HasPrevious: &hasPrev,
		}))
	}
	return httpapi.OK(data, opts...), nil
}

func getEnumsResult(resp *bdspropb.GetEnumsResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetEnums())
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	if !resp.GetLastest() {
		return httpapi.OK([]any{}, opts...), nil
	}
	data = append(data, services.DealSettingEnumListEntry())
	return httpapi.OK(data, opts...), nil
}

func protoListResult[T proto.Message](items []T, total int64) (*httpapi.Result, error) {
	data, err := protoMessages(items)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

// listItemResponseResult maps bdspropb.ListItemResponse to envelope data=[] + pagination.total.
func listItemResponseResult(resp *bdspropb.ListItemResponse, page, size uint32) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []any{}
	}
	p := &httpapi.Pagination{Total: int64(resp.GetTotal())}
	if size > 0 {
		p.Page = int32(page)
		p.Size = int32(size)
	}
	return httpapi.OK(data, httpapi.WithPagination(p)), nil
}

func appointmentListResult(resp *crmpb.ListAppointmentsResponse) (*httpapi.Result, error) {
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
