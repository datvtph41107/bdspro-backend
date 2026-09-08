package delivery

import "fmt"

// EffectFailureKind captures what the process can know about an external
// delivery failure. Unknown means the provider may have accepted the request
// before the response was lost, so retry intentionally has at-least-once
// semantics.
type EffectFailureKind string

const (
	EffectFailurePermanent EffectFailureKind = "permanent"
	EffectFailureRetryable EffectFailureKind = "retryable"
	EffectFailureUnknown   EffectFailureKind = "unknown"
)

type EffectError struct {
	Kind EffectFailureKind
	Err  error
}

func (e *EffectError) Error() string {
	if e == nil {
		return "delivery effect failed"
	}
	if e.Err == nil {
		return fmt.Sprintf("delivery effect failed: %s", e.Kind)
	}
	return e.Err.Error()
}

func (e *EffectError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
