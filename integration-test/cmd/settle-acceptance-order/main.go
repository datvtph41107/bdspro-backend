// Command settle-acceptance-order submits a provider-authenticated SePay event
// to the configured Gateway. Payment Service still owns
// reference matching, exact-money validation, settlement and fulfillment.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		gateway     = flag.String("gateway", envOr("QHPRO_GATEWAY_URL", "http://127.0.0.1:8000"), "Gateway base URL")
		reference   = flag.String("reference", "", "commercial order reference shown by the client")
		amountMinor = flag.Int64("amount-minor", 0, "exact frozen VND amount in minor units")
		transaction = flag.Uint64("transaction-id", uint64(time.Now().UnixNano()%int64(^uint32(0))), "unique provider transaction ID")
		providerKey = flag.String("provider-key", strings.TrimSpace(os.Getenv("SEPAY_API_KEY")), "acceptance SePay provider credential")
	)
	flag.Parse()

	ref := strings.TrimSpace(*reference)
	if ref == "" || *amountMinor <= 0 || *transaction == 0 || *transaction > uint64(^uint32(0)) || strings.TrimSpace(*providerKey) == "" {
		fatal("reference, positive amount-minor, transaction-id and SEPAY_API_KEY are required")
	}
	if *amountMinor > 9_000_000_000_000_000 {
		fatal("amount-minor is outside the supported exact JSON integer range")
	}

	now := time.Now().In(time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60))
	payload := map[string]any{
		"gateway":         "QHPRO_ACCEPTANCE",
		"transactionDate": now.Format("2006-01-02 15:04:05"),
		"accountNumber":   "ACCEPTANCE",
		"content":         ref,
		"transferType":    "in",
		"description":     "QHPRO acceptance provider settlement",
		"transferAmount":  *amountMinor,
		"referenceCode":   ref,
		"id":              *transaction,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		fatal("encode request: %v", err)
	}

	url := strings.TrimRight(*gateway, "/") + "/v2/payment/sepay/webhook"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		fatal("create request: %v", err)
	}
	req.Header.Set("Authorization", strings.TrimSpace(*providerKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fatal("submit provider event: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		fatal("read response: %v", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fatal("provider event rejected: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var envelope struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &envelope); err != nil {
			fatal("decode response: %v", err)
		}
		if envelope.Code != 0 {
			fatal("provider event rejected: code=%d message=%s", envelope.Code, envelope.Message)
		}
	}

	fmt.Printf("Settlement evidence accepted: reference=%s transaction_id=%d amount_minor=%d\n", ref, *transaction, *amountMinor)
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "settle-acceptance-order: "+format+"\n", args...)
	os.Exit(1)
}
