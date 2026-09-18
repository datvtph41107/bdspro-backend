package quota

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAcceptanceContextHonorsLeaseAndRequestDeadline(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		name        string
		reservation Reservation
		denied      bool
	}{
		{"missing expiry", Reservation{Required: true, State: StateReserved}, true},
		{"expired", Reservation{Required: true, State: StateReserved, ExpiresAt: now.Add(-time.Second)}, true},
		{"safety gap", Reservation{Required: true, State: StateReserved, ExpiresAt: now.Add(5 * time.Second)}, true},
		{"canceled", Reservation{Required: true, State: StateCanceled, ExpiresAt: now.Add(time.Hour)}, true},
		{"live", Reservation{Required: true, State: StateReserved, ExpiresAt: now.Add(time.Hour)}, false},
		{"committed replay", Reservation{Required: true, State: StateCommitted, ExpiresAt: now.Add(-time.Hour)}, false},
		{"unmetered", Reservation{}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parent, cancelParent := context.WithTimeout(context.Background(), time.Second)
			defer cancelParent()
			ctx, cancel, err := AcceptanceContext(parent, tc.reservation, now)
			if tc.denied {
				if !errors.Is(err, ErrAcceptanceLease) {
					t.Fatalf("error=%v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer cancel()
			want, _ := parent.Deadline()
			got, ok := ctx.Deadline()
			if !ok || !got.Equal(want) {
				t.Fatalf("deadline=%v, want %v", got, want)
			}
		})
	}
}
