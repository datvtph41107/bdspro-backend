package _utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Diffie-Hellman parameters (sử dụng RFC 3526 - 2048-bit MODP Group)
var (
	// Prime modulus (2048-bit)
	dhPrime = mustParseBigInt("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF", 16)

	// Generator
	dhGenerator = big.NewInt(2)
)

// mustParseBigInt parses a hex string to big.Int, panic if failed
func mustParseBigInt(s string, base int) *big.Int {
	n := new(big.Int)
	n.SetString(s, base)
	return n
}

// DHKeyPair represents a Diffie-Hellman key pair
type DHKeyPair struct {
	PrivateKey *big.Int
	PublicKey  *big.Int
}

// GenerateDHKeyPair generates a new Diffie-Hellman key pair
func GenerateDHKeyPair() (*DHKeyPair, error) {
	// Generate random private key (256 bits for security)
	privateKey, err := rand.Int(rand.Reader, dhPrime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Calculate public key: g^private mod p
	publicKey := new(big.Int).Exp(dhGenerator, privateKey, dhPrime)

	return &DHKeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// ComputeSharedSecret computes the shared secret using private key and peer's public key
func ComputeSharedSecret(privateKey *big.Int, peerPublicKey *big.Int) (*big.Int, error) {
	if peerPublicKey.Cmp(big.NewInt(1)) <= 0 || peerPublicKey.Cmp(dhPrime) >= 0 {
		return nil, fmt.Errorf("invalid peer public key")
	}

	// Calculate shared secret: peerPublicKey^private mod p
	sharedSecret := new(big.Int).Exp(peerPublicKey, privateKey, dhPrime)
	return sharedSecret, nil
}

// EncodePublicKey encodes a public key to base64 string
func EncodePublicKey(publicKey *big.Int) string {
	return base64.StdEncoding.EncodeToString(publicKey.Bytes())
}

// DecodePublicKey decodes a base64 string to public key
func DecodePublicKey(encoded string) (*big.Int, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	publicKey := new(big.Int).SetBytes(bytes)
	return publicKey, nil
}

// EncodePrivateKey encodes a private key to base64 string
func EncodePrivateKey(privateKey *big.Int) string {
	return base64.StdEncoding.EncodeToString(privateKey.Bytes())
}

// DecodePrivateKey decodes a base64 string to private key
func DecodePrivateKey(encoded string) (*big.Int, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	privateKey := new(big.Int).SetBytes(bytes)
	return privateKey, nil
}

// GenerateAuthKey generates a random auth key (32 bytes)
func GenerateAuthKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate auth key: %w", err)
	}
	return hex.EncodeToString(key), nil
}

// DeriveEncryptionKey derives an AES-256 key from shared secret
func DeriveEncryptionKey(sharedSecret *big.Int) []byte {
	hash := sha256.Sum256(sharedSecret.Bytes())
	return hash[:]
}

// EncryptAuthKey encrypts the auth key using AES-256-GCM with the shared secret

// EncryptAuthKeyXOR: XOR authKey với sharedSecretHex (đã là 32-byte SHA256)
func EncryptAuthKeyXOR(authKeyHex string, sharedSecretHex string) (string, error) {
	// Convert authKey từ hex → bytes
	authBytes, err := hex.DecodeString(authKeyHex)
	if err != nil {
		return "", err
	}

	// Convert sharedSecret từ hex → bytes (SHA256 output, 32 bytes)
	keyBytes, err := hex.DecodeString(sharedSecretHex)
	if err != nil {
		return "", err
	}

	// XOR byte từng byte
	result := make([]byte, len(authBytes))
	for i := 0; i < len(authBytes); i++ {
		result[i] = authBytes[i] ^ keyBytes[i%len(keyBytes)]
	}

	// Encode Base64 → return
	return base64.StdEncoding.EncodeToString(result), nil
}

// Giải mã cũng giống hệt vì XOR đảo ngược chính nó
func DecryptAuthKeyXOR(encryptedBase64 string, sharedSecretHex string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	keyBytes, err := hex.DecodeString(sharedSecretHex)
	if err != nil {
		return "", err
	}

	result := make([]byte, len(cipherBytes))
	for i := 0; i < len(cipherBytes); i++ {
		result[i] = cipherBytes[i] ^ keyBytes[i%len(keyBytes)]
	}

	return hex.EncodeToString(result), nil
}
