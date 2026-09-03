package config

import "testing"

func TestExecutionModeSelectsServiceDNS(t *testing.T) {
	tests := map[string]bool{
		"":          false,
		"host":      false,
		"container": true,
	}
	for mode, want := range tests {
		if got := usesServiceDNS(mode); got != want {
			t.Fatalf("usesServiceDNS(%q) = %t, want %t", mode, got, want)
		}
	}
}

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

func TestRuntimeEndpointUsesServingOwnerAlias(t *testing.T) {
	cfg := Runtime{UseServiceDNS: true, Service: ServiceRuntime{
		Domain: "localhost",
		Host:   map[string]string{"auth": "user"},
		Port:   map[string]int{"auth": 8201},
	}}
	got, err := cfg.Endpoint("auth")
	if err != nil {
		t.Fatal(err)
	}
	if got != "user:8201" {
		t.Fatalf("endpoint = %q", got)
	}
}

func TestRuntimeEndpointRejectsMissingPort(t *testing.T) {
	cfg := Runtime{Service: ServiceRuntime{Domain: "localhost", Port: map[string]int{}}}
	if _, err := cfg.Endpoint("hub"); err == nil {
		t.Fatal("expected missing port error")
	}
}

func TestRuntimeHTTPEndpointCannotFallBackToGRPCPort(t *testing.T) {
	cfg := Runtime{Service: ServiceRuntime{Domain: "localhost", Port: map[string]int{"file": 8203}}}
	if _, err := cfg.HTTPEndpoint("file"); err == nil {
		t.Fatal("HTTP endpoint incorrectly reused the gRPC port")
	}
}

func TestRuntimeHTTPEndpointUsesHTTPServingProcess(t *testing.T) {
	cfg := Runtime{UseServiceDNS: true, Service: ServiceRuntime{
		Domain:   "localhost",
		HTTPHost: map[string]string{"file": "file"},
		HTTPPort: map[string]int{"file": 8002},
	}}
	got, err := cfg.HTTPEndpoint("file")
	if err != nil {
		t.Fatal(err)
	}
	if got != "file:8002" {
		t.Fatalf("HTTP endpoint = %q", got)
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
