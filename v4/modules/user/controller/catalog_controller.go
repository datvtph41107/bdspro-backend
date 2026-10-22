package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/user/dto"
	"mdn4/modules/user/services"
	userpb "pb/types/user"
)

// CatalogController owns main-area / purpose-use / profession HTTP routes.
type CatalogController struct {
	svc services.CatalogService
}

func NewCatalogController(svc services.CatalogService) *CatalogController {
	return &CatalogController{svc: svc}
}

func (c *CatalogController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/user/admin/main-area", c.getMainAreas)
	httpapi.POST(mux, "/v2/user/admin/main-area", c.createMainArea)
	httpapi.PUT(mux, "/v2/user/admin/main-area/{id}", c.updateMainArea)
	httpapi.DELETE(mux, "/v2/user/admin/main-area/{id}", c.deleteMainArea)
	httpapi.GET(mux, "/v2/user/main-area", c.getMainAreaItems)

	httpapi.GET(mux, "/v2/user/admin/purpose-use", c.getPurposeUses)
	httpapi.POST(mux, "/v2/user/admin/purpose-use", c.createPurposeUse)
	httpapi.PUT(mux, "/v2/user/admin/purpose-use/{id}", c.updatePurposeUse)
	httpapi.DELETE(mux, "/v2/user/admin/purpose-use/{id}", c.deletePurposeUse)
	httpapi.GET(mux, "/v2/user/purpose-use", c.getPurposeUseItems)

	httpapi.GET(mux, "/v2/user/admin/profession", c.listProfessions)
	httpapi.POST(mux, "/v2/user/admin/profession", c.createProfession)
	httpapi.PUT(mux, "/v2/user/admin/profession/{id}", c.updateProfession)
	httpapi.DELETE(mux, "/v2/user/admin/profession/{id}", c.deleteProfession)
	httpapi.GET(mux, "/v2/user/professions", c.getProfessionItems)
}

func (c *CatalogController) getMainAreas(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.svc.GetMainAreas(ctx)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) createMainArea(ctx context.Context, req dto.MainAreaWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.CreateMainArea(ctx, &userpb.MainArea{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) updateMainArea(ctx context.Context, req dto.MainAreaWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.UpdateMainArea(ctx, &userpb.MainArea{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) deleteMainArea(ctx context.Context, req dto.CatalogIDPath) (*httpapi.Result, error) {
	result, err := c.svc.DeleteMainArea(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) getMainAreaItems(ctx context.Context, req dto.CatalogSearchRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetMainAreaItems(ctx, req.Text, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) getPurposeUses(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.svc.GetPurposeUses(ctx)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) createPurposeUse(ctx context.Context, req dto.PurposeUseWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.CreatePurposeUse(ctx, &userpb.PurposeUse{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) updatePurposeUse(ctx context.Context, req dto.PurposeUseWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.UpdatePurposeUse(ctx, &userpb.PurposeUse{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) deletePurposeUse(ctx context.Context, req dto.CatalogIDPath) (*httpapi.Result, error) {
	result, err := c.svc.DeletePurposeUse(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) getPurposeUseItems(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.svc.GetPurposeUseItems(ctx)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) listProfessions(ctx context.Context, req dto.CatalogPageRequest) (*httpapi.Result, error) {
	result, err := c.svc.ListProfessions(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) createProfession(ctx context.Context, req dto.ProfessionWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.CreateProfession(ctx, professionFromWrite(req))
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) updateProfession(ctx context.Context, req dto.ProfessionWriteRequest) (*httpapi.Result, error) {
	pb := professionFromWrite(req)
	pb.Id = req.ID
	result, err := c.svc.UpdateProfession(ctx, pb)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) deleteProfession(ctx context.Context, req dto.CatalogIDPath) (*httpapi.Result, error) {
	result, err := c.svc.DeleteProfession(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *CatalogController) getProfessionItems(ctx context.Context, req dto.CatalogPageRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetProfessionItems(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func professionFromWrite(req dto.ProfessionWriteRequest) *userpb.Profession {
	return &userpb.Profession{
		Name:           req.Name,
		Issuer:         req.Issuer,
		IssueDate:      req.IssueDate,
		VerifiedStatus: req.VerifiedStatus,
		ApproveUserId:  req.ApproveUserID,
		IsActive:       req.IsActive,
	}
}
