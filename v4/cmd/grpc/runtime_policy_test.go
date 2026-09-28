package main

import "testing"

func TestLoadRuntimePolicyDefaultsAuxiliaryActorsOff(t *testing.T) {
	t.Setenv(envModuleWarmupEnabled, "")
	t.Setenv(envDashboardStatsEnabled, "")
	t.Setenv(envZNSSchedulerEnabled, "")

	policy, err := loadRuntimePolicy()
	if err != nil {
		t.Fatalf("loadRuntimePolicy() error = %v", err)
	}

	if policy.ModuleWarmup ||
		policy.DashboardStats ||
		policy.ZNSScheduler {
		t.Fatalf("policy = %+v, want all defaults disabled", policy)
	}
}

func TestLoadRuntimePolicyCanEnableAuxiliaryActors(t *testing.T) {
	t.Setenv(envModuleWarmupEnabled, "true")
	t.Setenv(envDashboardStatsEnabled, "true")
	t.Setenv(envZNSSchedulerEnabled, "true")

	policy, err := loadRuntimePolicy()
	if err != nil {
		t.Fatalf("loadRuntimePolicy() error = %v", err)
	}

	if !policy.ModuleWarmup ||
		!policy.DashboardStats ||
		!policy.ZNSScheduler {
		t.Fatalf("policy = %+v, want all enabled", policy)
	}
}

func TestLoadRuntimePolicyRejectsInvalidValue(t *testing.T) {
	t.Setenv(envModuleWarmupEnabled, "sometimes")

	if _, err := loadRuntimePolicy(); err == nil {
		t.Fatal("loadRuntimePolicy() error = nil, want invalid config error")
	}
}
