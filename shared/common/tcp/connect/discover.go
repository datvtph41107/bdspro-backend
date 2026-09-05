package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DiscoverFromGateway gọi GET {gatewayURL}/internal/services để lấy danh sách
// upstream đang up/down theo gateway.
func DiscoverFromGateway(ctx context.Context, gatewayURL string) (*DiscoverResponse, error) {
	gatewayURL = strings.TrimRight(strings.TrimSpace(gatewayURL), "/")
	if gatewayURL == "" {
		return nil, fmt.Errorf("tcp/connect: empty gateway URL")
	}

	url := gatewayURL + "/internal/services"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("tcp/connect: build discover request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	if deadline, ok := ctx.Deadline(); ok {
		client.Timeout = time.Until(deadline)
		if client.Timeout <= 0 {
			client.Timeout = time.Second
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tcp/connect: discover from gateway: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("tcp/connect: read discover body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tcp/connect: gateway discover status %d: %s", resp.StatusCode, string(body))
	}

	var out DiscoverResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("tcp/connect: decode discover response: %w", err)
	}
	return &out, nil
}

// UpServices lọc các service đang up (optional filter theo tên).
func UpServices(all []ServiceStatus, wanted []string) []ServiceStatus {
	wantSet := map[string]struct{}{}
	for _, w := range wanted {
		w = strings.TrimSpace(w)
		if w != "" {
			wantSet[w] = struct{}{}
		}
	}

	out := make([]ServiceStatus, 0, len(all))
	for _, s := range all {
		if !s.Up || s.Address == "" {
			continue
		}
		if len(wantSet) > 0 {
			if _, ok := wantSet[s.Name]; !ok {
				continue
			}
		}
		out = append(out, s)
	}
	return out
}
