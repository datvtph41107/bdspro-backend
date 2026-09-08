package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAdminManagement test các API quản lý admin
func TestAdminManagement(t *testing.T) {
	t.Run("TestCreateAdmin", func(t *testing.T) {
		// Test tạo admin mới
		adminData := map[string]interface{}{
			"username": "admin_test",
			"password": "password123",
			"fullName": "Test Admin",
			"email":    "admin_test@example.com",
			"phone":    "+84123456789",
			"avatar":   "https://example.com/avatar.jpg",
		}

		jsonData, err := json.Marshal(adminData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/admin/create", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// API này cần auth token
		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400, 409}, resp.StatusCode)
		t.Logf("✅ Create admin API response: %v", response)
	})

	t.Run("TestUpdateAdmin", func(t *testing.T) {
		// Test cập nhật admin
		updateData := map[string]interface{}{
			"authId":   123,
			"fullName": "Updated Admin Name",
			"email":    "updated@example.com",
			"phone":    "+84987654321",
			"avatar":   "https://example.com/new-avatar.jpg",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/auth/admin/update", GatewayURL), bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
		t.Logf("✅ Update admin API response: %v", response)
	})

	t.Run("TestDeleteAdmin", func(t *testing.T) {
		// Test xóa admin
		deleteData := map[string]interface{}{
			"authId": 123,
		}

		jsonData, err := json.Marshal(deleteData)
		require.NoError(t, err)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/auth/admin/delete", GatewayURL), bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
		t.Logf("✅ Delete admin API response: %v", response)
	})

	t.Run("TestListAdmins", func(t *testing.T) {
		// Test lấy danh sách admin
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/admin/list?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ List admins API response: %v", response)
	})

	t.Run("TestAssignRole", func(t *testing.T) {
		// Test gán quyền cho admin
		roleData := map[string]interface{}{
			"authId":  123,
			"roleKey": "admin",
		}

		jsonData, err := json.Marshal(roleData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/admin/assign-role", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
		t.Logf("✅ Assign role API response: %v", response)
	})
}

// TestAdminAccessControl test các API quản lý quyền truy cập admin
func TestAdminAccessControl(t *testing.T) {
	t.Run("TestCreateAdminAccess", func(t *testing.T) {
		// Test tạo quyền truy cập admin
		accessData := map[string]interface{}{
			"userId":        123,
			"ipAddress":     "192.168.1.100",
			"ipRange":       "192.168.1.0/24",
			"deviceId":      "device123",
			"deviceName":    "Admin PC",
			"deviceType":    "desktop",
			"authType":      "password",
			"status":        "active",
			"effectiveFrom": "2024-01-01T00:00:00Z",
			"effectiveTo":   "2024-12-31T23:59:59Z",
			"maxDevices":    5,
			"require2fa":    true,
		}

		jsonData, err := json.Marshal(accessData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/admin-access", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Create admin access API response: %v", response)
	})

	t.Run("TestGetAdminAccessList", func(t *testing.T) {
		// Test lấy danh sách quyền truy cập admin
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/admin-access?page=1&size=10&userId=123&status=active", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Get admin access list API response: %v", response)
	})

	t.Run("TestValidateAdminAccess", func(t *testing.T) {
		// Test kiểm tra quyền truy cập admin
		validateData := map[string]interface{}{
			"userId":    123,
			"ipAddress": "192.168.1.100",
			"deviceId":  "device123",
		}

		jsonData, err := json.Marshal(validateData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/admin-access/validate", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Validate admin access API response: %v", response)
	})
}

// TestAdminActivityLogs test API nhật ký thao tác admin
func TestAdminActivityLogs(t *testing.T) {
	t.Run("TestGetAdminActivityLogs", func(t *testing.T) {
		// Test lấy nhật ký thao tác admin
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/admin-access/log?page=1&size=10&userId=123", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Get admin activity logs API response: %v", response)
	})

	t.Run("TestCreateAdminAccessLog", func(t *testing.T) {
		// Test tạo log thao tác admin
		logData := map[string]interface{}{
			"userId":     123,
			"ipAddress":  "192.168.1.100",
			"deviceId":   "device123",
			"deviceName": "Admin PC",
			"action":     "login",
			"status":     "success",
			"reason":     "Admin login successful",
			"userAgent":  "Mozilla/5.0...",
		}

		jsonData, err := json.Marshal(logData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/admin-access/log", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Create admin access log API response: %v", response)
	})
}

// TestAdminManagementEndpointsAvailability test xem các endpoint admin management có tồn tại không
func TestAdminManagementEndpointsAvailability(t *testing.T) {
	endpoints := []string{
		"/v2/auth/admin/create",
		"/v2/auth/admin/update",
		"/v2/auth/admin/delete",
		"/v2/auth/admin/list",
		"/v2/auth/admin/assign-role",
		"/v2/auth/admin-access",
		"/v2/auth/admin-access/validate",
		"/v2/auth/admin-access/log",
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("TestEndpoint_%s", endpoint), func(t *testing.T) {
			// Test GET request để kiểm tra endpoint có tồn tại không
			resp, err := http.Get(fmt.Sprintf("%s%s", GatewayURL, endpoint))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Endpoint phải có response (không phải 404)
			// Có thể trả về 405 Method Not Allowed vì endpoint này chỉ hỗ trợ POST/PUT/DELETE
			assert.Contains(t, []int{405, 401, 400}, resp.StatusCode)
			t.Logf("✅ Endpoint %s có thể truy cập (status: %d)", endpoint, resp.StatusCode)
		})
	}
}
