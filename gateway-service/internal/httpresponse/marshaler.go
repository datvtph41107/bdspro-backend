package httpresponse

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	sharepb "pb/types/shared"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Marshaler struct{}

type customJSONDecoder struct {
	decoder *json.Decoder
}

func (d *customJSONDecoder) Decode(v interface{}) error {
	var raw json.RawMessage
	if err := d.decoder.Decode(&raw); err != nil {
		return err
	}

	if msg, ok := v.(proto.Message); ok {
		return (protojson.UnmarshalOptions{
			DiscardUnknown: false,
		}).Unmarshal(raw, msg)
	}

	return json.Unmarshal(raw, v)
}

// PaginationInfo chứa thông tin phân trang
type PaginationInfo struct {
	Page    int32 `json:"page"`    // Trang hiện tại
	Size    int32 `json:"size"`    // Kích thước trang
	Total   int64 `json:"total"`   // Tổng số bản ghi (totalElements)
	Pages   int32 `json:"pages"`   // Tổng số trang
	HasNext bool  `json:"hasNext"` // Có trang tiếp không
}

// StandardResponse với hỗ trợ cả format cũ và mới
type StandardResponse struct {
	Code    int         `json:"code"`           // 0 = success, >0 = error code
	Message string      `json:"message"`        // success message or error message
	Data    interface{} `json:"data,omitempty"` // Giữ nguyên field data (backward compatible)
	// Các field cũ -
	TotalElements int64 `json:"totalElements,omitempty"`
	Total         int64 `json:"total,omitempty"`
	// pagination object
	Pagination *PaginationInfo `json:"pagination,omitempty"` // Thông tin phân trang đầy đủ
}

