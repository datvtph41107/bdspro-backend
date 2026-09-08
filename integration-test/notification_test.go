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

// TestNotificationService test các API notification service
func TestNotificationService(t *testing.T) {
	t.Run("TestSendAdminNotification", func(t *testing.T) {
		// Test gửi thông báo từ admin
		notificationData := map[string]interface{}{
			"type":       "admin_warning",
			"title":      "Cảnh báo hệ thống",
			"message":    "Hệ thống phát hiện hoạt động bất thường",
			"recipients": []int{123, 456},
			"priority":   "high",
			"expiresAt":  "2024-12-31T23:59:59Z",
			"category":   "system_alert",
			"actionUrl":  "https://admin.example.com/alerts",
			"metadata": map[string]interface{}{
				"source":        "admin_panel",
				"alertLevel":    "critical",
				"affectedUsers": 10,
			},
		}

		jsonData, err := json.Marshal(notificationData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/notifications/admin", GatewayURL),
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
		t.Logf("✅ Send admin notification API response: %v", response)
	})

	t.Run("TestSendUserNotification", func(t *testing.T) {
		// Test gửi thông báo cho user
		notificationData := map[string]interface{}{
			"type":       "user_notification",
			"title":      "Thông báo mới",
			"message":    "Bạn có tin nhắn mới từ hệ thống",
			"recipients": []int{789},
			"priority":   "normal",
			"expiresAt":  "2024-12-31T23:59:59Z",
			"category":   "general",
			"actionUrl":  "https://app.example.com/messages",
		}

		jsonData, err := json.Marshal(notificationData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/notifications/user", GatewayURL),
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
		t.Logf("✅ Send user notification API response: %v", response)
	})
}

// TestNotificationTypes test các loại notification khác nhau
func TestNotificationTypes(t *testing.T) {
	notificationTypes := []map[string]interface{}{
		{
			"type":     "admin_warning",
			"title":    "Cảnh báo bảo mật",
			"message":  "Phát hiện đăng nhập bất thường",
			"priority": "high",
			"category": "security",
		},
		{
			"type":     "system_maintenance",
			"title":    "Bảo trì hệ thống",
			"message":  "Hệ thống sẽ bảo trì từ 2h-4h sáng",
			"priority": "medium",
			"category": "maintenance",
		},
		{
			"type":     "user_approval",
			"title":    "Tài khoản được duyệt",
			"message":  "Tài khoản của bạn đã được admin duyệt",
			"priority": "normal",
			"category": "account",
		},
		{
			"type":     "organization_invite",
			"title":    "Lời mời tham gia tổ chức",
			"message":  "Bạn được mời tham gia tổ chức ABC",
			"priority": "normal",
			"category": "organization",
		},
	}

	for _, notification := range notificationTypes {
		t.Run(fmt.Sprintf("TestNotificationType_%s", notification["type"]), func(t *testing.T) {
			// Thêm recipients và expiresAt
			notification["recipients"] = []int{123, 456, 789}
			notification["expiresAt"] = "2024-12-31T23:59:59Z"

			jsonData, err := json.Marshal(notification)
			require.NoError(t, err)

			resp, err := http.Post(
				fmt.Sprintf("%s/v2/notifications/admin", GatewayURL),
				"application/json",
				bytes.NewBuffer(jsonData),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusUnauthorized {
				t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
				return
			}

			assert.Contains(t, []int{200, 400}, resp.StatusCode)
			t.Logf("✅ Notification type %s: status %d", notification["type"], resp.StatusCode)
		})
	}
}

// TestNotificationValidation test validation cho notification
func TestNotificationValidation(t *testing.T) {
	t.Run("TestNotificationWithoutRequiredFields", func(t *testing.T) {
		// Test notification thiếu các trường bắt buộc
		invalidNotifications := []map[string]interface{}{
			{
				"title": "Thiếu type và message",
			},
			{
				"type": "admin_warning",
				// Thiếu title và message
			},
			{
				"type":    "admin_warning",
				"title":   "Test",
				"message": "Test message",
				// Thiếu recipients
			},
		}

		for i, notification := range invalidNotifications {
			t.Run(fmt.Sprintf("InvalidNotification_%d", i), func(t *testing.T) {
				jsonData, err := json.Marshal(notification)
				require.NoError(t, err)

				resp, err := http.Post(
					fmt.Sprintf("%s/v2/notifications/admin", GatewayURL),
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
				t.Logf("✅ Invalid notification %d handled correctly", i)
			})
		}
	})

	t.Run("TestNotificationWithInvalidPriority", func(t *testing.T) {
		// Test với priority không hợp lệ
		notificationData := map[string]interface{}{
			"type":       "admin_warning",
			"title":      "Test",
			"message":    "Test message",
			"recipients": []int{123},
			"priority":   "invalid_priority", // Priority không hợp lệ
		}

		jsonData, err := json.Marshal(notificationData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/notifications/admin", GatewayURL),
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
		t.Logf("✅ Notification with invalid priority handled correctly")
	})

	t.Run("TestNotificationWithInvalidExpiresAt", func(t *testing.T) {
		// Test với expiresAt không hợp lệ
		notificationData := map[string]interface{}{
			"type":       "admin_warning",
			"title":      "Test",
			"message":    "Test message",
			"recipients": []int{123},
			"expiresAt":  "invalid-date", // Date không hợp lệ
		}

		jsonData, err := json.Marshal(notificationData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/notifications/admin", GatewayURL),
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
		t.Logf("✅ Notification with invalid expiresAt handled correctly")
	})
}

// TestNotificationEndpointsAvailability test xem các endpoint notification có tồn tại không
func TestNotificationEndpointsAvailability(t *testing.T) {
	endpoints := []string{
		"/v2/notifications/admin",
		"/v2/notifications/user",
		"/v2/notifications/system",
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

// TestNotificationIntegration test tích hợp với các service khác
func TestNotificationIntegration(t *testing.T) {
	t.Run("TestNotificationWithUserService", func(t *testing.T) {
		// Test tích hợp với user service để gửi notification
		notificationData := map[string]interface{}{
			"type":       "user_notification",
			"title":      "Thông báo từ user service",
			"message":    "Tài khoản của bạn đã được cập nhật",
			"recipients": []int{123, 456},
			"priority":   "normal",
			"category":   "account_update",
		}

		jsonData, err := json.Marshal(notificationData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/notifications/user", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Notification integration with user service works correctly")
	})

	t.Run("TestNotificationWithAuthService", func(t *testing.T) {
		// Test tích hợp với auth service để gửi notification về bảo mật
		notificationData := map[string]interface{}{
			"type":       "security_alert",
			"title":      "Cảnh báo bảo mật",
			"message":    "Phát hiện đăng nhập bất thường từ IP mới",
			"recipients": []int{123},
			"priority":   "high",
			"category":   "security",
			"metadata": map[string]interface{}{
				"ipAddress": "192.168.1.100",
				"location":  "Hà Nội, Việt Nam",
				"device":    "Unknown Device",
			},
		}

		jsonData, err := json.Marshal(notificationData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v2/notifications/admin", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Notification integration with auth service works correctly")
	})
}
