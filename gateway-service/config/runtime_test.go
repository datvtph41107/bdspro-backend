package config

import "testing"

func TestRuntimeEndpointUsesLocalDomainWithoutRuntimeEnv(t *testing.T) {
	cfg := Runtime{Service: ServiceRuntime{Domain: "localhost", Port: map[string]int{"hub": 8280}}}
	got, err := cfg.Endpoint("hub")
	if err != nil {
		t.Fatal(err)
	}
	if got != "localhost:8280" {
		t.Fatalf("endpoint = %q", got)
	}
}

func TestRuntimeEndpointUsesServiceDNSWhenEnabled(t *testing.T) {
	cfg := Runtime{UseServiceDNS: true, Service: ServiceRuntime{Domain: "localhost", Port: map[string]int{"tqd": 8219}}}
	got, err := cfg.Endpoint("tqd")
	if err != nil {
		t.Fatal(err)
	}
	if got != "tqd:8219" {
		t.Fatalf("endpoint = %q", got)
	}
}

func TestRuntimeEndpointRejectsMissingPort(t *testing.T) {
	cfg := Runtime{Service: ServiceRuntime{Domain: "localhost", Port: map[string]int{}}}
	if _, err := cfg.Endpoint("hub"); err == nil {
		t.Fatal("expected missing port error")
	}
}

func TestRuntimeValidateRequiresJWTVerificationKey(t *testing.T) {
	cfg := Runtime{
		Server: ServerRuntime{Port: 8000},
		Service: ServiceRuntime{Domain: "localhost", Port: map[string]int{
			"notification": 8204, "payment": 8205, "bdspro": 8202, "crm": 8206,
			"social": 8209, "user": 8201, "chat": 8208, "auth": 8201,
			"assistant": 8218, "hub": 8280, "tqd": 8219,
		}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() accepted missing JWT verification key")
	}
}
