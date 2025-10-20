# Recent Changes & Updates

## 📋 Change Log

**Last Updated**: October 17, 2025
**Status**: ✅ Active Development

---

## 🔄 Latest Changes

### ✅ October 17, 2025 - Enhanced OTP Request API

**Feature**: Add `existed` field to OTP request response

**Affected Service**: Auth Service

**Changes Made**:

#### 1. Protobuf Schema Update
**File**: `shared/protobuf/schema/auth/auth_method.proto`

```protobuf
message RequestOTPResponse {
    uint64 authId = 1;
    string message = 2;
    bool existed = 3;  // NEW: true if user exists, false if new registration
}
```

#### 2. DTO Update
**File**: `auth-service/internal/dto/auth_extended.go`

```go
type RequestLoginResponse struct {
	AuthID   uint64 `json:"authId"`
	Provider string `json:"provider"`
	OAuthID  string `json:"oauthId"`
	Existed  bool   `json:"existed"` // NEW field
}
```

#### 3. Usecase Logic Update
**File**: `auth-service/internal/usecase/auth_user_usecase.go`

```go
func (s *AuthUsecase) RequestOtp(c context.Context, otp dto.OtpRequest) (*dto.RequestLoginResponse, error) {
    // ... validation ...
    
    existed := false // Track user existence
    
    oauth, err := s.AuthMethodRepo.FindByPhone(c, otp.Phone)
    if err == nil && oauth != nil {
        // User already exists
        existed = true
        // ... send OTP ...
    } else {
        // New user registration
        existed = false
        // ... create new auth method ...
    }
    
    return &dto.RequestLoginResponse{
        AuthID:   oauth.ID,
        Provider: "PHONE",
        OAuthID:  otp.Phone,
        Existed:  existed,  // Return existence status
    }, nil
}
```

#### 4. Handler Update
**File**: `auth-service/infra/handler/auth_handler.go`

```go
func (h *AuthHandler) RequestOTP(ctx context.Context, req *authpb.RequestOTPRequest) (*authpb.RequestOTPResponse, error) {
    result, err := h.AuthUsecase.RequestOtp(ctx, *dto)
    if err != nil {
        return nil, err
    }

    response := &authpb.RequestOTPResponse{
        AuthId:  result.AuthID,
        Existed: result.Existed,  // NEW field
    }

    return response, nil
}
```

**Impact**:
- ✅ Frontend can now detect new vs. returning users
- ✅ Can show different UI flow for registration vs. login
- ✅ Better UX with conditional messaging
- ✅ Analytics can track new user registrations

**API Endpoint**: `POST /v2/auth/otp/request`

**Request**:
```json
{
  "phone": "0901234567",
  "fullname": "Nguyen Van A"  // Required for new users
}
```

**Response** (Existing User):
```json
{
  "authId": 123,
  "message": "OTP đã được gửi",
  "existed": true
}
```

**Response** (New User):
```json
{
  "authId": 456,
  "message": "Tài khoản mới đã được tạo",
  "existed": false
}
```

**Frontend Usage**:
```typescript
const response = await requestOTP(phone, fullname);

if (response.existed) {
    // Show: "Chào mừng quay lại! Vui lòng nhập OTP"
    showLoginFlow();
} else {
    // Show: "Chào mừng bạn đến với BDSPro! Vui lòng nhập OTP để hoàn tất đăng ký"
    showRegistrationFlow();
}
```

**Testing**:
```bash
# Test with existing user
curl -X POST http://localhost:8080/v2/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567", "fullname": "Test User"}'

# Expected: existed = true (if user exists)
# Expected: existed = false (if new user)
```

---

## 🚀 Deployment Notes

**Steps to Deploy**:
1. ✅ Protobuf schema updated
2. ✅ Protobuf code generated
3. ✅ DTO updated
4. ✅ Usecase logic updated
5. ✅ Handler updated
6. ⏳ Ready for deployment

**No Breaking Changes**: This is a backwards-compatible addition (new field only)

**Migration Required**: No

**Database Changes**: No

**Service Restart Required**: Yes (Auth Service only)

---

**Change Type**: Feature Enhancement
**Priority**: Medium
**Risk Level**: Low
**Backward Compatible**: Yes ✅
**Tested**: Pending
**Deployed**: Pending

