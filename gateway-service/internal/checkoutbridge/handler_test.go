package checkoutbridge

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHandlerRendersSignedPOSTToOfficialSandbox(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v2/payment/sepay/checkout", Handler())
	query := url.Values{
		"checkout_url":         {"https://pay-sandbox.sepay.vn/v1/checkout/init"},
		"merchant":             {"merchant-1"},
		"operation":            {"PURCHASE"},
		"payment_method":       {"BANK_TRANSFER"},
		"order_invoice_number": {"QHP-1"},
		"order_amount":         {"99000"},
		"currency":             {"VND"},
		"signature":            {"signed-value"},
	}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/payment/sepay/checkout?"+query.Encode(), nil)
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Security-Policy"), "form-action https://pay-sandbox.sepay.vn")
	require.Contains(t, recorder.Body.String(), `action="https://pay-sandbox.sepay.vn/v1/checkout/init" method="post"`)
	require.Contains(t, recorder.Body.String(), `name="signature" value="signed-value"`)
}

func TestHandlerRejectsArbitraryCheckoutTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v2/payment/sepay/checkout", Handler())
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/payment/sepay/checkout?checkout_url="+url.QueryEscape("https://attacker.example/collect")+"&merchant=m&operation=PURCHASE&payment_method=CARD&order_invoice_number=o&order_amount=1&currency=VND&signature=s", nil)
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.False(t, strings.Contains(recorder.Body.String(), "attacker.example"))
}
