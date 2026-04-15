// common/quota/mock/data.go
package mock

import (
    "context"
    "errors"
    "sync"
)

type FakeParcel struct {
    ID          string
    Coordinates string
    Area        float64
}

type FakeProfile struct {
    UserID uint64
    Name   string
    Phone  string
}

var (
    parcels = map[string]*FakeParcel{
        "PARCEL_001": {ID: "001", Coordinates: "10.0,20.0", Area: 100.5},
        "PARCEL_002": {ID: "002", Coordinates: "10.1,20.1", Area: 200.0},
        "PARCEL_003": {ID: "003", Coordinates: "10.2,20.2", Area: 150.0},
    }
    profiles = map[uint64]*FakeProfile{
        1: {UserID: 1, Name: "Nguyen Van A", Phone: "0901234567"},
        2: {UserID: 2, Name: "Tran Thi B", Phone: "0901234568"},
    }
    mu sync.RWMutex
)

// GetParcel trả về parcel hoặc lỗi mô phỏng
func GetParcel(ctx context.Context, id string) (*FakeParcel, error) {
    mu.RLock()
    defer mu.RUnlock()
    // Mô phỏng lỗi DB
    if id == "DB_ERROR" {
        return nil, errors.New("database connection timeout")
    }
    p, ok := parcels[id]
    if !ok {
        return nil, nil // not found
    }
    return p, nil
}

// GetProfile tương tự
func GetProfile(ctx context.Context, userID uint64) (*FakeProfile, error) {
    mu.RLock()
    defer mu.RUnlock()
    if userID == 999 {
        return nil, nil
    }
    p, ok := profiles[userID]
    if !ok {
        return nil, nil
    }
    return p, nil
}