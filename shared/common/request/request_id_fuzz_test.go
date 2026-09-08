package request

import "testing"

func FuzzResolveRequestID(fuzzer *testing.F) {
	seedValues := []string{
		"",
		"client-request-123",
		" unsafe ",
		"request\nid",
		"request/id",
	}

	for _, seed := range seedValues {
		fuzzer.Add(seed)
	}

	fuzzer.Fuzz(func(t *testing.T, raw string) {
		decision := ResolveRequestID([]string{raw})

		if !IsValidRequestID(decision.ID) {
			t.Fatalf(
				"canonical request ID invalid: %q",
				decision.ID,
			)
		}

		normalized := NormalizeRequestID(raw)

		if IsValidRequestID(normalized) {
			if decision.ID != normalized {
				t.Fatalf(
					"valid value changed: got %q, want %q",
					decision.ID,
					normalized,
				)
			}

			if decision.Source != RequestIDSourceClient {
				t.Fatalf(
					"valid client value source = %q",
					decision.Source,
				)
			}
		}
	})
}
