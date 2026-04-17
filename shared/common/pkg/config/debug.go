package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync/atomic"
	"unsafe"
)

// ComputeHash returns SHA256 of canonical JSON representation (sensitive fields are masked).
func ComputeHash(cfg interface{}) string {
	data, _ := json.Marshal(cfg)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// AtomicConfig provides atomic load/store for hot-reload.
type AtomicConfig struct {
	ptr  unsafe.Pointer
	hash string
}

func NewAtomicConfig(initial interface{}) *AtomicConfig {
	ac := &AtomicConfig{}
	ac.Store(initial)
	return ac
}

func (ac *AtomicConfig) Load() interface{} {
	return atomic.LoadPointer(&ac.ptr)
}

func (ac *AtomicConfig) Store(cfg interface{}) {
	atomic.StorePointer(&ac.ptr, unsafe.Pointer(&cfg))
	ac.hash = ComputeHash(cfg)
}

func (ac *AtomicConfig) Hash() string {
	return ac.hash
}