---

## ✅ October 17, 2025 - Enhanced OTP Verify Logic

**Feature**: Mandatory fullname validation for new user registration

**Affected Service**: Auth Service

**Changes Made**:

#### 1. Usecase Logic Update
**File**: `auth-service/internal/usecase/auth_user_usecase.go`

```go
func (s *AuthUsecase) VerifyOtp(c context.Context, otp dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
    // ... OTP validation logic ...
    
    infoEntity, err := s.ProfileProvider.GetByProfileID(c, e.UserID)
    isNewUser := false
    if err != nil || infoEntity == nil || infoEntity.ProfileID == 0 {
        // NEW: Kiểm tra nếu là user mới thì bắt buộc phải có họ tên
        if otp.Fullname == "" {
            return nil, _errors.ReturnError(400, "Vui lòng nhập họ tên khi đăng ký tài khoản mới")
        }
        
        // Tạo profile với fullname từ request
        infoEntity, err = s.ProfileProvider.CreateProfile(c, e.ID, s.PreePlanID, otp.Fullname, e.Email, e.AuthName, e.Avatar)
        if err != nil {
            return nil, err
        }

        e.UserID = infoEntity.ProfileID
        s.AuthMethodRepo.Update(c, e)
        isNewUser = true
    }
    // ... rest of logic ...
}
```

**Impact**:
- ✅ New users MUST provide fullname during OTP verification
- ✅ Existing users can login without fullname (backward compatible)
- ✅ Better data quality for new registrations
- ✅ Clear error message for missing fullname

**API Endpoint**: `POST /v2/auth/otp/verify`

**Request** (New User):
```json
{
  "phone": "0901234567",
  "otp": "123456",
  "authId": 456,
  "fullname": "Nguyen Van A"  // REQUIRED for new users
}
```

**Request** (Existing User):
```json
{
  "phone": "0901234567", 
  "otp": "123456",
  "authId": 123
  // fullname not required for existing users
}
```

**Error Response** (Missing Fullname for New User):
```json
{
  "code": 400,
  "message": "Vui lòng nhập họ tên khi đăng ký tài khoản mới"
}
```

**Frontend Logic**:
```typescript
// After OTP request, check existed field
const otpRequest = await requestOTP(phone, fullname);

if (otpRequest.existed) {
    // Existing user - fullname not required
    const verifyRequest = {
        phone,
        otp,
        authId: otpRequest.authId
    };
} else {
    // New user - fullname required
    const verifyRequest = {
        phone,
        otp, 
        authId: otpRequest.authId,
        fullname: fullname  // REQUIRED
    };
}

const result = await verifyOTP(verifyRequest);
```

**Testing**:
```bash
# Test new user without fullname (should fail)
curl -X POST http://localhost:8080/v2/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567", "otp": "123456", "authId": 456}'

# Expected: 400 error with message about missing fullname

# Test new user with fullname (should succeed)
curl -X POST http://localhost:8080/v2/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567", "otp": "123456", "authId": 456, "fullname": "Nguyen Van A"}'

# Expected: Success with access token
```

**Change Type**: Logic Enhancement
**Priority**: High
**Risk Level**: Low
**Backward Compatible**: Yes ✅
**Tested**: Pending
**Deployed**: Pending

---

## ✅ October 17, 2025 - OAuth Callback Enhancement

**Feature**: Return fullname from OAuth provider in callback response

**Affected Service**: Auth Service

**Changes Made**:

#### 1. DTO Update
**File**: `auth-service/internal/dto/auth.go`

```go
// AuthLoginResponse - Added FullName field
type AuthLoginResponse struct {
    AccessToken    string `json:"accessToken"`
    AuthID         uint64 `json:"authId"`
    Role           string `json:"role,omitempty"`
    Username       string `json:"username,omitempty"`
    ProfileID      uint64 `json:"profileId"`
    RefreshToken   string `json:"refreshToken"`
    OrganizationID uint64 `json:"organizationId,omitempty"`
    FullName       string `json:"fullName"` // NEW: OAuth fullname
}

// MakeAuthLoginResponse - Set fullname from OAuth
func MakeAuthLoginResponse(accessToken string, authEntity *domain.AuthMethod, role string) *AuthLoginResponse {
    return &AuthLoginResponse{
        AccessToken: accessToken,
        Role:      role,
        Username:  authEntity.AuthName,
        AuthID:    authEntity.ID,
        ProfileID: authEntity.UserID,
        FullName:  authEntity.FullName, // NEW: OAuth fullname
    }
}
```

