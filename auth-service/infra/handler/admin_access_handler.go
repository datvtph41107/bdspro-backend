package handler

import (
	"auth/internal/usecase"
	authpb "pb/types/auth"
)

type AdminAccessHandler struct {
	AdminAccessUsecase *usecase.AdminAccessUsecase
	authpb.UnimplementedAuthServiceServer
}

func NewAdminAccessHandler(adminAccessUsecase *usecase.AdminAccessUsecase) *AdminAccessHandler {
	return &AdminAccessHandler{
		AdminAccessUsecase: adminAccessUsecase,
	}
}
