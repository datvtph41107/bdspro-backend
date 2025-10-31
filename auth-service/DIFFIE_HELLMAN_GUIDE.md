# Hướng dẫn sử dụng Diffie-Hellman Key Exchange trong Auth Service

## Tổng quan

Auth Service đã được tích hợp cơ chế **Diffie-Hellman Key Exchange** để bảo mật `auth_key` khi xác thực OTP. Khi user đăng ký và verify thành công, server sẽ sinh ra:
- **Cặp khóa DH** (private key + public key) 
- **Auth key** (khóa xác thực duy nhất)

Khi verify OTP, client gửi public key của mình lên, server sẽ:
1. Tính **shared secret** từ server private key và client public key
2. Mã hóa **auth_key** bằng shared secret
3. Trả về **encrypted auth_key** và **server public key** cho client

Client có thể dùng private key của mình và server public key để tính lại shared secret và giải mã auth_key.

---

## Luồng hoạt động

### 1. Đăng ký (RequestOTP)

```
Client                                    Server
  |                                          |
  |-------- POST /v2/auth/otp/request ----->|
  |  { phone: "+84901234567" }               |
  |                                          |
  |                                          | - Tạo AuthMethod mới
  |                                          | - Sinh DH key pair (private, public)
  |                                          | - Sinh auth_key
  |                                          | - Lưu vào database
  |                                          | - Gửi OTP
  |                                          |
  |<------ 200 OK --------------------------| 
  |  { authId: 123, message: "..." }         |
```

### 2. Xác thực (VerifyOTP) - Với Diffie-Hellman

```
Client                                             Server
  |                                                   |
  | 1. Generate DH key pair (client_private, client_public)
  |                                                   |
  |-------- POST /v2/auth/otp/verify -------------->|
  |  {                                                |
  |    phone: "+84901234567",                         |
  |    otp: "123456",                                 |
  |    clientPublicKey: "base64_encoded_key",         |
  |    platform: "MOBILE",                            |
  |    version: "1.0.0",                              |
  |    os: "iOS"                                      |
  |  }                                                |
  |                                                   |
  |                                                   | 2. Validate OTP
  |                                                   | 3. Decode client_public_key
  |                                                   | 4. Compute shared_secret = 
  |                                                   |    DH(server_private, client_public)
  |                                                   | 5. Encrypt auth_key with shared_secret
  |                                                   |
  |<------ 200 OK -----------------------------------| 
  |  {                                                |
  |    accessToken: "...",                            |
  |    refreshToken: "...",                           |
  |    encryptedAuthKey: "base64_encrypted_key",      |
  |    serverPublicKey: "base64_server_public_key",   |
  |    ...                                            |
  |  }                                                |
  |                                                   |
  | 6. Compute shared_secret = 
  |    DH(client_private, server_public)
  | 7. Decrypt auth_key using shared_secret
  | 8. Store auth_key securely in local storage
```

### 3. Đăng nhập lại (LoginWithKey) - Không cần OTP

```
Client                                             Server
  |                                                   |
  | 1. Read stored auth_key from secure storage
  | 2. Generate NEW DH key pair (client_private_new, client_public_new)
  | 3. Compute shared_secret_old = DH(client_private_old, server_public_old)
  | 4. Encrypt auth_key with shared_secret_old
  |                                                   |
  |-------- POST /v2/auth/login/key ---------------->|
  |  {                                                |
  |    phone: "+84901234567",                         |
  |    encryptedAuthKey: "base64_encrypted_key",      |
  |    clientPublicKey: "new_base64_client_public",   |
  |    platform: "MOBILE",                            |
  |    version: "1.0.0",                              |
  |    os: "iOS"                                      |
  |  }                                                |
  |                                                   |
  |                                                   | 5. Find auth_method by phone
  |                                                   | 6. Decode client_public_key_new
  |                                                   | 7. Compute shared_secret_old = 
  |                                                   |    DH(server_private, client_public_old)
  |                                                   | 8. Decrypt auth_key
  |                                                   | 9. Validate auth_key == stored_auth_key
  |                                                   | 10. Compute shared_secret_new =
  |                                                   |     DH(server_private, client_public_new)
  |                                                   | 11. Encrypt auth_key with shared_secret_new
  |                                                   |
  |<------ 200 OK -----------------------------------| 
  |  {                                                |
  |    accessToken: "...",                            |
  |    refreshToken: "...",                           |
  |    encryptedAuthKey: "new_base64_encrypted_key",  |
  |    serverPublicKey: "base64_server_public_key",   |
  |    ...                                            |
  |  }                                                |
  |                                                   |
  | 12. Compute shared_secret_new = 
  |     DH(client_private_new, server_public)
  | 13. Decrypt NEW encrypted auth_key
  | 14. Store new encrypted auth_key for next login
```

