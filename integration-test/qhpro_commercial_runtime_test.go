package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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

const runtimeSmokeEnabled = "QHPRO_RUNTIME_SMOKE"

const (
	acceptanceAdminUsername = "qhpro.acceptance.admin"
	acceptanceUserUsername  = "qhpro.acceptance.user"
)

// TestQHPROCommercialRuntime verifies an already-provisioned commercial state
// through Gateway. It is opt-in because it requires the Compose stack and the
// durable checkout/report acceptance data.
func TestQHPROCommercialRuntime(t *testing.T) {
	if strings.TrimSpace(os.Getenv(runtimeSmokeEnabled)) != "1" {
		t.Skip(runtimeSmokeEnabled + "=1 is required")
	}

	gatewayURL := strings.TrimRight(envOrDefault("QHPRO_GATEWAY_URL", GatewayURL), "/")
	secret := strings.TrimSpace(os.Getenv("JWT_KEY_GENERATE"))
	if secret == "" {
		t.Fatal("JWT_KEY_GENERATE is required")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	status, _ := runtimeJSONRequest(t, client, gatewayURL+"/v2/user/commercial/profile", "")
	if status != http.StatusUnauthorized {
		t.Fatalf("commercial profile without token: status=%d, want=%d", status, http.StatusUnauthorized)
	}
	planStatus, publicPlans := runtimeJSONRequest(t, client, gatewayURL+"/v2/user/commercial/plans?productCode=qhpro&subjectKind=profile", "")
	if planStatus != http.StatusOK || len(jsonArray(publicPlans, "plans")) == 0 {
		t.Fatalf("public commercial catalog: status=%d body=%v", planStatus, publicPlans)
	}

	token := signRuntimeAccessToken(t, secret, 1)
	adminStatus, adminBody := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/user/admin/commercial/subscriptions", token)
	if adminStatus != http.StatusForbidden {
		t.Fatalf("commercial admin IAM boundary: status=%d, want=%d body=%v", adminStatus, http.StatusForbidden, adminBody)
	}

	_, profile := runtimeJSONRequest(t, client, gatewayURL+"/v2/user/commercial/profile", token)
	if jsonNumber(profile, "profileId", "profile_id") != 1 || !jsonBool(profile, "hasSubscription", "has_subscription") {
		t.Fatalf("commercial profile fixture is not active: %v", profile)
	}

	_, usage := runtimeJSONRequest(t, client, gatewayURL+"/v2/tqd/commercial/generated-report-usage", token)
	if jsonNumber(usage, "durableUsed", "durable_used") < 2 || jsonNumber(usage, "remaining") < 0 {
		t.Fatalf("usage projection lost durable evidence: %v", usage)
	}

	_, reports := runtimeJSONRequest(t, client, gatewayURL+"/v2/tqd/map-workspace/reports?page=0&size=20", token)
	reportItems := jsonArray(reports, "data")
	if len(reportItems) < 2 {
		t.Fatalf("generated report fixture is incomplete: %v", reports)
	}

	pdfURL := ""
	for _, raw := range reportItems {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		reportMeta, _ := item["reportMeta"].(map[string]any)
		if reportMeta == nil {
			reportMeta, _ = item["report_meta"].(map[string]any)
		}
		format := jsonString(item, "format")
		if format == "" {
			format = jsonString(reportMeta, "format")
		}
		if !strings.EqualFold(format, "pdf") {
			continue
		}
		pdfURL = jsonString(item, "pdfUrl", "pdf_url")
		if assets, ok := item["assets"].(map[string]any); ok && pdfURL == "" {
			pdfURL = jsonString(assets, "pdfUrl", "pdf_url")
		}
		if pdfURL != "" {
			break
		}
	}
	if pdfURL == "" {
		t.Fatalf("no PDF artifact found: %v", reports)
	}
	assertPDF(t, client, pdfURL)
}

// TestQHPROCommercialAdminRuntime verifies acceptance IAM data through the same
// Gateway contracts consumed by React Admin. Production migrations
// intentionally never seed an operator account.
func TestQHPROCommercialAdminRuntime(t *testing.T) {
	if strings.TrimSpace(os.Getenv(runtimeSmokeEnabled)) != "1" {
		t.Skip(runtimeSmokeEnabled + "=1 is required")
	}

	gatewayURL := strings.TrimRight(envOrDefault("QHPRO_GATEWAY_URL", GatewayURL), "/")
	password := strings.TrimSpace(os.Getenv("QHPRO_ACCEPTANCE_ADMIN_PASSWORD"))
	if password == "" {
		t.Fatal("QHPRO_ACCEPTANCE_ADMIN_PASSWORD is required")
	}
	client := &http.Client{Timeout: 15 * time.Second}
	loginURL := gatewayURL + "/v2/auth/admin/login"
	loginPayload := map[string]any{
		"username":   acceptanceAdminUsername,
		"password":   password,
		"version":    "runtime-smoke",
		"platform":   "WEB",
		"os":         "Linux",
		"deviceName": "QHPRO Runtime Smoke",
	}

	invalidPayload := make(map[string]any, len(loginPayload))
	for key, value := range loginPayload {
		invalidPayload[key] = value
	}
	invalidPayload["password"] = "invalid-acceptance-password"
	invalidStatus, invalidLogin := runtimeJSONPostAny(t, client, loginURL, invalidPayload, "")
	// Legacy shared.ErrorResponse compatibility intentionally uses HTTP 200;
	// requestApi rejects the envelope code. Lock both code and absence of a
	// credential so an invalid login can never be mistaken for authentication.
	if invalidStatus != http.StatusOK || jsonNumber(invalidLogin, "code") != http.StatusUnauthorized || jsonString(invalidLogin, "accessToken", "access_token") != "" {
		t.Fatalf("admin invalid credential boundary: status=%d body=%v", invalidStatus, invalidLogin)
	}

	loginStatus, login := runtimeJSONPostAny(t, client, loginURL, loginPayload, "")
	if loginStatus != http.StatusOK {
		t.Fatalf("admin login: status=%d body=%v", loginStatus, login)
	}
	token := jsonString(login, "accessToken", "access_token")
	if token == "" || jsonNumber(login, "profileId", "profile_id") == 0 || jsonString(login, "role") != "ROLE_ADMIN" {
		t.Fatalf("admin login response is incomplete: %v", login)
	}
	permissionStatus, permissions := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/auth/permission/me", token)
	if permissionStatus != http.StatusOK || len(jsonArray(permissions, "keys")) != 8 {
		t.Fatalf("admin permission snapshot: status=%d body=%v", permissionStatus, permissions)
	}

	checks := []struct {
		name       string
		path       string
		arrayKey   string
		requireRow bool
	}{
		{name: "catalog", path: "/v2/user/admin/catalog/plan-versions?page=1&pageSize=20", arrayKey: "plan_versions", requireRow: true},
		{name: "subscriptions", path: "/v2/user/admin/commercial/subscriptions?page=1&pageSize=20", arrayKey: "subscriptions", requireRow: true},
		{name: "usage", path: "/v2/tqd/admin/commercial/usage?page=1&pageSize=20", arrayKey: "usage", requireRow: true},
		{name: "reports", path: "/v2/tqd/admin/commercial/generated-reports?page=1&pageSize=20", arrayKey: "reports", requireRow: true},
		{name: "orders", path: "/v2/payment/admin/commercial/orders?page=1&pageSize=20", arrayKey: "orders", requireRow: true},
		{name: "fulfillments", path: "/v2/payment/admin/commercial/fulfillments?page=1&pageSize=20", arrayKey: "fulfillments", requireRow: true},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			status, payload := runtimeJSONRequestAny(t, client, gatewayURL+check.path, token)
			if status != http.StatusOK {
				t.Fatalf("status=%d body=%v", status, payload)
			}
			if check.requireRow && len(jsonArray(payload, check.arrayKey)) == 0 {
				t.Fatalf("%s projection is empty: %v", check.arrayKey, payload)
			}
			if check.name == "usage" && strings.TrimSpace(os.Getenv("QHPRO_EXPECT_USAGE_RUNTIME_UNAVAILABLE")) == "1" {
				foundDurableUsage := false
				for _, raw := range jsonArray(payload, check.arrayKey) {
					usage, ok := raw.(map[string]any)
					if !ok || jsonNumber(usage, "durableUsed", "durable_used") <= 0 {
						continue
					}
					foundDurableUsage = true
					if jsonBool(usage, "runtimeAvailable", "runtime_available") || jsonString(usage, "reconciliationHealth", "reconciliation_health", "reconciliation") != "runtime_unavailable" {
						t.Fatalf("usage did not preserve durable truth while Redis was unavailable: %v", usage)
					}
				}
				if !foundDurableUsage {
					t.Fatalf("usage projection contains no durable accounting evidence: %v", payload)
				}
			}
		})
	}

	ordersStatus, orders := runtimeJSONRequestAny(
		t, client, gatewayURL+"/v2/payment/admin/commercial/orders?page=1&pageSize=20", token,
	)
	orderRows := jsonArray(orders, "orders")
	if ordersStatus != http.StatusOK || len(orderRows) == 0 {
		t.Fatalf("payment order evidence: status=%d body=%v", ordersStatus, orders)
	}
	var ownerOrder map[string]any
	var ownerProfileID uint64
	for _, row := range orderRows {
		order, ok := row.(map[string]any)
		if !ok || jsonString(order, "subjectKind", "subject_kind") != "profile" {
			continue
		}
		profileID, err := strconv.ParseUint(jsonString(order, "subjectId", "subject_id"), 10, 64)
		if err != nil || profileID == 0 {
			continue
		}
		ownerOrder = order
		ownerProfileID = profileID
		break
	}
	if ownerOrder == nil {
		t.Fatalf("payment owner progress requires a profile-owned order: %v", orders)
	}
	orderID := jsonNumber(ownerOrder, "orderId", "order_id")
	secret := strings.TrimSpace(os.Getenv("JWT_KEY_GENERATE"))
	if orderID == 0 || secret == "" {
		t.Fatalf("payment progress fixture is incomplete: order=%v jwt_key_present=%t", ownerOrder, secret != "")
	}
	ownerToken := signRuntimeAccessToken(t, secret, ownerProfileID)
	progressStatus, progress := runtimeJSONRequestAny(
		t,
		client,
		fmt.Sprintf("%s/v2/payment/commercial/orders/%d", gatewayURL, orderID),
		ownerToken,
	)
	if progressStatus != http.StatusOK || jsonNumber(progress, "orderId", "order_id") != orderID || jsonString(progress, "progressState", "progress_state") == "" {
		t.Fatalf("payment owner progress: status=%d body=%v", progressStatus, progress)
	}

	projectionStatus, projection := runtimeJSONRequestAny(
		t, client, fmt.Sprintf("%s/v2/user/admin/commercial/users/%d", gatewayURL, ownerProfileID), token,
	)
	if projectionStatus != http.StatusOK || jsonNumber(projection, "profileId", "profile_id") != int64(ownerProfileID) {
		t.Fatalf("user commercial projection: status=%d body=%v", projectionStatus, projection)
	}
}

