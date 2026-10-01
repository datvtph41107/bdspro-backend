package httpapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Empty is a no-payload request for handlers with no path/query/body input.
type Empty struct{}

// PathUint64 is a convenience for handlers that only need a path id.
type PathUint64 struct {
	ID uint64 `path:"id"`
}

// DecodeRequest binds path params (`path` tags) then query (`form` tags) into dst.
// Path tag may list alternatives separated by "|" (e.g. path:"id|dealId").
func DecodeRequest(r *http.Request, dst any) error {
	if r == nil {
		return fmt.Errorf("request is required")
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a non-nil pointer")
	}
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("dst must point to a struct")
	}
	if err := bindPath(elem, r); err != nil {
		return err
	}
	return DecodeQuery(r, dst)
}

// DecodePathJSON binds path params then JSON body into dst.
func DecodePathJSON(r *http.Request, dst any) error {
	if r == nil {
		return fmt.Errorf("request is required")
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a non-nil pointer")
	}
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("dst must point to a struct")
	}
	if err := bindPath(elem, r); err != nil {
		return err
	}
	if r.Body == nil || r.Body == http.NoBody {
		return nil
	}
	// Allow empty body for DELETE-like calls that still use JSON binder.
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(buf))) == 0 {
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(buf)))
	return DecodeJSON(r, dst)
}

func bindPath(v reflect.Value, r *http.Request) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		valueField := v.Field(i)
		if !valueField.CanSet() {
			continue
		}
		tag := field.Tag.Get("path")
		if tag == "" || tag == "-" {
			continue
		}
		raw := ""
		for _, key := range strings.Split(tag, "|") {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}
			if v := strings.TrimSpace(r.PathValue(key)); v != "" {
				raw = v
				break
			}
		}
		if raw == "" {
			return fmt.Errorf("%s is required", primaryPathKey(tag))
		}
		if err := setQueryField(valueField, raw, []string{raw}); err != nil {
			return fmt.Errorf("invalid %s: %w", primaryPathKey(tag), err)
		}
		switch valueField.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if valueField.Uint() == 0 {
				return fmt.Errorf("%s is required", primaryPathKey(tag))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if valueField.Int() == 0 {
				return fmt.Errorf("%s is required", primaryPathKey(tag))
			}
		}
	}
	return nil
}

func primaryPathKey(tag string) string {
	parts := strings.Split(tag, "|")
	if len(parts) == 0 {
		return "id"
	}
	return strings.TrimSpace(parts[0])
}

// Result is the success payload written by Bind handlers.
type Result struct {
	Data any
	Opts []OKOption
}

// OK builds a success Result for Bind handlers.
func OK(data any, opts ...OKOption) *Result {
	return &Result{Data: data, Opts: opts}
}

// Handler is a typed request handler: bind Req, then run business logic.
type Handler[Req any] func(ctx context.Context, req Req) (*Result, error)

func wrap[Req any](decode func(*http.Request, any) error, h Handler[Req]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := decode(r, &req); err != nil {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := h(r.Context(), req)
		if err != nil {
			WriteServiceError(w, err)
			return
		}
		if result == nil {
			WriteOK(w, nil)
			return
		}
		WriteOK(w, result.Data, result.Opts...)
	}
}

// Bind wraps a typed Handler with automatic path+query decode.
func Bind[Req any](h Handler[Req]) http.HandlerFunc {
	return wrap(DecodeRequest, h)
}

// BindJSON wraps a typed Handler with automatic path+JSON body decode.
func BindJSON[Req any](h Handler[Req]) http.HandlerFunc {
	return wrap(DecodePathJSON, h)
}

func normalize(method, pattern string) string {
	method = strings.ToUpper(strings.TrimSpace(method))
	pattern = strings.TrimSpace(pattern)
	if strings.HasPrefix(pattern, method+" ") {
		return pattern
	}
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	return method + " " + pattern
}

// GET registers a path+query bound handler.
func GET[Req any](mux *http.ServeMux, pattern string, h Handler[Req]) {
	mux.HandleFunc(normalize("GET", pattern), Bind(h))
}

// POST registers a path+JSON bound handler.
func POST[Req any](mux *http.ServeMux, pattern string, h Handler[Req]) {
	mux.HandleFunc(normalize("POST", pattern), BindJSON(h))
}

// PUT registers a path+JSON bound handler.
func PUT[Req any](mux *http.ServeMux, pattern string, h Handler[Req]) {
	mux.HandleFunc(normalize("PUT", pattern), BindJSON(h))
}

// PATCH registers a path+JSON bound handler.
func PATCH[Req any](mux *http.ServeMux, pattern string, h Handler[Req]) {
	mux.HandleFunc(normalize("PATCH", pattern), BindJSON(h))
}

// DELETE registers a path(+query) bound handler.
func DELETE[Req any](mux *http.ServeMux, pattern string, h Handler[Req]) {
	mux.HandleFunc(normalize("DELETE", pattern), Bind(h))
}

// ParsePathUint64 reads a required positive uint64 path value.
func ParsePathUint64(r *http.Request, key string) (uint64, error) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("%s is required", key)
	}
	return id, nil
}