---

## API Endpoints

### 1. Request OTP (Đăng ký hoặc đăng nhập)

**Endpoint:** `POST /v2/auth/otp/request`

**Request Body:**
```json
{
  "phone": "+84901234567",
  "fullname": "Nguyễn Văn A" // Optional, chỉ cần khi đăng ký
}
```

**Response:**
```json
{
  "authId": 123,
  "message": "OTP đã được gửi",
  "existed": true // true: user đã tồn tại, false: user mới
}
```

**Server tạo:**
- Cặp khóa DH (private_key, public_key)
- Auth key ngẫu nhiên (32 bytes)
- Lưu vào database cột: `private_key`, `public_key`, `auth_key`

---

### 2. Verify OTP (Xác thực)

**Endpoint:** `POST /v2/auth/otp/verify`

**Request Body:**
```json
{
  "phone": "+84901234567",
  "otp": "123456",
  "authId": 123,
  "platform": "MOBILE",
  "version": "1.0.0",
  "os": "iOS",
  "deviceName": "iPhone 13 Pro",
  "clientPublicKey": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMII..." // Base64 encoded
}
```

**Response:**
```json
{
  "authId": 123,
  "profileId": 456,
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "fullName": "Nguyễn Văn A",
  "role": "USER",
  "encryptedAuthKey": "nK7x9pQ2vL4mR8tY6wZ3sA==", // Base64 encoded
  "serverPublicKey": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMII...", // Base64 encoded
  "referralCode": "ABC123",
  "requireFullname": false
}
```

**Giải thích các trường mới:**
- `clientPublicKey`: Client gửi lên public key của mình (base64 encoded)
- `encryptedAuthKey`: Server trả về auth_key đã được mã hóa bằng shared secret (AES-256-GCM)
- `serverPublicKey`: Public key của server để client tính shared secret

---

### 3. Login with Key (Đăng nhập bằng Auth Key)

**Endpoint:** `POST /v2/auth/login/key`

**Mục đích:** Đăng nhập nhanh không cần OTP bằng auth_key đã lưu

**Request Body:**
```json
{
  "phone": "+84901234567",
  "encryptedAuthKey": "nK7x9pQ2vL4mR8tY6wZ3sA==", 
  "clientPublicKey": "NEW_MIIBIjANBgkqhkiG9w0...", 
  "platform": "MOBILE",
  "version": "1.0.0",
  "os": "iOS",
  "deviceName": "iPhone 13 Pro"
}
```

**Response:**
```json
{
  "authId": 123,
  "profileId": 456,
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "fullName": "Nguyễn Văn A",
  "role": "USER",
  "encryptedAuthKey": "pX9tQ3rL5nR7sY8wZ4aB==", 
  "serverPublicKey": "MIIBIjANBgkqhkiG9w0...", 
  "referralCode": "ABC123",
  "requireFullname": false
}
```

**Giải thích trường:**
- **Request `encryptedAuthKey`**: Auth key đã được mã hóa từ lần verify OTP trước hoặc lần login trước
- **Request `clientPublicKey`**: Public key **MỚI** sinh ra cho lần login này  
- **Response `encryptedAuthKey`**: Auth key được mã hóa **MỚI** với `clientPublicKey` mới để dùng cho lần login sau
- **Response `serverPublicKey`**: Server public key (giống lần verify, dùng để verify signature nếu cần)

**Lưu ý quan trọng:**
- ✅ Mỗi lần login, client phải sinh DH key pair **MỚI** 
- ✅ Client dùng shared secret **CŨ** (từ lần trước) để mã hóa auth_key gửi lên
- ✅ Server trả về auth_key được mã hóa bằng shared secret **MỚI**
- ✅ Client lưu encrypted auth_key MỚI này cho lần login sau
- ⚠️ Không lưu trữ private key - chỉ giữ trong memory khi sử dụng

---

## Cách tích hợp cho Client

### A. Client-side Implementation (JavaScript/TypeScript)

#### 1. Generate Client DH Key Pair

```typescript
import { subtle } from 'crypto';

// Diffie-Hellman parameters (RFC 3526 - 2048-bit MODP Group)
const DH_PRIME = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF";
const DH_GENERATOR = 2;

async function generateDHKeyPair() {
  // Generate random private key
  const privateKeyBytes = new Uint8Array(32);
  crypto.getRandomValues(privateKeyBytes);
  const privateKey = BigInt('0x' + Array.from(privateKeyBytes).map(b => b.toString(16).padStart(2, '0')).join(''));
  
  // Calculate public key: g^private mod p
  const prime = BigInt('0x' + DH_PRIME);
  const generator = BigInt(DH_GENERATOR);
  const publicKey = modPow(generator, privateKey, prime);
  
  return {
    privateKey: privateKey.toString(16),
    publicKey: publicKey.toString(16),
  };
}

function modPow(base: bigint, exp: bigint, mod: bigint): bigint {
  let result = 1n;
  base = base % mod;
  while (exp > 0n) {
    if (exp % 2n === 1n) result = (result * base) % mod;
    exp = exp >> 1n;
    base = (base * base) % mod;
  }
  return result;
}

// Encode to base64
function encodePublicKey(publicKey: string): string {
  const bytes = hexToBytes(publicKey);
  return btoa(String.fromCharCode(...bytes));
}
```

