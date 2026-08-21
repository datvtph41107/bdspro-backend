package cmd

import (
	"context"
	"testing"
)

func TestStartReportTargetRequiresCanonicalDependencies(t *testing.T) {
	t.Parallel()

	closeTarget, err := startReportTarget(
		context.Background(),
		nil,
		nil,
		nil,
	)
	if err == nil {
		t.Fatal("startReportTarget() must fail when canonical dependencies are missing")
	}
	if closeTarget != nil {
		t.Fatal("failed target startup must not return a live cleanup handle")
	}
}
