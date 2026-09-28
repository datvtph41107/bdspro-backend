package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	envModuleWarmupEnabled = "QHPRO_USER_MODULE_WARMUP_ENABLED"

	envDashboardStatsEnabled = "QHPRO_USER_DASHBOARD_STATS_ENABLED"

	envZNSSchedulerEnabled = "QHPRO_USER_ZNS_SCHEDULER_ENABLED"
)

// runtimePolicy controls process-owned auxiliary work.
//
// Request-serving capabilities are not feature-flagged here. These switches
// only decide whether optional startup/background actors participate in this
// process instance.
type runtimePolicy struct {
	ModuleWarmup   bool
	DashboardStats bool
	ZNSScheduler   bool
}

func loadRuntimePolicy() (runtimePolicy, error) {
	// v4 defaults auxiliary actors off so a side-by-side process does not
	// double-run jobs against the shared user database.
	moduleWarmup, err :=
		boolEnv(envModuleWarmupEnabled, false)
	if err != nil {
		return runtimePolicy{}, err
	}

	dashboardStats, err :=
		boolEnv(envDashboardStatsEnabled, false)
	if err != nil {
		return runtimePolicy{}, err
	}

	znsScheduler, err :=
		boolEnv(envZNSSchedulerEnabled, false)
	if err != nil {
		return runtimePolicy{}, err
	}

	return runtimePolicy{
		ModuleWarmup:   moduleWarmup,
		DashboardStats: dashboardStats,
		ZNSScheduler:   znsScheduler,
	}, nil
}

func boolEnv(name string, defaultValue bool) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf(
			"parse %s=%q as boolean: %w",
			name,
			raw,
			err,
		)
	}

	return value, nil
}
