package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	acceptanceUserUsername  = "qhpro_acceptance_0000000000"
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
	if permissionStatus != http.StatusOK || len(jsonArray(permissions, "keys")) < 10 {
		t.Fatalf("admin permission snapshot: status=%d body=%v", permissionStatus, permissions)
	}

	t.Run("IAM role permission matrix", func(t *testing.T) {
		statusCode, roles := runtimeJSONRequestAny(
			t,
			client,
			gatewayURL+"/v2/auth/role/list?page=0&size=100",
			token,
		)
		if statusCode != http.StatusOK {
			t.Fatalf("list IAM roles: status=%d body=%v", statusCode, roles)
		}
		var acceptanceRole map[string]any
		for _, raw := range jsonArray(roles, "data") {
			role, ok := raw.(map[string]any)
			if ok && jsonString(role, "key") == "QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN" {
				acceptanceRole = role
				break
			}
		}
		if acceptanceRole == nil {
			t.Fatalf("acceptance IAM role is absent: %v", roles)
		}
		roleID := jsonNumber(acceptanceRole, "id")
		permissionIDs := jsonArray(acceptanceRole, "permissionIds")
		if len(permissionIDs) == 0 {
			permissionIDs = jsonArray(acceptanceRole, "permission_ids")
		}
		if roleID == 0 || len(permissionIDs) < 15 {
			t.Fatalf("acceptance IAM role is incomplete: %v", acceptanceRole)
		}
		statusCode, saved := runtimeJSONPostAny(
			t,
			client,
			gatewayURL+"/v2/auth/role/permissions",
			map[string]any{"roleId": roleID, "permissionIds": permissionIDs},
			token,
		)
		if statusCode != http.StatusOK || jsonNumber(saved, "id") != roleID {
			t.Fatalf("replace IAM role permissions: role_id=%d permission_count=%d status=%d body=%v role=%v", roleID, len(permissionIDs), statusCode, saved, acceptanceRole)
		}

		groupListStatus, groups := runtimeJSONRequestAny(
			t, client, gatewayURL+"/v2/auth/role-group/list?page=0&size=1000", token,
		)
		if groupListStatus != http.StatusOK {
			t.Fatalf("list IAM role groups: status=%d body=%v", groupListStatus, groups)
		}
		// Deleted groups are soft-deleted and therefore still hold their unique
		// key in PostgreSQL. Start from wall-clock seconds so repeated acceptance
		// runs do not reuse a key hidden from the normal list query.
		groupKeyValue := time.Now().Unix()
		for _, raw := range jsonArray(groups, "data") {
			group, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if key := jsonNumber(group, "key"); key >= groupKeyValue {
				groupKeyValue = key + 1
			}
		}
		if groupKeyValue > int64(^uint32(0)) {
			t.Fatalf("no role-group key remains in uint32 range")
		}
		groupKey := uint32(groupKeyValue)
		groupName := fmt.Sprintf("NT %d", groupKey)
		statusCode, createdGroup := runtimeJSONPostAny(
			t, client, gatewayURL+"/v2/auth/role-group",
			map[string]any{"groupName": groupName, "description": "Nhom quyen runtime", "key": groupKey}, token,
		)
		groupID := jsonNumber(createdGroup, "id")
		if statusCode != http.StatusOK || groupID == 0 {
			t.Fatalf("create IAM role group: status=%d body=%v", statusCode, createdGroup)
		}
		defer func() {
			if groupID != 0 {
				_, _ = runtimeJSONMethod(client, http.MethodDelete, fmt.Sprintf("%s/v2/auth/role-group/%d", gatewayURL, groupID), nil, token)
			}
		}()

		updatedName := fmt.Sprintf("Sua %d", groupKey)
		statusCode, updatedGroup := runtimeJSONMethodAny(
			t, client, http.MethodPut, fmt.Sprintf("%s/v2/auth/role-group/%d", gatewayURL, groupID),
			map[string]any{"id": groupID, "groupName": updatedName, "description": "Da kiem chung", "key": groupKey}, token,
		)
		if statusCode != http.StatusOK || jsonString(updatedGroup, "groupName", "group_name") != updatedName {
			t.Fatalf("update IAM role group: status=%d body=%v", statusCode, updatedGroup)
		}

		statusCode, savedGroupPermissions := runtimeJSONMethodAny(
			t, client, http.MethodPut, gatewayURL+"/v2/auth/role-group/permissions",
			map[string]any{"groupId": groupID, "permissionIds": permissionIDs[:2]}, token,
		)
		if statusCode != http.StatusOK || jsonNumber(savedGroupPermissions, "id") != groupID {
			t.Fatalf("replace IAM role group permissions: status=%d body=%v", statusCode, savedGroupPermissions)
		}
		statusCode, groupPermissions := runtimeJSONRequestAny(
			t, client, fmt.Sprintf("%s/v2/auth/role-group/%d/permissions", gatewayURL, groupKey), token,
		)
		if statusCode != http.StatusOK || len(jsonArray(groupPermissions, "data")) != 2 {
			t.Fatalf("read IAM role group permissions: status=%d body=%v", statusCode, groupPermissions)
		}
		statusCode, deletedGroup := runtimeJSONMethodAny(
			t, client, http.MethodDelete, fmt.Sprintf("%s/v2/auth/role-group/%d", gatewayURL, groupID), nil, token,
		)
		if statusCode != http.StatusOK || jsonNumber(deletedGroup, "id") != groupID {
			t.Fatalf("delete IAM role group: status=%d body=%v", statusCode, deletedGroup)
		}
		groupID = 0
	})

	if strings.TrimSpace(os.Getenv("QHPRO_CATALOG_LIFECYCLE_SMOKE")) == "1" {
		t.Run("catalog draft publish retire lifecycle", func(t *testing.T) {
			runtimeCatalogLifecycle(t, client, gatewayURL, token)
		})
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
		{name: "users", path: "/v2/user/admin/users?page=1&size=20", arrayKey: "data", requireRow: true},
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
	userAdminStatus, userAdminBody := runtimeJSONRequestAny(
		t, client, gatewayURL+"/v2/user/admin/users?page=1&size=1", ownerToken,
	)
	if userAdminStatus != http.StatusForbidden {
		t.Fatalf("ordinary user reached User Admin boundary: status=%d body=%v", userAdminStatus, userAdminBody)
	}
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
	userDetailStatus, userDetail := runtimeJSONRequestAny(
		t, client, fmt.Sprintf("%s/v2/user/admin/users/%d", gatewayURL, ownerProfileID), token,
	)
	if userDetailStatus != http.StatusOK || jsonNumber(userDetail, "profileId", "profile_id") != int64(ownerProfileID) {
		t.Fatalf("user admin detail: status=%d body=%v", userDetailStatus, userDetail)
	}

	t.Run("user lifecycle", func(t *testing.T) {
		suffix := strconv.FormatInt(time.Now().UnixNano(), 36)
		username := "qhpro.admin.created." + suffix
		password := "Qhpro-Admin-Created-" + suffix
		createStatus, created := runtimeJSONMethodAny(t, client, http.MethodPost,
			gatewayURL+"/v2/user/admin/user/new", map[string]any{
				"username": username,
				"password": password,
				"fullName": "QHPRO Admin-created User",
				"email":    username + "@example.invalid",
			}, token)
		profileID := jsonNumber(created, "profileId", "profile_id")
		if createStatus != http.StatusOK || profileID == 0 || jsonNumber(created, "authId", "auth_id") == 0 {
			t.Fatalf("create user: status=%d body=%v", createStatus, created)
		}

		loginStatus, createdLogin := runtimeJSONPostAny(t, client, gatewayURL+"/v2/auth/password/login", map[string]any{
			"username": username, "password": password, "version": "runtime-smoke",
			"platform": "WEB", "os": "Linux", "deviceName": "Admin User Lifecycle",
		}, "")
		if loginStatus != http.StatusOK || jsonString(createdLogin, "accessToken", "access_token") == "" || jsonNumber(createdLogin, "profileId", "profile_id") != profileID {
			t.Fatalf("created user credential: status=%d body=%v", loginStatus, createdLogin)
		}

		updateStatus, updated := runtimeJSONMethodAny(t, client, http.MethodPut,
			fmt.Sprintf("%s/v2/user/admin/users/%d", gatewayURL, profileID), map[string]any{
				"profileId": profileID,
				"fullName":  "QHPRO Admin-updated User",
				"email":     username + "@example.invalid",
				"status":    10,
			}, token)
		if updateStatus != http.StatusOK || jsonString(updated, "fullName", "full_name") != "QHPRO Admin-updated User" {
			t.Fatalf("update user: status=%d body=%v", updateStatus, updated)
		}

		deleteStatus, deleted := runtimeJSONMethodAny(t, client, http.MethodDelete,
			fmt.Sprintf("%s/v2/user/admin/users/%d?reason=acceptance-cleanup", gatewayURL, profileID), nil, token)
		if deleteStatus != http.StatusOK {
			t.Fatalf("delete user: status=%d body=%v", deleteStatus, deleted)
		}
		detailStatus, detail := runtimeJSONRequestAny(t, client,
			fmt.Sprintf("%s/v2/user/admin/users/%d", gatewayURL, profileID), token)
		if detailStatus == http.StatusOK {
			t.Fatalf("deleted user remained visible: body=%v", detail)
		}
	})
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
	username := envOrDefault("QHPRO_ACCEPTANCE_USERNAME", acceptanceUserUsername)
	loginPayload := map[string]any{
		"username":   username,
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
	return runtimeJSONMethodAny(t, client, http.MethodPost, url, value, token)
}

func runtimeJSONMethodAny(t *testing.T, client *http.Client, method, url string, value map[string]any, token string) (int, map[string]any) {
	t.Helper()
	var requestBody io.Reader
	if value != nil {
		raw, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		requestBody = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(t.Context(), method, url, requestBody)
	if err != nil {
		t.Fatal(err)
	}
	if value != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer resp.Body.Close()

	var responseBody map[string]any
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&responseBody); err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("%s %s: decode status %d: %v", method, url, resp.StatusCode, err)
	}
	if responseBody == nil {
		responseBody = map[string]any{}
	}
	payload := responseBody
	if data, ok := responseBody["data"].(map[string]any); ok {
		payload = data
	}
	return resp.StatusCode, payload
}

