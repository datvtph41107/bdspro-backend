package main

import (
	_utils "common/utils"
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
