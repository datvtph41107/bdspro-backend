package controller

import (
	"context"
	"net/http"
	"time"

	"common/httpapi"
	"mdn4/modules/user/dto"
	"mdn4/modules/user/services"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

// AdminUserController owns admin-users + admin-profile HTTP routes.
type AdminUserController struct {
	svc services.AdminUserService
}

func NewAdminUserController(svc services.AdminUserService) *AdminUserController {
	return &AdminUserController{svc: svc}
}

func (c *AdminUserController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/user/admin/users", c.listAllUsers)
	httpapi.POST(mux, "/v2/user/admin/users/{profileId}/lock", c.lockUser)
	httpapi.POST(mux, "/v2/user/admin/users/{profileId}/unlock", c.unlockUser)
	httpapi.POST(mux, "/v2/user/admin/users/{profileId}/approve", c.approveUser)
	httpapi.GET(mux, "/v2/user/stats", c.getUserStats)
	httpapi.GET(mux, "/v2/user/admin/users/dashboard/daily", c.getUserGrowthDailyStats)
	httpapi.POST(mux, "/v2/user/admin/user/new", c.createUser)
	httpapi.PUT(mux, "/v2/user/admin/users/{profileId}", c.updateUser)
	mux.HandleFunc("DELETE /v2/user/admin/users/{profileId}", httpapi.BindJSON(c.deleteUser))
	httpapi.GET(mux, "/v2/user/admin/users/{profileId}", c.getUserDetail)

	httpapi.POST(mux, "/v2/user/admin-profile/create", c.createAdmin)
	httpapi.PUT(mux, "/v2/user/admin-profile/update/{id}", c.updateAdmin)
	httpapi.DELETE(mux, "/v2/user/admin-profile/delete/{id}", c.deleteAdmin)
	httpapi.GET(mux, "/v2/user/admin-profile/list", c.listAdmins)
	httpapi.GET(mux, "/v2/user/admin-profile/detail/{id}", c.getAdminDetail)
}

func (c *AdminUserController) listAllUsers(ctx context.Context, req dto.ListAllUsersRequest) (*httpapi.Result, error) {
	result, err := c.svc.ListAllUsers(ctx, &userpb.ListAllUsersRequest{
		Search:    req.Search,
		Page:      req.Page,
		Size:      req.Size,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		RoleType:  req.RoleType,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Bookmark:  req.Bookmark,
		RoleId:    req.RoleID,
		Sort:      req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) lockUser(ctx context.Context, req dto.LockUserRequest) (*httpapi.Result, error) {
	if err := c.svc.LockUser(ctx, &userpb.LockUserRequest{
		ProfileId: req.ProfileID,
		Reason:    req.Reason,
		LockType:  uint32(req.LockType),
	}); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *AdminUserController) unlockUser(ctx context.Context, req dto.UnLockUserRequest) (*httpapi.Result, error) {
	if err := c.svc.UnLockUser(ctx, &userpb.UnLockUserRequest{
		ProfileId: req.ProfileID,
		Reason:    req.Reason,
	}); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *AdminUserController) approveUser(ctx context.Context, req dto.ApproveUserRequest) (*httpapi.Result, error) {
	if err := c.svc.ApproveUser(ctx, &userpb.ApproveUserRequest{
		ProfileId: req.ProfileID,
		Note:      req.Note,
	}); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *AdminUserController) getUserStats(ctx context.Context, req dto.GetUserStatsRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetUserStats(ctx, &sharepb.GetStatsRequest{
		FirstDate:  req.FirstDate,
		SecondDate: req.SecondDate,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) getUserGrowthDailyStats(ctx context.Context, req dto.GetUserGrowthDailyStatsRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetUserGrowthDailyStats(ctx, &userpb.GetUserGrowthDailyStatsRequest{
		LastDates: req.LastDates,
		Gap:       req.Gap,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) createUser(ctx context.Context, req userpb.CreateUserRequest) (*httpapi.Result, error) {
	result, err := c.svc.CreateUser(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) updateUser(ctx context.Context, req dto.UpdateUserRequest) (*httpapi.Result, error) {
	var roleID uint64
	if req.RoleID != nil {
		roleID = *req.RoleID
	}
	var gender uint32
	if req.Gender != nil {
		gender = *req.Gender
	}
	var birth string
	if req.Birth != nil {
		birth = req.Birth.Format(time.RFC3339)
	}
	result, err := c.svc.UpdateUser(ctx, &userpb.UpdateUserRequest{
		ProfileId:    req.ProfileID,
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Avatar:       req.Avatar,
		RoleId:       roleID,
		Gender:       gender,
		Address:      req.Address,
		Birth:        birth,
		Status:       uint32(req.Status),
		IdVerifiedAt: req.IDVerifiedAt,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) deleteUser(ctx context.Context, req dto.DeleteUserRequest) (*httpapi.Result, error) {
	if err := c.svc.DeleteUser(ctx, &userpb.DeleteUserRequest{
		ProfileId: req.ProfileID,
		Reason:    req.Reason,
	}); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *AdminUserController) getUserDetail(ctx context.Context, req dto.AdminProfileIDPath) (*httpapi.Result, error) {
	result, err := c.svc.GetUserDetail(ctx, req.ProfileID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) createAdmin(ctx context.Context, req dto.AdminProfileWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.CreateAdmin(ctx, adminProfileFromWrite(req))
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) updateAdmin(ctx context.Context, req dto.AdminProfileWriteRequest) (*httpapi.Result, error) {
	pb := adminProfileFromWrite(req)
	pb.Id = req.ID
	result, err := c.svc.UpdateAdmin(ctx, pb)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) deleteAdmin(ctx context.Context, req dto.AdminIDPath) (*httpapi.Result, error) {
	result, err := c.svc.DeleteAdmin(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) listAdmins(ctx context.Context, req dto.AdminListRequest) (*httpapi.Result, error) {
	result, err := c.svc.ListAdmins(ctx, &userpb.AdminListRequest{
		Page:   req.Page,
		Size:   req.Size,
		RoleId: req.RoleID,
		Name:   req.Name,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *AdminUserController) getAdminDetail(ctx context.Context, req dto.AdminIDPath) (*httpapi.Result, error) {
	result, err := c.svc.GetAdminDetail(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func adminProfileFromWrite(req dto.AdminProfileWriteRequest) *userpb.AdminProfile {
	return &userpb.AdminProfile{
		FullName:         req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		Avatar:           req.Avatar,
		Address:          req.Address,
		Gender:           req.Gender,
		Birth:            req.Birth,
		Status:           req.Status,
		RoleType:         req.RoleType,
		RoleRealEstate:   req.RoleRealEstate,
		Position:         req.Position,
		DepartmentId:     req.DepartmentID,
		Bookmark:         req.Bookmark,
		TickVerified:     req.TickVerified,
		Attachments:      req.Attachments,
		InternalNotes:    req.InternalNotes,
		SendNotification: req.SendNotification,
		JobTitle:         req.JobTitle,
		WorkAt:           req.WorkAt,
		RoleKey:          req.RoleKey,
		RoleId:           req.RoleID,
		Username:         req.Username,
		Password:         req.Password,
		SuperAdmin:       req.SuperAdmin,
	}
}
