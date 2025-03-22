package services

import (
	_utils "common/utils"
	"fmt"
	"payment/data"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type ApiClient struct {
	Client        *resty.Client
	Key           string
	MemVerifyUrl  string
	MemVerifyPath string
}

func NewApiClient() *ApiClient {
	return &ApiClient{
		Client:        resty.New(),
		Key:           viper.GetString("payment.secretKey"),
		MemVerifyUrl:  viper.GetString("membership.url"),
		MemVerifyPath: viper.GetString("membership.path"),
	}
}

func (client *ApiClient) CompletedPlan(data data.NotiComplete) {
	targetURL := fmt.Sprintf("http://%s%s", client.MemVerifyUrl, client.MemVerifyPath)
	plantText := strconv.FormatUint(data.TargetId, 10) + ":" + strconv.FormatInt(int64(data.Type), 10)
	signature := _utils.GenerateHmacSHA256(plantText, client.Key)
	data.Signature = signature
	client.Client.R().
		// SetHeaderMultiValues(c.Request.Header).
		SetBody(data).
		Execute("POST", targetURL)
	// if err != nil {
	// 	return nil, err
	// }
}
