package controller

import (
	"context"
	"net/http"

	base_enum "base/enum"
	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/enums"
	"userv4/modules/crm/services"
)

// ContactController owns /v2/crm/contact/* routes.
type ContactController struct {
	local services.ContactService
	uc    *services.ContactUsecase
}

func NewContactController(local services.ContactService, uc *services.ContactUsecase) *ContactController {
	return &ContactController{local: local, uc: uc}
}

func (c *ContactController) RegisterHTTP(mux *http.ServeMux) {
	// Paths follow shared/protobuf/schema/crm/contact.proto (frontend contract).
	httpapi.GET(mux, "/v2/crm/contact/list/me", c.listMe)
	httpapi.GET(mux, "/v2/crm/contact/list/user/{ownerId}", c.listUser)
	httpapi.GET(mux, "/v2/crm/contact/list/{ownerOf}/{ownerId}", c.search)
	httpapi.GET(mux, "/v2/crm/contact/detail/{id}/products", c.products)
	httpapi.GET(mux, "/v2/crm/contact/detail/{id}", c.detail)
	httpapi.POST(mux, "/v2/crm/contact", c.create)
	// PUT edit/{id}, note/{id} and {contactId}/products overlap (e.g. /edit/products),
	// which ServeMux rejects, so a single pattern dispatches them.
	mux.HandleFunc("PUT /v2/crm/contact/{first}/{second}", c.putDispatch(
		httpapi.BindJSON(c.update),
		httpapi.BindJSON(c.note),
		httpapi.BindJSON(c.updateProducts),
	))
	httpapi.DELETE(mux, "/v2/crm/contact/{id}", c.delete)
	httpapi.POST(mux, "/v2/crm/contact/sync", c.sync)
	httpapi.GET(mux, "/v2/crm/contact/relation-ship/{id}", c.relationShip)
	httpapi.GET(mux, "/v2/crm/contact/phone/{ownerOf}/{ownerId}/{phone}", c.checkPhone)
	httpapi.GET(mux, "/v2/crm/contact/by-product/{id}", c.byProduct)
	httpapi.GET(mux, "/v2/crm/contact/with-member", c.withMember)
	httpapi.GET(mux, "/v2/crm/contact/friends", c.friends)
	httpapi.GET(mux, "/v2/crm/contact/interested/product/{productId}", c.interested)
	httpapi.POST(mux, "/v2/crm/contact/pin", c.pin)
}

func (c *ContactController) putDispatch(edit, note, products http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		first, second := r.PathValue("first"), r.PathValue("second")
		switch {
		case first == "edit":
			r.SetPathValue("id", second)
			edit.ServeHTTP(w, r)
		case first == "note":
			r.SetPathValue("id", second)
			note.ServeHTTP(w, r)
		case second == "products":
			r.SetPathValue("id", first)
			products.ServeHTTP(w, r)
		default:
			httpapi.WriteError(w, http.StatusNotFound, "route not found")
		}
	}
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
	return contactListResult(result)
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
	items, total, err := c.uc.Search(ctx, base_enum.EOwnerOf(req.OwnerOf), dto.ContactSearchDTO{
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
	return contactDetailResult(out)
}

func (c *ContactController) create(ctx context.Context, req dto.ContactCreateHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireUC(); err != nil {
		return nil, err
	}
	out, err := c.uc.Create(ctx, &domain.ContactEntity{
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
	out, err := c.uc.Update(ctx, req.ID, &domain.ContactEntity{
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
	entities := make([]domain.ContactEntity, 0, len(req.Contacts))
	for _, item := range req.Contacts {
		entities = append(entities, domain.ContactEntity{
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
	out, _, err := c.uc.SearchByPhone(ctx, base_enum.EOwnerOf(req.OwnerOf), req.OwnerID, req.Phone)
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
