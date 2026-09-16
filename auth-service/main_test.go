package main

import (
	"errors"
	"log/slog"
	"testing"
)

func TestRunProcessReturnsFailureForRunError(t *testing.T) {
	t.Setenv("QHPRO_LOG_OUTPUT", "stdout")

	original := slog.Default()
	defer slog.SetDefault(original)

	if got := runProcess(func() error {
		return errors.New("boom")
	}); got != 1 {
		t.Fatalf("runProcess() = %d, want 1", got)
	}
}

func TestRunProcessReturnsSuccess(t *testing.T) {
	t.Setenv("QHPRO_LOG_OUTPUT", "stdout")

	original := slog.Default()
	defer slog.SetDefault(original)

	called := false

	got := runProcess(func() error {
		called = true
		return nil
	})

	if got != 0 {
		t.Fatalf("runProcess() = %d, want 0", got)
	}

	if !called {
		t.Fatal("run function was not called")
	}
}

func TestRunProcessStopsWhenLoggingConfigurationFails(t *testing.T) {
	t.Setenv("QHPRO_LOG_OUTPUT", "unsupported")

	original := slog.Default()
	defer slog.SetDefault(original)

	called := false

	got := runProcess(func() error {
		called = true
		return nil
	})

	if got != 1 {
		t.Fatalf("runProcess() = %d, want 1", got)
	}

	if called {
		t.Fatal("run function called after logging configuration failure")
	}
}
