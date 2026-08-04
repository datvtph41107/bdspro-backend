package operationbinding

import "testing"

func TestTransportValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		transport Transport
		wantError bool
	}{
		{name: "unary", transport: GRPCUnary("/tqdpb.ParcelService/GetParcelPlanning")},
		{name: "stream", transport: GRPCServerStream("/tqdpb.TileService/BuildPMTiles")},
		{name: "route", transport: HTTPRoute("get", "/v2/tqd/parcels/{parcelId}/planning")},
		{name: "missing slash", transport: GRPCUnary("tqdpb.ParcelService/GetParcelPlanning"), wantError: true},
		{name: "route query", transport: HTTPRoute("GET", "/v2/tqd/parcels?view=detail"), wantError: true},
		{name: "unsupported method", transport: HTTPRoute("TRACE", "/v2/test"), wantError: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.transport.Validate()
			if (err != nil) != test.wantError {
				t.Fatalf("Validate() error = %v, wantError %v", err, test.wantError)
			}
		})
	}
}
