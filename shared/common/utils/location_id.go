package _utils

import (
	"errors"
	"strconv"
	"strings"
)

const (
	LocationIDKindProvince uint32 = 1
	LocationIDKindWard     uint32 = 2
)

var (
	errLocationIDRequired = errors.New("location id is required")
	errLocationIDInvalid  = errors.New("location id must be numeric with last digit 1 (province) or 2 (ward)")
)

func EncodeLocationID(sourceID uint64, kind uint32) (uint64, error) {
	if sourceID == 0 || (kind != LocationIDKindProvince && kind != LocationIDKindWard) {
		return 0, errLocationIDRequired
	}
	return (sourceID * 10) + uint64(kind), nil
}

func DecodeLocationID(encodedID uint64) (sourceID uint64, kind uint32, err error) {
	if encodedID == 0 {
		return 0, 0, errLocationIDRequired
	}
	kind = uint32(encodedID % 10)
	sourceID = encodedID / 10
	if sourceID == 0 || (kind != LocationIDKindProvince && kind != LocationIDKindWard) {
		return 0, 0, errLocationIDInvalid
	}
	return sourceID, kind, nil
}

func DecodeLocationIDFromString(value string) (sourceID uint64, kind uint32, err error) {
	encodedID, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, 0, errLocationIDInvalid
	}
	return DecodeLocationID(encodedID)
}

func GetLocationIDKind(encodedID uint64) (uint32, error) {
	_, kind, err := DecodeLocationID(encodedID)
	return kind, err
}

func IsEncodedLocationID(value uint64) bool {
	_, _, err := DecodeLocationID(value)
	return err == nil
}