#### 2. Send VerifyOTP Request

```typescript
async function verifyOTP(phone: string, otp: string, clientPublicKey: string) {
  const response = await fetch('/v2/auth/otp/verify', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      phone,
      otp,
      platform: 'MOBILE',
      version: '1.0.0',
      os: 'iOS',
      deviceName: 'iPhone 13 Pro',
      clientPublicKey, // Base64 encoded
    }),
  });
  
  return await response.json();
}
```

#### 3. Decrypt Auth Key

```typescript
import CryptoJS from 'crypto-js';

async function decryptAuthKey(
  encryptedAuthKey: string,
  clientPrivateKey: string,
  serverPublicKey: string
): Promise<string> {
  // 1. Compute shared secret
  const serverPubKeyBigInt = BigInt('0x' + base64ToHex(serverPublicKey));
  const clientPrivKeyBigInt = BigInt('0x' + clientPrivateKey);
  const prime = BigInt('0x' + DH_PRIME);
  
  const sharedSecret = modPow(serverPubKeyBigInt, clientPrivKeyBigInt, prime);
  
  // 2. Derive encryption key (SHA-256 of shared secret)
  const sharedSecretHex = sharedSecret.toString(16);
  const encryptionKey = CryptoJS.SHA256(CryptoJS.enc.Hex.parse(sharedSecretHex));
  
  // 3. Decrypt using AES-256-GCM
  const ciphertext = CryptoJS.enc.Base64.parse(encryptedAuthKey);
  const nonce = CryptoJS.lib.WordArray.create(ciphertext.words.slice(0, 3)); // 12 bytes nonce
  const encrypted = CryptoJS.lib.WordArray.create(ciphertext.words.slice(3));
  
  const decrypted = CryptoJS.AES.decrypt(
    { ciphertext: encrypted },
    encryptionKey,
    {
      iv: nonce,
      mode: CryptoJS.mode.GCM,
      padding: CryptoJS.pad.NoPadding,
    }
  );
  
  return decrypted.toString(CryptoJS.enc.Utf8);
}
```

#### 4. Complete Flow

```typescript
async function loginWithOTP(phone: string, otp: string) {
  // 1. Generate DH key pair
  const { privateKey, publicKey } = await generateDHKeyPair();
  const clientPublicKeyBase64 = encodePublicKey(publicKey);
  
  // 2. Verify OTP
  const response = await verifyOTP(phone, otp, clientPublicKeyBase64);
  
  // 3. Decrypt auth key
  const authKey = await decryptAuthKey(
    response.encryptedAuthKey,
    privateKey,
    response.serverPublicKey
  );
  
  // 4. Store securely
  localStorage.setItem('accessToken', response.accessToken);
  localStorage.setItem('refreshToken', response.refreshToken);
  localStorage.setItem('authKey', authKey); // Store encrypted in production
  
  console.log('Login successful, auth_key:', authKey);
}
```

---

### B. Client-side Implementation (Flutter/Dart)

#### 1. Add dependencies

```yaml
dependencies:
  pointycastle: ^3.7.3 # For DH key exchange
  encrypt: ^5.0.1      # For AES encryption
```

#### 2. Generate DH Key Pair

```dart
import 'dart:typed_data';
import 'dart:convert';
import 'package:pointycastle/export.dart';

class DHKeyExchange {
  static final BigInt _dhPrime = BigInt.parse(
    'FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF',
    radix: 16,
  );
  static final BigInt _dhGenerator = BigInt.from(2);

  static Future<Map<String, String>> generateKeyPair() async {
    // Generate random private key
    final random = SecureRandom('Fortuna');
    final seed = Uint8List.fromList(List.generate(32, (_) => random.nextUint8()));
    random.seed(KeyParameter(seed));
    
    final privateKey = random.nextBigInteger(256);
    
    // Calculate public key: g^private mod p
    final publicKey = _dhGenerator.modPow(privateKey, _dhPrime);
    
    return {
      'privateKey': privateKey.toRadixString(16),
      'publicKey': publicKey.toRadixString(16),
    };
  }

  static String encodePublicKey(String publicKeyHex) {
    final bytes = _hexToBytes(publicKeyHex);
    return base64Encode(bytes);
  }

  static BigInt computeSharedSecret(String privateKeyHex, String serverPublicKeyBase64) {
    final privateKey = BigInt.parse(privateKeyHex, radix: 16);
    final serverPublicKeyBytes = base64Decode(serverPublicKeyBase64);
    final serverPublicKey = BigInt.parse(
      serverPublicKeyBytes.map((e) => e.toRadixString(16).padLeft(2, '0')).join(),
      radix: 16,
    );
    
    return serverPublicKey.modPow(privateKey, _dhPrime);
  }

  static Uint8List _hexToBytes(String hex) {
    final result = Uint8List(hex.length ~/ 2);
    for (int i = 0; i < hex.length; i += 2) {
      result[i ~/ 2] = int.parse(hex.substring(i, i + 2), radix: 16);
    }
    return result;
  }
}
```