// runtimeCatalogLifecycle proves the exact browser contract: draft is Admin
// only, publication makes the version discoverable by clients, and retirement
// removes it from new sales without deleting its durable record.
func runtimeCatalogLifecycle(t *testing.T, client *http.Client, gatewayURL, token string) {
	t.Helper()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)
	planCode := "qhpro.acceptance." + suffix
	priceCode := planCode + ".monthly"
	draft := map[string]any{
		"productCode": "qhpro", "productDisplayName": "QHPro",
		"planCode": planCode, "version": "1.0.0", "displayName": "Gói nghiệm thu " + suffix,
		"subjectScope": "profile", "tierRank": 900, "subscriptionTermDays": 30,
		"entitlements": []map[string]any{
			{"code": "workspace.report.feature", "kind": "feature_access", "featureCode": "workspace.report.generate", "amount": 0, "unlimited": false, "period": "none"},
			{"code": "workspace.report.allowance", "kind": "usage_allowance", "meterCode": "workspace.report_generation.accepted", "amount": 7, "unlimited": false, "period": "subscription_cycle"},
		},
		"operations": []map[string]any{{"code": "workspace.report.generate", "featureCode": "workspace.report.generate", "meterCode": "workspace.report_generation.accepted", "unitsPerAction": 1}},
		"prices":     []map[string]any{{"code": priceCode, "kind": "recurring", "currency": "VND", "amountMinor": 125000, "billingUnit": "subscription_cycle", "quantity": 1}},
	}
	statusCode, created := runtimeJSONMethodAny(t, client, http.MethodPost, gatewayURL+"/v2/user/admin/catalog/plan-versions", map[string]any{"draft": draft}, token)
	planVersionID := jsonNumber(created, "id")
	if statusCode != http.StatusOK || planVersionID == 0 || jsonString(created, "status") != "draft" {
		t.Fatalf("create catalog draft: status=%d body=%v", statusCode, created)
	}
	completed := false
	defer func() {
		if completed {
			return
		}
		// Best-effort cleanup: delete if it is still a draft, otherwise retire
		// it before leaving the failed acceptance run.
		_, _ = runtimeJSONMethod(client, http.MethodPost, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d:retire", gatewayURL, planVersionID), map[string]any{"planVersionId": planVersionID}, token)
		_, _ = runtimeJSONMethod(client, http.MethodDelete, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d", gatewayURL, planVersionID), nil, token)
	}()
	assertPublicPlanVisibility(t, client, gatewayURL, planCode, false, 0, 0)

	// Stable product/plan metadata comes back from the owner. A draft version
	// may only replace its own terms, never mutate that shared parent state.
	draft["productDisplayName"] = jsonString(created, "productDisplayName", "product_display_name")
	draft["tierRank"] = jsonNumber(created, "tierRank", "tier_rank")
	draft["displayName"] = "Gói nghiệm thu đã cập nhật " + suffix
	draft["entitlements"].([]map[string]any)[1]["amount"] = 12
	draft["prices"].([]map[string]any)[0]["amountMinor"] = 149000
	statusCode, updated := runtimeJSONMethodAny(t, client, http.MethodPut, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d", gatewayURL, planVersionID), map[string]any{"planVersionId": planVersionID, "draft": draft}, token)
	if statusCode != http.StatusOK || jsonString(updated, "displayName", "display_name") != draft["displayName"] {
		t.Fatalf("update catalog draft: status=%d body=%v", statusCode, updated)
	}
	statusCode, validated := runtimeJSONMethodAny(t, client, http.MethodPost, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d:validate", gatewayURL, planVersionID), map[string]any{"planVersionId": planVersionID}, token)
	if statusCode != http.StatusOK || !jsonBool(validated, "valid") || len(jsonString(validated, "termsChecksum", "terms_checksum")) != 64 {
		t.Fatalf("validate catalog draft: status=%d body=%v", statusCode, validated)
	}

	effectiveFrom := time.Now().UTC().Add(2 * time.Second)
	statusCode, published := runtimeJSONMethodAny(t, client, http.MethodPost, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d:publish", gatewayURL, planVersionID), map[string]any{"planVersionId": planVersionID, "effectiveFrom": effectiveFrom.Format(time.RFC3339Nano)}, token)
	publishedVersion := jsonMap(published, "planVersion", "plan_version")
	if statusCode != http.StatusOK || jsonString(publishedVersion, "status") != "active" {
		t.Fatalf("publish catalog version: status=%d body=%v", statusCode, published)
	}
	time.Sleep(time.Until(effectiveFrom) + 150*time.Millisecond)
	assertPublicPlanVisibility(t, client, gatewayURL, planCode, true, 149000, 12)

	statusCode, retired := runtimeJSONMethodAny(t, client, http.MethodPost, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d:retire", gatewayURL, planVersionID), map[string]any{"planVersionId": planVersionID}, token)
	if statusCode != http.StatusOK || jsonString(retired, "status") != "retired" {
		t.Fatalf("retire catalog version: status=%d body=%v", statusCode, retired)
	}
	assertPublicPlanVisibility(t, client, gatewayURL, planCode, false, 0, 0)

	// Deletion is deliberately a different rule: only a never-published draft
	// may be removed. This second version proves that bounded cleanup path.
	draft["version"] = "1.0.1"
	draft["displayName"] = "Bản nháp có thể xóa " + suffix
	statusCode, removable := runtimeJSONMethodAny(t, client, http.MethodPost, gatewayURL+"/v2/user/admin/catalog/plan-versions", map[string]any{"draft": draft}, token)
	removableID := jsonNumber(removable, "id")
	if statusCode != http.StatusOK || removableID == 0 {
		t.Fatalf("create removable catalog draft: status=%d body=%v", statusCode, removable)
	}
	statusCode, deleted := runtimeJSONMethodAny(t, client, http.MethodDelete, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d", gatewayURL, removableID), nil, token)
	if statusCode != http.StatusOK {
		t.Fatalf("delete catalog draft: status=%d body=%v", statusCode, deleted)
	}
	statusCode, deleted = runtimeJSONRequestAny(t, client, fmt.Sprintf("%s/v2/user/admin/catalog/plan-versions/%d", gatewayURL, removableID), token)
	if statusCode == http.StatusOK {
		t.Fatalf("deleted catalog draft remained visible: %v", deleted)
	}
	completed = true
}

func assertPublicPlanVisibility(t *testing.T, client *http.Client, gatewayURL, planCode string, wantVisible bool, wantPrice, wantQuota int64) {
	t.Helper()
	statusCode, payload := runtimeJSONRequestAny(t, client, gatewayURL+"/v2/user/commercial/plans?productCode=qhpro&subjectKind=profile", "")
	if statusCode != http.StatusOK {
		t.Fatalf("public catalog: status=%d body=%v", statusCode, payload)
	}
	for _, raw := range jsonArray(payload, "plans") {
		item, ok := raw.(map[string]any)
		if !ok || jsonString(item, "planCode", "plan_code") != planCode {
			continue
		}
		if !wantVisible {
			t.Fatalf("plan %s unexpectedly visible to client: %v", planCode, item)
		}
		prices := jsonArray(item, "prices")
		entitlements := jsonArray(item, "entitlements")
		if len(prices) == 0 || jsonNumber(prices[0].(map[string]any), "amountMinor", "amount_minor") != wantPrice {
			t.Fatalf("published price mismatch: %v", item)
		}
		for _, entitlementRaw := range entitlements {
			entitlement, ok := entitlementRaw.(map[string]any)
			if ok && jsonString(entitlement, "meterCode", "meter_code") == "workspace.report_generation.accepted" && jsonNumber(entitlement, "amount") == wantQuota {
				return
			}
		}
		t.Fatalf("published quota mismatch: %v", item)
	}
	if wantVisible {
		t.Fatalf("plan %s is not visible in client catalog: %v", planCode, payload)
	}
}

func runtimeJSONMethod(client *http.Client, method, url string, value map[string]any, token string) (int, error) {
	var requestBody io.Reader
	if value != nil {
		raw, err := json.Marshal(value)
		if err != nil {
			return 0, err
		}
		requestBody = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return 0, err
	}
	if value != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	return resp.StatusCode, nil
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

func jsonMap(value map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		if result, ok := value[key].(map[string]any); ok {
			return result
		}
	}
	return map[string]any{}
}

func Example_runtimeSmoke() {
	fmt.Println("QHPRO_RUNTIME_SMOKE=1 JWT_KEY_GENERATE=<local-key> go test -run TestQHPROCommercialRuntime")
	// Output: QHPRO_RUNTIME_SMOKE=1 JWT_KEY_GENERATE=<local-key> go test -run TestQHPROCommercialRuntime
}
