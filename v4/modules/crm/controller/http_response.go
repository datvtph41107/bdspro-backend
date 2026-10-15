package controller

import (
	"encoding/json"
	"net/http"

	crmpb "pb/types/crm"
	"userv4/internal/httpapi"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/mapper"

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

func protoDataResult(msg proto.Message) (*httpapi.Result, error) {
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
}

// contactDetailResult maps ContactEntity to crmpb.ContactDTO JSON (gateway contract)
// and always includes productIds (protojson omits empty repeated fields).
func contactDetailResult(entity *domain.ContactEntity) (*httpapi.Result, error) {
	if entity == nil {
		return httpapi.OK(nil), nil
	}
	pb := (&mapper.ContactMapper{}).ContactToPb(entity)
	data, err := protoToAny(pb)
	if err != nil {
		return nil, err
	}
	m, ok := data.(map[string]any)
	if !ok {
		return httpapi.OK(data), nil
	}
	ids := entity.ProductIDs
	if ids == nil {
		ids = []uint64{}
	}
	m["productIds"] = ids
	if entity.ProfileID != nil {
		m["profileId"] = *entity.ProfileID
	}
	return httpapi.OK(m), nil
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

func submitResult(id uint64, message string) (*httpapi.Result, error) {
	return httpapi.OK(map[string]any{
		"id":      id,
		"message": message,
	}), nil
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

func contactListResult(resp *crmpb.ContactListDTO) (*httpapi.Result, error) {
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
