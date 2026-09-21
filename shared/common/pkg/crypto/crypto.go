// shared/common/crypto/manager.go
package crypto

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type DefaultKeyManager struct {
	activeKey     []byte
	activeKeyID   string
	activeVersion string
	mu            sync.RWMutex
}

var (
	globalManager     *DefaultKeyManager
	globalManagerOnce sync.Once
)

// GetGlobalKeyManager returns singleton key manager
func GetGlobalKeyManager() *DefaultKeyManager {
	globalManagerOnce.Do(func() {
		globalManager = &DefaultKeyManager{
			activeVersion: "v1",
		}
	})
	return globalManager
}

// InitializeFromConfig khởi tạo từ config
func (m *DefaultKeyManager) InitializeFromConfig(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	passphrase := viper.GetString("encryption.aes_passphrase")
	if passphrase == "" {
		passphrase = os.Getenv("CONFIG_ENCRYPTION_PASSPHRASE")
	}
	if passphrase == "" {
		return errors.New("encryption passphrase is required (set in config or env)")
	}

	hash := sha256.Sum256([]byte(passphrase))
	m.activeKey = hash[:]
	m.activeKeyID = generateKeyID(m.activeKey)
	return nil
}

// InitializeWithPassphrase khởi tạo với passphrase trực tiếp
func (m *DefaultKeyManager) InitializeWithPassphrase(ctx context.Context, passphrase string) error {
	if passphrase == "" {
		return errors.New("passphrase is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	hash := sha256.Sum256([]byte(passphrase))
	m.activeKey = hash[:]
	m.activeKeyID = generateKeyID(m.activeKey)
	return nil
}

// GetActiveKey implements KeyManager
func (m *DefaultKeyManager) GetActiveKey(ctx context.Context) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeKey == nil {
		return nil, errors.New("no active encryption key")
	}
	return m.activeKey, nil
}

// GetActiveEncryptor returns encryptor with active key
func (m *DefaultKeyManager) GetActiveEncryptor(ctx context.Context) (*AESEncryptor, error) {
	key, err := m.GetActiveKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get active encryption key: %w", err)
	}
	encryptor, err := NewAESEncryptorFromKey(key)
	if err != nil {
		return nil, fmt.Errorf("create active encryptor: %w", err)
	}
	return encryptor, nil
}

// RotateKey implements KeyManager
func (m *DefaultKeyManager) RotateKey(ctx context.Context, newKey []byte) error {
	if len(newKey) != AESKeySize {
		return fmt.Errorf("new key must be %d bytes, got %d", AESKeySize, len(newKey))
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO: Lưu key cũ vào backup để decrypt data cũ nếu cần
	m.activeKey = newKey
	m.activeKeyID = generateKeyID(newKey)
	return nil
}

func (m *DefaultKeyManager) GetKeyVersion() string {
	return m.activeVersion
}
