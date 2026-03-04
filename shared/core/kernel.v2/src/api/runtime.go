package api

import (
	"context"
	"core/kernel.v2/src/internal/runtime/signal"
	"reflect"
	"syscall"
)

// Run tự động phát hiện các component trong app và chạy.
// app: struct chứa các field implement runtime.Component
func Run(ctx context.Context, app interface{}) error {
	orch := orchestrator.New()

	// Dùng reflection để tìm tất cả component
	if err := registerComponents(orch, app); err != nil {
		return err
	}

	if err := orch.Start(ctx); err != nil {
		return err
	}

	waiter := signal.NewWaiter(syscall.SIGINT, syscall.SIGTERM)
	return waiter.Wait(ctx, orch.Shutdown)
}

// registerComponents đệ quy tìm các field implement Component
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

		// Nếu field implement Component
		if comp, ok := field.Interface().(runtime.Component); ok {
			if err := orch.RegisterComponent(comp); err != nil {
				return err
			}
			continue
		}

		// Nếu field là struct, đệ quy vào trong
		if field.Kind() == reflect.Struct {
			if err := registerComponents(orch, field.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}
