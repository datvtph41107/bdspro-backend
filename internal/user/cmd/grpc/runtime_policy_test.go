package main

import "testing"

func TestLoadRuntimePolicyDefaultsAuxiliaryActorsOn(t *testing.T) {
	t.Setenv(envModuleWarmupEnabled, "")
	t.Setenv(envDashboardStatsEnabled, "")
	t.Setenv(envZNSSchedulerEnabled, "")

	policy, err := loadRuntimePolicy()
	if err != nil {
		t.Fatalf("loadRuntimePolicy() error = %v", err)
	}

	if !policy.ModuleWarmup ||
		!policy.DashboardStats ||
		!policy.ZNSScheduler {
		t.Fatalf("policy = %+v, want all defaults enabled", policy)
	}
}

func TestLoadRuntimePolicyCanDisableAuxiliaryActors(t *testing.T) {
	t.Setenv(envModuleWarmupEnabled, "false")
	t.Setenv(envDashboardStatsEnabled, "false")
	t.Setenv(envZNSSchedulerEnabled, "false")

	policy, err := loadRuntimePolicy()
	if err != nil {
		t.Fatalf("loadRuntimePolicy() error = %v", err)
	}

	if policy.ModuleWarmup ||
		policy.DashboardStats ||
		policy.ZNSScheduler {
		t.Fatalf("policy = %+v, want all disabled", policy)
	}
}

func TestLoadRuntimePolicyRejectsInvalidValue(t *testing.T) {
	t.Setenv(envModuleWarmupEnabled, "sometimes")

	if _, err := loadRuntimePolicy(); err == nil {
		t.Fatal("loadRuntimePolicy() error = nil, want invalid config error")
	}
}
