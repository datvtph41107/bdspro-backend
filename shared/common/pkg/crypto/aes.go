// shared/common/crypto/aes.go
package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
)

// AESEncryptor implements Encryptor interface with AES-GCM
type AESEncryptor struct {
	key        []byte
	keyID      string
	keyVersion string
	mu         sync.RWMutex
}

// NewAESEncryptorFromPassphrase tạo encryptor từ passphrase
func NewAESEncryptorFromPassphrase(passphrase string) *AESEncryptor {
	hash := sha256.Sum256([]byte(passphrase))
	return &AESEncryptor{
		key:        hash[:],
		keyID:      generateKeyID(hash[:]),
		keyVersion: KeyVersion,
	}
}

// NewAESEncryptorFromKey tạo encryptor từ key có sẵn (32 bytes)
func NewAESEncryptorFromKey(key []byte) (*AESEncryptor, error) {
	if len(key) != AESKeySize {
		return nil, fmt.Errorf("key must be %d bytes, got %d", AESKeySize, len(key))
	}
	return &AESEncryptor{
		key:        key,
		keyID:      generateKeyID(key),
		keyVersion: KeyVersion,
	}, nil
}

// Encrypt implements Encryptor.Encrypt
func (e *AESEncryptor) Encrypt(ctx context.Context, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher for encryption: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create AES-GCM for encryption: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate AES-GCM nonce: %w", err)
	}

	version := []byte{0x01}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	result := append(version, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt implements Encryptor.Decrypt
func (e *AESEncryptor) Decrypt(ctx context.Context, encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("decode encrypted value from base64: %w", err)
	}
	if len(data) < 1 {
		return "", errors.New("encrypted data is too short")
	}

	version := data[0]
	if version != 0x01 {
		return "", fmt.Errorf("unsupported encryption version: %d", version)
	}
	ciphertext := data[1:]

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher for decryption: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create AES-GCM for decryption: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("encrypted ciphertext is too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt AES-GCM ciphertext: %w", err)
	}
	return string(plaintext), nil
}

func (e *AESEncryptor) EncryptBytes(ctx context.Context, plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher for byte encryption: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM for byte encryption: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate AES-GCM byte nonce: %w", err)
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (e *AESEncryptor) DecryptBytes(ctx context.Context, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher for byte decryption: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM for byte decryption: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("encrypted byte ciphertext is too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt AES-GCM byte ciphertext: %w", err)
	}
	return plaintext, nil
}

// EncryptBytesCTR mã hóa binary bằng AES-256-CTR: [nonce 16 bytes][ciphertext].
func (e *AESEncryptor) EncryptBytesCTR(ctx context.Context, plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher for CTR encryption: %w", err)
	}

	nonce := make([]byte, CTRNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate AES-CTR nonce: %w", err)
	}

	stream := cipher.NewCTR(block, nonce)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	out := make([]byte, CTRNonceSize+len(ciphertext))
	copy(out, nonce)
	copy(out[CTRNonceSize:], ciphertext)
	return out, nil
}

// DecryptBytesCTR giải mã binary AES-256-CTR: [nonce 16 bytes][ciphertext].
func (e *AESEncryptor) DecryptBytesCTR(ctx context.Context, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, nil
	}
	if len(ciphertext) < CTRNonceSize {
		return nil, errors.New("encrypted CTR ciphertext is too short")
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher for CTR decryption: %w", err)
	}

	nonce := ciphertext[:CTRNonceSize]
	data := ciphertext[CTRNonceSize:]

	stream := cipher.NewCTR(block, nonce)
	plaintext := make([]byte, len(data))
	stream.XORKeyStream(plaintext, data)
	return plaintext, nil
}

func (e *AESEncryptor) GetKeyID() string {
	return e.keyID
}

func generateKeyID(key []byte) string {
	hash := sha256.Sum256(key)
	return fmt.Sprintf("%x", hash[:8])
}