// TestQHPROCommercialClientLoginRuntime proves that the acceptance identity
// uses the real public password-login transport and receives an ordinary user
// JWT. Production authentication remains OTP/OAuth or an explicitly configured
// password credential created by User Service.
func TestQHPROCommercialClientLoginRuntime(t *testing.T) {
	if strings.TrimSpace(os.Getenv(runtimeSmokeEnabled)) != "1" {
		t.Skip(runtimeSmokeEnabled + "=1 is required")
	}

	gatewayURL := strings.TrimRight(envOrDefault("QHPRO_GATEWAY_URL", GatewayURL), "/")
	password := strings.TrimSpace(os.Getenv("QHPRO_ACCEPTANCE_USER_PASSWORD"))
	if password == "" {
		t.Fatal("QHPRO_ACCEPTANCE_USER_PASSWORD is required")
	}
	client := &http.Client{Timeout: 15 * time.Second}
	loginPayload := map[string]any{
		"username":   acceptanceUserUsername,
		"password":   password,
		"version":    "runtime-smoke",
		"platform":   "WEB",
		"os":         "Linux",
		"deviceName": "QHPRO Client Acceptance",
	}

	loginStatus, login := runtimeJSONPostAny(
		t, client, gatewayURL+"/v2/auth/password/login", loginPayload, "",
	)
	token := jsonString(login, "accessToken", "access_token")
	profileID := jsonNumber(login, "profileId", "profile_id")
	if loginStatus != http.StatusOK || token == "" || profileID == 0 || jsonString(login, "role") == "ROLE_ADMIN" {
		t.Fatalf("acceptance client login: status=%d body=%v", loginStatus, login)
	}

	profileStatus, profile := runtimeJSONRequestAny(
		t, client, gatewayURL+"/v2/user/commercial/profile", token,
	)
	if profileStatus != http.StatusOK || jsonNumber(profile, "profileId", "profile_id") != profileID {
		t.Fatalf("acceptance client projection: status=%d body=%v", profileStatus, profile)
	}
}

