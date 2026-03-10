package api

import (
	"context"
	"reflect"

	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"core/kernel.v2/src/internal/registry"
	"core/kernel.v2/src/internal/runtime/orchestrator"
	"core/kernel.v2/src/internal/runtime/signal"
	"core/kernel.v2/src/modules/database"
	"core/kernel.v2/src/modules/http"
	"core/kernel.v2/src/modules/logger"

	"github.com/google/martian/v3/h2/grpc"
)

var (
	procRegistry contract_runtime.ProcessorRegistry
)

func init() {
	procRegistry = registry.NewProcessorRegistry()
	// Đăng ký tất cả processors từ các modules
	registerProcessors()
}

func registerProcessors() {
	// Logger
	procRegistry.Register(&logger.Processor{})
	// Tracer
	procRegistry.Register(&tracer.Processor{})
	// Metrics
	procRegistry.Register(&metrics.Processor{})
	// Database
	procRegistry.Register(&database.Processor{})
	// HTTP
	procRegistry.Register(&http.Processor{})
	// gRPC
	procRegistry.Register(&grpc.Processor{})
	// Job
	procRegistry.Register(&job.Processor{})
}

// Run tự động phát hiện các component trong app và chạy.
func Run(ctx context.Context, app interface{}) error {
	orch := orchestrator.New(procRegistry)

	// Dùng reflection để tìm tất cả các field implement Component
	if err := registerComponents(orch, app); err != nil {
		return err
	}

	if err := orch.Start(ctx); err != nil {
		return err
	}

	waiter := signal.NewWaiter()
	return waiter.Wait(ctx, orch.Shutdown)
}

func registerComponents(orch runtime.Orchestrator, app interface{}) error {
	val := reflect.ValueOf(app)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsNil() {
			continue
		}
		// Nếu field implement runtime.Component
		if comp, ok := field.Interface().(runtime.Component); ok {
			if err := orch.RegisterComponent(comp); err != nil {
				return err
			}
		}
		// Nếu field là struct, đệ quy (hỗ trợ lồng)
		if field.Kind() == reflect.Struct {
			if err := registerComponents(orch, field.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}
