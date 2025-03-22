package services

import (
	"auth/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CookieService struct {
	PathRefresh string
}

func NewCookieService() *CookieService {
	return &CookieService{
		PathRefresh: "/user/auth/token/refresh",
	}
}

// Đặt session browser (có thể bổ sung logic xử lý session tại đây)
func (s *CookieService) SetSessionBrowser(w http.ResponseWriter, refreshToken string) {
	// Thêm logic lưu session nếu cần
}

// Đặt refresh token trong cookie
func (s *CookieService) SetRefreshToken(w http.ResponseWriter, refreshToken string) {
	cookie := &http.Cookie{
		Name:     "rt",
		Value:    refreshToken,
		HttpOnly: true,                                                 // Ngăn JavaScript đọc cookie
		Secure:   true,                                                 // Chỉ gửi cookie qua HTTPS
		MaxAge:   int(config.AppProperties.Jwt.RefreshExpMinutes * 60), // Cookie tồn tại trong 1 ngày
		Path:     s.PathRefresh,                                        // Chỉ áp dụng cho đường dẫn refresh token
		// Domain: "yourdomain.com", // Uncomment nếu backend & frontend chạy trên domain khác nhau
	}

	http.SetCookie(w, cookie)
}

// Set cookie
func (s *CookieService) Set(c *gin.Context, key string, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		HttpOnly: true,         // Ngăn JavaScript đọc cookie
		Secure:   true,         // Chỉ gửi cookie qua HTTPS
		MaxAge:   24 * 60 * 60, // Cookie tồn tại trong 1 ngày
		Path:     "/",          // Chỉ áp dụng cho đường dẫn refresh token
		// Domain: "yourdomain.com", // Uncomment nếu backend & frontend chạy trên domain khác nhau
	}

	http.SetCookie(c.Writer, cookie)
}

// Xóa refresh token khỏi cookie
func (s *CookieService) DeleteRefreshToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "rt",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   0, // Hết hạn ngay lập tức
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}
