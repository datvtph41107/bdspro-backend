package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// This proof uses real HTTP owners and worker/File output. Only login
// credentials are provisioned by the development fixture; no commercial or
// GIS business rows are fabricated. Bank events are explicitly simulated.
func TestQHPROSubscriptionLifecycleRuntime(t *testing.T) {
	if os.Getenv(runtimeSmokeEnabled) != "1" {
		t.Skip("local runtime required")
	}
	gateway := strings.TrimRight(envOrDefault("QHPRO_GATEWAY_URL", GatewayURL), "/")
	client := &http.Client{Timeout: 20 * time.Second}
	key := strings.TrimSpace(os.Getenv("SEPAY_API_KEY"))
	if key == "" {
		t.Fatal("SEPAY_API_KEY required")
	}
	statusCode, login := runtimeJSONPostWithHeaders(t, client, gateway+"/v2/auth/password/login", map[string]any{
		"username": os.Getenv("QHPRO_ACCEPTANCE_USERNAME"), "password": os.Getenv("QHPRO_ACCEPTANCE_USER_PASSWORD"),
		"version": "lifecycle-acceptance", "platform": "WEB", "os": "Linux", "deviceName": "Local lifecycle acceptance",
	}, nil)
	token := jsonString(login, "accessToken", "access_token")
	if statusCode != 200 || token == "" {
		t.Fatal("local acceptance login failed")
	}
	auth := map[string]string{"Authorization": "Bearer " + token}
	runKey := "lifecycle-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	proof := map[string]any{"bankEvidence": "authenticated local webhook simulation", "profileId": jsonNumber(login, "profileId", "profile_id"), "username": os.Getenv("QHPRO_ACCEPTANCE_USERNAME")}
	defer func() {
		if dir := os.Getenv("QHPRO_ACCEPTANCE_OUTPUT_DIR"); dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Error(err)
				return
			}
			proof["passed"] = !t.Failed()
			b, _ := json.MarshalIndent(proof, "", "  ")
			if err := os.WriteFile(filepath.Join(dir, "subscription-lifecycle.json"), b, 0644); err != nil {
				t.Error(err)
			}
		}
	}()
	base, amount, retireBase := publishRuntimePlan(t, client, gateway, 100, 2, true)
	defer retireBase()
	upgrade, upgradeAmount, retireUpgrade := publishRuntimePlan(t, client, gateway, 200, 4, true)
	defer retireUpgrade()
	defer runtimeJSONMethod(client, http.MethodPost, gateway+"/v2/user/commercial/subscription:cancel", map[string]any{"productCode": "qhpro", "reason": "local acceptance cleanup"}, token)
	proof["plans"] = map[string]any{"base": base, "upgrade": upgrade}
	get := func(path string) map[string]any {
		s, b := runtimeJSONRequestAny(t, client, gateway+path, token)
		if s != 200 {
			t.Fatalf("GET %s: status=%d body=%v", path, s, b)
		}
		return b
	}
	post := func(path string, body map[string]any, command string) (int, map[string]any) {
		h := cloneHeaders(auth)
		if command != "" {
			h["Idempotency-Key"] = runKey + "-" + command
		}
		return runtimeJSONPostWithHeaders(t, client, gateway+path, body, h)
	}
	reportBody := map[string]any{"reportType": 2, "entityType": 2, "title": "PDF nghiệm thu vòng đời gói",
		"location": map[string]any{"address": "Hà Nội — báo cáo vị trí local"},
		"spatial":  map[string]any{"centroid": map[string]any{"lat": 21.03145, "lon": 105.84123}}}
	assertDeniedReport := func(command string) {
		s, b := post("/v2/tqd/map-workspace/reports", reportBody, command)
		if s != http.StatusForbidden {
			t.Fatalf("new report without access: status=%d body=%v", s, b)
		}
	}
	if jsonBool(get("/v2/user/commercial/profile"), "hasSubscription", "has_subscription") {
		t.Fatal("acceptance user must begin free")
	}
	assertDeniedReport("free-report")
	proof["freeReportDenied"] = true
	checkout := func(plan, command string) map[string]any {
		s, b := post("/v2/user/checkout", map[string]any{"planCode": plan, "subjectKind": "profile"}, command)
		if s != 200 || jsonNumber(b, "orderId", "order_id") == 0 {
			t.Fatalf("checkout: status=%d body=%v", s, b)
		}
		return b
	}
	order := checkout(base, "base-checkout")
	replay := checkout(base, "base-checkout")
	if jsonNumber(order, "orderId", "order_id") != jsonNumber(replay, "orderId", "order_id") {
		t.Fatal("checkout replay duplicated order")
	}
	proof["checkoutReplaySameOrder"] = true
	createAttempt := func(order map[string]any, command string) {
		s, b := post("/v2/user/checkout/payment-attempts", map[string]any{"orderId": jsonNumber(order, "orderId", "order_id"), "subjectKind": "profile", "method": "bank_transfer"}, command)
		if s != 200 || !strings.HasPrefix(jsonString(jsonObject(b, "nextAction"), "qrPayload", "qr_payload"), "000201") {
			// Protobuf JSON may use snake_case on the public Gateway.
			if s != 200 || !strings.HasPrefix(jsonString(jsonObject(b, "next_action"), "qrPayload", "qr_payload"), "000201") {
				t.Fatalf("QR attempt: status=%d body=%v", s, b)
			}
		}
	}
	createAttempt(order, "base-attempt")
	assertDeniedReport("unpaid-report")
	proof["unpaidReportDenied"] = true
	webhookBody := func(order map[string]any, value int64) map[string]any {
		return map[string]any{"gateway": "LOCAL_LIFECYCLE_TEST", "transactionDate": time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05"),
			"accountNumber": "LOCAL_TEST", "content": jsonString(order, "reference"), "transferType": "in", "transferAmount": value,
			"description": "Local authenticated bank simulation", "referenceCode": runKey, "id": time.Now().UnixNano()%4_000_000_000 + 1}
	}
	s, _ := runtimeJSONPostWithHeaders(t, client, gateway+"/v2/payment/sepay/webhook", webhookBody(order, amount), map[string]string{"Authorization": "Apikey invalid"})
	if s != 401 {
		t.Fatalf("invalid webhook credential status=%d", s)
	}
	proof["invalidWebhookDenied"] = true
	wrong := webhookBody(order, amount-1)
	s, _ = runtimeJSONPostWithHeaders(t, client, gateway+"/v2/payment/sepay/webhook", wrong, map[string]string{"Authorization": "Apikey " + key})
	if s != 200 || jsonBool(get("/v2/user/commercial/profile"), "hasSubscription", "has_subscription") {
		t.Fatal("wrong amount activated subscription")
	}
	proof["wrongAmountDidNotActivate"] = true
	settle := func(order map[string]any, value int64) {
		body := webhookBody(order, value)
		for range 2 {
			s, b := runtimeJSONPostWithHeaders(t, client, gateway+"/v2/payment/sepay/webhook", body, map[string]string{"Authorization": "Apikey " + key})
			if s != 200 {
				t.Fatalf("settlement: status=%d body=%v", s, b)
			}
		}
		waitFor(t, 45*time.Second, "active payment", func() (bool, string) {
			b := get(fmt.Sprintf("/v2/payment/commercial/orders/%d", jsonNumber(order, "orderId", "order_id")))
			return jsonString(b, "progressState", "progress_state") == "ACTIVE", fmt.Sprint(b)
		})
	}
	settle(order, amount)
	baseProfile := get("/v2/user/commercial/profile")
	baseSubscription := jsonMap(baseProfile, "effectiveSubscription", "effective_subscription")
	proof["activated"] = baseSubscription
	usage := func() map[string]any { return get("/v2/tqd/commercial/generated-report-usage") }
	if jsonNumber(usage(), "limit") != 2 {
		t.Fatal("base quota limit differs from catalog")
	}
	badBody := map[string]any{"reportType": 999, "entityType": 2}
	s, _ = post("/v2/tqd/map-workspace/reports", badBody, "invalid-report")
	if s != 400 || jsonNumber(usage(), "durableUsed", "durable_used") != 0 {
		t.Fatal("invalid report spent quota")
	}
	proof["invalidReportNoCharge"] = true
	createReport := func(command string) map[string]any {
		s, b := post("/v2/tqd/map-workspace/reports", reportBody, command)
		if s != 200 {
			t.Fatalf("create report: status=%d body=%v", s, b)
		}
		return jsonObject(b, "item")
	}
	first := createReport("report-1")
	firstReplay := createReport("report-1")
	reportID := jsonNumber(first, "reportId", "report_id", "id")
	if reportID == 0 || reportID != jsonNumber(firstReplay, "reportId", "report_id", "id") {
		t.Fatal("report replay duplicated output")
	}
	createReport("report-2")
	s, _ = post("/v2/tqd/map-workspace/reports", reportBody, "report-over-limit")
	if s != 429 || jsonNumber(usage(), "durableUsed", "durable_used") != 2 {
		t.Fatal("exhausted quota admitted another report")
	}
	proof["exhaustion"] = usage()
	reportPath := fmt.Sprintf("/v2/tqd/map-workspace/reports/%d", reportID)
	waitPDF := func() string {
		url := ""
		waitFor(t, 45*time.Second, "PDF ready", func() (bool, string) {
			b := get(reportPath)
			url = jsonString(jsonObject(b, "assets"), "pdfUrl", "pdf_url")
			return jsonNumber(b, "status") == 10 && url != "", fmt.Sprint(b)
		})
		assertPDF(t, client, url)
		return url
	}
	url := waitPDF()
	proof["reportId"], proof["pdfUrl"] = reportID, url
	fileResponse, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	fileBytes, err := io.ReadAll(io.LimitReader(fileResponse.Body, 32<<20))
	fileResponse.Body.Close()
	if err != nil || len(fileBytes) < 1000 {
		t.Fatal("PDF artifact incomplete")
	}
	fileHash := sha256.Sum256(fileBytes)
	proof["pdfSHA256"] = hex.EncodeToString(fileHash[:])
	if dir := os.Getenv("QHPRO_ACCEPTANCE_OUTPUT_DIR"); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "subscription-lifecycle-report.pdf"), fileBytes, 0644); err != nil {
			t.Fatal(err)
		}
	}
	for range 2 {
		s, b := post(reportPath+":regenerate", map[string]any{"reportId": reportID}, "")
		if s != 200 {
			t.Fatalf("regenerate: %d %v", s, b)
		}
	}
	if waitPDF() != url {
		t.Fatal("regenerate replaced immutable output authority")
	}
	if jsonNumber(usage(), "durableUsed", "durable_used") != 2 {
		t.Fatal("regenerate charged again")
	}
	proof["regenerateSameReportNoCharge"] = true
	s, share := post(reportPath+":share", map[string]any{"reportId": reportID}, "")
	shareURL := jsonString(share, "shareUrl", "share_url")
	if s != 200 || !strings.HasPrefix(shareURL, "http") {
		t.Fatalf("share missing real output: %d %v", s, share)
	}
	assertPDF(t, client, shareURL)
	proof["shareReturnsRealPDF"] = true
	// Real local GIS may be empty. These failure cases must still leave usage at zero.
	readStatus, err := runtimeJSONMethod(client, http.MethodGet, gateway+"/v2/tqd/parcels/9000000000/detail", nil, token)
	if err != nil {
		t.Fatal(err)
	}
	proof["dossierWithoutCommandStatus"] = readStatus
	if readStatus != 400 {
		t.Fatalf("missing metering key status=%d", readStatus)
	}
	request, _ := http.NewRequest(http.MethodGet, gateway+"/v2/tqd/parcels/9000000000/detail", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Idempotency-Key", runKey+"-missing-parcel")
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 404 {
		t.Fatalf("nonexistent parcel status=%d", response.StatusCode)
	}
	proof["missingParcelNoResult"] = true
	s, _ = post("/v2/tqd/parcels/by-polygon", map[string]any{"coordinates": []any{}, "timestamp": 1}, "invalid-area")
	if s != 400 {
		t.Fatalf("invalid area status=%d", s)
	}
	// The same stale timestamp with a new intent must not bypass validation/charge.
	s, _ = post("/v2/tqd/parcels/by-polygon", map[string]any{"coordinates": []any{}, "timestamp": 1}, "invalid-area-2")
	if s != 400 {
		t.Fatalf("stale timestamp skipped area validation: status=%d", s)
	}
	proof["invalidAreaNoCharge"] = true
	for _, raw := range jsonArray(get("/v2/tqd/commercial/usage"), "usage") {
		row, _ := raw.(map[string]any)
		if jsonString(row, "operation") != "workspace.report.generate" && jsonNumber(row, "durableUsed", "durable_used") != 0 {
			t.Fatalf("failed GIS action charged: %v", row)
		}
	}
	for _, plan := range []string{base} {
		s, _ = post("/v2/user/checkout", map[string]any{"planCode": plan, "subjectKind": "profile"}, "same-plan")
		if s != http.StatusPreconditionFailed {
			t.Fatalf("same plan checkout status=%d", s)
		}
	}
	proof["sameTierCheckoutStatus"] = s
	upOrder := checkout(upgrade, "upgrade-checkout")
	createAttempt(upOrder, "upgrade-attempt")
	if jsonNumber(usage(), "limit") != 2 {
		t.Fatal("unpaid upgrade changed quota")
	}
	settle(upOrder, upgradeAmount)
	upSubscription := jsonMap(get("/v2/user/commercial/profile"), "effectiveSubscription", "effective_subscription")
	if jsonNumber(upSubscription, "subscriptionId", "subscription_id") != jsonNumber(baseSubscription, "subscriptionId", "subscription_id") ||
		fmt.Sprint(jsonMap(upSubscription, "currentPeriodEnd", "current_period_end")) != fmt.Sprint(jsonMap(baseSubscription, "currentPeriodEnd", "current_period_end")) {
		t.Fatal("upgrade reset subscription period")
	}
	upUsage := usage()
	if jsonNumber(upUsage, "limit") != 4 || jsonNumber(upUsage, "durableUsed", "durable_used") != 2 || jsonNumber(upUsage, "remaining") != 2 {
		t.Fatalf("upgrade lost usage: %v", upUsage)
	}
	proof["upgradePreservedUsageAndPeriod"] = upUsage
	s, _ = post("/v2/user/checkout", map[string]any{"planCode": base, "subjectKind": "profile"}, "downgrade")
	if s != http.StatusPreconditionFailed {
		t.Fatalf("downgrade checkout status=%d", s)
	}
	proof["downgradeCheckoutStatus"] = s
	proof["downgradeDeniedBeforePayment"] = true
	createReport("report-after-upgrade")
	if jsonNumber(usage(), "durableUsed", "durable_used") != 3 {
		t.Fatal("upgraded quota did not account report")
	}
	s, canceled := post("/v2/user/commercial/subscription:cancel", map[string]any{"productCode": "qhpro", "reason": "local lifecycle acceptance"}, "")
	if s != 200 || !jsonBool(canceled, "changed") {
		t.Fatalf("cancel: %d %v", s, canceled)
	}
	assertDeniedReport("after-cancel")
	if jsonBool(get("/v2/user/commercial/profile"), "hasSubscription", "has_subscription") {
		t.Fatal("canceled access remained active")
	}
	if jsonNumber(createReport("report-1"), "reportId", "report_id", "id") != reportID {
		t.Fatal("accepted report lost after cancel")
	}
	assertPDF(t, client, url)
	s, cancelReplay := post("/v2/user/commercial/subscription:cancel", map[string]any{"productCode": "qhpro", "reason": "local replay"}, "")
	if s != 200 || jsonBool(cancelReplay, "changed") {
		t.Fatal("cancel replay changed state")
	}
	proof["cancelImmediateReplaySafeHistoryReadable"] = true
	newOrder := checkout(base, "repurchase")
	createAttempt(newOrder, "repurchase-attempt")
	settle(newOrder, amount)
	newSubscription := jsonMap(get("/v2/user/commercial/profile"), "effectiveSubscription", "effective_subscription")
	if jsonNumber(newSubscription, "subscriptionId", "subscription_id") == jsonNumber(baseSubscription, "subscriptionId", "subscription_id") || jsonNumber(usage(), "durableUsed", "durable_used") != 0 {
		t.Fatal("repurchase reused old quota period")
	}
	createReport("new-period-report")
	proof["repurchaseNewPeriod"] = usage()
	t.Logf("PASS lifecycle: paid/replay/wrong-amount/QR/exhaustion/regenerate/share/upgrade/cancel/repurchase; profile=%d report=%d", jsonNumber(login, "profileId", "profile_id"), reportID)
}
