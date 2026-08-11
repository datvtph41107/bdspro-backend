package _jwt

import (
	_callercontext "common/callercontext"

	"context"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrAuthenticationRequired = errors.New("authentication is required")
	ErrCallerClassification   = errors.New("caller classification is invalid")
)

/**
 * CallerClassificationInput contains already-authenticated edge evidence.
 * Principal must only be supplied after JWT verification succeeds.
 * APIKey must only be supplied after the API-key authority verifies the secret.
 */
type CallerClassificationInput struct {
	Principal      *Principal
	APIKey         *VerifiedAPIKey
	AllowAnonymous bool
}

/**
 * CallerResolution contains the non-secret identity produced at the HTTP edge.
 */
type CallerResolution struct {
	Caller    _callercontext.Caller
	APIKey    VerifiedAPIKey
	HasAPIKey bool
}

/**
 * ResolveHTTPCaller parses credentials, verifies API keys through their authority,
 * canonicalizes the compatibility header, and classifies the logical caller.
 */
func ResolveHTTPCaller(
	ctx context.Context,
	r *http.Request,
	principal *Principal,
	allowAnonymous bool,
	verifier APIKeyVerifier,
) (CallerResolution, error) {
	rawAPIKey, hasAPIKey, err := APIKeyFromRequest(r)
	if err != nil {
		return CallerResolution{}, err
	}

	var verifiedAPIKey *VerifiedAPIKey
	if hasAPIKey {
		if verifier == nil {
			return CallerResolution{}, ErrAPIKeyUnavailable
		}
		result, verifyErr := verifier.VerifyAPIKey(ctx, rawAPIKey)
		if verifyErr != nil {
			return CallerResolution{}, verifyErr
		}
		verifiedAPIKey = &result
	}

	caller, err := ClassifyCaller(CallerClassificationInput{
		Principal:      principal,
		APIKey:         verifiedAPIKey,
		AllowAnonymous: allowAnonymous,
	})
	if err != nil {
		return CallerResolution{}, err
	}

	resolution := CallerResolution{Caller: caller, HasAPIKey: hasAPIKey}
	if verifiedAPIKey != nil {
		resolution.APIKey = VerifiedAPIKey{
			ID:      verifiedAPIKey.ID,
			AppName: strings.TrimSpace(verifiedAPIKey.AppName),
		}
		if r != nil {
			r.Header.Del(LegacyAPIKeyHeader)
			r.Header.Set(APIKeyHeader, rawAPIKey)
		}
	}
	return resolution, nil
}

/**
 * ClassifyCaller projects verified edge evidence into a non-secret caller identity.
 */
func ClassifyCaller(input CallerClassificationInput) (_callercontext.Caller, error) {
	if input.Principal != nil && !ActorFromPrincipal(input.Principal).IsValid() {
		return _callercontext.Caller{}, ErrCallerClassification
	}

	verifiedAPIKey, hasAPIKey, err := normalizeVerifiedAPIKey(input.APIKey)
	if err != nil {
		return _callercontext.Caller{}, err
	}

	caller := _callercontext.Caller{}
	switch {
	case input.Principal != nil:
		caller.Kind = _callercontext.KindUser
	case hasAPIKey:
		caller.Kind = _callercontext.KindAPIKey
	case input.AllowAnonymous:
		caller.Kind = _callercontext.KindAnonymous
	default:
		return _callercontext.Caller{}, ErrAuthenticationRequired
	}

	if hasAPIKey {
		caller.APIKeyID = verifiedAPIKey.ID
		caller.APIKeyApp = verifiedAPIKey.AppName
	}
	if !caller.IsValid() {
		return _callercontext.Caller{}, ErrCallerClassification
	}
	return caller, nil
}

func normalizeVerifiedAPIKey(key *VerifiedAPIKey) (VerifiedAPIKey, bool, error) {
	if key == nil {
		return VerifiedAPIKey{}, false, nil
	}
	appName := strings.TrimSpace(key.AppName)
	if key.ID == 0 || appName == "" || len(appName) > 128 || containsControl(appName) {
		return VerifiedAPIKey{}, false, ErrCallerClassification
	}
	return VerifiedAPIKey{ID: key.ID, AppName: appName}, true, nil
}

func containsControl(value string) bool {
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}
