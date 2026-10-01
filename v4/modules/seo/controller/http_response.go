package controller

import (
	"net/http"

	"userv4/internal/httpapi"
)

func dataResult(v any) (*httpapi.Result, error) {
	return httpapi.OK(v), nil
}

func listResult(items any, total int64) (*httpapi.Result, error) {
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func writeXML(w http.ResponseWriter, xml string, status int) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(xml))
}
