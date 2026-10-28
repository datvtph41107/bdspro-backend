package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/crm/dto"
	"mdn4/modules/crm/enums"
	"mdn4/modules/crm/services"
)

type LeadController struct{ uc *services.LeadUsecase }

func NewLeadController(uc *services.LeadUsecase) *LeadController { return &LeadController{uc: uc} }

func (c *LeadController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/lead/detail/{id}", c.detail)
	httpapi.GET(mux, "/v2/crm/lead/by-contact/{id}", c.byContact)
	httpapi.POST(mux, "/v2/crm/lead/{ownerOf}/{ownerId}", c.create)
	httpapi.PUT(mux, "/v2/crm/lead/{id}", c.update)
	httpapi.DELETE(mux, "/v2/crm/lead/{id}", c.delete)
	httpapi.POST(mux, "/v2/crm/lead/assign-person", c.assignPerson)
	httpapi.POST(mux, "/v2/crm/lead/assign-crm", c.assignCRM)
	httpapi.GET(mux, "/v2/crm/lead/manage/{ownerOf}/{ownerId}/{pipelineId}", c.manage)
	httpapi.GET(mux, "/v2/crm/lead/care-expired", c.careExpired)
	httpapi.PUT(mux, "/v2/crm/lead/stage", c.switchStage)
}

func (c *LeadController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "lead service not configured")
	}
	return nil
}

type leadCreateReq struct {
	OwnerOf string `path:"ownerOf"`
	OwnerID uint64 `path:"ownerId"`
	dto.LeadSaveDTO
}
type leadManageReq struct {
	OwnerOf    string `path:"ownerOf"`
	OwnerID    uint64 `path:"ownerId"`
	PipelineID uint64 `path:"pipelineId"`
	Page       uint32 `form:"page"`
	Size       uint32 `form:"size"`
}
// leadAssignPersonReq mirrors crmpb.LeadAssignPersonRequest.
type leadAssignPersonReq struct {
	LeadID           uint64 `json:"leadId"`
	OwnerID          uint64 `json:"ownerId"`
	OwnerType        int32  `json:"ownerType"`
	ChargePersonID   uint64 `json:"chargePersonId"`
	ChargePersonType int32  `json:"chargePersonType"`
	NoteAssign       string `json:"noteAssign"`
}
type leadAssignCRMReq struct {
	ContactID      uint64 `json:"contactId"`
	PersonChargeID uint64 `json:"personChargeId"`
	StageID        uint64 `json:"stageId"`
	Note           string `json:"note"`
}
// leadSwitchStageReq mirrors crmpb.LeadStageRequest.
type leadSwitchStageReq struct {
	LeadID    uint64 `json:"leadId"`
	StageID   uint64 `json:"stageId"`
	StageNote string `json:"stageNote"`
}
type leadUpdateReq struct {
	ID uint64 `path:"id" json:"id"`
	dto.LeadSaveDTO
}

func (c *LeadController) detail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *LeadController) byContact(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.GetByContact(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *LeadController) create(ctx context.Context, req leadCreateReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	ownerOf, err := dto.ParseOwnerOf(req.OwnerOf)
	if err != nil {
		return nil, err
	}
	in := req.LeadSaveDTO
	in.OwnerID = req.OwnerID
	in.OwnerOf = ownerOf
	out, err := c.uc.Create(ctx, &in)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *LeadController) update(ctx context.Context, req leadUpdateReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Update(ctx, req.ID, &req.LeadSaveDTO)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *LeadController) delete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "deleted")
}
func (c *LeadController) assignPerson(ctx context.Context, req leadAssignPersonReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Assign(ctx, req.LeadID, &dto.LeadDTO{
		OwnerID:        &req.OwnerID,
		OwnerType:      enums.EOwnerOf(req.OwnerType),
		ChargePersonId: &req.ChargePersonID,
		ChargeType:     enums.EOwnerOf(req.ChargePersonType),
		AssignNote:     req.NoteAssign,
	})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *LeadController) assignCRM(ctx context.Context, req leadAssignCRMReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.AssignCrm(ctx, req.ContactID, req.PersonChargeID, req.StageID, req.Note)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *LeadController) manage(ctx context.Context, req leadManageReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	ownerOf, err := dto.ParseOwnerOf(req.OwnerOf)
	if err != nil {
		return nil, err
	}
	items, total, err := c.uc.GetManager(ctx, dto.LeadSearchDTO{
		OwnerOf: ownerOf, OwnerID: req.OwnerID, PipelineID: &req.PipelineID,
	})
	if err != nil {
		return nil, err
	}
	_ = req.Page
	_ = req.Size
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
func (c *LeadController) careExpired(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	n, err := c.uc.GetCareExpiredTime(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]any{"careExpired": n}), nil
}
func (c *LeadController) switchStage(ctx context.Context, req leadSwitchStageReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if req.LeadID == 0 || req.StageID == 0 {
		return nil, _errors.BadRequestException("leadId and stageId is required")
	}
	out, err := c.uc.SwitchStage(ctx, req.LeadID, &req.StageID, req.StageNote)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
