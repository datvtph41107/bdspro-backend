package context

import (
	"context"
	"time"
)

// platform/core/context/values.go

// ─── Route Service ──────────────────────────────────────────────────────────
func WithRouteService(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, keyRouteService, service)
}

func RouteService(ctx context.Context) string {
	v, _ := ctx.Value(keyRouteService).(string)
	return v
}

// ─── Route Timeout ──────────────────────────────────────────────────────────
func WithRouteTimeout(ctx context.Context, timeout time.Duration) context.Context {
	return context.WithValue(ctx, keyRouteTimeout, timeout)
}

func RouteTimeout(ctx context.Context) time.Duration {
	v, _ := ctx.Value(keyRouteTimeout).(time.Duration)
	return v
}

// ─── Route Requires Auth ────────────────────────────────────────────────────
func WithRouteRequiresAuth(ctx context.Context, requires bool) context.Context {
	return context.WithValue(ctx, keyRouteRequiresAuth, requires)
}

func RouteRequiresAuth(ctx context.Context) bool {
	v, _ := ctx.Value(keyRouteRequiresAuth).(bool)
	return v
}

// ─── Route Required Role ────────────────────────────────────────────────────
func WithRouteRequiredRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, keyRouteRequiredRole, role)
}

func RouteRequiredRole(ctx context.Context) string {
	v, _ := ctx.Value(keyRouteRequiredRole).(string)
	return v
}

// ─── Service Instances ──────────────────────────────────────────────────────
func WithServiceInstances(ctx context.Context, instances []string) context.Context {
	return context.WithValue(ctx, keyServiceInstances, instances)
}

func ServiceInstances(ctx context.Context) []string {
	v, _ := ctx.Value(keyServiceInstances).([]string)
	return v
}
