package _utils

import (
	"fmt"
	"testing"
)

func TestWithCustomPublicPrivateKeys(t *testing.T) {
	// TRUYỀN VÀO 2 BIẾN NÀY
	var myPrivateKeyEncoded string = "r7PNrlR2bJxGfnDh+p0mojad4GErP7YVRQmGrjLzofHRvhCaWN7WwdLU25mPOeQkNzSweafH3t9a1Q3JZ5ULE09pyf5DRmxrBhj57Lpzj8oFrw4o0QEA4RPdO4RYKjyrceJzWNa6wJ/9oFH+J7NXSe2xWDSGppt4AhM1NbWjulcrqF3ol6oS1BBWf5hzGuZ57ndw/hy+TPMVN7eQUZUIBvMD9kVbI7OsVpIPHROiiDP6pw3brZIE7BbiwJR34RKMYw/rV+chHrIkyOtAQl7OBaZRUB1SUIpwQS9n7m9P7St5xIdqecvJhQEpYVDkMonGJv2lahN/cIxz3+xrbZEZ/A=="  // Điền private key (base64) của bạn vào đây
	var peerPublicKeyEncoded string = "F2picKa/hfpiEEglej2U5ZBFjevnc2LoDJ+AI82CZc38HCIrE9OKripc4JR6fz14f6xsFhuJUDp+/aFg/KFA95rpqM2TiZY0wNlDNV/sb84F1B4d5dXu+6dvtnIYh33rDqq+niMVsO09k5fNVbh6jebIjhkZw3w5KJ/2sYKXbvxv3V9Wpd7AzWpPJB1nbyDT6POCu0GKCsc3UAsQqP66unsieXOol0OZrFLa9raEcMqbDc0/t97sdupfAx3IJBqSOa0vlC9SOh0L46LZ/2x2How+N7av9+NeYrmzcs+s0YJYUczK12QGEMQdbbUV3L3bzr2IuVvn9o7hHSsOcDB0Yg==" // Điền public key (base64) của peer vào đây
	var authKey string = "659c1838fa48bc0d104b09bea37fc0c9f302c4c7a43b8f090913a69ef40fcec8"

	// Decode keys
	myPrivateKey, err := DecodePrivateKey(myPrivateKeyEncoded)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	fmt.Printf("\n=====My Private Key: %s\n", myPrivateKey.Text(16))

	peerPublicKey, err := DecodePublicKey(peerPublicKeyEncoded)
	if err != nil {
		t.Fatalf("Failed to decode peer public key: %v", err)
	}

	fmt.Printf("\n=====Peer Public Key: %s\n", peerPublicKey.Text(16))
	if err != nil {
		t.Fatalf("Failed to decode peer public key: %v", err)
	}

	// Compute shared secret
	sharedSecret, err := ComputeSharedSecret(myPrivateKey, peerPublicKey)
	if err != nil {
		t.Fatalf("Failed to compute shared secret: %v", err)
	}

	// Derive encryption key
	// encryptionKey := sharedSecret
	sharedKeyStr := sharedSecret.Text(16)
	fmt.Printf("\n=======Encryption Key: %x\n\n", sharedSecret)

	// Mã hóa authKey bằng shared secret
	fmt.Printf("Original Auth Key: %s\n", authKey)
	encryptedAuthKey, err := EncryptAuthKeyXOR(authKey, sharedKeyStr)
	if err != nil {
		t.Fatalf("Failed to encrypt auth key: %v", err)
	}
	fmt.Printf("Encrypted Auth Key: %s\n\n", encryptedAuthKey)

	// Giải mã để verify
	decryptedAuthKey, err := DecryptAuthKeyXOR(encryptedAuthKey, sharedKeyStr)
	if err != nil {
		t.Fatalf("Failed to decrypt auth key: %v", err)
	}
	fmt.Printf("Decrypted Auth Key: %s\n", decryptedAuthKey)

	// Verify
	if authKey == decryptedAuthKey {
		fmt.Println("✅ Encryption/Decryption successful!")
	} else {
		t.Fatal("❌ Decryption failed!")
	}
}