func runtimeJSONRequest(t *testing.T, client *http.Client, url, token string) (int, map[string]any) {
	t.Helper()
	status, body := runtimeJSONRequestAny(t, client, url, token)
	if token != "" && status != http.StatusOK {
		t.Fatalf("GET %s: status=%d body=%v", url, status, body)
	}
	return status, body
}

func runtimeJSONRequestAny(t *testing.T, client *http.Client, url, token string) (int, map[string]any) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&body); err != nil {
		t.Fatalf("GET %s: decode status %d: %v", url, resp.StatusCode, err)
	}
	payload := body
	if data, ok := body["data"].(map[string]any); ok {
		payload = data
	}
	return resp.StatusCode, payload
}

func runtimeJSONPostAny(t *testing.T, client *http.Client, url string, value map[string]any, token string) (int, map[string]any) {
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
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&body); err != nil {
		t.Fatalf("POST %s: decode status %d: %v", url, resp.StatusCode, err)
	}
	payload := body
	if data, ok := body["data"].(map[string]any); ok {
		payload = data
	}
	return resp.StatusCode, payload
}

func signRuntimeAccessToken(t *testing.T, secret string, profileID uint64) string {
	t.Helper()
	now := time.Now()
	header := map[string]any{"alg": "HS256", "typ": "JWT"}
	claims := map[string]any{
		"sub": profileID, "profile": profileID, "origin": profileID,
		"type": "ACCESS", "role": "ROLE_USER", "roleIds": []uint64{},
		"session": 1, "iat": now.Unix(), "exp": now.Add(10 * time.Minute).Unix(),
	}
	encodedHeader := encodeJWTPart(t, header)
	encodedClaims := encodeJWTPart(t, claims)
	unsigned := encodedHeader + "." + encodedClaims
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encodeJWTPart(t *testing.T, value any) string {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(raw)
}

func assertPDF(t *testing.T, client *http.Client, url string) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Range", "bytes=0-7")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("download PDF: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("download PDF: status=%d", resp.StatusCode)
	}
	header := make([]byte, 8)
	if _, err := io.ReadFull(resp.Body, header); err != nil {
		t.Fatalf("read PDF header: %v", err)
	}
	// PDF producers legitimately emit different specification revisions
	// (for example 1.4 or 1.7). The acceptance contract is a real PDF
	// artifact, not a renderer-specific version string.
	if !bytes.HasPrefix(header, []byte("%PDF-")) {
		t.Fatalf("PDF header=%q", header)
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func jsonArray(value map[string]any, key string) []any {
	items, _ := value[key].([]any)
	return items
}

func jsonNumber(value map[string]any, keys ...string) int64 {
	for _, key := range keys {
		if number, ok := value[key].(float64); ok {
			return int64(number)
		}
	}
	return 0
}

func jsonBool(value map[string]any, keys ...string) bool {
	for _, key := range keys {
		if result, ok := value[key].(bool); ok {
			return result
		}
	}
	return false
}

func jsonString(value map[string]any, keys ...string) string {
	for _, key := range keys {
		if result, ok := value[key].(string); ok {
			return result
		}
	}
	return ""
}

func Example_runtimeSmoke() {
	fmt.Println("QHPRO_RUNTIME_SMOKE=1 JWT_KEY_GENERATE=<local-key> go test -run TestQHPROCommercialRuntime")
	// Output: QHPRO_RUNTIME_SMOKE=1 JWT_KEY_GENERATE=<local-key> go test -run TestQHPROCommercialRuntime
}