func (c *Marshaler) Marshal(v interface{}) ([]byte, error) {
	// Default response
	resp := StandardResponse{
		Code:    0,
		Message: "success",
	}

	// Try to convert to proto.Message
	msg, ok := v.(proto.Message)
	if !ok {
		return json.Marshal(v)
	}

	// HttpBody: trả raw JSON/binary (public TQD projections, …)
	if hb, ok := msg.(*httpbody.HttpBody); ok {
		if hb == nil {
			return []byte("null"), nil
		}
		return hb.GetData(), nil
	}

	// BytesValue được các gRPC handler dùng để vận chuyển JSON đã marshal.
	// Phải trả raw bytes; nếu để encoding/json xử lý, []byte sẽ bị base64.
	if bv, ok := msg.(*wrapperspb.BytesValue); ok {
		if bv == nil {
			return []byte("null"), nil
		}
		return bv.GetValue(), nil
	}

	val := reflect.ValueOf(msg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Kiểm tra field Code
	codeField := val.FieldByName("Code")
	if codeField.IsValid() && codeField.CanInterface() {
		var code int64
		switch codeField.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			code = codeField.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			code = int64(codeField.Uint())
		}

		if code == 8500 {
			dataField := val.FieldByName("Data")
			if dataField.IsValid() && dataField.CanInterface() {
				return json.Marshal(dataField.Interface())
			}
			return json.Marshal(v)
		}

		if code == 8501 {
			dataField := val.FieldByName("Data")
			if dataField.IsValid() && dataField.CanInterface() {
				if b, ok := dataField.Interface().([]byte); ok {
					return b, nil
				}
			}
			return json.Marshal(v)
		}
	}

	// Kiểm tra field Data
	dataField := val.FieldByName("Data")
	if dataField.IsValid() && dataField.CanInterface() {
		// Một số response (AdminOpportunity, list+funnel, …) có sibling fields
		// (events, commercial, funnel, …) ngoài Data/Total. Custom envelope cũ
		// chỉ lấy Data → FE mất nhật ký / commercial / phễu. Khi có sibling,
		// nhét cả message vào data để giữ nguyên các field đó.
		if hasProtoSiblingPayload(val) {
			raw, err := protojson.MarshalOptions{
				UseProtoNames:   false,
				EmitUnpopulated: false,
			}.Marshal(msg)
			if err != nil {
				return nil, err
			}
			var envelope any
			if err := json.Unmarshal(raw, &envelope); err != nil {
				return nil, err
			}
			resp.Data = envelope
		} else {
			resp.Data = dataField.Interface()
		}

		// Khởi tạo pagination info
		pagination := &PaginationInfo{}

		// Tìm trực tiếp Total / TotalElements, Page, Size fields
		// Nhiều service dùng TotalElements (RoleListResponse, PermissionListResponse…)
		totalField := val.FieldByName("Total")
		if !totalField.IsValid() {
			totalField = val.FieldByName("TotalElements")
		}
		if totalField.IsValid() && totalField.CanInterface() {
			switch totalField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				total := totalField.Int() // luôn trả về int64
				resp.Total = total
				resp.TotalElements = total
				pagination.Total = total

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				total := int64(totalField.Uint())
				resp.Total = total
				resp.TotalElements = total
				pagination.Total = total
			}
		}

		pageField := val.FieldByName("Page")
		if pageField.IsValid() && pageField.CanInterface() {
			switch p := pageField.Interface().(type) {
			case int32:
				pagination.Page = p
			case uint32:
				pagination.Page = int32(p)
			}
		}

		sizeField := val.FieldByName("Size")
		if sizeField.IsValid() && sizeField.CanInterface() {
			switch s := sizeField.Interface().(type) {
			case int32:
				pagination.Size = s
			case uint32:
				pagination.Size = int32(s)
			}
		}

		totalPagesField := val.FieldByName("TotalPages")
		if totalPagesField.IsValid() && totalPagesField.CanInterface() {
			switch tp := totalPagesField.Interface().(type) {
			case int32:
				pagination.Pages = tp
			case uint32:
				pagination.Pages = int32(tp)
			case int64:
				pagination.Pages = int32(tp)
			}
		}

		// Pagination key
		paginationField := val.FieldByName("Pagination")
		if paginationField.IsValid() && paginationField.CanInterface() {
			if paginationField.Kind() == reflect.Ptr {
				paginationField = paginationField.Elem()
			}

			if pageField := paginationField.FieldByName("Page"); pageField.IsValid() {
				if page, ok := pageField.Interface().(int32); ok {
					pagination.Page = page
				}
			}
			if sizeField := paginationField.FieldByName("Size"); sizeField.IsValid() {
				if size, ok := sizeField.Interface().(int32); ok {
					pagination.Size = size
				}
			}
			if totalField := paginationField.FieldByName("Total"); totalField.IsValid() {
				if total, ok := totalField.Interface().(int64); ok {
					pagination.Total = total
					resp.Total = total
					resp.TotalElements = total
				}
			}
			if pagesField := paginationField.FieldByName("Pages"); pagesField.IsValid() {
				if pages, ok := pagesField.Interface().(int32); ok {
					pagination.Pages = pages
				}
			}
			if hasNextField := paginationField.FieldByName("HasNext"); hasNextField.IsValid() {
				if hasNext, ok := hasNextField.Interface().(bool); ok {
					pagination.HasNext = hasNext
				}
			}
		}

		if pagination.Size > 0 && pagination.Total > 0 && pagination.Pages == 0 {
			pagination.Pages = int32((pagination.Total + int64(pagination.Size) - 1) / int64(pagination.Size))
			pagination.HasNext = pagination.Page < pagination.Pages
		}

		if pagination.Page > 0 || pagination.Size > 0 || pagination.Total > 0 {
			resp.Pagination = pagination
		}

		return json.Marshal(resp)
	}

	resp.Data = v
	return json.Marshal(resp)
}

func (c *Marshaler) Unmarshal(data []byte, v interface{}) error {
	if msg, ok := v.(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, msg)
	}
	return json.Unmarshal(data, v)
}

// protoJSONDecoder — grpc-gateway Decode body vào proto phải dùng protojson
// (encoding/json không map đúng mọi field proto → assigneeId/qaStatus bị = 0).
type protoJSONDecoder struct {
	r io.Reader
}

func (d *protoJSONDecoder) Decode(v interface{}) error {
	data, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return io.EOF
	}
	if msg, ok := v.(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, msg)
	}
	return json.Unmarshal(data, v)
}

func (c *Marshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &protoJSONDecoder{r: r}
}

func (c *Marshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}

