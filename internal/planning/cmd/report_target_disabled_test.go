package cmd

import (
	"context"
	"log/slog"
	"testing"

	"tqd/config"
)

func TestStartReportTargetRequiresCanonicalDependencies(t *testing.T) {
	t.Parallel()

	closeTarget, err := startReportTarget(
		context.Background(),
		nil,
		nil,
		nil,
		nil,
		config.ReportRuntimeConfig{},
		slog.Default(),
	)
	if err == nil {
		t.Fatal("startReportTarget() must fail when canonical dependencies are missing")
	}
	if closeTarget != nil {
		t.Fatal("failed target startup must not return a live cleanup handle")
	}
}
