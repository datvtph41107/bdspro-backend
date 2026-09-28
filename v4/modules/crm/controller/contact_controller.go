package controller

import (
	"context"
	"net/http"

	crmpb "pb/types/crm"
	"userv4/internal/httpapi"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/services"
)

// ContactController owns GET /v2/crm/contact/list/user/{ownerId}.
type ContactController struct {
	contact services.ContactService
}

func NewContactController(contact services.ContactService) *ContactController {
	return &ContactController{contact: contact}
}

func (c *ContactController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/contact/list/user/{ownerId}", c.listUser)
}

func (c *ContactController) listUser(ctx context.Context, req dto.ContactListUserRequest) (*httpapi.Result, error) {
	search := dto.ContactSearchRequest{
		Page:      req.Page,
		Size:      req.Size,
		Text:      req.Text,
		Timestamp: req.Timestamp,
	}
	result, err := c.contact.ListUser(ctx, req.OwnerID, search)
	if err != nil {
		return nil, err
	}
	return contactListResult(result)
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
