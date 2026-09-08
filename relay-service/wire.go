//go:build wireinject
// +build wireinject

package main

import (
	_db "common/config"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_db.NewDB,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() error {
	wire.Build(wireSet)
	return nil
}
