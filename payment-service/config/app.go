package config

import (
	"payment/external/mod/gate/services"
	"payment/external/mod/wallet/router"
	grpc_server "payment/external/server"

	"gorm.io/gorm"
)

var PublicRoutes = []string{
	"/payment/swagger",
	"/payment/payos/webhook/verify",
	"/payment/wallet/account/new",
}

type ConfigApp struct {
	DB                *gorm.DB
	PaymentService    *services.PaymentService
	AccountRouter     *router.AccountRouter
	TransactionRouter *router.TransactionRouter
	GrpcServer        *grpc_server.GrpcServer
}

// var RootApp *ConfigApp

func NewConfigApp(db *gorm.DB,
	fileService *services.PaymentService,
	AccountRouter *router.AccountRouter,
	TransactionRouter *router.TransactionRouter,
	GrpcServer *grpc_server.GrpcServer,
) *ConfigApp {
	return &ConfigApp{
		DB:                db,
		PaymentService:    fileService,
		AccountRouter:     AccountRouter,
		TransactionRouter: TransactionRouter,
		GrpcServer:        GrpcServer,
	}
}
