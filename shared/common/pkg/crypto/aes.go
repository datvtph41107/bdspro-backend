// shared/common/crypto/aes.go
package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"sync"

	_errors "common/errors" // import errors package
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
		return nil, _errors.ReturnError(ErrCodeInvalidKeySize,
			fmt.Sprintf("key must be %d bytes, got %d", AESKeySize, len(key)))
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
		return "", _errors.ReturnError(ErrCodeEncryptionFailed,
			fmt.Sprintf("failed to create cipher: %v", err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", _errors.ReturnError(ErrCodeEncryptionFailed,
			fmt.Sprintf("failed to create GCM: %v", err))
	}

	// Tạo nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", _errors.ReturnError(ErrCodeEncryptionFailed,
			"failed to generate nonce")
	}

	// Mã hóa: format = [version(1 byte)][nonce(12 bytes)][ciphertext+tag]
	version := []byte{0x01} // version 1
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Gộp version + ciphertext
	result := append(version, ciphertext...)

	// Encode Base64
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt implements Encryptor.Decrypt
func (e *AESEncryptor) Decrypt(ctx context.Context, encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	// Decode Base64
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", _errors.ReturnError(ErrCodeDecodingFailed,
			fmt.Sprintf("base64 decode failed: %v", err))
	}

	if len(data) < 1 {
		return "", _errors.ReturnError(ErrCodeInvalidCiphertext,
			"data too short")
	}

	// Check version
	version := data[0]
	if version != 0x01 {
		return "", _errors.ReturnError(ErrCodeUnsupportedVersion,
			fmt.Sprintf("unsupported version: %d", version))
	}

	ciphertext := data[1:]

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", _errors.ReturnError(ErrCodeDecryptionFailed,
			fmt.Sprintf("failed to create cipher: %v", err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", _errors.ReturnError(ErrCodeDecryptionFailed,
			fmt.Sprintf("failed to create GCM: %v", err))
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", _errors.ReturnError(ErrCodeInvalidCiphertext,
			"ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", _errors.ReturnError(ErrCodeDecryptionFailed,
			fmt.Sprintf("decryption failed: %v", err))
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
		return nil, _errors.ReturnError(ErrCodeEncryptionFailed, err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, _errors.ReturnError(ErrCodeEncryptionFailed, err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, _errors.ReturnError(ErrCodeEncryptionFailed, "failed to generate nonce")
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
		return nil, _errors.ReturnError(ErrCodeDecryptionFailed, err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, _errors.ReturnError(ErrCodeDecryptionFailed, err.Error())
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, _errors.ReturnError(ErrCodeInvalidCiphertext, "ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
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
		return nil, _errors.ReturnError(ErrCodeEncryptionFailed, err.Error())
	}

	nonce := make([]byte, CTRNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, _errors.ReturnError(ErrCodeEncryptionFailed, "failed to generate ctr nonce")
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
		return nil, _errors.ReturnError(ErrCodeInvalidCiphertext, "ciphertext too short")
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, _errors.ReturnError(ErrCodeDecryptionFailed, err.Error())
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
	return fmt.Sprintf("%x", hash[:8]) // lấy 8 bytes đầu
}
