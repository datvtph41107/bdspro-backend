package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
	"mdn4/modules/crm/services"
)

type SharingController struct {
	uc *services.SharingAccessUsecase
}

func NewSharingController(uc *services.SharingAccessUsecase) *SharingController {
	return &SharingController{uc: uc}
}

func (c *SharingController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/crm/sharing/share", c.share)
	httpapi.GET(mux, "/v2/crm/sharing/list/{contactId}", c.list)
}

func (c *SharingController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "sharing service not configured")
	}
	return nil
}

type sharingListReq struct {
	ContactID uint64 `path:"contactId"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
}
type sharingShareReq struct {
	ContactID   uint64                 `json:"contactId"`
	Permissions []int32                `json:"permissions"`
	Note        string                 `json:"note"`
	OwnerType   int32                  `json:"ownerType"`
	DeletedIDs  []uint64               `json:"deletedIds"`
	Receivers   []domain.SharingEntity `json:"receivers"`
}

func (c *SharingController) list(ctx context.Context, req sharingListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.List(ctx, req.ContactID, &dto.SharingAccessSearchDTO{})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
func (c *SharingController) share(ctx context.Context, req sharingShareReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.BulkShare(ctx, req.ContactID, req.Permissions, req.Note, enums.EOwnerOf(req.OwnerType), req.DeletedIDs, req.Receivers); err != nil {
		return nil, err
	}
	return submitResult(req.ContactID, "shared")
}
