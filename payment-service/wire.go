//go:build wireinject
// +build wireinject

package main

import (
	_db "common/db"
	"payment/repo"
	"payment/services"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_db.NewDB,
	repo.NewPaymentRepo,
	NewConfigApp,
	services.NewPaymentService,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
