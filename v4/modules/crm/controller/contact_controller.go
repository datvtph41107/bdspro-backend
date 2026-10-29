package controller

import (
	"context"
	"net/http"

	base_enum "base/enum"
	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
	"mdn4/modules/crm/provider"
	"mdn4/modules/crm/services"
)

// ContactController owns /v2/crm/contact/* routes.
type ContactController struct {
	local     services.ContactService
	uc        *services.ContactUsecase
	histories provider.ContactHistoryReader
}

func NewContactController(local services.ContactService, uc *services.ContactUsecase, histories provider.ContactHistoryReader) *ContactController {
	return &ContactController{local: local, uc: uc, histories: histories}
}

// contactHistoryReq is GET /v2/crm/contact/{id}/history (was /v2/notification/history/contact/{id}).
type contactHistoryReq struct {
	ID   uint64 `path:"id"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

func (c *ContactController) history(ctx context.Context, req contactHistoryReq) (*httpapi.Result, error) {
	if c.histories == nil {
		return nil, _errors.ReturnError(501, "contact history not configured")
	}
	items, total, err := c.histories.ContactHistory(ctx, req.ID, _dto.Pagable{Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	return protoListWithTotalResult(items, total)
}

func (c *ContactController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/contact/list/me", c.listMe)
	httpapi.GET(mux, "/v2/crm/contact/list/user/{ownerId}", c.listUser)
	httpapi.GET(mux, "/v2/crm/contact/list/{ownerOf}/{ownerId}", c.search)
	httpapi.GET(mux, "/v2/crm/contact/phone/{ownerOf}/{ownerId}/{phone}", c.checkPhone)
	httpapi.GET(mux, "/v2/crm/contact/with-member", c.withMember)
	httpapi.GET(mux, "/v2/crm/contact/friends", c.friends)
	httpapi.GET(mux, "/v2/crm/contact/interested/product/{productId}", c.interested)
	httpapi.POST(mux, "/v2/crm/contact", c.create)
	httpapi.POST(mux, "/v2/crm/contact/sync", c.sync)
	httpapi.POST(mux, "/v2/crm/contact/pin", c.pin)

	// Two-segment routes must all be {id}/<action>: ServeMux panics on
	// overlapping patterns such as <action>/{id} vs {id}/<action>.
	httpapi.GET(mux, "/v2/crm/contact/{id}/detail", c.detail)
	httpapi.GET(mux, "/v2/crm/contact/{id}/products", c.products)
	httpapi.GET(mux, "/v2/crm/contact/{id}/history", c.history)
	httpapi.GET(mux, "/v2/crm/contact/{id}/relation-ship", c.relationShip)
	httpapi.GET(mux, "/v2/crm/contact/{id}/by-product", c.byProduct)
	httpapi.PUT(mux, "/v2/crm/contact/{id}/edit", c.update)
	httpapi.PUT(mux, "/v2/crm/contact/{id}/note", c.note)
	httpapi.PUT(mux, "/v2/crm/contact/{id}/products", c.updateProducts)
	httpapi.DELETE(mux, "/v2/crm/contact/{id}", c.delete)
}

func (c *ContactController) requireLocal() error {
	if c == nil || c.local == nil {
		return _errors.ReturnError(501, "contact service not configured")
	}
	return nil
}

func (c *ContactController) requireUC() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "contact usecase not configured")
	}
	return nil
}

func (c *ContactController) listUser(ctx context.Context, req dto.ContactListUserRequest) (*httpapi.Result, error) {
	if err := c.requireLocal(); err != nil {
		return nil, err
	}
	search := dto.ContactSearchRequest{Page: req.Page, Size: req.Size, Text: req.Text, Timestamp: req.Timestamp}
	result, err := c.local.ListUser(ctx, req.OwnerID, search)
	if err != nil {
		return nil, err
	}
	return httpapi.OKResponse(&httpapi.ResultResponse{
		Data:  result.Data,
		Total: int64(result.Total),
	})
}

func (c *ContactController) listMe(ctx context.Context, req dto.ContactListMeRequest) (*httpapi.Result, error) {
	if err := c.requireLocal(); err != nil {
		return nil, err
	}
	result, err := c.local.ListUser(ctx, 0, dto.ContactSearchRequest{
		Page: req.Page, Size: req.Size, Text: req.Text, Timestamp: req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	return contactListResult(result)
}

func (c *ContactController) search(ctx context.Context, req dto.ContactSearchHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	ownerOf, err := dto.ParseOwnerOf(req.OwnerOf)
	if err != nil {
		return nil, err
	}
	items, total, err := c.uc.Search(ctx, ownerOf, dto.ContactSearchDTO{
		Pagable: _dto.Pagable{Page: req.Page, Size: req.Size},
		Text:    req.Text,
	})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *ContactController) detail(ctx context.Context, req dto.ContactDetailHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	out, err := c.uc.Detail(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	// return contactDetailResult(out)
	return httpapi.OKResponse(&httpapi.ResultResponse{
		Data: out,
	})
}

func (c *ContactController) create(ctx context.Context, req dto.ContactCreateHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	out, err := c.uc.Create(ctx, &domain.Contact{
		Phone: req.Phone, FullName: req.FullName, Email: req.Email, Zalo: req.Zalo,
		Company: req.Company, Note: req.Note, Avatar: req.Avatar,
		OwnerOf: base_enum.EOwnerOf(req.OwnerOf), OwnerID: req.OwnerID,
		ProductIDs: req.ProductIds,
	})
	if err != nil {
		return nil, err
	}
	detail, err := c.uc.Detail(ctx, out.ID)
	if err != nil {
		return contactDetailResult(out)
	}
	return contactDetailResult(detail)
}

func (c *ContactController) update(ctx context.Context, req dto.ContactSaveHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	out, err := c.uc.Update(ctx, req.ID, &domain.Contact{
		Phone: req.Phone, FullName: req.FullName, Email: req.Email, Zalo: req.Zalo,
		Company: req.Company, Note: req.Note, Avatar: req.Avatar,
		ProductIDs: req.ProductIds,
	})
	if err != nil {
		return nil, err
	}
	detail, err := c.uc.Detail(ctx, out.ID)
	if err != nil {
		return contactDetailResult(out)
	}
	return contactDetailResult(detail)
}

func (c *ContactController) delete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "deleted")
}

func (c *ContactController) sync(ctx context.Context, req dto.ContactSyncHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	entities := make([]domain.Contact, 0, len(req.Contacts))
	for _, item := range req.Contacts {
		entities = append(entities, domain.Contact{
			Phone: item.Phone, FullName: item.FullName, Email: item.Email, Avatar: item.Avatar,
		})
	}
	out, err := c.uc.Sync(ctx, base_enum.EOwnerOf(req.OwnerType), entities)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}

func (c *ContactController) relationShip(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	out := c.uc.RelationShip(ctx, req.ID)
	return httpapi.OK(out), nil
}

func (c *ContactController) note(ctx context.Context, req dto.ContactNoteHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	if err := c.uc.Note(ctx, req.ID, req.Note); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "updated")
}

func (c *ContactController) checkPhone(ctx context.Context, req dto.ContactCheckPhoneHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	ownerOf, err := dto.ParseOwnerOf(req.OwnerOf)
	if err != nil {
		return nil, err
	}
	out, _, err := c.uc.SearchByPhone(ctx, ownerOf, req.OwnerID, req.Phone)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}

func (c *ContactController) byProduct(ctx context.Context, req dto.ContactByProductHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.GetListContactByProductId(ctx, req.ID, dto.ContactSearchDTO{
		Pagable: _dto.Pagable{Page: req.Page, Size: req.Size},
	})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *ContactController) updateProducts(ctx context.Context, req dto.ContactUpdateProductsHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	if err := c.uc.UpdateProductIds(ctx, req.ContactID, req.ProductIDs); err != nil {
		return nil, err
	}
	return submitResult(req.ContactID, "updated")
}

func (c *ContactController) products(ctx context.Context, req dto.ContactProductsHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	ids, total, err := c.uc.GetContactProducts(ctx, req.ID, _dto.Pagable{Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(ids, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *ContactController) withMember(ctx context.Context, req dto.ContactWithMemberHTTPRequest) (*httpapi.Result, error) {
	// Handled historically via friend/contact join; expose list/me shape for now.
	return c.listMe(ctx, dto.ContactListMeRequest{Page: req.Page, Size: req.Size, Text: req.Text})
}

func (c *ContactController) friends(ctx context.Context, req dto.ContactWithMemberHTTPRequest) (*httpapi.Result, error) {
	return c.withMember(ctx, req)
}

func (c *ContactController) interested(ctx context.Context, req dto.ContactInterestedHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.GetInterestedContactsByProduct(ctx, &dto.ContactProductInterestedFilter{
		ProductID: req.ProductID,
		Pagable:   _dto.Pagable{Page: req.Page, Size: req.Size},
	})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *ContactController) pin(ctx context.Context, req dto.ContactPinHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	if err := c.uc.PinContact(ctx, &dto.PinRequestDTO{
		Target: enums.PinTarget(req.Target), TargetIDs: req.TargetIDs, NameSpace: req.Namespace, Action: req.Action,
	}); err != nil {
		return nil, err
	}
	return submitResult(0, "ok")
}
