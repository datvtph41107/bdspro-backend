package controller

import (
	"net/http"

	bdspropb "pb/types/bdspro"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// PublicController owns public bdspro HTTP routes under /v2/bdspro/public.
type PublicController struct {
	public services.PublicService
}

func NewPublicController(public services.PublicService) *PublicController {
	return &PublicController{public: public}
}

func (c *PublicController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/bdspro/public/enums", c.httpGetEnums)
	mux.HandleFunc("GET /v2/bdspro/public/enums/{text}", c.httpGetEnumByName)
}

func (c *PublicController) httpGetEnums(w http.ResponseWriter, r *http.Request) {
	var req dto.GetEnumsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.public.GetEnums(r.Context(), req.Timestamp)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeGetEnumsResponse(w, result)
}

func (c *PublicController) httpGetEnumByName(w http.ResponseWriter, r *http.Request) {
	text := r.PathValue("text")
	if text == "" {
		writeHTTPError(w, http.StatusBadRequest, "text is required")
		return
	}
	result, err := c.public.GetEnumByName(r.Context(), text)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func writeGetEnumsResponse(w http.ResponseWriter, resp *bdspropb.GetEnumsResponse) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(resp.GetEnums())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	// Unchanged catalog: empty list + timestamp (lastest=false).
	if !resp.GetLastest() {
		writeOK(w, []any{}, opts...)
		return
	}
	writeOK(w, data, opts...)
}
