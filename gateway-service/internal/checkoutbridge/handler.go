// Package checkoutbridge converts a safe GET navigation from a browser or
// native application into the signed HTML POST required by SePay Checkout.
// It never receives or stores the merchant secret.
package checkoutbridge

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var checkoutHosts = map[string]struct{}{
	"https://pay-sandbox.sepay.vn/v1/checkout/init": {},
	"https://pay.sepay.vn/v1/checkout/init":         {},
}

var fieldOrder = []string{
	"merchant", "operation", "payment_method", "order_invoice_number",
	"order_amount", "currency", "order_description", "customer_id",
	"success_url", "error_url", "cancel_url", "custom_data", "signature",
}

type formField struct {
	Name  string
	Value string
}

type pageData struct {
	Action template.URL
	Nonce  string
	Fields []formField
}

var checkoutPage = template.Must(template.New("sepay-checkout").Parse(`<!doctype html>
<html lang="vi">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Đang chuyển đến SePay</title>
  <style nonce="{{.Nonce}}">body{font-family:system-ui,sans-serif;background:#f8fafc;color:#0f172a;display:grid;min-height:100vh;margin:0;place-items:center}.card{background:#fff;border:1px solid #e2e8f0;border-radius:16px;padding:28px;text-align:center;box-shadow:0 12px 30px rgba(15,23,42,.08)}button{background:#0f766e;border:0;border-radius:10px;color:#fff;font-weight:700;padding:12px 18px}</style>
</head>
<body>
  <main class="card">
    <p>Đang chuyển đến cổng thanh toán SePay…</p>
    <form id="sepay-checkout" action="{{.Action}}" method="post" referrerpolicy="no-referrer">
      {{range .Fields}}<input type="hidden" name="{{.Name}}" value="{{.Value}}">{{end}}
      <noscript><button type="submit">Tiếp tục thanh toán</button></noscript>
    </form>
  </main>
  <script nonce="{{.Nonce}}">document.getElementById('sepay-checkout').submit();</script>
</body>
</html>`))

// Handler renders only official SePay endpoints and an allowlisted field set.
// SePay remains the authority that verifies the merchant HMAC signature.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		action := strings.TrimSpace(c.Query("checkout_url"))
		if _, ok := checkoutHosts[action]; !ok || !validRequest(c.Request.URL.Query()) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": "INVALID_SEPAY_CHECKOUT", "message": "Thông tin khởi tạo thanh toán không hợp lệ",
			})
			return
		}

		nonceBytes := make([]byte, 18)
		if _, err := rand.Read(nonceBytes); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		nonce := base64.RawStdEncoding.EncodeToString(nonceBytes)
		fields := make([]formField, 0, len(fieldOrder))
		for _, name := range fieldOrder {
			if value := c.Query(name); value != "" {
				fields = append(fields, formField{Name: name, Value: value})
			}
		}

		c.Header("Cache-Control", "no-store")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Content-Security-Policy", "default-src 'none'; style-src 'nonce-"+nonce+"'; script-src 'nonce-"+nonce+"'; form-action https://pay-sandbox.sepay.vn https://pay.sepay.vn; base-uri 'none'; frame-ancestors 'none'")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)
		if err := checkoutPage.Execute(c.Writer, pageData{Action: template.URL(action), Nonce: nonce, Fields: fields}); err != nil {
			_ = c.Error(err)
		}
		c.Abort()
	}
}

// Middleware avoids a conflicting concrete Gin route under the existing
// /v2/*any gRPC-Gateway catch-all. Non-checkout requests pass through unchanged.
func Middleware() gin.HandlerFunc {
	handler := Handler()
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet || c.Request.URL.Path != "/v2/payment/sepay/checkout" {
			c.Next()
			return
		}
		handler(c)
	}
}

func validRequest(values url.Values) bool {
	allowed := map[string]struct{}{"checkout_url": {}}
	for _, field := range fieldOrder {
		allowed[field] = struct{}{}
	}
	for name, entries := range values {
		if _, ok := allowed[name]; !ok || len(entries) != 1 {
			return false
		}
	}
	amount, err := strconv.ParseInt(values.Get("order_amount"), 10, 64)
	if err != nil || amount <= 0 || values.Get("merchant") == "" || values.Get("signature") == "" ||
		values.Get("operation") != "PURCHASE" || values.Get("currency") != "VND" ||
		values.Get("order_invoice_number") == "" {
		return false
	}
	switch values.Get("payment_method") {
	case "BANK_TRANSFER", "NAPAS_BANK_TRANSFER", "CARD":
		return true
	default:
		return false
	}
}
