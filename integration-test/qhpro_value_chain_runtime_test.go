package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestQHPROValueChainRuntime is the release-level proof for the first QHPRO
// sellable capability. It crosses only public boundaries; no subscription,
// quota, usage, report or notification business row is inserted by the test.
func TestQHPROValueChainRuntime(t *testing.T) {
	if strings.TrimSpace(os.Getenv(runtimeSmokeEnabled)) != "1" {
		t.Skip(runtimeSmokeEnabled + "=1 is required")
	}

	gatewayURL := strings.TrimRight(envOrDefault("QHPRO_GATEWAY_URL", GatewayURL), "/")
	username := strings.TrimSpace(os.Getenv("QHPRO_ACCEPTANCE_USERNAME"))
	password := strings.TrimSpace(os.Getenv("QHPRO_ACCEPTANCE_USER_PASSWORD"))
	providerKey := strings.TrimSpace(os.Getenv("SEPAY_API_KEY"))
	if username == "" || password == "" || providerKey == "" {
		t.Fatal("QHPRO_ACCEPTANCE_USERNAME, QHPRO_ACCEPTANCE_USER_PASSWORD and SEPAY_API_KEY are required")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	loginStatus, login := runtimeJSONPostWithHeaders(t, client, gatewayURL+"/v2/auth/password/login", map[string]any{
		"username": username, "password": password, "version": "acceptance",
		"platform": "WEB", "os": "Linux", "deviceName": "QHPRO Value Chain Acceptance",
	}, nil)
	token := jsonString(login, "accessToken", "access_token")
	profileID := jsonNumber(login, "profileId", "profile_id")
	if loginStatus != http.StatusOK || token == "" || profileID <= 0 {
		t.Fatalf("login transport: status=%d body=%v", loginStatus, login)
	}
	auth := map[string]string{"Authorization": "Bearer " + token}

	plansStatus, plansPayload := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/user/commercial/plans?productCode=qhpro&subjectKind=profile", "")
	if plansStatus != http.StatusOK {
		t.Fatalf("catalog: status=%d body=%v", plansStatus, plansPayload)
	}
	planCode, amountMinor := paidPlan(plansPayload, "qhpro.basic")
	if planCode == "" || amountMinor <= 0 {
		t.Fatalf("canonical paid plan is unavailable: %v", plansPayload)
	}

	runKey := "accept-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	checkoutHeaders := cloneHeaders(auth)
	checkoutHeaders["Idempotency-Key"] = runKey + "-checkout"
	checkoutStatus, checkout := runtimeJSONPostWithHeaders(t, client, gatewayURL+"/v2/user/checkout", map[string]any{
		"planCode": planCode, "subjectKind": "profile",
	}, checkoutHeaders)
	orderID := jsonNumber(checkout, "orderId", "order_id")
	reference := jsonString(checkout, "reference")
	if checkoutStatus != http.StatusOK || orderID <= 0 || reference == "" {
		t.Fatalf("checkout: status=%d body=%v", checkoutStatus, checkout)
	}

	attemptHeaders := cloneHeaders(auth)
	attemptHeaders["Idempotency-Key"] = runKey + "-attempt"
	attemptStatus, attempt := runtimeJSONPostWithHeaders(t, client, gatewayURL+"/v2/user/checkout/payment-attempts", map[string]any{
		"orderId": orderID, "subjectKind": "profile", "method": "bank_transfer",
	}, attemptHeaders)
	if attemptStatus != http.StatusOK || jsonNumber(attempt, "attemptId", "attempt_id") <= 0 || jsonString(attempt, "provider") != "sepay" {
		t.Fatalf("payment attempt: status=%d body=%v", attemptStatus, attempt)
	}

	transactionID := uint64(time.Now().UTC().UnixNano()%4_000_000_000) + 1
	settlementStatus, settlement := runtimeJSONPostWithHeaders(t, client, gatewayURL+"/v2/payment/sepay/webhook", map[string]any{
		"gateway": "QHPRO_ACCEPTANCE", "transactionDate": time.Now().Format("2006-01-02 15:04:05"),
		"accountNumber": "ACCEPTANCE", "content": reference, "transferType": "in",
		"description": "QHPRO canonical value-chain acceptance", "transferAmount": amountMinor,
		"referenceCode": reference, "id": transactionID,
	}, map[string]string{"Authorization": providerKey})
	if settlementStatus != http.StatusOK || jsonNumber(settlement, "code") != 0 {
		t.Fatalf("provider settlement: status=%d body=%v", settlementStatus, settlement)
	}

	waitFor(t, 45*time.Second, "Payment fulfillment and User subscription", func() (bool, string) {
		status, progress := runtimeJSONRequestAny(t, client, fmt.Sprintf("%s/v2/payment/commercial/orders/%d", gatewayURL, orderID), token)
		if status != http.StatusOK {
			return false, fmt.Sprintf("status=%d body=%v", status, progress)
		}
		return jsonString(progress, "progressState", "progress_state") == "ACTIVE", fmt.Sprint(progress)
	})

	profileStatus, profile := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/user/commercial/profile", token)
	if profileStatus != http.StatusOK || !jsonBool(profile, "hasSubscription", "has_subscription") {
		t.Fatalf("subscription projection: status=%d body=%v", profileStatus, profile)
	}

	usageStatus, usageBefore := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/tqd/commercial/generated-report-usage", token)
	if usageStatus != http.StatusOK || !jsonBool(usageBefore, "allowed") {
		t.Fatalf("entitlement/quota projection: status=%d body=%v", usageStatus, usageBefore)
	}
	usedBefore := jsonNumber(usageBefore, "durableUsed", "durable_used")
	reportHeaders := cloneHeaders(auth)
	reportHeaders["Idempotency-Key"] = runKey + "-report"
	reportBody := map[string]any{
		// A fresh database intentionally contains no fabricated planning parcel.
		// A location snapshot is a real sellable report target and proves the
		// same entitlement, quota, durable acceptance, renderer and File path.
		"reportType": 2, "entityType": 2, "entityId": profileID,
		"title": "Báo cáo nghiệm thu chuỗi giá trị", "subtitle": username,
		"location": map[string]any{"address": "QHPRO acceptance", "province": "Hà Nội"},
		"spatial": map[string]any{
			"centroid": map[string]any{"lat": 21.0278, "lon": 105.8342},
			"bounds":   map[string]any{"minLon": 105.83, "minLat": 21.02, "maxLon": 105.84, "maxLat": 21.03},
		},
	}
	reportStatus, report := runtimeJSONPostWithHeaders(t, client, gatewayURL+"/v2/tqd/map-workspace/reports", reportBody, reportHeaders)
	item := jsonObject(report, "item")
	reportID := jsonNumber(item, "reportId", "report_id")
	if reportStatus != http.StatusOK || reportID <= 0 {
		t.Fatalf("report acceptance: status=%d body=%v", reportStatus, report)
	}

	replayStatus, replay := runtimeJSONPostWithHeaders(t, client, gatewayURL+"/v2/tqd/map-workspace/reports", reportBody, reportHeaders)
	if replayStatus != http.StatusOK || jsonNumber(jsonObject(replay, "item"), "reportId", "report_id") != reportID {
		t.Fatalf("report idempotent replay: status=%d body=%v", replayStatus, replay)
	}

	var pdfURL string
	waitFor(t, 60*time.Second, "report worker and File artifact", func() (bool, string) {
		status, current := runtimeJSONRequestAny(t, client, fmt.Sprintf("%s/v2/tqd/map-workspace/reports/%d", gatewayURL, reportID), token)
		assets := jsonObject(current, "assets")
		pdfURL = jsonString(assets, "pdfUrl", "pdf_url")
		return status == http.StatusOK && pdfURL != "", fmt.Sprintf("status=%d body=%v", status, current)
	})
	assertPDF(t, client, pdfURL)

	usageAfterStatus, usageAfter := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/tqd/commercial/generated-report-usage", token)
	if usageAfterStatus != http.StatusOK || jsonNumber(usageAfter, "durableUsed", "durable_used") != usedBefore+1 {
		t.Fatalf("durable usage/idempotency: before=%v after=%v", usageBefore, usageAfter)
	}

	waitFor(t, 45*time.Second, "Payment outbox, broker Inbox and customer notification", func() (bool, string) {
		status, notifications := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/notification/list?page=0&size=20&timestamp=0", token)
		for _, raw := range jsonArray(notifications, "data") {
			row, _ := raw.(map[string]any)
			if jsonString(row, "title") == "Thanh toán thành công" && jsonNumber(row, "targetId", "target_id") == orderID {
				return status == http.StatusOK, fmt.Sprint(notifications)
			}
		}
		return false, fmt.Sprintf("status=%d body=%v", status, notifications)
	})

	t.Logf("PASS value chain: profile=%d plan=%s order=%d report=%d used=%d pdf=%s", profileID, planCode, orderID, reportID, usedBefore+1, pdfURL)
}

func runtimeJSONPostWithHeaders(t *testing.T, client *http.Client, url string, value map[string]any, headers map[string]string) (int, map[string]any) {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		t.Fatalf("POST %s: read response: %v", url, err)
	}
	body := map[string]any{}
	if len(bytes.TrimSpace(bodyBytes)) != 0 {
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			t.Fatalf("POST %s: decode status %d: %v body=%q", url, resp.StatusCode, err, bodyBytes)
		}
	}
	payload := body
	if data, ok := body["data"].(map[string]any); ok {
		payload = data
	}
	return resp.StatusCode, payload
}

func paidPlan(payload map[string]any, desired string) (string, int64) {
	for _, raw := range jsonArray(payload, "plans") {
		plan, _ := raw.(map[string]any)
		if jsonString(plan, "planCode", "plan_code") != desired {
			continue
		}
		for _, priceRaw := range jsonArray(plan, "prices") {
			price, _ := priceRaw.(map[string]any)
			if jsonString(price, "kind") == "recurring" {
				return desired, jsonNumber(price, "amountMinor", "amount_minor")
			}
		}
	}
	return "", 0
}

func cloneHeaders(source map[string]string) map[string]string {
	result := make(map[string]string, len(source)+1)
	for key, value := range source {
		result[key] = value
	}
	return result
}

func jsonObject(value map[string]any, key string) map[string]any {
	object, _ := value[key].(map[string]any)
	return object
}

func waitFor(t *testing.T, timeout time.Duration, name string, check func() (bool, string)) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	last := "not observed"
	for time.Now().Before(deadline) {
		ok, detail := check()
		last = detail
		if ok {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s: %s", name, last)
}