#### 3. Decrypt Auth Key

```dart
import 'package:encrypt/encrypt.dart';
import 'package:crypto/crypto.dart';

Future<String> decryptAuthKey(
  String encryptedAuthKeyBase64,
  String clientPrivateKeyHex,
  String serverPublicKeyBase64,
) async {
  // 1. Compute shared secret
  final sharedSecret = DHKeyExchange.computeSharedSecret(
    clientPrivateKeyHex,
    serverPublicKeyBase64,
  );
  
  // 2. Derive encryption key (SHA-256 of shared secret)
  final sharedSecretBytes = _bigIntToBytes(sharedSecret);
  final encryptionKey = sha256.convert(sharedSecretBytes).bytes;
  
  // 3. Decrypt using AES-256-GCM
  final ciphertext = base64Decode(encryptedAuthKeyBase64);
  final nonce = ciphertext.sublist(0, 12); // First 12 bytes
  final encrypted = ciphertext.sublist(12);
  
  final key = Key(Uint8List.fromList(encryptionKey));
  final iv = IV(Uint8List.fromList(nonce));
  final encrypter = Encrypter(AES(key, mode: AESMode.gcm));
  
  final decrypted = encrypter.decrypt(Encrypted(encrypted), iv: iv);
  return decrypted;
}

Uint8List _bigIntToBytes(BigInt bigInt) {
  final hex = bigInt.toRadixString(16);
  final paddedHex = hex.length % 2 == 0 ? hex : '0$hex';
  return Uint8List.fromList(
    List.generate(
      paddedHex.length ~/ 2,
      (i) => int.parse(paddedHex.substring(i * 2, i * 2 + 2), radix: 16),
    ),
  );
}
```

---

## Bảo mật

### 1. Lưu trữ Private Key
- **Không bao giờ** lưu private key lên server
- Lưu trong bộ nhớ tạm (memory) hoặc secure storage (Keychain/Keystore)
- Xóa sau khi tính shared secret

### 2. Lưu trữ Auth Key
- Mã hóa bằng device-specific key trước khi lưu vào local storage
- Sử dụng Keychain (iOS) hoặc EncryptedSharedPreferences (Android)

### 3. Shared Secret
- Không lưu trữ shared secret
- Tính toán mỗi khi cần dùng

### 4. Transport Layer
- Luôn dùng HTTPS/TLS để bảo vệ dữ liệu truyền tải
- Diffie-Hellman là lớp bảo mật bổ sung, không thay thế TLS

---

## Testing

### Test Diffie-Hellman Flow

```bash
# 1. Request OTP
curl -X POST http://localhost:8080/v2/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"phone": "+84901234567"}'

# Response: { "authId": 123, "message": "OTP đã được gửi" }

# 2. Generate client key pair (using utility script)
# clientPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMII..."

# 3. Verify OTP with client public key
curl -X POST http://localhost:8080/v2/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+84901234567",
    "otp": "123456",
    "platform": "MOBILE",
    "version": "1.0.0",
    "os": "iOS",
    "clientPublicKey": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMII..."
  }'

# Response sẽ có: encryptedAuthKey và serverPublicKey
```

---

## Troubleshooting

### Lỗi "Failed to decode client public key"
- Kiểm tra client public key có đúng format base64 không
- Đảm bảo sử dụng đúng DH parameters (RFC 3526 - 2048-bit MODP Group)

### Lỗi "Failed to compute shared secret"
- Kiểm tra client public key nằm trong khoảng hợp lệ (1 < key < prime)
- Đảm bảo private key được sinh ngẫu nhiên đúng cách

### Lỗi "Failed to decrypt auth key"
- Kiểm tra shared secret được tính đúng cả client và server
- Đảm bảo sử dụng đúng thuật toán AES-256-GCM
- Kiểm tra nonce được extract đúng (12 bytes đầu tiên)

---

## Liên hệ

Nếu có thắc mắc hoặc cần hỗ trợ, vui lòng liên hệ team Backend.

