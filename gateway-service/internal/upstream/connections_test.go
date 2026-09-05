package upstream

import (
	"common/rpc"
	gatewayconfig "gateway/config"
	"testing"
)

func TestOpenBuildsAndClosesExplicitGatewayTopology(t *testing.T) {
	ports := map[string]int{
		"notification": 8204,
		"payment":      8205,
		"bdspro":       8202,
		"crm":          8206,
		"organization": 8207,
		"social":       8209,
		"user":         8201,
		"chat":         8208,
		"auth":         8201,
		"assistant":    8218,
		"hub":          8280,
		"tqd":          8219,
	}
	cfg := gatewayconfig.Runtime{
		Service: gatewayconfig.ServiceRuntime{Domain: "127.0.0.1", Port: ports},
		RPCTransport: rpc.TransportConfig{
			ServiceAssertion: rpc.ServiceAssertionConfig{
				ServiceID: "gateway-service",
			},
		},
	}
	connections, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if connections.Hub == nil || connections.TQD == nil || connections.User == nil {
		t.Fatal("expected explicit Gateway connections")
	}
	if err := connections.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
