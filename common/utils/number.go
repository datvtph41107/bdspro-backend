package _utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// isValidPhoneNumber kiểm tra số điện thoại có hợp lệ hay không
func IsValidPhoneNumber(phoneNumber string) bool {
	phoneRegex := `^\+?\d{1,3}?[-.\s]?\(?\d{1,4}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phoneNumber)
}

// GenerateOTP tạo mã OTP gồm 6 chữ số
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(1000000)
	return fmt.Sprintf("%06d", otp) // Định dạng thành 6 chữ số
}

// TimeNow trả về thời gian hiện tại dưới dạng chuỗi "yyyy-MM-dd HH:mm:ss"
func TimeNow() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05") // Định dạng giống Java
}

func ValidatePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^(?:\+?84|0)(3|5|7|8|9)[0-9]{8}$`)
	return re.MatchString(phone)
}

func IsParsableToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