// hasProtoSiblingPayload reports whether a response message carries payload
// fields other than Data + pagination (e.g. events, commercial, funnel).
func hasProtoSiblingPayload(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return false
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}
	skip := map[string]struct{}{
		"Data": {}, "Total": {}, "TotalElements": {}, "Page": {}, "Size": {}, "TotalPages": {}, "Pagination": {},
		"Code": {}, "Message": {}, "state": {}, "sizeCache": {}, "unknownFields": {},
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		sf := typ.Field(i)
		if !sf.IsExported() {
			continue
		}
		if strings.HasPrefix(sf.Name, "XXX_") {
			continue
		}
		if _, ok := skip[sf.Name]; ok {
			continue
		}
		f := val.Field(i)
		if !f.IsValid() || !f.CanInterface() {
			continue
		}
		switch f.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
			if f.IsNil() || (f.Kind() == reflect.Slice && f.Len() == 0) {
				continue
			}
		case reflect.String:
			if f.String() == "" {
				continue
			}
		case reflect.Bool:
			if !f.Bool() {
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if f.Int() == 0 {
				continue
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if f.Uint() == 0 {
				continue
			}
		case reflect.Float32, reflect.Float64:
			if f.Float() == 0 {
				continue
			}
		}
		return true
	}
	return false
}

func (c *Marshaler) ContentType(v interface{}) string {
	if msg, ok := v.(proto.Message); ok {
		if hb, ok := msg.(*httpbody.HttpBody); ok && hb != nil {
			if ct := strings.TrimSpace(hb.GetContentType()); ct != "" {
				return ct
			}
			return "application/json"
		}
		if cr, ok := msg.(*sharepb.CommonResponse); ok && cr.Code == 8501 {
			if contentType := rawResponseContentType(cr.Message); contentType != "" {
				return contentType
			}
			return "application/x-protobuf"
		}

		val := reflect.ValueOf(msg)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		codeField := val.FieldByName("Code")
		if codeField.IsValid() && codeField.CanInterface() {
			var code int64
			switch codeField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				code = codeField.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				code = int64(codeField.Uint())
			}
			if code == 8501 {
				return "application/x-protobuf"
			}
		}
	}
	return "application/json"
}

func rawResponseContentType(message string) string {
	contentType := strings.TrimSpace(message)
	switch strings.ToLower(contentType) {
	case "application/xml", "application/xml; charset=utf-8", "text/xml", "text/xml; charset=utf-8":
		return contentType
	default:
		return ""
	}
}

// ForwardMetadata converts transport-neutral gRPC response
// metadata into the real HTTP response. This keeps status/header ownership in
// gateway-service while CRM/TQD/Hub remain gRPC-only.
func ForwardMetadata(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	serverMetadata, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	for _, key := range []string{
		"cache-control",
		"etag",
		"x-robots-tag",
		"x-content-type-options",
		"retry-after",
		"content-disposition",
	} {
		values := serverMetadata.HeaderMD.Get(key)
		if len(values) == 0 || strings.TrimSpace(values[0]) == "" {
			continue
		}
		w.Header().Set(http.CanonicalHeaderKey(key), strings.TrimSpace(values[0]))
	}

	values := serverMetadata.HeaderMD.Get("x-http-status")
	if len(values) == 0 {
		return nil
	}
	statusCode, err := strconv.Atoi(strings.TrimSpace(values[0]))
	if err != nil || statusCode < 100 || statusCode > 599 {
		return nil
	}
	w.WriteHeader(statusCode)
	return nil
}

// ForwardCommon applies CommonResponse HTTP metadata without requiring the original request.
func ForwardCommon(_ context.Context, w http.ResponseWriter, resp proto.Message) error {
	cr, ok := resp.(*sharepb.CommonResponse)
	if !ok || cr == nil {
		return nil
	}
	for _, ck := range cr.Cookies {
		http.SetCookie(w, &http.Cookie{
			Name: ck.Name, Value: ck.Value, Path: ck.Path, Domain: ck.Domain,
			HttpOnly: ck.HttpOnly, Secure: ck.Secure, MaxAge: int(ck.MaxAge),
		})
	}
	if location := strings.TrimSpace(cr.RedirectUrl); location != "" {
		// grpc-gateway ForwardResponseOption has no *http.Request; writing Location
		// directly avoids the nil-request panic path in http.Redirect.
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusFound)
	}
	return nil
}
