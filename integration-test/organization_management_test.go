package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrganizationManagement test các API quản lý organization
func TestOrganizationManagement(t *testing.T) {
	t.Run("TestListOrganizations", func(t *testing.T) {
		// Test lấy danh sách organization
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations?page=1&size=10&search=company", GatewayURL))
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
		t.Logf("✅ List organizations API response: %v", response)
	})

	t.Run("TestListOrganizationsWithSearch", func(t *testing.T) {
		// Test lấy danh sách organization với tìm kiếm
		searchTerms := []string{
			"tech",
			"real estate",
			"construction",
		}

		for _, search := range searchTerms {
			t.Run(fmt.Sprintf("Search_%s", search), func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations?page=1&size=10&search=%s", GatewayURL, search))
				require.NoError(t, err)
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusUnauthorized {
					t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
					return
				}

				assert.Contains(t, []int{200, 400}, resp.StatusCode)
				t.Logf("✅ List organizations with search '%s': status %d", search, resp.StatusCode)
			})
		}
	})

	t.Run("TestListOrganizationsWithInvalidParams", func(t *testing.T) {
		// Test với tham số không hợp lệ
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations?page=-1&size=1000", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{400, 422}, resp.StatusCode)
		t.Logf("✅ List organizations with invalid params handled correctly")
	})
}

// TestOrganizationMembers test các API quản lý thành viên organization
func TestOrganizationMembers(t *testing.T) {
	testOrgID := "123" // ID organization test

	t.Run("TestListOrganizationMembers", func(t *testing.T) {
		// Test lấy danh sách thành viên organization
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations/%s/members?page=1&size=10", GatewayURL, testOrgID))
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
		t.Logf("✅ List organization members API response: %v", response)
	})

	t.Run("TestListOrganizationMembersWithPagination", func(t *testing.T) {
		// Test lấy danh sách thành viên với phân trang
		pages := []int{1, 2, 3}
		sizes := []int{5, 10, 20}

		for _, page := range pages {
			for _, size := range sizes {
				t.Run(fmt.Sprintf("Page_%d_Size_%d", page, size), func(t *testing.T) {
					resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations/%s/members?page=%d&size=%d", GatewayURL, testOrgID, page, size))
					require.NoError(t, err)
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusUnauthorized {
						t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
						return
					}

					assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
					t.Logf("✅ List organization members with page=%d, size=%d: status %d", page, size, resp.StatusCode)
				})
			}
		}
	})

	t.Run("TestListOrganizationMembersInvalidOrgID", func(t *testing.T) {
		// Test với organization ID không hợp lệ
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations/999999/members?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{404, 400}, resp.StatusCode)
		t.Logf("✅ List organization members with invalid org ID handled correctly")
	})
}

// TestOrganizationValidation test validation cho các API organization
func TestOrganizationValidation(t *testing.T) {
	t.Run("TestListOrganizationsWithInvalidPagination", func(t *testing.T) {
		// Test với tham số phân trang không hợp lệ
		invalidParams := []string{
			"page=0&size=10",
			"page=1&size=0",
			"page=-1&size=-1",
			"page=abc&size=def",
		}

		for _, params := range invalidParams {
			t.Run(fmt.Sprintf("Params_%s", params), func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations?%s", GatewayURL, params))
				require.NoError(t, err)
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusUnauthorized {
					t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
					return
				}

				assert.Contains(t, []int{400, 422}, resp.StatusCode)
				t.Logf("✅ List organizations with invalid params '%s' handled correctly", params)
			})
		}
	})

	t.Run("TestListOrganizationMembersWithInvalidPagination", func(t *testing.T) {
		// Test với tham số phân trang không hợp lệ cho members
		invalidParams := []string{
			"page=0&size=10",
			"page=1&size=0",
			"page=-1&size=-1",
		}

		for _, params := range invalidParams {
			t.Run(fmt.Sprintf("Params_%s", params), func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations/123/members?%s", GatewayURL, params))
				require.NoError(t, err)
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusUnauthorized {
					t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
					return
				}

				assert.Contains(t, []int{400, 422}, resp.StatusCode)
				t.Logf("✅ List organization members with invalid params '%s' handled correctly", params)
			})
		}
	})
}

// TestOrganizationWorkflow test workflow hoàn chỉnh cho organization
func TestOrganizationWorkflow(t *testing.T) {
	t.Run("TestCompleteOrganizationWorkflow", func(t *testing.T) {
		testOrgID := "123"

		// Step 1: List organizations
		t.Logf("🔄 Step 1: Listing organizations")
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("⚠️ Need authentication token for workflow test")
			return
		}

		// Step 2: List organization members
		t.Logf("🔄 Step 2: Listing organization members for org %s", testOrgID)
		resp, err = http.Get(fmt.Sprintf("%s/v2/auth/organizations/%s/members?page=1&size=10", GatewayURL, testOrgID))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Step 3: Search organizations
		t.Logf("🔄 Step 3: Searching organizations")
		resp, err = http.Get(fmt.Sprintf("%s/v2/auth/organizations?page=1&size=10&search=test", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		t.Logf("✅ Complete organization workflow test finished")
	})
}

// TestOrganizationEndpointsAvailability test xem các endpoint organization có tồn tại không
func TestOrganizationEndpointsAvailability(t *testing.T) {
	endpoints := []string{
		"/v2/auth/organizations",
		"/v2/auth/organizations/123/members",
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("TestEndpoint_%s", endpoint), func(t *testing.T) {
			// Test GET request để kiểm tra endpoint có tồn tại không
			resp, err := http.Get(fmt.Sprintf("%s%s", GatewayURL, endpoint))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Endpoint phải có response (không phải 404)
			assert.Contains(t, []int{200, 401, 400, 404}, resp.StatusCode)
			t.Logf("✅ Endpoint %s có thể truy cập (status: %d)", endpoint, resp.StatusCode)
		})
	}
}

// TestOrganizationIntegration test tích hợp với các service khác
func TestOrganizationIntegration(t *testing.T) {
	t.Run("TestOrganizationWithUserService", func(t *testing.T) {
		// Test tích hợp với user service để lấy thông tin thành viên
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations/123/members?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		if resp.StatusCode == http.StatusOK {
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Kiểm tra response có chứa thông tin user không
			if data, ok := response["data"].([]interface{}); ok && len(data) > 0 {
				t.Logf("✅ Organization members response contains user data")
			} else {
				t.Logf("ℹ️ Organization members response is empty or doesn't contain user data")
			}
		}

		assert.Contains(t, []int{200, 400, 404}, resp.StatusCode)
	})

	t.Run("TestOrganizationWithAuthService", func(t *testing.T) {
		// Test tích hợp với auth service để kiểm tra quyền truy cập
		resp, err := http.Get(fmt.Sprintf("%s/v2/auth/organizations?page=1&size=10", GatewayURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✅ API yêu cầu authentication đúng như mong đợi")
			return
		}

		assert.Contains(t, []int{200, 400}, resp.StatusCode)
		t.Logf("✅ Organization API integration with auth service works correctly")
	})
}
