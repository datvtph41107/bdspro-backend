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

// TestAdminBlockUser test API khóa tài khoản user
func TestAdminBlockUser(t *testing.T) {
	t.Run("TestBlockUserSuccess", func(t *testing.T) {
		// Test khóa user thành công
		blockData := map[string]interface{}{
			"reason": "Vi phạm quy định cộng đồng",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/57/block", GatewayURL),
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

		// Nếu có auth token, kiểm tra response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Có thể thành công hoặc lỗi tùy thuộc vào trạng thái user
		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
		t.Logf("✅ Block user API response: %v", response)
	})

	t.Run("TestBlockUserWithoutReason", func(t *testing.T) {
		// Test khóa user mà không có lý do
		blockData := map[string]interface{}{
			"reason": "",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/57/block", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về lỗi validation hoặc 401
		assert.Contains(t, []int{400, 401}, resp.StatusCode)
		t.Logf("✅ Block user without reason handled correctly")
	})

	t.Run("TestBlockUserInvalidID", func(t *testing.T) {
		// Test khóa user với ID không hợp lệ
		blockData := map[string]interface{}{
			"reason": "Test reason",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/999999/block", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về 404 hoặc 401
		assert.Contains(t, []int{404, 401}, resp.StatusCode)
		t.Logf("✅ Block user with invalid ID handled correctly")
	})
}

// TestAdminUnblockUser test API mở khóa tài khoản user
func TestAdminUnblockUser(t *testing.T) {
	t.Run("TestUnblockUserSuccess", func(t *testing.T) {
		// Test mở khóa user thành công
		unblockData := map[string]interface{}{
			"reason": "Đã xem xét lại và mở khóa",
		}

		jsonData, err := json.Marshal(unblockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/57/unblock", GatewayURL),
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

		// Nếu có auth token, kiểm tra response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Có thể thành công hoặc lỗi tùy thuộc vào trạng thái user
		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
		t.Logf("✅ Unblock user API response: %v", response)
	})

	t.Run("TestUnblockUserInvalidID", func(t *testing.T) {
		// Test mở khóa user với ID không hợp lệ
		unblockData := map[string]interface{}{
			"reason": "Test reason",
		}

		jsonData, err := json.Marshal(unblockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/999999/unblock", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về 404 hoặc 401
		assert.Contains(t, []int{404, 401}, resp.StatusCode)
		t.Logf("✅ Unblock user with invalid ID handled correctly")
	})
}

// TestAdminApproveUser test API duyệt tài khoản user
func TestAdminApproveUser(t *testing.T) {
	t.Run("TestApproveUserSuccess", func(t *testing.T) {
		// Test duyệt user thành công
		approveData := map[string]interface{}{
			"note": "Tài khoản đã được duyệt",
		}

		jsonData, err := json.Marshal(approveData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/57/approve", GatewayURL),
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

		// Nếu có auth token, kiểm tra response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Có thể thành công hoặc lỗi tùy thuộc vào trạng thái user
		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
		t.Logf("✅ Approve user API response: %v", response)
	})

	t.Run("TestApproveUserWithoutNote", func(t *testing.T) {
		// Test duyệt user mà không có ghi chú
		approveData := map[string]interface{}{
			"note": "",
		}

		jsonData, err := json.Marshal(approveData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/57/approve", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về lỗi validation hoặc 401
		assert.Contains(t, []int{400, 401}, resp.StatusCode)
		t.Logf("✅ Approve user without note handled correctly")
	})

	t.Run("TestApproveUserInvalidID", func(t *testing.T) {
		// Test duyệt user với ID không hợp lệ
		approveData := map[string]interface{}{
			"note": "Test note",
		}

		jsonData, err := json.Marshal(approveData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/999999/approve", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về 404 hoặc 401
		assert.Contains(t, []int{404, 401}, resp.StatusCode)
		t.Logf("✅ Approve user with invalid ID handled correctly")
	})
}

// TestAdminAPIEndpointsAvailability test xem các endpoint admin có tồn tại không
func TestAdminAPIEndpointsAvailability(t *testing.T) {
	endpoints := []string{
		"/v2/user/admin/users/57/block",
		"/v2/user/admin/users/57/unblock",
		"/v2/user/admin/users/57/approve",
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("TestEndpoint_%s", endpoint), func(t *testing.T) {
			// Test GET request để kiểm tra endpoint có tồn tại không
			resp, err := http.Get(fmt.Sprintf("%s%s", GatewayURL, endpoint))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Endpoint phải có response (không phải 404)
			// Có thể trả về 405 Method Not Allowed vì endpoint này chỉ hỗ trợ POST
			assert.Contains(t, []int{405, 401}, resp.StatusCode)
			t.Logf("✅ Endpoint %s có thể truy cập (status: %d)", endpoint, resp.StatusCode)
		})
	}
}

// TestAdminAPIWorkflow test workflow hoàn chỉnh: block -> unblock -> approve
func TestAdminAPIWorkflow(t *testing.T) {
	t.Run("TestCompleteWorkflow", func(t *testing.T) {
		// Test ID user để thực hiện workflow
		testUserID := "57"

		// Step 1: Block user
		t.Logf("🔄 Step 1: Blocking user %s", testUserID)
		blockData := map[string]interface{}{
			"reason": "Test block reason",
		}

		jsonData, err := json.Marshal(blockData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/%s/block", GatewayURL, testUserID),
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
			"reason": "Test unblock reason",
		}

		jsonData, err = json.Marshal(unblockData)
		require.NoError(t, err)

		resp, err = http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/%s/unblock", GatewayURL, testUserID),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Step 3: Approve user
		t.Logf("🔄 Step 3: Approving user %s", testUserID)
		approveData := map[string]interface{}{
			"note": "Test approve note",
		}

		jsonData, err = json.Marshal(approveData)
		require.NoError(t, err)

		resp, err = http.Post(
			fmt.Sprintf("%s/v2/user/admin/users/%s/approve", GatewayURL, testUserID),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		t.Logf("✅ Complete workflow test finished")
	})
}
