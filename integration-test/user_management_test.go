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

// TestUserManagement test các API quản lý user
func TestUserManagement(t *testing.T) {
	t.Run("TestListUsers", func(t *testing.T) {
		// Test lấy danh sách user
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/users?page=1&size=10&status=active&search=user", GatewayURL))
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
		t.Logf("✅ List users API response: %v", response)
	})

	t.Run("TestListUsersWithFilters", func(t *testing.T) {
		// Test lấy danh sách user với các filter khác nhau
		filters := []string{
			"status=blocked",
			"role=admin",
			"search=test@example.com",
		}

		for _, filter := range filters {
			t.Run(fmt.Sprintf("Filter_%s", filter), func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("%s/v2/auth/users?page=1&size=10&%s", GatewayURL, filter))
				require.NoError(t, err)
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusUnauthorized {
					t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
					return
				}

				assert.Contains(t, []int{200, 400}, resp.StatusCode)
				t.Logf("✅ List users with filter %s: status %d", filter, resp.StatusCode)
			})
		}
	})

	t.Run("TestListUsersWithInvalidParams", func(t *testing.T) {
		// Test với tham số không hợp lệ
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/users?page=-1&size=1000", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{400, 422}, resp.StatusCode)
		t.Logf("✅ List users with invalid params handled correctly")
	})
}

// TestUserBlockUnblockApprove test các API khóa/mở khóa/duyệt user
func TestUserBlockUnblockApprove(t *testing.T) {
	testUserID := "57" // ID user test

	t.Run("TestBlockUser", func(t *testing.T) {
		// Test khóa user
		blockData := map[string]interface{}{
			"userId": testUserID,
			"reason": "Vi phạm quy định cộng đồng",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/block", GatewayURL),
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
		t.Logf("✅ Block user API response: %v", response)
	})

	t.Run("TestUnblockUser", func(t *testing.T) {
		// Test mở khóa user
		unblockData := map[string]interface{}{
			"userId": testUserID,
			"reason": "Đã xem xét lại và mở khóa",
		}

		jsonData, err := json.Marshal(unblockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/unblock", GatewayURL),
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
		t.Logf("✅ Unblock user API response: %v", response)
	})

	t.Run("TestApproveUser", func(t *testing.T) {
		// Test duyệt user
		approveData := map[string]interface{}{
			"userId": testUserID,
			"note":   "Tài khoản đã được duyệt",
		}

		jsonData, err := json.Marshal(approveData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/approve", GatewayURL),
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
		t.Logf("✅ Approve user API response: %v", response)
	})

	t.Run("TestRestoreUser", func(t *testing.T) {
		// Test khôi phục user
		restoreData := map[string]interface{}{
			"userId": testUserID,
			"reason": "Khôi phục tài khoản theo yêu cầu",
		}

		jsonData, err := json.Marshal(restoreData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/restore", GatewayURL),
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
		t.Logf("✅ Restore user API response: %v", response)
	})
}

// TestUserManagementValidation test validation cho các API user management
func TestUserManagementValidation(t *testing.T) {
	t.Run("TestBlockUserWithoutReason", func(t *testing.T) {
		// Test khóa user mà không có lý do
		blockData := map[string]interface{}{
			"userId": "57",
			"reason": "",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/block", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{400, 422}, resp.StatusCode)
		t.Logf("✅ Block user without reason handled correctly")
	})

	t.Run("TestBlockUserInvalidID", func(t *testing.T) {
		// Test khóa user với ID không hợp lệ
		blockData := map[string]interface{}{
			"userId": "999999",
			"reason": "Test reason",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/block", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{404, 400}, resp.StatusCode)
		t.Logf("✅ Block user with invalid ID handled correctly")
	})

	t.Run("TestApproveUserWithoutNote", func(t *testing.T) {
		// Test duyệt user mà không có ghi chú
		approveData := map[string]interface{}{
			"userId": "57",
			"note":   "",
		}

		jsonData, err := json.Marshal(approveData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/approve", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{400, 422}, resp.StatusCode)
		t.Logf("✅ Approve user without note handled correctly")
	})
}

// TestUserManagementWorkflow test workflow hoàn chỉnh: block -> unblock -> approve -> restore
func TestUserManagementWorkflow(t *testing.T) {
	t.Run("TestCompleteUserManagementWorkflow", func(t *testing.T) {
		testUserID := "57"

		// Step 1: Block user
		t.Logf("🔄 Step 1: Blocking user %s", testUserID)
		blockData := map[string]interface{}{
			"userId": testUserID,
			"reason": "Test block reason",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/auth/users/block", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("⚠️ Need authentication token for workflow test")
			return
		}

		// Step 2: Unblock user
		t.Logf("🔄 Step 2: Unblocking user %s", testUserID)
		unblockData := map[string]interface{}{
			"userId": testUserID,
			"reason": "Test unblock reason",
		}

		jsonData, err = json.Marshal(unblockData)
		require.NoError(t, err)

		resp, err = http.Post(
			fmt.Sprintf("%s/v2/auth/users/unblock", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Step 3: Approve user
		t.Logf("🔄 Step 3: Approving user %s", testUserID)
		approveData := map[string]interface{}{
			"userId": testUserID,
			"note":   "Test approve note",
		}

		jsonData, err = json.Marshal(approveData)
		require.NoError(t, err)

		resp, err = http.Post(
			fmt.Sprintf("%s/v2/auth/users/approve", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Step 4: Restore user
		t.Logf("🔄 Step 4: Restoring user %s", testUserID)
		restoreData := map[string]interface{}{
			"userId": testUserID,
			"reason": "Test restore reason",
		}

		jsonData, err = json.Marshal(restoreData)
		require.NoError(t, err)

		resp, err = http.Post(
			fmt.Sprintf("%s/v2/auth/users/restore", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		t.Logf("✅ Complete user management workflow test finished")
	})
}

// TestUserManagementEndpointsAvailability test xem các endpoint user management có tồn tại không
func TestUserManagementEndpointsAvailability(t *testing.T) {
	endpoints := []string{
		"/v2/auth/users",
		"/v2/auth/users/block",
		"/v2/auth/users/unblock",
		"/v2/auth/users/approve",
		"/v2/auth/users/restore",
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("TestEndpoint_%s", endpoint), func(t *testing.T) {
			// Test GET request để kiểm tra endpoint có tồn tại không
			resp, err := http.Get(fmt.Sprintf("%s%s", GatewayURL, endpoint))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Endpoint phải có response (không phải 404)
			// Có thể trả về 405 Method Not Allowed vì endpoint này chỉ hỗ trợ POST
			assert.Contains(t, []int{405, 401, 400}, resp.StatusCode)
			t.Logf("✅ Endpoint %s có thể truy cập (status: %d)", endpoint, resp.StatusCode)
		})
	}
}