#### 2. Handler Update
**File**: `auth-service/infra/handler/oauth_handler.go`

```go
func (h *OAuthHandler) HandleCallback(ctx context.Context, req *authpb.HandleCallbackRequest) (*authpb.HandleCallbackResponse, error) {
    var result *dto.AuthLoginResponse
    var err error

    switch req.Provider {
    case "zalo":
        result, err = h.ZaloAuthService.HandleCallback(ginCtx, req.Code, req.State)
    case "google":
        result, err = h.GoogleAuthService.HandleCallback(ginCtx, req.Code, req.State)
    case "facebook":
        result, err = h.FacebookAuthService.HandleCallback(ginCtx, req.Code, req.Platform, req.State)
    }

    return &authpb.HandleCallbackResponse{
        AuthId:       result.AuthID,
        ProfileId:    result.ProfileID,
        AccessToken:  result.AccessToken,
        RefreshToken: result.RefreshToken,
        FullName:     result.FullName, // NEW: Return fullname from OAuth
        Provider:     req.Provider,
    }, nil
}
```

**Impact**:
- ✅ Frontend receives fullname immediately after OAuth login
- ✅ No need for additional API call to get user info
- ✅ Better UX with instant profile display
- ✅ Consistent with OTP flow

**OAuth Providers Supported**:
- **Google**: Gets fullname from `userInfo["name"]`
- **Zalo**: Gets fullname from `userInfo["name"]`
- **Facebook**: Gets fullname from user profile

**API Endpoint**: `POST /v2/auth/oauth/callback`

**Request**:
```json
{
  "provider": "google",
  "code": "4/0AY0e-g7...",
  "state": "random_state"
}
```

**Response** (Before):
```json
{
  "authId": 123,
  "profileId": 456,
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "provider": "google"
}
```

**Response** (After):
```json
{
  "authId": 123,
  "profileId": 456,
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "fullName": "Nguyen Van A",  // ← NEW
  "provider": "google"
}
```

**Frontend Usage**:
```typescript
const oauthCallback = await handleOAuthCallback(provider, code, state);

// Immediately display user info
console.log(`Welcome, ${oauthCallback.fullName}!`);

// Save to state
setUser({
    authId: oauthCallback.authId,
    profileId: oauthCallback.profileId,
    fullName: oauthCallback.fullName,
    accessToken: oauthCallback.accessToken
});

// No need for additional API call
// Before: await getUserProfile(profileId)
// After: Already have fullName
```

**Testing**:
```bash
# Test Google OAuth
curl -X POST http://localhost:8080/v2/auth/oauth/callback \
  -H "Content-Type: application/json" \
  -d '{"provider": "google", "code": "4/0AY0e-g7...", "state": "random_state"}'

# Expected: Response includes fullName field

# Test Zalo OAuth
curl -X POST http://localhost:8080/v2/auth/oauth/callback \
  -H "Content-Type: application/json" \
  -d '{"provider": "zalo", "code": "abc123...", "state": "random_state"}'

# Expected: Response includes fullName from Zalo
```

**Change Type**: Feature Enhancement
**Priority**: Medium
**Risk Level**: Low
**Backward Compatible**: Yes ✅ (Adding new field)
**Tested**: Pending
**Deployed**: Pending

---

## ✅ October 17, 2025 - OTP Rate Limiting Update

**Feature**: Reduce OTP request rate limit from 5 minutes to 1 minute

**Affected Service**: Auth Service

**Changes Made**:

#### 1. Configuration Update
**Files**: 
- `auth-service/config/local.yml`
- `auth-service/config/develop.yml`

```yaml
otp:
  max-times-resend: 5
  max-times-otp-enter: 3 #số lần nhập OTP
  expired-after-seconds: 600 #10p hết hạn otp
  lock-verify-after: 900 # khóa tài khoản 15p nếu sai OTP/gửi OTP 5 lần
  lock-send-after: 60 # lần gửi OTP tiếp theo trong request hiện tại sau 60s
  lock-request-after: 60 # CHANGED: từ 300s (5p) → 60s (1p)
```

