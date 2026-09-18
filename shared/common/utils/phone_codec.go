package _utils

import (
	"errors"
	"strings"
)

// PhoneDigitKey maps index 0..9 -> substitute char for that digit.
// Shared with clients that encode the phone before putting it in the URL.
// Characters are URL-path-safe and exclude digits so a raw phone cannot pass as encoded.
const PhoneDigitKey = "kQmZnXpWvR"

var (
	errPhoneCodecEmpty   = errors.New("phone codec: empty value")
	errPhoneCodecInvalid = errors.New("phone codec: invalid encoded phone")
	errPhoneCodecDigit   = errors.New("phone codec: phone must contain only digits")
)

var phoneDecodeMap = func() map[rune]byte {
	m := make(map[rune]byte, len(PhoneDigitKey))
	for i, ch := range PhoneDigitKey {
		m[ch] = byte('0' + i)
	}
	return m
}()

// EncodePhone replaces each digit 0-9 with the shared 10-character key.
// Example with PhoneDigitKey: 0961745160 -> kRpQWnXQpk.
func EncodePhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "", errPhoneCodecEmpty
	}

	var b strings.Builder
	b.Grow(len(phone))
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return "", errPhoneCodecDigit
		}
		b.WriteByte(PhoneDigitKey[ch-'0'])
	}
	return b.String(), nil
}

// DecodePhone reverses EncodePhone and returns the raw digit phone number.
func DecodePhone(encoded string) (string, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return "", errPhoneCodecEmpty
	}

	var b strings.Builder
	b.Grow(len(encoded))
	for _, ch := range encoded {
		digit, ok := phoneDecodeMap[ch]
		if !ok {
			return "", errPhoneCodecInvalid
		}
		b.WriteByte(digit)
	}
	return b.String(), nil
}
