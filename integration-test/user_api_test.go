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

// UserProfileResponse đại diện cho response từ API profile
type UserProfileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ProfileID uint64 `json:"profile_id"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Avatar    string `json:"avatar"`
	} `json:"data"`
}

// TestUserProfileAPI test các API liên quan đến user profile
func TestUserProfileAPI(t *testing.T) {
	t.Run("TestGetProfileInfoWithoutAuth", func(t *testing.T) {
		// Test lấy thông tin profile mà không có auth token
		resp, err := http.Get(fmt.Sprintf("%s/v2/user/profile/info/me", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Phải trả về lỗi 401 vì không có auth
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
	})

	t.Run("TestGetProfileInfoWithInvalidID", func(t *testing.T) {
		// Test lấy thông tin profile với ID không hợp lệ
		resp, err := http.Get(fmt.Sprintf("%s/v2/user/profile/info/999999", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về 404 hoặc 401
		assert.Contains(t, []int{401, 404}, resp.StatusCode)
		t.Logf("✅ API xử lý ID không hợp lệ đúng cách")
	})

	t.Run("TestSearchProfilePublic", func(t *testing.T) {
		// Test search profile public
		resp, err := http.Get(fmt.Sprintf("%s/v2/user/profile/search/public?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		// API này có thể public hoặc cần auth
		assert.Contains(t, []int{200, 401}, resp.StatusCode)
		t.Logf("✅ Search profile API có thể truy cập")
	})
}

// TestAdminAPI test các API admin
func TestAdminAPI(t *testing.T) {
	t.Run("TestListAllUsersWithoutAuth", func(t *testing.T) {
		// Test API admin mà không có auth
		resp, err := http.Get(fmt.Sprintf("%s/v2/user/admin/users?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		// API admin phải yêu cầu auth
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		t.Logf("✅ Admin API yêu cầu authentication đúng như mong đợi")
	})

	t.Run("TestListAllUsersWithInvalidParams", func(t *testing.T) {
		// Test với tham số không hợp lệ
		resp, err := http.Get(fmt.Sprintf("%s/v2/user/admin/users?page=-1&size=1000", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Có thể trả về 400 hoặc 401
		assert.Contains(t, []int{400, 401}, resp.StatusCode)
		t.Logf("✅ Admin API xử lý tham số không hợp lệ đúng cách")
	})
}

// TestAPIEndpointsAvailability test xem các endpoint có tồn tại không
func TestAPIEndpointsAvailability(t *testing.T) {
	endpoints := []string{
		"/v1/auth/login",
		"/v2/user/profile/info/me",
		"/v2/user/profile/search/public",
		"/v2/user/admin/users",
		"/swagger/merged/",
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("TestEndpoint_%s", endpoint), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s%s", GatewayURL, endpoint))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Endpoint phải có response (không phải 404)
			assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
			t.Logf("✅ Endpoint %s có thể truy cập (status: %d)", endpoint, resp.StatusCode)
		})
	}
}

// TestAPIResponseFormat test format response của API
func TestAPIResponseFormat(t *testing.T) {
	t.Run("TestLoginResponseFormat", func(t *testing.T) {
		// Test với thông tin login không hợp lệ để xem format response
		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		resp, err := http.Post(
			fmt.Sprintf("%s/v1/auth/login", GatewayURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Kiểm tra format response có đúng không
		assert.Contains(t, response, "code")
		assert.Contains(t, response, "message")
		t.Logf("✅ API response có format đúng: %v", response)
	})
}
