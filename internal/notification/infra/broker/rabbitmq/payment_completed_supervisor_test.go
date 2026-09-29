package rabbitmq

import (
	"testing"
	"time"
)

func TestNextConsumerRetryDelayIsBounded(t *testing.T) {
	maximum := 30 * time.Second
	cases := []struct {
		current time.Duration
		want    time.Duration
	}{
		{current: time.Second, want: 2 * time.Second},
		{current: 16 * time.Second, want: maximum},
		{current: maximum, want: maximum},
	}
	for _, tc := range cases {
		if got := nextConsumerRetryDelay(tc.current, maximum); got != tc.want {
			t.Fatalf("nextConsumerRetryDelay(%s)=%s, want %s", tc.current, got, tc.want)
		}
	}
}
