//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	"map/config"
	"map/infra/db"
	"map/infra/handler"
	"map/infra/postgre"
	"map/initial"
	"map/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var WireSet = wire.NewSet(
	// Config
	config.LoadConfig,

	// Database
	_db.NewDB,
	db.Migrate,

	// Repository
	postgre.NewMapPointPostgres,

	// Usecase
	usecase.NewMapPointUsecase,

	// Handler
	handler.NewMapPointHandler,

	// Router
	NewRouter,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*gin.Engine, error) {
	wire.Build(WireSet)
	return nil, nil
}

// InitializeRuntime khởi tạo runtime
func InitializeRuntime() error {
	return initial.InitializeRuntime()
}
