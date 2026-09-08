package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// This module is a black-box runtime suite, not a unit-test package. Keep
	// `go test ./...` safe on a clean clone and require an owning root target to
	// opt in after it has started and checked the backend topology.
	if strings.TrimSpace(os.Getenv("QHPRO_INTEGRATION_TEST")) != "1" &&
		strings.TrimSpace(os.Getenv(runtimeSmokeEnabled)) != "1" {
		fmt.Println("integration runtime not requested; use a root Make acceptance target")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
