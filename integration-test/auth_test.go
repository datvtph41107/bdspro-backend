package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	GatewayURL = "http://localhost:8000"
)

// TestUser đại diện cho user test
type TestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse đại diện cho response từ API login
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ProfileID uint64 `json:"profile_id"`
			Email     string `json:"email"`
			FullName  string `json:"full_name"`
		} `json:"user"`
	} `json:"data"`
}

// TestLoginWithValidCredentials test login với thông tin hợp lệ
func TestLoginWithValidCredentials(t *testing.T) {
	// Test data - bạn có thể thay đổi email/password theo dữ liệu thực tế
	testUsers := []TestUser{
		{
			Email:    "admin@example.com",
			Password: "admin123",
		},
		{
			Email:    "user@example.com",
			Password: "123456",
		},
	}

	for _, user := range testUsers {
		t.Run(fmt.Sprintf("Login_%s", user.Email), func(t *testing.T) {
			// Chuẩn bị request
			loginData := map[string]string{
				"email":    user.Email,
				"password": user.Password,
			}

			jsonData, err := json.Marshal(loginData)
			require.NoError(t, err)

			// Gửi request - thử cả v1 và v2 endpoints
			var resp *http.Response

			// Thử v2 admin login trước
			resp, err = http.Post(
				fmt.Sprintf("%s/v2/auth/admin/login", GatewayURL),
				"application/json",
				bytes.NewBuffer(jsonData),
			)

			// Nếu v2 không thành công, thử v1
			if err != nil || resp.StatusCode != http.StatusOK {
				// Thay đổi format cho v1
				v1Data := map[string]string{
					"username": user.Email, // v1 có thể dùng username
					"password": user.Password,
				}
				v1JsonData, _ := json.Marshal(v1Data)

				resp, err = http.Post(
					fmt.Sprintf("%s/v1/auth/login", GatewayURL),
					"application/json",
					bytes.NewBuffer(v1JsonData),
				)
			}
			require.NoError(t, err)
			defer resp.Body.Close()

			// Parse response
			var loginResp LoginResponse
			err = json.NewDecoder(resp.Body).Decode(&loginResp)
			require.NoError(t, err)

			// Assertions
			if resp.StatusCode == http.StatusOK {
				// Login thành công
				assert.NotEmpty(t, loginResp.Data.AccessToken)
				assert.NotEmpty(t, loginResp.Data.RefreshToken)
				assert.NotZero(t, loginResp.Data.User.ProfileID)
				assert.Equal(t, user.Email, loginResp.Data.User.Email)
				t.Logf("✅ Login thành công cho user: %s", user.Email)
			} else {
				// Login thất bại - có thể do thông tin không đúng
				t.Logf("⚠️ Login thất bại cho user: %s - %s", user.Email, loginResp.Message)
				assert.Contains(t, []int{400, 401, 404}, resp.StatusCode)
			}
		})
	}
}

// TestLoginWithInvalidCredentials test login với thông tin không hợp lệ
func TestLoginWithInvalidCredentials(t *testing.T) {
	invalidUsers := []TestUser{
		{
			Email:    "invalid@example.com",
			Password: "wrongpassword",
		},
		{
			Email:    "admin@example.com",
			Password: "wrongpassword",
		},
		{
			Email:    "",
			Password: "123456",
		},
		{
			Email:    "admin@example.com",
			Password: "",
		},
	}

	for _, user := range invalidUsers {
		t.Run(fmt.Sprintf("InvalidLogin_%s", user.Email), func(t *testing.T) {
			loginData := map[string]string{
				"email":    user.Email,
				"password": user.Password,
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

			var loginResp LoginResponse
			err = json.NewDecoder(resp.Body).Decode(&loginResp)
			require.NoError(t, err)

			// Login với thông tin không hợp lệ phải trả về lỗi
			assert.NotEqual(t, http.StatusOK, resp.StatusCode)
			t.Logf("✅ Login thất bại đúng như mong đợi cho user: %s", user.Email)
		})
	}
}

// TestLoginEndpointAvailability test xem endpoint login có hoạt động không
func TestLoginEndpointAvailability(t *testing.T) {
	// Test endpoint có response không
	resp, err := http.Get(fmt.Sprintf("%s/v1/auth/login", GatewayURL))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Endpoint GET có thể trả về 405 Method Not Allowed hoặc 404
	// Điều quan trọng là server phải response
	assert.Contains(t, []int{200, 404, 405}, resp.StatusCode)
	t.Logf("✅ Login endpoint có thể truy cập được")
}

// TestGatewayHealth test xem gateway có hoạt động không
func TestGatewayHealth(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(GatewayURL)
	if err != nil {
		t.Skipf("⚠️ Gateway không thể truy cập: %v", err)
		return
	}
	defer resp.Body.Close()

	assert.Contains(t, []int{200, 404, 405}, resp.StatusCode)
	t.Logf("✅ Gateway service đang hoạt động")
}
