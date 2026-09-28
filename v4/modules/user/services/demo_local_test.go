package services

import (
	"context"
	"testing"

	_enums "common/domain/enum"
)

func TestLocalDemoService_SetDemoMode_RequiresClients(t *testing.T) {
	svc := NewLocalDemoService(nil, nil, nil, nil, nil, nil, nil, nil)
	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(1))

	err := svc.SetDemoMode(ctx, true)
	if err == nil {
		t.Fatal("expected error when product service missing")
	}

	err = svc.SetDemoMode(ctx, false)
	if err == nil {
		t.Fatal("expected error when product service missing on disable")
	}
}
