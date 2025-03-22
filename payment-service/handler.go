package main

import (
	_utils "common/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

type PathPrefix struct {
	Prefix string
	Path   string
}

type SignatureDecode struct {
	Signature   string
	ProfileID   int64
	ExpiredTime int64
}

func SplitPathEncode(input string) PathPrefix {
	if len(input) == 0 {
		return PathPrefix{"", ""}
	}
	return PathPrefix{
		Prefix: string(input[0]),
		Path:   input[1:],
	}
}

func SignDecode(key, signEncode string) *SignatureDecode {
	xorDecoded, err := _utils.XorDecode(signEncode, key)
	if err != nil {

	}
	split := strings.Split(xorDecoded, ":")

	if len(split) < 3 {
		return nil
	}

	profileID, err1 := strconv.ParseInt(split[1], 10, 64)
	expiredTime, err2 := strconv.ParseInt(split[2], 10, 64)
	if err1 != nil || err2 != nil {
		return nil
	}

	return &SignatureDecode{
		Signature:   split[0],
		ProfileID:   profileID,
		ExpiredTime: expiredTime,
	}
}

func SignEncode(key string, profileID, expiredTime int64, signature string) string {
	data := signature + ":" + strconv.FormatInt(profileID, 10) + ":" + strconv.FormatInt(expiredTime, 10)
	return _utils.XorEncode(data, key)
}

// XorEncode và XorDecode sẽ được định nghĩa sau

func sortObjDataByKey(data map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedData := make(map[string]interface{})
	for _, k := range keys {
		sortedData[k] = data[k]
	}

	return sortedData
}

func convertObjToQueryStr(data map[string]interface{}) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryStr string
	for i, k := range keys {
		v := data[k]
		var valueStr string
		switch v := v.(type) {
		case string:
			valueStr = v
		case float64:
			valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			valueStr = strconv.Itoa(v)
		case nil:
			valueStr = ""
		default:
			jsonVal, _ := json.Marshal(v)
			valueStr = string(jsonVal)
		}
		if i > 0 {
			queryStr += "&"
		}
		queryStr += k + "=" + valueStr
	}
	return queryStr
}

func IsValidData(data map[string]interface{}, currentSignature, key string) bool {
	sortedData := sortObjDataByKey(data)
	dataQueryStr := convertObjToQueryStr(sortedData)

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(dataQueryStr))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature == currentSignature
}
