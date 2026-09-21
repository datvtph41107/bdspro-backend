// shared/common/crypto/interface.go
package crypto

import "context"

// Encryptor interface định nghĩa các phương thức mã hóa/giải mã
type Encryptor interface {
	// Encrypt mã hóa plaintext thành encrypted string (Base64)
	Encrypt(ctx context.Context, plaintext string) (string, error)

	// Decrypt giải mã encrypted string thành plaintext
	Decrypt(ctx context.Context, encrypted string) (string, error)

	// EncryptBytes mã hóa byte slice
	EncryptBytes(ctx context.Context, plaintext []byte) ([]byte, error)

	// DecryptBytes giải mã byte slice
	DecryptBytes(ctx context.Context, ciphertext []byte) ([]byte, error)

	// EncryptBytesCTR / DecryptBytesCTR — AES-256-CTR cho binary (tile); format: [nonce 16][ciphertext].
	EncryptBytesCTR(ctx context.Context, plaintext []byte) ([]byte, error)
	DecryptBytesCTR(ctx context.Context, ciphertext []byte) ([]byte, error)

	// GetKeyID trả về ID của key đang dùng (cho audit)
	GetKeyID() string
}

// KeyManager interface quản lý encryption keys
type KeyManager interface {
	// GetActiveKey lấy key đang active
	GetActiveKey(ctx context.Context) ([]byte, error)

	// RotateKey thay đổi key (cho phép rotate mà không downtime)
	RotateKey(ctx context.Context, newKey []byte) error

	// GetKeyVersion lấy version hiện tại
	GetKeyVersion() string
}

const (
	AESKeySize   = 32 // AES-256
	NonceSize    = 12 // GCM recommended nonce size
	CTRNonceSize = 16 // AES block size — IV cho CTR (tile)
	KeyVersion   = "v1"
)