#### 2. Logic Explanation
**File**: `auth-service/internal/usecase/otp_usecase.go`

```go
case enums.LIMIT_REQUEST_TIME:
    if otpEntity.OTPDate != nil {
        nextTime := otpEntity.OTPDate.Add(time.Duration(s.properties.LockRequestAfter) * time.Second)
        if time.Now().Before(nextTime) {
            return _errors.ReturnError(
                int32(code),
                "Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau "+fmt.Sprintf("%.0f", time.Until(nextTime).Seconds())+"s",
            )
        }
    }
```

**Before**: `LockRequestAfter = 300 seconds (5 phút)`
**After**: `LockRequestAfter = 60 seconds (1 phút)`

#### 3. OTP Flow Comparison

**Before (5 minutes)**:
```
Time 00:00 - Request OTP ✅
Time 00:30 - Request OTP ❌ (Phải đợi 4.5 phút)
Time 05:00 - Request OTP ✅
Time 10:00 - Request OTP ✅
```

**After (1 minute)**:
```
Time 00:00 - Request OTP ✅
Time 00:30 - Request OTP ❌ (Phải đợi 30s)
Time 01:00 - Request OTP ✅
Time 02:00 - Request OTP ✅
```

**Impact**:
- ✅ Better UX: Users can request OTP more frequently
- ✅ Less frustration: 1 minute wait time vs 5 minutes
- ⚠️ Security consideration: More frequent requests allowed
- ✅ Still protected: Max 5 requests before account lock

**Rate Limiting Rules**:
1. **REQUEST OTP** (`POST /v2/auth/otp/request`): 
   - Minimum 1 minute between requests (changed from 5 minutes)
   - Error message: "Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau Xs"

2. **RESEND OTP** (`POST /v2/auth/otp/resend`):
   - Minimum 1 minute between resends (unchanged)
   - Error message: "Hãy thử lại sau Xs"

3. **MAX ATTEMPTS**:
   - Max 5 OTP send attempts
   - Max 3 OTP verification attempts
   - Account locked for 15 minutes after exceeding limits

**Security Measures Still in Place**:
- ✅ Max 5 OTP requests before account lock
- ✅ Max 3 verification attempts before account lock
- ✅ 15-minute account lock after exceeding limits
- ✅ OTP expires after 10 minutes
- ✅ Device and IP rate limiting

**API Error Messages**:

**Before (5 minutes)**:
```json
{
  "code": 1017,
  "message": "Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau 298s"
}
```

**After (1 minute)**:
```json
{
  "code": 1017,
  "message": "Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau 58s"
}
```

**Frontend Implementation**:
```typescript
const requestOTP = async (phone: string) => {
    try {
        const response = await api.post('/v2/auth/otp/request', { phone });
        return response.data;
    } catch (error) {
        if (error.code === 1017) {
            // Extract wait time from message
            const waitTime = extractWaitTime(error.message);
            
            // Show countdown timer
            showCountdown(waitTime);
            
            // Before: Wait up to 5 minutes
            // After: Wait up to 1 minute
        }
    }
};
```

**Testing**:
```bash
# Test 1: First OTP request
curl -X POST http://localhost:8080/v2/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567"}'
# Expected: Success

# Test 2: Immediate second request (should fail)
curl -X POST http://localhost:8080/v2/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567"}'
# Expected: Error 1017, wait ~60s

# Test 3: After 1 minute (should succeed)
sleep 60
curl -X POST http://localhost:8080/v2/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567"}'
# Expected: Success
```

**Deployment Notes**:
- ⚠️ Config change only - no code changes
- ✅ Service restart required to apply new config
- ✅ No database migration needed
- ✅ Backward compatible

**Recommendations**:
1. Monitor OTP request frequency after deployment
2. Consider adjusting based on abuse patterns
3. Add analytics to track user retry behavior
4. Consider A/B testing 1min vs 2min vs 3min

**Change Type**: Configuration Change
**Priority**: Medium
**Risk Level**: Low
**Backward Compatible**: Yes ✅
**Tested**: Pending
**Deployed**: Pending
**Restart Required**: Yes (Auth Service)

